/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	// "errors"
	"fmt"
	// "strconv"

	// "github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"gothub-team/gotiac/features"
	"gothub-team/gotiac/util"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new gotiac project.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(util.Logo)
		fmt.Println(`Welcome to gotiac! Let's create the infrastructure for your got project.`)
		fmt.Println()

		items, err := util.SelectItems(`Select the features you want to initiate.`, 0, []*util.Item{
			{ID: "CDN"},
			{ID: "Storage"},
		})

		// Label:     "Select the features you want to initiate.",

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		for _, item := range items {
			switch item.ID {
			case "CDN":
				features.Cdn()
			case "Storage":
				features.Storage()
			default:
				fmt.Printf("Unknown feature: %s\n", item.ID)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
