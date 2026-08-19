package handler

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/mythenox/bingetracker/internal/app"
	"github.com/mythenox/bingetracker/internal/database"
)

// make sure to require ffprobe, comes with mpv but not vlc

type ShowExistsError struct{}

func (*ShowExistsError) Error() string {
	return "The specified show has already been initialized"
}

// initiate a show to be tracked. requires show name, season to be added, and dirpath of the season.
// ex: bingetracker init "twin peaks" s01 videos/twin-peaks-s01
// automatically adds all episodes found in given directory

func HandlerInit(c context.Context,
	s *app.State,
	showTitle string,
	seasonNumber int,
	seasonDirPath string,
) error {
	dirEntries, err := os.ReadDir(seasonDirPath)
	if err != nil {
		return err
	}

	// check if show already exists in database

	_, err = s.Q.GetShow(context.Background(), showTitle)
	if err == nil {
		return new(ShowExistsError)
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

	_, err = s.Q.AddShow(context.Background(), database.AddShowParams{
		Title:         showTitle,
		TotalEpisodes: int64(episodeCount - 1),
	})
	if err != nil {
		return fmt.Errorf("error adding show to database: %w", err)
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

	fmt.Printf("Successfully added %d episodes of season %d of %s.\n", episodeCount-1, season.SeasonNumber, season.ShowTitle)

	return nil
}

// use ffprobe to get video runtime

func getVideoDuration(filePath string) (float64, error) {
	args := []string{
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	}

	cmd := exec.Command("ffprobe", args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0.0, fmt.Errorf("failed to run ffprobe: %w", err)
	}

	durationStr := strings.TrimSpace(out.String())
	runtime, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0.0, fmt.Errorf("failed to parse duration string: %w", err)
	}

	return runtime, nil
}
