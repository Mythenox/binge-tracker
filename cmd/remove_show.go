/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/mythenox/bingetracker/internal/handler"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var removeShowCmd = &cobra.Command{
	Use:   "show <show title>",
	Short: "Removes the specified show from the database.",
	Long: `This command removes the specified show from the database, along with all of its seasons and episodes.
Usage example: 'bingetracker remove show Twin Peaks'`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		showTitle := strings.ToLower(strings.Join(args, " "))

		return handler.HandlerRemoveShow(cmd.Context(), s, showTitle)
	},
}

func init() {
	removeCmd.AddCommand(removeShowCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
