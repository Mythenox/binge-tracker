package handler

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/mythenox/binge-tracker/internal/app"
	"github.com/mythenox/binge-tracker/internal/database"
	"github.com/mythenox/binge-tracker/internal/listen"
)

type playerArgs struct {
	title   string
	season  int64
	episode int64
	restart bool
	flags   []string
}

// play [mpv, vlc] <show_name> <season_number> <episode_number> [optional extra mpv/vlc flags]

// RunE func(cmd *Command, args []string) error

func HandlerPlayMPV(
	c context.Context,
	s *app.State,
	showTitle string,
	seasonNumber, episodeNumber int,
	restart bool,
	playerArgs []string) error {
	// add functionality to automatically resume from last watched spot

	episode, err := s.DB.GetEpisode(context.Background(), database.GetEpisodeParams{
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

	_, err = s.DB.UpdateEpisodeStats(context.Background(), database.UpdateEpisodeStatsParams{
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

func handlerPlayVLC(s *app.State, p playerArgs) {}
