/*
Copyright © 2024 Zhang Zhihui <ZhangZhihuiAAA@126.com>
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
        deleteRun()
    },
}

func init() {
    rootCmd.AddCommand(deleteCmd)

    deleteCmd.Flags().StringVarP(&key, "key", "k", "", "Key to delete")
    deleteCmd.MarkFlagRequired("key")
}

func deleteRun() {
    _, ok := index[key]
    if ok {
        fmt.Println("Found key:", key)
        delete(index, key)
    } else {
        fmt.Println("Key not found:", key)
        return
    }

    // Now, delete it from data
    for i, k := range data {
        if k.Filename == key {
            data = slices.Delete(data, i, i + 1)
            break
        }
    }

    err := saveJSONFile(JSONFILE)
    if err != nil {
        fmt.Println("Error saving data:", err)
    }

    fmt.Println("Deleted key:", key)
}