package handler

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strconv"

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

func HandlerPlay(s *app.State, cmd app.Command) error {
	// add functionality to automatically resume from last watched spot
	// add optional flag to restart from beginning

	flagSimple := flag.Bool("s", false, "declares simple episode identifier input (allows for inputs such as s03e04 instead of two separate arguments)")
	flagRestart := flag.Bool("r", false, "plays video from the beginning instead of resuming from where playback was last stopped")

	flag.Parse()

	var videoPlayer, showTitle, episodeArg, seasonArg string

	btArgs := cmd.Args
	var playerFlags []string

	doubleDashIdx := slices.Index(cmd.Args, "--")
	if doubleDashIdx != -1 {
		playerFlags = append(playerFlags, cmd.Args[doubleDashIdx+1:]...)
		btArgs = cmd.Args[:doubleDashIdx]
	}

	videoPlayer, showTitle = btArgs[0], btArgs[1]

	if *flagSimple {
		if len(btArgs) != 3 {
			return fmt.Errorf("Usage: %s (mpv | vlc) <show title> <episode identifier> [optional flags]", cmd.Name)
		}
		episodeArg = btArgs[2]

		episodeIdx := -1
		for i := range episodeArg {
			if string(episodeArg[i]) == "e" {
				episodeIdx = i
				break
			}
		}

		// if the string does not contain an e or ends in e (in the latter case an episode number was not given)
		if episodeIdx == -1 || episodeIdx == len(episodeArg)-1 {
			return fmt.Errorf("Usage: %s (mpv | vlc) <show title> <episode identifier> [optional flags]", cmd.Name)
		}

		seasonArg = CleanInput(episodeArg[0:episodeIdx], "s")
		episodeArg = CleanInput(episodeArg[episodeIdx:], "e")
	} else if len(cmd.Args) != 4 {
		return fmt.Errorf("Usage: %s (mpv | vlc) <show title> <season number> <episode number> [optional flags]", cmd.Name)
	} else {
		seasonArg = CleanInput(btArgs[2], "s")
		episodeArg = CleanInput(btArgs[3], "e")
	}

	seasonNumber, err := strconv.Atoi(seasonArg)
	if err != nil {
		return err
	}
	episodeNumber, err := strconv.Atoi(episodeArg)
	if err != nil {
		return err
	}

	args := playerArgs{
		title:   showTitle,
		season:  int64(seasonNumber),
		episode: int64(episodeNumber),
		restart: *flagRestart,
		flags:   playerFlags,
	}

	switch videoPlayer {
	case "mpv":
		return HandlerPlayMPV(s, args)
	case "vlc":
		return errors.New("vlc not implemented yet")
	default:
		return errors.New("unsupported player")
	}
}

func HandlerPlayMPV(s *app.State, p playerArgs) error {

	playerFlags := []string{
		fmt.Sprintf("--input-ipc-server=%s", s.Cfg.SocketPath),
		fmt.Sprintf("--script=%s", s.Cfg.ScriptPath),
	}
	playerFlags = append(playerFlags, p.flags...)

	episode, err := s.DB.GetEpisode(context.Background(), database.GetEpisodeParams{
		ShowTitle:     p.title,
		SeasonNumber:  p.season,
		EpisodeNumber: p.episode,
	})
	if err != nil {
		return err
	}

	// decide what timestamp to start the file at

	if p.restart {
		playerFlags = append(playerFlags, "--rebase-start-time=yes")
	} else if episode.Viewtime != 0.0 {
		startFlag := fmt.Sprintf("--start=+%d", int(episode.Viewtime))
		playerFlags = append(playerFlags, startFlag)
	}

	playerCmd := exec.Command("mpv", playerFlags...)

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
		fmt.Printf("You finished watching s%de%d of %s.\n", p.season, p.episode, p.title)
	} else {
		fmt.Printf("You did not finish watching s%de%d of %s.\n", p.season, p.episode, p.title)
	}

	return nil
}

func handlerPlayVLC(s *app.State, p playerArgs) {}

func CleanInput(inputStr, prefix string) string {
	// ex: CleanInput("e03", "e") == "3", CleanInput("s8", "s") == "8"
	if string(inputStr[0]) == prefix {
		if len(inputStr) > 2 && string(inputStr[1]) == "0" {
			// if input is of the form s0n (n >= 0), converts to just n
			return inputStr[2:]
		} else {
			// if input is of the form sn (n >= 0), converts to just n
			return inputStr[1:]
		}
	}
	return inputStr
}
