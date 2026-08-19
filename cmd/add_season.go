/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/mythenox/bingetracker/internal/handler"
	"github.com/spf13/cobra"
)

// bingetracker add season Twin Peaks s02 <dirpath>

// seasonCmd represents the season command
var addSeasonCmd = &cobra.Command{
	Use:     "season <show title> <season identifier> <season directory path>",
	Short:   "Add a new season to an existing show",
	Example: "add season Twin Peaks s02 '~/Downloads/Twin Peaks Season 2'",
	Args:    cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		showTitleWords := args[:len(args)-2]
		seasonIdentifier := args[len(args)-2]
		seasonDirPath := args[len(args)-1]

		showTitle := strings.ToLower(strings.Join(showTitleWords, " "))

		seasonNumber, err := extractFromSeasonIdentifier(seasonIdentifier)
		if err != nil {
			return err
		}

		return handler.HandlerAddSeason(cmd.Context(), s,
			seasonNumber, showTitle, seasonDirPath)
	},
}

func init() {
	addCmd.AddCommand(addSeasonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seasonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seasonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
