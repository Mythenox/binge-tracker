/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/mythenox/bingetracker/internal/handler"
	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets the entire database",
	Long:  `This command removes all shows from the database, along with their respective seasons and episodes.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Print("The following command will completely reset the database, including all seasons and shows. Continue? (Y/n) ")
		userInput := ""
		fmt.Scanln(&userInput)

		if userInput == "Y" || userInput == "y" {
			err := handler.HandlerReset(cmd.Context(), s)
			if err != nil {
				fmt.Println("Error resetting database:", err)
				return err
			}
		} else {
			fmt.Println("Reset aborted.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// resetCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// resetCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
