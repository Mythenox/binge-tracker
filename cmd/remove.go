/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var removeCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Removes a show, season, or episode",
	Long: `This command does nothing on its own, and must be used with the 'show', 'season', or 'episode' subcommands.
See the help commands for the respective subcommands for more information their usage.`,
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// ask for confirmation?

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
