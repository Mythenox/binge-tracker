/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mythenox/binge-tracker/internal/handler"
	"github.com/spf13/cobra"
)

// bingetracker add season Twin Peaks s02 <dirpath>

// seasonCmd represents the season command
var addSeasonCmd = &cobra.Command{
	Use:   "season",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		showTitle := args[0]
		seasonIdentifier := args[1]
		seasonDirPath := args[2]

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
