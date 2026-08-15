package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/mythenox/binge-tracker/internal/app"
	"github.com/mythenox/binge-tracker/internal/database"
)

type ShowNotInitializedErr struct{}

func (*ShowNotInitializedErr) Error() string {
	return "The provided show has not been initialized."
}

func HandlerListEpisodes(
	cmdContext context.Context,
	s *app.State,
	showTitle string,
	seasonNumber int,
) error {

	episodes, err := s.Q.GetAllEpisodesFromSeason(cmdContext, database.GetAllEpisodesFromSeasonParams{
		ShowTitle:    showTitle,
		SeasonNumber: int64(seasonNumber),
	})
	if err != nil {
		return err
	}

	if len(episodes) == 0 {
		_, err := s.Q.GetShow(cmdContext, showTitle)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return new(ShowNotInitializedErr)
			}
			return err
		}
		return errors.New("The provided season has not been added.")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	for _, episode := range episodes {
		viewtime, err := formatDurationHMMSS(int64(episode.Viewtime))
		if err != nil {
			return fmt.Errorf("Unable to format episode viewtime: %w", err)
		}

		runtime, err := formatDurationHMMSS(int64(episode.Runtime))
		if err != nil {
			return fmt.Errorf("Unable to format episode runtime: %w", err)
		}

		episodeIdentifier := formatEpisodeIdentifier(episode.SeasonNumber, episode.EpisodeNumber)

		if episode.Watched {
			line := fmt.Sprintf("%s\t%s/%s\t(watched)",
				episodeIdentifier, viewtime, runtime)
			fmt.Fprintln(w, line)
		} else {
			line := fmt.Sprintf("%s\t%s/%s",
				episodeIdentifier, viewtime, runtime)
			fmt.Fprintln(w, line)
		}
	}

	fmt.Printf("%s Season %d episodes:\n", showTitle, seasonNumber)
	w.Flush()

	return nil
}

func HandlerListSeasons(cmdContext context.Context, s *app.State, showTitle string) error {
	// TODO: add verbose flag to allow user to see directory of season?
	seasons, err := s.Q.GetSeasonsFromShow(cmdContext, showTitle)
	if err != nil {
		return err
	}

	if len(seasons) == 0 {
		return new(ShowNotInitializedErr)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	for _, season := range seasons {
		seasonIdentifier := formatSeasonIdentifier(season.SeasonNumber)
		watchedEpisodes := season.TotalEpisodes - season.UnwatchedEpisodes

		line := fmt.Sprintf("%s\t%d/%d eps watched",
			seasonIdentifier, watchedEpisodes, season.TotalEpisodes)

		fmt.Fprintln(w, line)
	}

	fmt.Printf("%s seasons:\n", showTitle)
	w.Flush()

	return nil
}

func HandlerListShows(cmdContext context.Context, s *app.State) error {
	shows, err := s.Q.GetAllShows(cmdContext)
	if err != nil {
		return err
	}

	if len(shows) == 0 {
		return errors.New("No shows have been initialized.")
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	// VDM 2 seasons 1/300 eps watched
	for _, show := range shows {
		watchedEpisodes := show.TotalEpisodes - show.UnwatchedEpisodes

		line := fmt.Sprintf("%s\t%d season(s) added\t%d/%d eps watched",
			show.Title, show.Seasons, watchedEpisodes, show.TotalEpisodes)

		fmt.Fprintln(w, line)
	}

	w.Flush()
	fmt.Println("\nFor more information about a show's seasons, try 'bingetracker list seasons <show title>'.")

	return nil
}

func formatDurationHMMSS(seconds int64) (string, error) {
	durationString := fmt.Sprintf("%ds", seconds)

	d, err := time.ParseDuration(durationString)
	if err != nil {
		return "", err
	}

	t := time.Unix(0, 0).UTC().Add(d)

	return t.Format("04:05"), nil
}

func formatEpisodeIdentifier(seasonNumber, episodeNumber int64) string {
	// ex: (4,5) -> s04e05
	var season, episode string

	if seasonNumber < 10 {
		season = fmt.Sprintf("s0%d", seasonNumber)
	} else {
		season = fmt.Sprintf("s%d", seasonNumber)
	}

	if episodeNumber < 10 {
		episode = fmt.Sprintf("e0%d", episodeNumber)
	} else {
		episode = fmt.Sprintf("e%d", episodeNumber)
	}

	episodeIdentifier := season + episode
	return episodeIdentifier
}

func formatSeasonIdentifier(seasonNumber int64) string {
	// ex: 4 -> s04
	var season string

	if seasonNumber < 10 {
		season = fmt.Sprintf("s0%d", seasonNumber)
	} else {
		season = fmt.Sprintf("s%d", seasonNumber)
	}

	return season
}
