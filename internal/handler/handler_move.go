package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mythenox/bingetracker/internal/app"
	"github.com/mythenox/bingetracker/internal/database"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func HandlerMoveSeason(cmdContext context.Context, s *app.State, seasonNumber int,
	showTitle, seasonDirPath string) error {
	// warn user if episode count differs between new and old directories

	dirEntries, err := os.ReadDir(seasonDirPath)
	if err != nil {
		return err
	}

	// check if show already exists in database

	show, err := s.Q.GetShow(cmdContext, showTitle)
	if err != nil {
		return new(ShowNotFoundErr)
	}

	season, err := s.Q.GetSeason(cmdContext, database.GetSeasonParams{
		ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(SeasonNotFoundErr)
		}
		return err
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

	var overwrite bool

	if episodeCount == 1 {
		return errors.New("No valid files found in given directory")
	} else if episodeCount-1 != int(season.TotalEpisodes) {
		var userResponse string

		caser := cases.Title(language.English)

		fmt.Scanf(`WARNING: The number of episodes differs between the new and old
directories of Season %d of %s. Old episode viewing history must be overwritten
in order to use this new directory. Proceed? (Y/n): `, seasonNumber, caser.String(showTitle))
		fmt.Scanln(&userResponse)
		if strings.ToLower(userResponse) == "n" {
			fmt.Println("Move operation aborted.")
			return nil
		}
		overwrite = true
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	if overwrite {
		updatedSeason, err := qtx.UpdateSeason(cmdContext, database.UpdateSeasonParams{
			ShowTitle: season.ShowTitle, SeasonNumber: season.SeasonNumber,
			TotalEpisodes:     int64(episodeCount),
			UnwatchedEpisodes: int64(episodeCount),
		})
		if err != nil {
			return err
		}

		episodeDiff := updatedSeason.TotalEpisodes - season.TotalEpisodes

		_, err = qtx.UpdateShow(cmdContext, database.UpdateShowParams{
			ShowTitle: showTitle, TotalEpisodes: show.TotalEpisodes + episodeDiff,
			UnwatchedEpisodes: show.UnwatchedEpisodes - season.UnwatchedEpisodes + updatedSeason.TotalEpisodes,
		})

		for i := 1; i < episodeCount; i++ {
			_, err := qtx.OverwriteEpisode(cmdContext, database.OverwriteEpisodeParams{
				ShowTitle:     showTitle,
				SeasonNumber:  int64(seasonNumber),
				EpisodeNumber: int64(i),
				Runtime:       episodes[i].Runtime,
				Filepath:      episodes[i].Filepath,
			})
			if err != nil {
				return fmt.Errorf("error overwriting episode: %w", err)
			}
		}
	} else {
		for i := 1; i < episodeCount; i++ {
			err := qtx.UpdateEpisodePath(cmdContext, database.UpdateEpisodePathParams{
				ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
				EpisodeNumber: int64(i), Filepath: episodes[i].Filepath,
			})
			if err != nil {
				return fmt.Errorf("error updating episode filepath: %w", err)
			}
		}
	}

	caser := cases.Title(language.English)

	fmt.Printf("Successfully moved directory of season %d of %s.\n", seasonNumber, caser.String(showTitle))

	return tx.Commit()
}
