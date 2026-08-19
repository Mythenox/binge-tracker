/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mythenox/binge-tracker/internal/handler"
	"github.com/spf13/cobra"
)

// showsCmd represents the shows command
var listShowsCmd = &cobra.Command{
	Use:   "shows",
	Short: "Lists all initialized shows",
	Long:  `Lists all initialized shows, along with the number of seasons and total episode watch progress for each show.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return handler.HandlerListShows(cmd.Context(), s)
	},
}

func init() {
	listCmd.AddCommand(listShowsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
