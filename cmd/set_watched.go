/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/mythenox/bingetracker/internal/app"
	"github.com/mythenox/bingetracker/internal/handler"
	"github.com/spf13/cobra"
)

// bingetracker set watched <show title> {season or season range | episode or episode range}
// bingetracker set watched Twin Peaks s01-s02
// bingetracker set watched Twin Peaks s01e01-s01e06

// watchedCmd represents the watched command
var setWatchedCmd = &cobra.Command{
	Use:   "watched",
	Short: "Sets an episode or range of episodes as watched",
	Long: `This command sets an episode or range of episodes as watched. Can also input a season or range of seasons to set entire seasons as watched in bulk.
Usage example: 'bingetracker set watched Twin Peaks s01-s02'`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputIsRange := strings.Contains(args[len(args)-1], "-")
		episodeMode := strings.Contains(args[len(args)-1], "e")

		return setCompletion(cmd, args, s, inputIsRange, episodeMode, true)
	},
}

func init() {
	setCmd.AddCommand(setWatchedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func setCompletion(cmd *cobra.Command, args []string, s *app.State,
	inputIsRange, episodeMode, setWatched bool) error {
	showTitleWords := args[:len(args)-1]

	showTitle := strings.ToLower(strings.Join(showTitleWords, " "))

	if episodeMode {
		if inputIsRange {
			startEpisodeIdentifier := strings.Split(args[len(args)-1], "-")[0]
			endEpisodeIdentifier := strings.Split(args[len(args)-1], "-")[1]

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
			episodeIdentifier := args[len(args)-1]

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
			startSeasonIdentifier := strings.Split(args[len(args)-1], "-")[0]
			endSeasonIdentifier := strings.Split(args[len(args)-1], "-")[1]

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
			seasonIdentifier := args[len(args)-1]

			seasonNumber, err := extractFromSeasonIdentifier(seasonIdentifier)
			if err != nil {
				return err
			}

			return handler.HandlerSetSeasonCompletion(cmd.Context(), s, showTitle,
				seasonNumber, setWatched)
		}
	}
}
