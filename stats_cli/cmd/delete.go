/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"slices"

	"github.com/spf13/cobra"
)

var key string

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete command",
	Long: `A longer description of the delete command.`,
	Run: func(cmd *cobra.Command, args []string) {
		_, ok := index[key]
		if ok {
			logger.Info(fmt.Sprintf("Found key: %s", key))
			fmt.Println("Found key:", key)
			delete(index, key)
		} else {
			logger.Info(fmt.Sprintf("%s not found!", key))
			return
		}

		// Now, delete it from data
		fmt.Println(data)
		for i, k := range data {
			if k.Filename == key {
				data = slices.Delete(data, i, i + 1)
				break
			}
		}

		err := saveJSONFile(JSONFILE)
		if err != nil {
			logger.Warn("Error saving data:", err)
		}

		logger.Info(fmt.Sprintf("Deleted key %s:", key))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	deleteCmd.Flags().StringVarP(&key, "key", "k", "", "Key to delete")
	deleteCmd.MarkFlagRequired("key")
}
