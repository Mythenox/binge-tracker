package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/mythenox/binge-tracker/internal/app"
	"github.com/mythenox/binge-tracker/internal/database"
)

func HandlerSetSeasonFinished(cmdContext context.Context, s *app.State, showTitle string,
	seasonNumber int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	err = qtx.SetSeasonFinished(cmdContext, database.SetSeasonFinishedParams{
		ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
	})
	if err != nil {
		return err
	}

	err = qtx.SetSeasonEpisodesWatched(cmdContext, database.SetSeasonEpisodesWatchedParams{
		ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Season %d of %s has been set as finished.\n", seasonNumber, showTitle)
	return tx.Commit()
}

func HandlerSetSeasonRangeFinished(cmdContext context.Context, s *app.State, showTitle string,
	startSeasonNumber, endSeasonNumber int) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	for i := startSeasonNumber; i <= endSeasonNumber; i++ {
		err = qtx.SetSeasonFinished(cmdContext, database.SetSeasonFinishedParams{
			ShowTitle: showTitle, SeasonNumber: int64(i),
		})
		if err != nil {
			return err
		}

		err = qtx.SetSeasonEpisodesWatched(cmdContext, database.SetSeasonEpisodesWatchedParams{
			ShowTitle: showTitle, SeasonNumber: int64(i),
		})
		if err != nil {
			return err
		}
	}

	fmt.Printf("Seasons %d to %d of %s have been set as finished.\n", startSeasonNumber,
		endSeasonNumber, showTitle)
	return tx.Commit()
}

func HandlerSetEpisodeWatched(cmdContext context.Context, s *app.State, showTitle string,
	seasonNumber, episodeNumber int) error {
	err := s.Q.SetEpisodeWatched(cmdContext, database.SetEpisodeWatchedParams{
		ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
		EpisodeNumber: int64(episodeNumber),
	})
	if err != nil {
		return err
	}

	fmt.Printf("Episode %d of season %d of %s has been set as watched.\n", episodeNumber,
		seasonNumber, showTitle)

	return nil
}

func HandlerSetEpisodeRangeWatched(cmdContext context.Context, s *app.State, showTitle string,
	startNums, endNums []int) error {
	startSeasonNumber, startEpisodeNumber := startNums[0], startNums[1]
	endSeasonNumber, endEpisodeNumber := endNums[0], endNums[1]

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	if startSeasonNumber == endSeasonNumber {
		if startEpisodeNumber > endEpisodeNumber {
			return errors.New("Invalid episode range")
		}
		for i := startEpisodeNumber; i <= endEpisodeNumber; i++ {
			err = qtx.SetEpisodeWatched(cmdContext, database.SetEpisodeWatchedParams{
				ShowTitle: showTitle, SeasonNumber: int64(endSeasonNumber),
				EpisodeNumber: int64(i),
			})
			if err != nil {
				return err
			}
		}
	} else {
		if startSeasonNumber > endSeasonNumber {
			return errors.New("Invalid season range")
		}
		// just set full season as finished until i == endSeasonNumber
		for i := startSeasonNumber; i < endSeasonNumber; i++ {

			err = qtx.SetSeasonFinished(cmdContext, database.SetSeasonFinishedParams{
				ShowTitle: showTitle, SeasonNumber: int64(i),
			})
			if err != nil {
				return err
			}

			err = qtx.SetSeasonEpisodesWatched(cmdContext, database.SetSeasonEpisodesWatchedParams{
				ShowTitle: showTitle, SeasonNumber: int64(i),
			})
			if err != nil {
				return err
			}
		}

		for i := 1; i <= endEpisodeNumber; i++ {
			err = qtx.SetEpisodeWatched(cmdContext, database.SetEpisodeWatchedParams{
				ShowTitle: showTitle, SeasonNumber: int64(endSeasonNumber),
				EpisodeNumber: int64(i),
			})
			if err != nil {
				return err
			}
		}
	}

	start := formatEpisodeIdentifier(int64(startSeasonNumber), int64(startEpisodeNumber))
	end := formatEpisodeIdentifier(int64(endSeasonNumber), int64(endEpisodeNumber))

	fmt.Printf("All episodes between %s and %s of %s have been set as watched.\n", start,
		end, showTitle)
	return tx.Commit()
}
