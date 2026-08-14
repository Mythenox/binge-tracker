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
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	episodes, err := s.DB.GetAllSeasonEpisodes(cmdContext, database.GetAllSeasonEpisodesParams{
		ShowTitle:    showTitle,
		SeasonNumber: int64(seasonNumber),
	})
	if err != nil {
		return err
	}

	if len(episodes) == 0 {
		_, err := s.DB.GetShow(cmdContext, showTitle)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return new(ShowNotInitializedErr)
			}
			return err
		}
		return errors.New("The provided season has not been added.")
	}

	for _, episode := range episodes {
		viewtime := formatDurationHMMSS(time.Duration(episode.Viewtime))
		runtime := formatDurationHMMSS(time.Duration(episode.Runtime))
		episodeIdentifier := formatEpisodeIdentifier(episode.SeasonNumber, episode.EpisodeNumber)

		if episode.Watched {
			line := fmt.Sprintf("%s\t%s\t%s/%s\t(watched)",
				episode.ShowTitle, episodeIdentifier, viewtime, runtime)
			fmt.Fprintln(w, line)
		} else {
			line := fmt.Sprintf("%s\t%s\t%s/%s",
				episode.ShowTitle, episodeIdentifier, viewtime, runtime)
			fmt.Fprintln(w, line)
		}
	}

	w.Flush()
	return nil
}

func formatDurationHMMSS(d time.Duration) string {
	t := time.Unix(0, 0).UTC().Add(d)
	return t.Format("04:05")
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

func handlerListSeasons() {}

func handlerListShows() {}
