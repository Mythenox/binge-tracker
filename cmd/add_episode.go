/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/mythenox/bingetracker/internal/handler"
	"github.com/spf13/cobra"
)

// bingetracker add episode Twin Peaks s02e01 <path>

// episodeCmd represents the episode command
var addEpisodeCmd = &cobra.Command{
	Use:     "episode <show title> <episode identifier> <episode path>",
	Short:   "Add an episode to a season of a show",
	Example: "add episode Twin Peaks s02e01 ~/Downloads/twin_peaks_s02_e01.mp4",
	Args:    cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		showTitleWords := args[:len(args)-2]
		episodeIdentifier := args[len(args)-2]
		episodePath := args[len(args)-1]

		showTitle := strings.ToLower(strings.Join(showTitleWords, " "))

		nums, err := extractFromEpisodeIdentifier(episodeIdentifier)
		if err != nil {
			return err
		}

		seasonNumber, episodeNumber := nums[0], nums[1]

		return handler.HandlerAddEpisode(
			cmd.Context(),
			s,
			seasonNumber,
			episodeNumber,
			showTitle,
			episodePath,
		)
	},
}

func init() {
	addCmd.AddCommand(addEpisodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// episodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// episodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
