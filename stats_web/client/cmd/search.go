/*
Copyright © 2024 Zhang Zhihui <ZhangZhihuiAAA@126.com>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
    Use:   "search",
    Short: "Search for an entry",
    Long:  `This command searches for an entry by a given name.`,
    Run: func(cmd *cobra.Command, args []string) {
        server := viper.GetString("server")
        port := viper.GetString("port")

        dataset, _ := cmd.Flags().GetString("dataset")
        if dataset == "" {
            fmt.Println("Dataset name not provided!")
            return
        }

        url := "http://" + server + ":" + port + "/search/" + dataset
        resp, err := http.Get(url)
        if err != nil {
            fmt.Println(err)
            return
        }

        if resp.StatusCode != http.StatusOK {
            fmt.Println("Status code:", resp.StatusCode)
            return
        }

        data, err := io.ReadAll(resp.Body)
        if err != nil {
            fmt.Println(err)
            return
        }

        fmt.Print(string(data))
    },
}

func init() {
    rootCmd.AddCommand(searchCmd)
    searchCmd.Flags().StringP("dataset", "d", "", "Dataset name to search")
}
