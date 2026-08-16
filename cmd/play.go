/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mythenox/binge-tracker/internal/handler"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play <show name> <sXXeYY> [-- player_flags]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.Println("play called")
		dashIndex := cmd.ArgsLenAtDash()
		if dashIndex == -1 {
			dashIndex = len(args)
		}

		if dashIndex != 2 {
			return fmt.Errorf("accepts exactly 2 args before '--', received %d", dashIndex)
		}

		cmdArgs := args[:dashIndex]
		playerArgs := args[dashIndex:]

		restart, _ := cmd.Flags().GetBool("restart")

		showTitle := cmdArgs[0]
		episodeIdentifier := cmdArgs[1]

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
			return errors.New("vlc not implemented yet")
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
