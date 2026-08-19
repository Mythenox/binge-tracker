/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/mythenox/bingetracker/internal/handler"
	"github.com/spf13/cobra"
)

// bingetracker config set video_player=mpv

// configCmd represents the config command
var configSetCmd = &cobra.Command{
	Use:   "set <key>=<value>",
	Short: "Set value of certain config settings",
	Long: `At the moment, this command only allows you to change the values of the count_partial_progress and video_player settings.
Changing the values of the other settings will completely break the program, so please refrain from doing so.
The keys and values must be set in <key>=<value> format.
Usage example: 'bingetracker config set count_partial_progress=true'`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		mutableKeys := []string{"count_partial_progress", "video_player"}

		configChanges := make(map[string]string)
		for _, arg := range args {
			if !strings.Contains(arg, "=") {
				return errors.New("Input must be of the form key=value")
			}
			splitInput := strings.Split(arg, "=")
			key, value := strings.ToLower(splitInput[0]), strings.ToLower(splitInput[1])

			if !slices.Contains(mutableKeys, key) {
				return fmt.Errorf("The key %s cannot be altered.", key)
			}

			configChanges[key] = value
		}

		return handler.HandlerSetConfig(cmd.Context(), s, configChanges)
	},
}

func init() {
	configCmd.AddCommand(configSetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
