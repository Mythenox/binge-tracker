package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/mythenox/binge-tracker/internal/app"
	"github.com/mythenox/binge-tracker/internal/database"
	"github.com/mythenox/binge-tracker/internal/listen"
)

// play [mpv, vlc] <show_name> <season_number> <episode_number> [optional extra mpv/vlc flags]

func HandlerPlay(s *app.State, cmd app.Command) error {
	if len(cmd.Args) < 4 {
		return fmt.Errorf("Usage: %s (mpv | vlc) <show title> <season number> <episode number> [optional extra mpv/vlc flags]", cmd.Name)
	}

	videoPlayer, showTitle, seasonArg, episodeArg, otherArgs := cmd.Args[0], cmd.Args[1], cmd.Args[2], cmd.Args[3], cmd.Args[4:]
	mpvFlags := map[string]string{
		"ipc":    fmt.Sprintf("--input-ipc-server=%s", s.Cfg.SocketPath),
		"script": fmt.Sprintf("--script=%s", s.Cfg.ScriptPath),
	}

	seasonNumber, err := strconv.Atoi(seasonArg)
	if err != nil {
		return err
	}
	episodeNumber, err := strconv.Atoi(episodeArg)
	if err != nil {
		return err
	}

	episode, err := s.DB.GetEpisode(context.Background(), database.GetEpisodeParams{
		ShowTitle:     showTitle,
		SeasonNumber:  int64(seasonNumber),
		EpisodeNumber: int64(episodeNumber),
	})
	if err != nil {
		return err
	}

	if videoPlayer == "vlc" {
		return errors.New("vlc not yet supported")
	}

	_ = os.Remove(s.Cfg.SocketPath)

	playerArgs := []string{episode.Filepath, mpvFlags["ipc"], mpvFlags["script"]}
	playerArgs = append(playerArgs, otherArgs...)

	playerCmd := exec.Command(videoPlayer, playerArgs...)

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

	socketReady := false
	for range 50 { // try for up to 5 seconds
		if _, err := os.Stat(s.Cfg.SocketPath); err == nil {
			socketReady = true
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	if !socketReady {
		return fmt.Errorf("timed out waiting for mpv socket at %s", s.Cfg.SocketPath)
	}

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
