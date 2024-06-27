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

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
    Use:   "delete",
    Short: "Delete an entry",
    Long: `This commands deletes an existing entry from 
    the statistics web application given a name.`,
    Run: func(cmd *cobra.Command, args []string) {
        server := viper.GetString("server")
		port := viper.GetString("port")

		dataset, _ := cmd.Flags().GetString("dataset")
		if dataset == "" {
			fmt.Println("Dataset not provided!")
			return
		}

		// Create request
		url := "http://" + server + ":" + port + "/delete/" + dataset

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
    rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringP("dataset", "d", "", "Dataset name to delete")
}
