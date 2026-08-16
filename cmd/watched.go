/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/mythenox/binge-tracker/internal/app"
	"github.com/mythenox/binge-tracker/internal/handler"
	"github.com/spf13/cobra"
)

// bingetracker set watched <show title> {season or season range | episode or episode range}
// bingetracker set watched Twin Peaks s01-s02
// bingetracker set watched Twin Peaks s01e01-s01e06

// watchedCmd represents the watched command
var watchedCmd = &cobra.Command{
	Use:   "watched",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputIsRange := strings.Contains(args[1], "-")
		episodeMode := strings.Contains(args[1], "e")

		return setCompletion(cmd, args, s, inputIsRange, episodeMode, true)
	},
}

func setCompletion(cmd *cobra.Command, args []string, s *app.State,
	inputIsRange, episodeMode, setWatched bool) error {
	showTitle := args[0]

	if episodeMode {
		if inputIsRange {
			startEpisodeIdentifier := strings.Split(args[1], "-")[0]
			endEpisodeIdentifier := strings.Split(args[1], "-")[1]

			startNums, err := extractFromEpisodeIdentifier(startEpisodeIdentifier)
			if err != nil {
				return err
			}

			endNums, err := extractFromEpisodeIdentifier(endEpisodeIdentifier)
			if err != nil {
				return err
			}

			return handler.HandlerSetEpisodeRangeCompletion(cmd.Context(), s, showTitle,
				startNums, endNums, setWatched)
		} else {
			episodeIdentifier := args[1]

			nums, err := extractFromEpisodeIdentifier(episodeIdentifier)
			if err != nil {
				return err
			}

			seasonNumber, episodeNumber := nums[0], nums[1]

			return handler.HandlerSetEpisodeCompletion(cmd.Context(), s, showTitle,
				seasonNumber, episodeNumber, setWatched)
		}

	} else {
		if inputIsRange {
			startSeasonIdentifier := strings.Split(args[1], "-")[0]
			endSeasonIdentifier := strings.Split(args[1], "-")[1]

			startSeasonNumber, err := extractFromSeasonIdentifier(startSeasonIdentifier)
			if err != nil {
				return err
			}

			endSeasonNumber, err := extractFromSeasonIdentifier(endSeasonIdentifier)
			if err != nil {
				return err
			}

			return handler.HandlerSetSeasonRangeCompletion(cmd.Context(), s, showTitle,
				startSeasonNumber, endSeasonNumber, setWatched)
		} else {
			seasonIdentifier := args[1]

			seasonNumber, err := extractFromSeasonIdentifier(seasonIdentifier)
			if err != nil {
				return err
			}

			return handler.HandlerSetSeasonCompletion(cmd.Context(), s, showTitle,
				seasonNumber, setWatched)
		}
	}
}

func init() {
	setCmd.AddCommand(watchedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
