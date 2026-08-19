package handler

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/mythenox/binge-tracker/internal/app"
)

func HandlerSetConfig(cmdContext context.Context, s *app.State,
	configChanges map[string]string) error {

	for key, value := range configChanges {
		switch key {
		case "count_partial_progress":
			switch value {
			case "true":
				s.Cfg.CountPartialProgress = true
			case "false":
				s.Cfg.CountPartialProgress = false
			default:
				return errors.New("Invalid value for config variable count_partial_progress.")
			}
		case "video_player":
			supportedPlayers := []string{"mpv", "vlc"}
			if !slices.Contains(supportedPlayers, value) {
				return errors.New("Video players other than VLC and MPV are not supported.")
			} else {
				s.Cfg.VideoPlayer = value
			}
		default:
			return errors.New("Invalid input")
		}
	}

	err := s.WriteConfig()
	if err != nil {
		return err
	}

	fmt.Println("Successfully updated config.")

	return nil
}
