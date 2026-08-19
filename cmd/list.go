/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// bingetracker list shows
// bingetracker list episodes Twin Peaks s02
// bingetracker list seasons Twin Peaks

// ex output: Twin Peaks s01e01 85:32/89:00 (watched)
//			  Twin Peaks s01e02

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list {shows | seasons | episodes}",
	Short: "Lists shows, seasons, or episodes.",
	Long: `This command does nothing on its own, and must be used with the 'shows', 'seasons', or 'episodes' subcommands.
Usage example: 'bingetracker list seasons [arguments]'.
See the help commands for the respective subcommands for more information their usage.`,
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
