/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/mythenox/bingetracker/internal/handler"
	"github.com/spf13/cobra"
)

// bingetracker move Twin Peaks s02 <new dirpath>

// account for changing filepaths
// allows for changing of season dirpath?

// moveCmd represents the move command
var moveCmd = &cobra.Command{
	Use:   "move <show title> <season number> <new directory path>",
	Short: "Move the target directory of an existing season to a new path",
	Long: `This command allows for changing the target directory of an existing season.
If the new target directory has a different number of episodes, user confirmation is required to proceed, as the watch information for the season must be completely reset.
Usage example: 'bingetracker move Twin Peaks s02 '~/Downloads/Twin Peaks S02'`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		showTitleWords := args[:len(args)-2]
		seasonIdentifier := args[len(args)-2]
		seasonDirPath := args[len(args)-1]

		showTitle := strings.ToLower(strings.Join(showTitleWords, " "))

		seasonNumber, err := extractFromSeasonIdentifier(seasonIdentifier)
		if err != nil {
			return err
		}

		return handler.HandlerMoveSeason(cmd.Context(), s,
			seasonNumber, showTitle, seasonDirPath)
	},
}

func init() {
	rootCmd.AddCommand(moveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// moveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// moveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
