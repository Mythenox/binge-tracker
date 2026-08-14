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

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init <show name> <season> <season directory path>",
	Short: "Initiate a show to be tracked",
	Long: `initiate a show to be tracked. requires show name, season to be added, and dirpath of the season.
ex: bingetracker init "twin peaks" s01 videos/twin-peaks-s01
automatically adds all episodes found in given directory`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		// ex: bingetracker init "twin peaks" s01 videos/twin-peaks-s01
		seasonRegex := regexp.MustCompile(`(?i)^s(\d+)$`)

		showTitle := cmd.Flags().Args()[0]

		seasonIdentifier := cmd.Flags().Args()[1]
		matches := seasonRegex.FindStringSubmatch(seasonIdentifier)
		if len(matches) != 2 {
			return errors.New("Invalid format")
		}

		seasonNumber, err := strconv.Atoi(matches[1])
		if err != nil {
			return err
		}
		if seasonNumber < 0 {
			return errors.New("Season numbers must be non-negative")
		}

		seasonDirPath := cmd.Flags().Args()[2]

		return handler.HandlerInit(
			cmd.Context(),
			s,
			showTitle,
			seasonNumber,
			seasonDirPath,
		)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
