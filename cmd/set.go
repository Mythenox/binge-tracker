/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets an episode as watched or unwatched",
	Long: `This command does nothing on its own, and must be used with either the 'watched' or 'unwatched' subcommands.
Usage example: 'bingetracker set watched Twin Peaks s02 e01-e06'
See the help commands for the respective subcommands for more information their usage.`,
}

// bingetracker set watched twin peaks s02 e01-e06

func init() {
	rootCmd.AddCommand(setCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
