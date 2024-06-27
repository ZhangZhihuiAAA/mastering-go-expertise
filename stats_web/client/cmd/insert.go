/*
Copyright © 2024 Zhang Zhihui <ZhangZhihuiAAA@126.com>
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
    Use:   "insert",
    Short: "Insert a new entry",
    Long: `This command inserts a new entry to the 
    statistics web application.`,
    Run: func(cmd *cobra.Command, args []string) {
        server := viper.GetString("server")
		port := viper.GetString("port")

		dataset, _ := cmd.Flags().GetString("dataset")
		if dataset == "" {
			fmt.Println("Dataset not provided!")
			return
		}

		valuesStr, _ := cmd.Flags().GetString("values")
		if valuesStr == "" {
			fmt.Println("Data not provided!")
			return
		}

		values := strings.Split(valuesStr, ",")
		vSend := ""
		for _, v := range values {
			_, err := strconv.ParseFloat(v, 64)
			if err == nil {
				vSend = vSend + "/" + v
			}
		}

		url := "http://" + server + ":" + port + "/insert/" + dataset + "/" + vSend
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
    rootCmd.AddCommand(insertCmd)
	insertCmd.Flags().StringP("dataset", "d", "", "Dataset name")
	insertCmd.Flags().StringP("values", "v", "", "List of values separated by comma")
}
