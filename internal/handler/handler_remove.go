package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mythenox/bingetracker/internal/app"
	"github.com/mythenox/bingetracker/internal/database"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func HandlerRemoveSeason(cmdContext context.Context, s *app.State, showTitle string,
	seasonNumber int) error {
	// check if show already exists in database

	show, err := s.Q.GetShow(context.Background(), showTitle)
	if err != nil {
		return new(ShowNotFoundErr)
	}

	season, err := s.Q.GetSeason(cmdContext, database.GetSeasonParams{
		ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(ShowNotFoundErr)
		}
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	err = qtx.DeleteSeason(cmdContext, database.DeleteSeasonParams{
		ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
	})

	_, err = qtx.UpdateShow(cmdContext, database.UpdateShowParams{
		ShowTitle:         showTitle,
		Seasons:           show.Seasons - 1,
		TotalEpisodes:     show.TotalEpisodes - season.TotalEpisodes,
		UnwatchedEpisodes: show.UnwatchedEpisodes - season.UnwatchedEpisodes,
	})

	caser := cases.Title(language.English)

	fmt.Printf("Successfully removed season %d of %s.\n",
		seasonNumber, caser.String(showTitle))

	return tx.Commit()
}

// add episode, update corresponding total_episodes, unwatched_episodes, finished
// values of seasons, and total_episodes, unwatched_episodes values of shows
func HandlerRemoveEpisode(cmdContext context.Context, s *app.State,
	showTitle string, seasonNumber, episodeNumber int) error {
	// check if show is in database

	show, err := s.Q.GetShow(context.Background(), showTitle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(ShowNotFoundErr)
		}
		return err
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

	episode, err := s.Q.GetEpisode(cmdContext, database.GetEpisodeParams{
		ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
		EpisodeNumber: int64(episodeNumber),
	})

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	err = qtx.RemoveEpisode(cmdContext, database.RemoveEpisodeParams{
		ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
		EpisodeNumber: int64(episodeNumber),
	})

	var finished bool
	var seasonUnwatched, showUnwatched int64

	if episode.Watched {
		finished = season.UnwatchedEpisodes == season.TotalEpisodes-1
		seasonUnwatched = season.UnwatchedEpisodes
		showUnwatched = show.UnwatchedEpisodes
	} else {
		finished = season.Finished
		seasonUnwatched = season.UnwatchedEpisodes - 1
		showUnwatched = show.UnwatchedEpisodes - 1
	}

	_, err = s.Q.UpdateSeason(cmdContext, database.UpdateSeasonParams{
		ShowTitle:         showTitle,
		SeasonNumber:      int64(seasonNumber),
		TotalEpisodes:     season.TotalEpisodes - 1,
		UnwatchedEpisodes: seasonUnwatched,
		Finished:          finished,
	})

	_, err = s.Q.UpdateShow(cmdContext, database.UpdateShowParams{
		ShowTitle:         showTitle,
		Seasons:           show.Seasons,
		TotalEpisodes:     show.TotalEpisodes - 1,
		UnwatchedEpisodes: showUnwatched,
	})

	episodeIdentifier := formatEpisodeIdentifier(int64(seasonNumber), int64(episodeNumber))

	caser := cases.Title(language.English)

	fmt.Printf("Successfully removed %s of %s.\n", episodeIdentifier, caser.String(showTitle))

	return tx.Commit()
}

func HandlerRemoveShow(cmdContext context.Context, s *app.State, showTitle string) error {
	_, err := s.Q.GetShow(context.Background(), showTitle)
	if err != nil {
		return new(ShowNotFoundErr)
	}

	err = s.Q.RemoveShow(cmdContext, showTitle)
	if err != nil {
		return err
	}

	caser := cases.Title(language.English)

	fmt.Printf("Successfully removed %s.\n", caser.String(showTitle))

	return nil
}

func HandlerReset(cmdContext context.Context, s *app.State) error {
	shows, err := s.Q.GetAllShows(cmdContext)
	if err != nil {
		return err
	}

	if len(shows) == 0 {
		fmt.Println("No shows have been initialized.")
		return nil
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	for _, show := range shows {
		err = qtx.RemoveShow(cmdContext, show.Title)
		if err != nil {
			return err
		}
	}

	fmt.Println("Successfully reset database.")

	return tx.Commit()
}
