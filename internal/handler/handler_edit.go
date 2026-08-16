package handler

import (
	"context"
	"errors"
	"fmt"

	"github.com/mythenox/binge-tracker/internal/app"
	"github.com/mythenox/binge-tracker/internal/database"
)

func HandlerSetEpisodeRangeCompletion(cmdContext context.Context, s *app.State, showTitle string,
	startNums, endNums []int, setWatched bool) error {
	// watched == true -> set watched, == false -> set unwatched

	startSeasonNumber, startEpisodeNumber := startNums[0], startNums[1]
	endSeasonNumber, endEpisodeNumber := endNums[0], endNums[1]

	var setEpisodeCompletion func(q *database.Queries, ctx context.Context,
		showTitle string, seasonNumber int64, episodeNumber int64) error
	var setEpisodesForSeasonCompletion func(q *database.Queries, ctx context.Context,
		showTitle string, seasonNumber int64) error
	var setSeasonCompletion func(q *database.Queries, ctx context.Context,
		showTitle string, seasonNumber int64) error

	if setWatched {
		setEpisodeCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber, episodeNumber int64) error {
			return q.SetEpisodeWatched(ctx, database.SetEpisodeWatchedParams{
				ShowTitle:     showTitle,
				SeasonNumber:  seasonNumber,
				EpisodeNumber: episodeNumber,
			})
		}
		setEpisodesForSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.SetEpisodesForSeasonWatched(ctx, database.SetEpisodesForSeasonWatchedParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
		setSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.SetSeasonFinished(ctx, database.SetSeasonFinishedParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
	} else {
		setEpisodeCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber, episodeNumber int64) error {
			return q.ResetEpisodeStats(ctx, database.ResetEpisodeStatsParams{
				ShowTitle:     showTitle,
				SeasonNumber:  seasonNumber,
				EpisodeNumber: episodeNumber,
			})
		}
		setEpisodesForSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.ResetEpisodeStatsForSeason(ctx, database.ResetEpisodeStatsForSeasonParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
		setSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.ResetSeasonStats(ctx, database.ResetSeasonStatsParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
	}

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
			err = setEpisodeCompletion(qtx, cmdContext, showTitle,
				int64(endSeasonNumber), int64(i))
			if err != nil {
				return err
			}
		}
	} else {
		if startSeasonNumber > endSeasonNumber {
			return errors.New("Invalid season range")
		}
		// just set full season as finished/unfinished until i == endSeasonNumber
		for i := startSeasonNumber; i < endSeasonNumber; i++ {

			err = setSeasonCompletion(qtx, cmdContext, showTitle, int64(i))
			if err != nil {
				return err
			}

			err = setEpisodesForSeasonCompletion(qtx, cmdContext, showTitle, int64(i))
			if err != nil {
				return err
			}
		}

		for i := 1; i <= endEpisodeNumber; i++ {
			err = setEpisodeCompletion(qtx, cmdContext, showTitle,
				int64(endSeasonNumber), int64(i))
			if err != nil {
				return err
			}
		}
	}

	start := formatEpisodeIdentifier(int64(startSeasonNumber), int64(startEpisodeNumber))
	end := formatEpisodeIdentifier(int64(endSeasonNumber), int64(endEpisodeNumber))

	if setWatched {
		fmt.Printf("All episodes between %s and %s of %s have been set as watched.\n", start,
			end, showTitle)
	} else {
		fmt.Printf("All episodes between %s and %s of %s have been set as unwatched.\n", start,
			end, showTitle)
	}

	return tx.Commit()
}

func HandlerSetEpisodeCompletion(cmdContext context.Context, s *app.State, showTitle string,
	seasonNumber, episodeNumber int, setWatched bool) error {

	var setEpisodeCompletion func(q *database.Queries, ctx context.Context,
		showTitle string, seasonNumber int64, episodeNumber int64) error

	if setWatched {
		setEpisodeCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber, episodeNumber int64) error {
			return q.SetEpisodeWatched(ctx, database.SetEpisodeWatchedParams{
				ShowTitle:     showTitle,
				SeasonNumber:  seasonNumber,
				EpisodeNumber: episodeNumber,
			})
		}
	} else {
		setEpisodeCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber, episodeNumber int64) error {
			return q.ResetEpisodeStats(ctx, database.ResetEpisodeStatsParams{
				ShowTitle:     showTitle,
				SeasonNumber:  seasonNumber,
				EpisodeNumber: episodeNumber,
			})
		}
	}

	err := setEpisodeCompletion(s.Q, cmdContext, showTitle,
		int64(seasonNumber), int64(episodeNumber))
	if err != nil {
		return err
	}

	if setWatched {
		fmt.Printf("Episode %d of season %d of %s has been set as watched.\n", episodeNumber,
			seasonNumber, showTitle)
	} else {
		fmt.Printf("Episode %d of season %d of %s has been set as unwatched.\n", episodeNumber,
			seasonNumber, showTitle)
	}

	return nil
}

func HandlerSetSeasonRangeCompletion(cmdContext context.Context, s *app.State, showTitle string,
	startSeasonNumber, endSeasonNumber int, setWatched bool) error {
	var setEpisodesForSeasonCompletion func(q *database.Queries, ctx context.Context,
		showTitle string, seasonNumber int64) error
	var setSeasonCompletion func(q *database.Queries, ctx context.Context,
		showTitle string, seasonNumber int64) error

	if setWatched {
		setEpisodesForSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.SetEpisodesForSeasonWatched(ctx, database.SetEpisodesForSeasonWatchedParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
		setSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.SetSeasonFinished(ctx, database.SetSeasonFinishedParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
	} else {
		setEpisodesForSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.ResetEpisodeStatsForSeason(ctx, database.ResetEpisodeStatsForSeasonParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
		setSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.ResetSeasonStats(ctx, database.ResetSeasonStatsParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	for i := startSeasonNumber; i <= endSeasonNumber; i++ {
		err = setSeasonCompletion(qtx, cmdContext, showTitle, int64(i))
		if err != nil {
			return err
		}

		err = setEpisodesForSeasonCompletion(qtx, cmdContext, showTitle, int64(i))
		if err != nil {
			return err
		}
	}

	if setWatched {
		fmt.Printf("Seasons %d to %d of %s have been set as finished.\n", startSeasonNumber,
			endSeasonNumber, showTitle)
	} else {
		fmt.Printf("Seasons %d to %d of %s have been set as unfinished.\n", startSeasonNumber,
			endSeasonNumber, showTitle)
	}

	return tx.Commit()
}

func HandlerSetSeasonCompletion(cmdContext context.Context, s *app.State, showTitle string,
	seasonNumber int, setWatched bool) error {
	var setEpisodesForSeasonCompletion func(q *database.Queries, ctx context.Context,
		showTitle string, seasonNumber int64) error
	var setSeasonCompletion func(q *database.Queries, ctx context.Context,
		showTitle string, seasonNumber int64) error

	if setWatched {
		setEpisodesForSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.SetEpisodesForSeasonWatched(ctx, database.SetEpisodesForSeasonWatchedParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
		setSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.SetSeasonFinished(ctx, database.SetSeasonFinishedParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
	} else {
		setEpisodesForSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.ResetEpisodeStatsForSeason(ctx, database.ResetEpisodeStatsForSeasonParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
		setSeasonCompletion = func(q *database.Queries, ctx context.Context, showTitle string,
			seasonNumber int64) error {
			return q.ResetSeasonStats(ctx, database.ResetSeasonStatsParams{
				ShowTitle:    showTitle,
				SeasonNumber: seasonNumber,
			})
		}
	}

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	err = setSeasonCompletion(qtx, cmdContext, showTitle, int64(seasonNumber))
	if err != nil {
		return err
	}

	err = setEpisodesForSeasonCompletion(qtx, cmdContext, showTitle, int64(seasonNumber))
	if err != nil {
		return err
	}

	if setWatched {
		fmt.Printf("Season %d of %s has been set as finished.\n", seasonNumber, showTitle)
	} else {
		fmt.Printf("Season %d of %s has been set as unfinished.\n", seasonNumber, showTitle)
	}

	return tx.Commit()
}
