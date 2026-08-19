package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"

	"github.com/mythenox/bingetracker/internal/app"
	"github.com/mythenox/bingetracker/internal/database"
	"github.com/mythenox/bingetracker/internal/listen"
)

// play [mpv, vlc] <show_name> <season_number> <episode_number> [optional extra mpv/vlc flags]

// RunE func(cmd *Command, args []string) error

type SeasonWatchedErr struct{}

func (*SeasonWatchedErr) Error() string {
	return "All episodes of this season have been watched."
}

type ShowWatchedErr struct{}

func (*ShowWatchedErr) Error() string {
	return "All episodes of this show have been watched."
}

type EpisodeNotFoundErr struct{}

func (*EpisodeNotFoundErr) Error() string {
	return "The specified episode was not found in the database."
}

func HandlerPlayNext(
	cmdContext context.Context,
	s *app.State,
	showTitle string,
	seasonNumber int,
	verboseInput bool,
	playerArgs []string,
) error {
	// bingetracker play next Twin Peaks OR bingetracker play next Twin Peaks sXX
	nextUnfinishedSeason, err := s.Q.GetNextUnfinishedSeason(cmdContext, showTitle)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(ShowWatchedErr)
		}
		return err
	}

	// lowest season number with unwatched episodes
	currentSeasonNumber := nextUnfinishedSeason.SeasonNumber

	// if the user provided a season number, make sure that it doesn't point to a finished season
	// or use it if it's ahead of the currentSeasonNumber
	if verboseInput {
		if int64(seasonNumber) < currentSeasonNumber {
			return new(SeasonWatchedErr)
		} else if int64(seasonNumber) > currentSeasonNumber {
			// in this case, this means the user is watching a more recent season before finishing the older one(s)
			// this is fine, so just set the current season number to the user's input
			currentSeasonNumber = int64(seasonNumber)
		}
	}

	// at this point in the program, there must exist an unwatched episode in this season
	// otherwise it would have returned with a ShowWatchedErr already.
	nextEpisode, err := s.Q.GetNextUnwatchedEpisodeFromSeason(cmdContext, database.GetNextUnwatchedEpisodeFromSeasonParams{
		SeasonNumber: currentSeasonNumber,
		ShowTitle:    showTitle,
	})
	if err != nil {
		return err
	}

	if s.Cfg.VideoPlayer == "mpv" {
		return HandlerPlayMPV(
			cmdContext,
			s,
			showTitle,
			int(currentSeasonNumber),
			int(nextEpisode.EpisodeNumber),
			false,
			playerArgs,
		)
	} else {
		return HandlerPlayVLC(
			cmdContext,
			s,
			showTitle,
			int(currentSeasonNumber),
			int(nextEpisode.EpisodeNumber),
			false,
			playerArgs,
		)
	}
}

func HandlerPlayMPV(
	cmdContext context.Context,
	s *app.State,
	showTitle string,
	seasonNumber, episodeNumber int,
	restart bool,
	playerArgs []string) error {

	episode, err := s.Q.GetEpisode(cmdContext, database.GetEpisodeParams{
		ShowTitle:     showTitle,
		SeasonNumber:  int64(seasonNumber),
		EpisodeNumber: int64(episodeNumber),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(EpisodeNotFoundErr)
		}
		return err
	}

	// if the episode was previously fully watched, the user choosing to play it again
	// most likely wants to watch it again from the beginning

	if int(episode.Viewtime) == int(episode.Runtime) {
		restart = true
	}

	var connPath string
	operatingSystem := runtime.GOOS
	if operatingSystem == "windows" {
		connPath = s.Cfg.PipePath
	} else {
		connPath = s.Cfg.SocketPath
	}

	mpvFlags := []string{
		episode.Filepath,
		fmt.Sprintf("--input-ipc-server=%s", connPath),
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

	os.Remove(connPath)

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

	viewtime, err := listen.TrackViewTimeMPV(connPath)
	if err != nil {
		return err
	}

	err = updateDB(cmdContext, s, viewtime, episode)
	if err != nil {
		return err
	}

	return nil
}

// place lua file in ~/.local/share/vlc/lua/intf

func HandlerPlayVLC(
	cmdContext context.Context,
	s *app.State,
	showTitle string,
	seasonNumber, episodeNumber int,
	restart bool,
	playerArgs []string,
) error {
	episode, err := s.Q.GetEpisode(cmdContext, database.GetEpisodeParams{
		ShowTitle:     showTitle,
		SeasonNumber:  int64(seasonNumber),
		EpisodeNumber: int64(episodeNumber),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return new(EpisodeNotFoundErr)
		}
		return err
	}

	// if the episode was previously fully watched, the user choosing to play it again
	// most likely wants to watch it again from the beginning

	if int(episode.Viewtime) == int(episode.Runtime) {
		restart = true
	}

	vlcFlags := []string{
		episode.Filepath,
		"--extraintf=luaintf",
		"--lua-intf=tcptracker",
	}
	vlcFlags = append(vlcFlags, playerArgs...)

	// decide what timestamp to start the file at

	if restart {
		vlcFlags = append(vlcFlags, "--start-time=0")
	} else if episode.Viewtime != 0.0 {
		startFlag := fmt.Sprintf("--start-time=%f", episode.Viewtime)
		vlcFlags = append(vlcFlags, startFlag)
	}

	playerCmd := exec.Command("vlc", vlcFlags...)
	playerCmd.Env = append(os.Environ(),
		"XDG_DATA_HOME="+os.Getenv("HOME")+"/.local/share")

	fmt.Printf("Starting s%de%d of %s...\n", seasonNumber, episodeNumber, showTitle)

	viewtime, err := listen.TrackViewTimeVLC(playerCmd)
	if err != nil {
		return err
	}

	err = updateDB(cmdContext, s, viewtime, episode)
	if err != nil {
		return err
	}

	return nil
}

func updateDB(cmdContext context.Context, s *app.State,
	viewtime float64, episode database.Episode) error {
	var watched bool
	if s.Cfg.CountPartialProgress {
		watched = viewtime >= 0.0
	} else {
		watched = (viewtime/episode.Runtime >= 0.85)
	}

	_, err := s.Q.UpdateEpisodeStats(cmdContext, database.UpdateEpisodeStatsParams{
		Viewtime:      viewtime,
		Watched:       watched,
		ShowTitle:     episode.ShowTitle,
		SeasonNumber:  episode.SeasonNumber,
		EpisodeNumber: episode.EpisodeNumber,
	})
	if err != nil {
		return err
	}

	if watched {
		tx, err := s.DB.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()

		qtx := s.Q.WithTx(tx)

		cascadeWatchStatus(cmdContext, qtx, episode.ShowTitle,
			episode.SeasonNumber, 1, true)

		fmt.Printf("You finished watching s%de%d of %s.\n",
			episode.SeasonNumber, episode.EpisodeNumber, episode.ShowTitle)

		return tx.Commit()
	} else {
		fmt.Printf("You did not finish watching s%de%d of %s.\n",
			episode.SeasonNumber, episode.EpisodeNumber, episode.ShowTitle)
	}

	return nil
}
