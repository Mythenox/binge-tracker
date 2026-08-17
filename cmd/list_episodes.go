/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mythenox/binge-tracker/internal/handler"
	"github.com/spf13/cobra"
)

// bingetracker list seasons Twin Peaks
// verbose flag for listing filepath as well

// episodesCmd represents the episodes command
var episodesCmd = &cobra.Command{
	Use:   "episodes <show title> <season number>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		showTitle := args[0]
		seasonIdentifier := args[1]

		seasonNumber, err := extractFromSeasonIdentifier(seasonIdentifier)
		if err != nil {
			return err
		}

		return handler.HandlerListEpisodes(
			cmd.Context(),
			s,
			showTitle,
			seasonNumber,
		)
	},
}

func init() {
	listCmd.AddCommand(episodesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// episodesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// episodesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
