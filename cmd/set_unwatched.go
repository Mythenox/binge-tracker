/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// unwatchedCmd represents the unwatched command
var setUnwatchedCmd = &cobra.Command{
	Use:   "unwatched <show title> {season or season range | episode or episode range}",
	Short: "Sets an episode or range of episodes as unwatched",
	Long: `This command sets an episode or range of episodes as unwatched. Can also input a season or range of seasons to set entire seasons as unwatched in bulk.
Usage example: 'bingetracker set unwatched Twin Peaks s01-s02'`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputIsRange := strings.Contains(args[1], "-")
		episodeMode := strings.Contains(args[1], "e")

		return setCompletion(cmd, args, s, inputIsRange, episodeMode, false)
	},
}

func init() {
	setCmd.AddCommand(setUnwatchedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unwatchedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unwatchedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
