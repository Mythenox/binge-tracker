/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mythenox/binge-tracker/internal/handler"
	"github.com/spf13/cobra"
)

// seasonsCmd represents the seasons command
var listSeasonsCmd = &cobra.Command{
	Use:   "seasons <show title>",
	Short: "Lists all seasons of the specified show",
	Long:  `Lists all seasons of the specified show, as well as episode watch progress for each season.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		showTitle := args[0]
		return handler.HandlerListSeasons(cmd.Context(), s, showTitle)
	},
}

func init() {
	listCmd.AddCommand(listSeasonsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seasonsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seasonsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
