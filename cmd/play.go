/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/mythenox/bingetracker/internal/handler"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play <show name> <sXXeYY> [-- player_flags]",
	Short: "Play a specified episode",
	Long: `Plays the specified episode. Arguments can be passed through to the player by use of the '--' separator; all arguments past this separator will be passed to the player.
Usage example: 'bingetracker play Twin Peaks s01e01 -- --mute=yes'`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		dashIndex := cmd.ArgsLenAtDash()
		if dashIndex == -1 {
			dashIndex = len(args)
		}

		if dashIndex < 2 {
			return fmt.Errorf("needs at least 2 args before '--', received %d", dashIndex)
		}

		cmdArgs := args[:dashIndex]
		playerArgs := args[dashIndex:]

		restart, _ := cmd.Flags().GetBool("restart")

		showTitleWords := cmdArgs[:len(cmdArgs)-1]
		episodeIdentifier := cmdArgs[len(cmdArgs)-1]

		showTitle := strings.ToLower(strings.Join(showTitleWords, " "))

		nums, err := extractFromEpisodeIdentifier(episodeIdentifier)
		if err != nil {
			return err
		}

		seasonNumber, episodeNumber := nums[0], nums[1]

		switch s.Cfg.VideoPlayer {
		case "mpv":
			return handler.HandlerPlayMPV(
				cmd.Context(),
				s,
				showTitle,
				seasonNumber,
				episodeNumber,
				restart,
				playerArgs,
			)
		case "vlc":
			return handler.HandlerPlayVLC(
				cmd.Context(),
				s,
				showTitle,
				seasonNumber,
				episodeNumber,
				restart,
				playerArgs,
			)
		default:
			return errors.New("unsupported player")
		}
	},
}

// returns {XX, YY} from input sXXeYY.
func extractFromEpisodeIdentifier(episodeIdentifier string) ([]int, error) {
	episodeRegex := regexp.MustCompile(`(?i)^s(\d+)e(\d+)$`)

	matches := episodeRegex.FindStringSubmatch(episodeIdentifier)
	if len(matches) != 3 {
		return nil, errors.New("Invalid format")
	}

	seasonNumber, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, err
	}
	episodeNumber, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, err
	}

	nums := []int{seasonNumber, episodeNumber}

	return nums, nil
}

// returns XX from input sXX.
func extractFromSeasonIdentifier(seasonIdentifier string) (int, error) {
	seasonRegex := regexp.MustCompile(`(?i)^s(\d+)$`)

	matches := seasonRegex.FindStringSubmatch(seasonIdentifier)
	if len(matches) != 2 {
		return 0, errors.New("Invalid format")
	}

	seasonNumber, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0, err
	}

	return seasonNumber, nil
}

func init() {
	rootCmd.AddCommand(playCmd)
	playCmd.Flags().BoolP("restart", "r", false,
		"Restart from beginning of video instead of resuming from where playback was last stopped")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// playCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// playCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
