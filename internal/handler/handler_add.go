package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/mythenox/bingetracker/internal/app"
	"github.com/mythenox/bingetracker/internal/database"
)

type ShowNotFoundError struct{}

func (*ShowNotFoundError) Error() string {
	return "The specified show was not found in the database."
}

type SeasonNotFoundError struct{}

func (*SeasonNotFoundError) Error() string {
	return "The specified season was not found in the database."
}

type InvalidFileTypeError struct{}

func (*InvalidFileTypeError) Error() string {
	return "The specified file has an unsupported file type."
}

// add season, update corresponding season, total_episodes, unwatched_episodes values of show
func HandlerAddSeason(cmdContext context.Context, s *app.State, seasonNumber int,
	showTitle, seasonDirPath string) error {
	dirEntries, err := os.ReadDir(seasonDirPath)
	if err != nil {
		return err
	}

	// check if show already exists in database

	show, err := s.Q.GetShow(cmdContext, showTitle)
	if err != nil {
		return new(ShowNotFoundError)
	}

	var episodes []database.AddEpisodeParams
	supportedFormats := []string{".mkv", ".mp4", ".avi", ".mov", ".webm"}
	episodeCount := 1
	absDirPath, err := filepath.Abs(seasonDirPath)
	if err != nil {
		return err
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			continue
		}

		episodePath := filepath.Join(absDirPath, dirEntry.Name())

		ext := filepath.Ext(episodePath)
		if len(ext) == 0 || !slices.Contains(supportedFormats, ext) {
			continue
		}

		runtime, err := getVideoDuration(episodePath)
		if err != nil {
			return err
		}

		episode := database.AddEpisodeParams{
			EpisodeNumber: int64(episodeCount),
			Filepath:      episodePath,
			Runtime:       runtime,
			ShowTitle:     showTitle,
			SeasonNumber:  int64(seasonNumber),
		}

		episodes = append(episodes, episode)
		episodeCount++
	}

	if episodeCount == 1 {
		return errors.New("No valid files found in given directory")
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	newSeason, err := qtx.AddSeason(cmdContext, database.AddSeasonParams{
		SeasonNumber:  int64(seasonNumber),
		TotalEpisodes: int64(episodeCount - 1),
		ShowTitle:     showTitle,
	})
	if err != nil {
		return fmt.Errorf("error adding season to database: %w", err)
	}

	for _, episode := range episodes {
		_, err := qtx.AddEpisode(cmdContext, episode)
		if err != nil {
			return fmt.Errorf("error adding episode to database: %w", err)
		}
	}

	_, err = qtx.UpdateShow(cmdContext, database.UpdateShowParams{
		ShowTitle:         show.Title,
		TotalEpisodes:     show.TotalEpisodes + newSeason.TotalEpisodes,
		UnwatchedEpisodes: show.UnwatchedEpisodes + newSeason.TotalEpisodes,
		Seasons:           show.Seasons + 1,
	})

	fmt.Printf("Successfully added %d episodes of season %d of %s.\n",
		newSeason.TotalEpisodes, newSeason.SeasonNumber, newSeason.ShowTitle)

	return tx.Commit()
}

// add episode, update corresponding total_episodes, unwatched_episodes, finished
// values of seasons, and total_episodes, unwatched_episodes values of shows
func HandlerAddEpisode(cmdContext context.Context, s *app.State,
	seasonNumber, episodeNumber int, showTitle, episodePath string) error {
	// check if show is in database

	show, err := s.Q.GetShow(cmdContext, showTitle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(ShowNotFoundError)
		}
		return err
	}

	season, err := s.Q.GetSeason(cmdContext, database.GetSeasonParams{
		ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(SeasonNotFoundError)
		}
		return err
	}

	supportedFormats := []string{".mkv", ".mp4", ".avi", ".mov", ".webm"}

	ext := filepath.Ext(episodePath)
	if len(ext) == 0 || !slices.Contains(supportedFormats, ext) {
		return new(InvalidFileTypeError)
	}

	runtime, err := getVideoDuration(episodePath)
	if err != nil {
		return err
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	_, err = qtx.AddEpisode(cmdContext, database.AddEpisodeParams{
		EpisodeNumber: int64(episodeNumber),
		Filepath:      episodePath,
		Runtime:       runtime,
		ShowTitle:     showTitle,
		SeasonNumber:  int64(seasonNumber),
	})
	if err != nil {
		return fmt.Errorf("error adding episode to database: %w", err)
	}

	_, err = qtx.UpdateSeason(cmdContext, database.UpdateSeasonParams{
		ShowTitle:         showTitle,
		SeasonNumber:      int64(seasonNumber),
		TotalEpisodes:     season.TotalEpisodes + 1,
		UnwatchedEpisodes: season.UnwatchedEpisodes + 1,
		Finished:          false,
	})

	_, err = qtx.UpdateShow(cmdContext, database.UpdateShowParams{
		ShowTitle:         showTitle,
		Seasons:           show.Seasons,
		TotalEpisodes:     show.TotalEpisodes + 1,
		UnwatchedEpisodes: show.UnwatchedEpisodes + 1,
	})

	episodeIdentifier := formatEpisodeIdentifier(int64(seasonNumber), int64(episodeNumber))

	fmt.Printf("Successfully added %s of %s.\n", episodeIdentifier, showTitle)

	return tx.Commit()
}
