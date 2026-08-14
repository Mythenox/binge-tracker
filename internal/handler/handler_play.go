package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mythenox/binge-tracker/internal/app"
	"github.com/mythenox/binge-tracker/internal/database"
	"github.com/mythenox/binge-tracker/internal/listen"
)

// play [mpv, vlc] <show_name> <season_number> <episode_number> [optional extra mpv/vlc flags]

// RunE func(cmd *Command, args []string) error

type SeasonWatchedError struct{}

func (*SeasonWatchedError) Error() string {
	return "All episodes of this season have been watched."
}

type ShowWatchedError struct{}

func (*ShowWatchedError) Error() string {
	return "All episodes of this show have been watched."
}

type PlayNextArgs struct {
	VideoPlayer  string
	ShowTitle    string
	SeasonNumber int
}

func HandlerPlayNext(
	cmdContext context.Context,
	s *app.State,
	btArgs PlayNextArgs,
	skipInProgress bool,
	playerArgs []string,
) error {
	// bingetracker play next Twin Peaks OR bingetracker play next Twin Peaks sXX

	// skip_in_progress = false:
	// find first unfinished season -> currentSeason := GetUnfinishedSeasons[0]
	// find first unwatched episode in season -> GetUnwatchedSeasonEpisodes(currentSeason)[0]
	// any of the above don't exist -> error
	// else -> done

	if !skipInProgress {
		unfinishedSeasons, err := s.DB.GetUnfinishedSeasons(cmdContext, btArgs.ShowTitle)
		if err != nil {
			return err
		}

		// len == 0 means all seasons have been watched, i.e. the show has been fully watched.
		if len(unfinishedSeasons) == 0 {
			return new(ShowWatchedError)
		}

		currentSeasonNumber := unfinishedSeasons[0].SeasonNumber

		if btArgs.SeasonNumber != -1 && int64(btArgs.SeasonNumber) < currentSeasonNumber {
			return new(SeasonWatchedError)
		} else if int64(btArgs.SeasonNumber) > currentSeasonNumber {
			// in this case, this means the user is watching a more recent season before finishing the older one(s)
			// this is fine, so just set the current season number to the user's input
			currentSeasonNumber = int64(btArgs.SeasonNumber)
		}

		// at this point in the program, there must exist an unwatched episode in this season
		// otherwise it would have returned with a ShowWatchedError already.
		nextEpisode, err := s.DB.GetNextUnwatchedSeasonEpisode(cmdContext, database.GetNextUnwatchedSeasonEpisodeParams{
			SeasonNumber: currentSeasonNumber,
			ShowTitle:    btArgs.ShowTitle,
		})
		if err != nil {
			return err
		}

		if btArgs.VideoPlayer == "mpv" {
			return HandlerPlayMPV(
				cmdContext,
				s,
				btArgs.ShowTitle,
				int(currentSeasonNumber),
				int(nextEpisode.EpisodeNumber),
				false,
				playerArgs,
			)
		} else {
			// insert HandlerPlayVLC here
			return nil
		}
	} else {
		// skip_in_progress = true:
		// find season with least viewtime -> GetLeastViewedSeason[0] (ORDER BY total_viewtime ASC, season_number ASC)
		// find episode with 0 viewtime -> GetNoProgressSeasonEpisodes[0] (ORDER BY viewtime ASC, episode_number ASC)
		// episode with 0 viewtime doesn't exist -> error
		// else -> done

		leastViewedSeason, err := s.DB.GetLeastViewedSeason(cmdContext, btArgs.ShowTitle)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return new(ShowWatchedError)
			}
			return err
		}

		currentSeasonNumber := leastViewedSeason.SeasonNumber

		if btArgs.SeasonNumber != -1 && int64(btArgs.SeasonNumber) < currentSeasonNumber {
			return new(SeasonWatchedError)
		} else if int64(btArgs.SeasonNumber) > currentSeasonNumber {
			// in this case, this means the user is watching a more recent season before finishing the older one(s)
			// this is fine, so just set the current season number to the user's input
			currentSeasonNumber = int64(btArgs.SeasonNumber)
		}

		nextEpisode, err := s.DB.GetNextNoProgressSeasonEpisode(cmdContext, database.GetNextNoProgressSeasonEpisodeParams{
			ShowTitle:    btArgs.ShowTitle,
			SeasonNumber: currentSeasonNumber,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return new(SeasonWatchedError)
			}
			return err
		}

		if btArgs.VideoPlayer == "mpv" {
			return HandlerPlayMPV(
				cmdContext,
				s,
				btArgs.ShowTitle,
				int(currentSeasonNumber),
				int(nextEpisode.EpisodeNumber),
				false,
				playerArgs,
			)
		} else {
			// insert HandlerPlayVLC here
			return nil
		}
	}
}

func HandlerPlayMPV(
	cmdContext context.Context,
	s *app.State,
	showTitle string,
	seasonNumber, episodeNumber int,
	restart bool,
	playerArgs []string) error {

	episode, err := s.DB.GetEpisode(cmdContext, database.GetEpisodeParams{
		ShowTitle:     showTitle,
		SeasonNumber:  int64(seasonNumber),
		EpisodeNumber: int64(episodeNumber),
	})
	if err != nil {
		return err
	}

	mpvFlags := []string{
		episode.Filepath,
		fmt.Sprintf("--input-ipc-server=%s", s.Cfg.SocketPath),
		fmt.Sprintf("--script=%s", s.Cfg.ScriptPath),
	}
	mpvFlags = append(mpvFlags, playerArgs...)

	// decide what timestamp to start the file at

	if restart {
		mpvFlags = append(mpvFlags, "--rebase-start-time=yes")
	} else if episode.Viewtime != 0.0 {
		startFlag := fmt.Sprintf("--start=+%d", int(episode.Viewtime))
		mpvFlags = append(mpvFlags, startFlag)
	}

	playerCmd := exec.Command("mpv", mpvFlags...)

	_ = os.Remove(s.Cfg.SocketPath)

	fmt.Printf("Starting s%de%d of %s...\n", seasonNumber, episodeNumber, showTitle)

	err = playerCmd.Start()
	if err != nil {
		return fmt.Errorf("error launching video player: %v", err)
	}

	go func() {
		err := playerCmd.Wait()
		if err != nil {
			log.Printf("Process exited or was killed by user: %v", err)
			return
		}
		log.Println("Process exited gracefully on its own.")
	}()

	viewTime, err := listen.TrackViewTime(s.Cfg.SocketPath)
	if err != nil {
		return err
	}

	watched := (viewTime/episode.Runtime >= 0.85)

	_, err = s.DB.UpdateEpisodeStats(cmdContext, database.UpdateEpisodeStatsParams{
		Viewtime:      viewTime,
		Watched:       watched,
		ShowTitle:     episode.ShowTitle,
		SeasonNumber:  episode.SeasonNumber,
		EpisodeNumber: episode.EpisodeNumber,
	})
	if err != nil {
		return err
	}

	if watched {
		fmt.Printf("You finished watching s%de%d of %s.\n", seasonNumber, episodeNumber, showTitle)
	} else {
		fmt.Printf("You did not finish watching s%de%d of %s.\n", seasonNumber, episodeNumber, showTitle)
	}

	return nil
}

func handlerPlayVLC(
	c context.Context,
	s *app.State,
	showTitle string,
	seasonNumber, episodeNumber int,
	restart bool,
	playerArgs []string,
) {
}
