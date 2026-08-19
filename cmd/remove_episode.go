/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mythenox/binge-tracker/internal/handler"
	"github.com/spf13/cobra"
)

// bingetracker remove episode Twin Peaks s02e01

// episodeCmd represents the episode command
var removeEpisodeCmd = &cobra.Command{
	Use:   "episode <show title> <episode identifier>",
	Short: "Removes the specified episode from the database.",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		showTitle := args[0]
		episodeIdentifier := args[1]

		nums, err := extractFromEpisodeIdentifier(episodeIdentifier)
		if err != nil {
			return err
		}

		seasonNumber, episodeNumber := nums[0], nums[1]

		return handler.HandlerRemoveEpisode(cmd.Context(), s, showTitle, seasonNumber, episodeNumber)
	},
}

func init() {
	removeCmd.AddCommand(removeEpisodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// episodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// episodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
