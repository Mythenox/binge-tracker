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
var removeSeasonCmd = &cobra.Command{
	Use:   "season <show title> <season identifier>",
	Short: "Removes the specified season from the database.",
	Long: `This command removes the specified season from the database, along with all of its episodes.
Usage example: 'bingetracker remove season Twin Peaks s02'`,
	Args: cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		showTitleWords := args[:len(args)-1]
		seasonIdentifier := args[len(args)-1]

		showTitle := strings.ToLower(strings.Join(showTitleWords, " "))

		seasonNumber, err := extractFromSeasonIdentifier(seasonIdentifier)
		if err != nil {
			return err
		}

		return handler.HandlerRemoveSeason(cmd.Context(), s, showTitle, seasonNumber)
	},
}

func init() {
	removeCmd.AddCommand(removeSeasonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seasonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seasonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
