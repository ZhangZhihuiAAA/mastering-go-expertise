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

// listCmd represents the list command
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List all entries",
    Long:  `List all entries of the statistics web application.`,
    Run: func(cmd *cobra.Command, args []string) {
        server := viper.GetString("server")
        port := viper.GetString("port")

        // Create request
        url := "http://" + server + ":" + port + "/list"
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

        fmt.Println("List of entries:")
        fmt.Print(string(data))
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}
