package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/mythenox/binge-tracker/internal/app"
	"github.com/mythenox/binge-tracker/internal/database"
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

func HandlerAddSeason(cmdContext context.Context, s *app.State, seasonNumber int,
	showTitle, seasonDirPath string) error {
	dirEntries, err := os.ReadDir(seasonDirPath)
	if err != nil {
		return err
	}

	// check if show already exists in database

	_, err = s.Q.GetShow(context.Background(), showTitle)
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

	season, err := s.Q.AddSeason(context.Background(), database.AddSeasonParams{
		SeasonNumber:  int64(seasonNumber),
		TotalEpisodes: int64(episodeCount - 1),
		ShowTitle:     showTitle,
	})
	if err != nil {
		return fmt.Errorf("error adding season to database: %w", err)
	}

	for _, episode := range episodes {
		_, err := s.Q.AddEpisode(context.Background(), episode)
		if err != nil {
			return fmt.Errorf("error adding episode to database: %w", err)
		}
	}

	fmt.Printf("Successfully added %d episodes of season %d of %s.\n", season.TotalEpisodes, season.SeasonNumber, season.ShowTitle)

	return nil
}

func HandlerAddEpisode(cmdContext context.Context, s *app.State,
	seasonNumber, episodeNumber int, showTitle, episodePath string) error {
	// check if show is in database

	_, err := s.Q.GetShow(context.Background(), showTitle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(ShowNotFoundError)
		}
		return err
	}

	_, err = s.Q.GetSeason(cmdContext, database.GetSeasonParams{
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

	_, err = s.Q.AddEpisode(context.Background(), database.AddEpisodeParams{
		EpisodeNumber: int64(episodeNumber),
		Filepath:      episodePath,
		Runtime:       runtime,
		ShowTitle:     showTitle,
		SeasonNumber:  int64(seasonNumber),
	})
	if err != nil {
		return fmt.Errorf("error adding episode to database: %w", err)
	}

	episodeIdentifier := formatEpisodeIdentifier(int64(seasonNumber), int64(episodeNumber))

	fmt.Printf("Successfully added %s of %s.\n", episodeIdentifier, showTitle)

	return nil
}
