package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mythenox/bingetracker/internal/app"
	"github.com/mythenox/bingetracker/internal/database"
)

// s01e01-s02e06

func HandlerSetEpisodeRangeCompletion(cmdContext context.Context, s *app.State, showTitle string,
	startNums, endNums []int, setWatched bool) error {
	// watched == true -> set watched, == false -> set unwatched

	startSeason, startEpisode := startNums[0], startNums[1]
	endSeason, endEpisode := endNums[0], endNums[1]

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	var totalUpdatedEpisodeCount int64

	for i := startSeason; i <= endSeason; i++ {
		season, err := qtx.GetSeason(cmdContext, database.GetSeasonParams{
			ShowTitle: showTitle, SeasonNumber: int64(i),
		})
		if err != nil {
			return err
		}

		e := database.GetEpisodeParams{
			ShowTitle:    showTitle,
			SeasonNumber: int64(i),
		}

		var updates int64

		switch i {
		case startSeason:
			if startSeason == endSeason {
				updates, err = episodeLoop(cmdContext, qtx, e,
					int64(startEpisode), int64(endEpisode), setWatched)
			} else {
				updates, err = episodeLoop(cmdContext, qtx, e,
					int64(startEpisode), season.TotalEpisodes, setWatched)
			}
		case endSeason:
			updates, err = episodeLoop(cmdContext, qtx, e,
				1, int64(endEpisode), setWatched)
		default:
			updates, err = episodeLoop(cmdContext, qtx, e,
				1, season.TotalEpisodes, setWatched)
		}

		if err != nil {
			return err
		}

		err = updateSeasonWatchStatus(cmdContext, qtx, season,
			updates, setWatched)
		if err != nil {
			return err
		}

		totalUpdatedEpisodeCount += updates
	}

	err = updateShowWatchStatus(cmdContext, qtx, showTitle,
		totalUpdatedEpisodeCount, setWatched)
	if err != nil {
		return err
	}

	start := formatEpisodeIdentifier(int64(startSeason), int64(startEpisode))
	end := formatEpisodeIdentifier(int64(endSeason), int64(endEpisode))

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
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	e := database.GetEpisodeParams{ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
		EpisodeNumber: int64(episodeNumber)}

	updated, err := updateEpisodeWatchStatus(cmdContext, qtx,
		e, setWatched)
	if err != nil {
		return err
	}

	if updated {
		err := cascadeWatchStatus(cmdContext, qtx, showTitle, int64(seasonNumber), 1, setWatched)
		if err != nil {
			return err
		}
	}

	if setWatched {
		fmt.Printf("Episode %d of season %d of %s has been set as watched.\n", episodeNumber,
			seasonNumber, showTitle)
	} else {
		fmt.Printf("Episode %d of season %d of %s has been set as unwatched.\n", episodeNumber,
			seasonNumber, showTitle)
	}

	return tx.Commit()
}

func HandlerSetSeasonRangeCompletion(cmdContext context.Context, s *app.State, showTitle string,
	startSeason, endSeason int, setWatched bool) error {

	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	var totalUpdatedEpisodeCount int64

	for i := startSeason; i <= endSeason; i++ {
		season, err := qtx.GetSeason(cmdContext, database.GetSeasonParams{
			ShowTitle: showTitle, SeasonNumber: int64(i),
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return new(SeasonNotFoundError)
			}
			return err
		}

		e := database.GetEpisodeParams{ShowTitle: showTitle, SeasonNumber: int64(i)}

		updates, err := episodeLoop(cmdContext, qtx, e, 1, season.TotalEpisodes, setWatched)
		if err != nil {
			return err
		}

		err = updateSeasonWatchStatus(cmdContext, qtx, season, updates, setWatched)
		if err != nil {
			return err
		}

		totalUpdatedEpisodeCount += updates
	}

	err = updateShowWatchStatus(cmdContext, qtx, showTitle, totalUpdatedEpisodeCount, setWatched)
	if err != nil {
		return err
	}

	if setWatched {
		fmt.Printf("Seasons %d to %d of %s have been set as finished.\n", startSeason,
			endSeason, showTitle)
	} else {
		fmt.Printf("Seasons %d to %d of %s have been set as unfinished.\n", startSeason,
			endSeason, showTitle)
	}

	return tx.Commit()
}

func HandlerSetSeasonCompletion(cmdContext context.Context, s *app.State, showTitle string,
	seasonNumber int, setWatched bool) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := s.Q.WithTx(tx)

	season, err := qtx.GetSeason(cmdContext, database.GetSeasonParams{
		ShowTitle: showTitle, SeasonNumber: int64(seasonNumber),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(SeasonNotFoundError)
		}
		return err
	}

	e := database.GetEpisodeParams{ShowTitle: showTitle, SeasonNumber: int64(seasonNumber)}

	updates, err := episodeLoop(cmdContext, qtx, e, 1, season.TotalEpisodes, setWatched)
	if err != nil {
		return err
	}

	err = updateSeasonWatchStatus(cmdContext, qtx, season, updates, setWatched)
	if err != nil {
		return err
	}

	err = updateShowWatchStatus(cmdContext, qtx, showTitle, updates, setWatched)
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

func updateEpisodeWatchStatus(
	cmdContext context.Context, qtx *database.Queries,
	e database.GetEpisodeParams, setWatched bool) (bool, error) {
	episode, err := qtx.GetEpisode(cmdContext, e)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			episodeIdentifier := formatEpisodeIdentifier(e.SeasonNumber,
				e.EpisodeNumber)
			return false, fmt.Errorf("The episode %s of %s does was not found in the database.",
				episodeIdentifier, e.ShowTitle)
		}
		return false, err
	}
	// if episode.watched == true and setWatched == true, nothing needs to be done
	// same goes for if episode.watched == false and setWatched == false
	if episode.Watched == setWatched {
		return false, nil
	}

	var viewtime float64

	if setWatched {
		viewtime = episode.Runtime
	} else {
		viewtime = 0.0
	}

	_, err = qtx.UpdateEpisodeStats(cmdContext, database.UpdateEpisodeStatsParams{
		ShowTitle: e.ShowTitle, SeasonNumber: e.SeasonNumber,
		EpisodeNumber: e.EpisodeNumber,
		Viewtime:      viewtime,
		Watched:       setWatched,
	})
	if err != nil {
		return false, err
	}

	return true, nil
}

func episodeLoop(cmdContext context.Context, qtx *database.Queries,
	e database.GetEpisodeParams, startEpisodeNum, endEpisodeNum int64,
	setWatched bool) (int64, error) {
	var updates int64

	for i := startEpisodeNum; i <= endEpisodeNum; i++ {
		e.EpisodeNumber = i
		updated, err := updateEpisodeWatchStatus(cmdContext, qtx,
			e, setWatched)
		if err != nil {
			return 0, err
		}

		if updated {
			updates++
		}
	}

	return updates, nil
}

func updateSeasonWatchStatus(cmdContext context.Context, qtx *database.Queries,
	season database.Season, updatedEpisodeCount int64, setWatched bool) error {
	if setWatched {
		finished := season.UnwatchedEpisodes-updatedEpisodeCount <= 0
		_, err := qtx.UpdateSeason(cmdContext, database.UpdateSeasonParams{
			ShowTitle:         season.ShowTitle,
			SeasonNumber:      season.SeasonNumber,
			TotalEpisodes:     season.TotalEpisodes,
			UnwatchedEpisodes: season.UnwatchedEpisodes - updatedEpisodeCount,
			Finished:          finished,
		})
		if err != nil {
			return err
		}
	} else {
		finished := season.UnwatchedEpisodes+updatedEpisodeCount > 0
		_, err := qtx.UpdateSeason(cmdContext, database.UpdateSeasonParams{
			ShowTitle:         season.ShowTitle,
			SeasonNumber:      season.SeasonNumber,
			TotalEpisodes:     season.TotalEpisodes,
			UnwatchedEpisodes: season.UnwatchedEpisodes + updatedEpisodeCount,
			Finished:          finished,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func updateShowWatchStatus(cmdContext context.Context, qtx *database.Queries, showTitle string,
	updatedEpisodeCount int64, setWatched bool) error {
	show, err := qtx.GetShow(cmdContext, showTitle)
	if err != nil {
		return err
	}

	if setWatched {
		_, err = qtx.UpdateShow(cmdContext, database.UpdateShowParams{
			ShowTitle:         showTitle,
			Seasons:           show.Seasons,
			TotalEpisodes:     show.TotalEpisodes,
			UnwatchedEpisodes: show.UnwatchedEpisodes - updatedEpisodeCount,
		})
		if err != nil {
			return err
		}
	} else {
		_, err = qtx.UpdateShow(cmdContext, database.UpdateShowParams{
			ShowTitle:         showTitle,
			Seasons:           show.Seasons,
			TotalEpisodes:     show.TotalEpisodes,
			UnwatchedEpisodes: show.UnwatchedEpisodes + updatedEpisodeCount,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func cascadeWatchStatus(cmdContext context.Context, qtx *database.Queries, showTitle string,
	seasonNumber, updatedEpisodeCount int64, setWatched bool) error {
	season, err := qtx.GetSeason(cmdContext, database.GetSeasonParams{
		ShowTitle: showTitle, SeasonNumber: seasonNumber,
	})
	if err != nil {
		return err
	}

	err = updateSeasonWatchStatus(cmdContext, qtx, season, updatedEpisodeCount, setWatched)
	if err != nil {
		return err
	}

	err = updateShowWatchStatus(cmdContext, qtx, showTitle,
		updatedEpisodeCount, setWatched)
	if err != nil {
		return err
	}

	return nil
}
