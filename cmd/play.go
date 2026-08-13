/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"regexp"
	"strconv"

	"github.com/mythenox/binge-tracker/internal/handler"
	"github.com/spf13/cobra"
)

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play (mpv | vlc) <show name> <sXXeYY> [-- player_flags]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		episodeRegex := regexp.MustCompile(`(?i)^s(\d+)e(\d+)$`)

		btArgs := cmd.Flags().Args()

		episodeIdentifier := btArgs[2]
		matches := episodeRegex.FindStringSubmatch(episodeIdentifier)
		if len(matches) != 3 {
			return errors.New("Invalid format")
		}

		videoPlayer := btArgs[0]
		showTitle := btArgs[1]
		seasonNumber, err := strconv.Atoi(matches[1])
		if err != nil {
			return err
		}
		episodeNumber, err := strconv.Atoi(matches[2])
		if err != nil {
			return err
		}

		restart, _ := cmd.Flags().GetBool("restart")
		playerArgs := cmd.Flags().Args()[cmd.Flags().NArg():]

		switch videoPlayer {
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
