/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/mythenox/binge-tracker/internal/handler"
	"github.com/spf13/cobra"
)

// nextCmd represents the next command
// play next <show name> [-- player flags] OR play next <show name> <sXX> [-- player flags]
var nextCmd = &cobra.Command{
	Use:   "next <show name> [-- player flags] OR play next <show name> <sXX> [-- player flags]",
	Short: "Plays the next episode of a show",
	Long: `This command plays the next unwatched episode of the specified show.
If a season identifier is provided, the command will ignore unwatched episodes that precede that season.
Arguments can be passed through to the player by use of the '--' separator; all arguments past this separator will be passed to the player.
Usage example: 'bingetracker play next Twin Peaks -- --mute=yes'`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dashIndex := cmd.ArgsLenAtDash()
		if dashIndex == -1 {
			dashIndex = len(args)
		}

		if dashIndex > 2 {
			return fmt.Errorf("accepts at most 2 args before '--', received %d", dashIndex)
		}

		cmdArgs := args[:dashIndex]
		playerArgs := args[dashIndex:]

		showTitle := cmdArgs[0]
		verboseInput := dashIndex == 2
		var seasonNumber int
		// seasonNumber is ignored by handler if verboseInput == false

		if verboseInput {
			episodeRegex := regexp.MustCompile(`(?i)^s(\d+)$`)

			seasonIdentifier := cmdArgs[1]

			matches := episodeRegex.FindStringSubmatch(seasonIdentifier)
			if len(matches) != 2 {
				return errors.New("Invalid format")
			}

			var err error
			seasonNumber, err = strconv.Atoi(matches[1])
			if err != nil {
				return err
			}
		}

		return handler.HandlerPlayNext(cmd.Context(), s, showTitle, seasonNumber,
			verboseInput, playerArgs)
	},
}

func init() {
	playCmd.AddCommand(nextCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nextCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nextCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
