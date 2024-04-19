/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	// "errors"
	"fmt"
	// "strconv"

	"github.com/manifoldco/promptui"
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

		initProject()
		initFeatures()
	},
}

func initProject() {
	stagePrompt := promptui.Prompt{
		Label: "Enter the stage name you want to deploy to",
	}
	awsProfilePrompt := promptui.Prompt{
		Label: "Enter the AWS profile you want to use",
	}

	regions := []string{"us-east-1", "us-west-1", "us-west-2", "eu-west-1", "eu-central-1", "ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "sa-east-1"}
	awsRegionSelect := promptui.Select{
		Label:     "Enter the AWS region you want to use",
		Items:     regions,
		CursorPos: 4,
	}

	stage, _ := stagePrompt.Run()
	awsProfile, _ := awsProfilePrompt.Run()
	regionIdx, _, _ := awsRegionSelect.Run()
	awsRegion := regions[regionIdx]

	fmt.Printf("Initializing stage %s for profile %s in %s \n", stage, awsProfile, awsRegion)
}

func initFeatures() {
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
