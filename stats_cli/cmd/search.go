/*
Copyright © 2024 Zhang Zhihui <ZhangZhihuiAAA@126.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sid string

// searchCmd represents the search command
var searchCmd = &cobra.Command{
    Use:   "search",
    Short: "search command",
    Long:  `A longer description of the search command.`,
    Run: func(cmd *cobra.Command, args []string) {
        searchRun()
    },
}

func init() {
    rootCmd.AddCommand(searchCmd)

    searchCmd.Flags().StringVarP(&sid, "sid", "s", "", "Search key")
    searchCmd.MarkFlagRequired("sid")
}

func searchRun() {
    for i, k := range data {
        if k.Filename == sid {
            str, err := PrettyPrintJSONStream(data[i])
            if err != nil {
                fmt.Println(err)
            } else {
                fmt.Print(str)
            }
            return
        }
    }

    fmt.Printf("Key not found: %s\n", sid)
}