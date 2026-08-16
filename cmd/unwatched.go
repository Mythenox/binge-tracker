/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

// unwatchedCmd represents the unwatched command
var unwatchedCmd = &cobra.Command{
	Use:   "unwatched",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputIsRange := strings.Contains(args[1], "-")
		episodeMode := strings.Contains(args[1], "e")

		return setCompletion(cmd, args, s, inputIsRange, episodeMode, false)
	},
}

func init() {
	setCmd.AddCommand(unwatchedCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// unwatchedCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// unwatchedCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
