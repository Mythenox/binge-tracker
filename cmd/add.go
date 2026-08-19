/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new season or new episode",
	Long: `This command does nothing on its own; it must be used with either the 'season' or 'episode' subcommands.
For example: 'bingetracker add season [arguments]' or 'bingetracker add episode [arguments].
See the help commands for the respective 'season' and 'episode' subcommands for more information their usage.`,
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
