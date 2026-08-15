/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/mythenox/binge-tracker/internal/handler"
	"github.com/spf13/cobra"
)

// seasonsCmd represents the seasons command
var seasonsCmd = &cobra.Command{
	Use:   "seasons <show title>",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		showTitle := args[0]
		return handler.HandlerListSeasons(cmd.Context(), s, showTitle)
	},
}

func init() {
	listCmd.AddCommand(seasonsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// seasonsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// seasonsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
