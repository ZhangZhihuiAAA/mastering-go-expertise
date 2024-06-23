/*
Copyright © 2024 Zhang Zhihui <ZhangZhihuiAAA@126.com>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
    Use:   "list",
    Short: "list command",
    Long:  `A longer description of the list command`,
    Run: func(cmd *cobra.Command, args []string) {
        list()
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}

func list() {
    sort.Sort(data)
    text, err := PrettyPrintJSONStream(data)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println(text)

    logger.Info(fmt.Sprintf("%d records in total.", len(data)))
}

// PrettyPrintJSONStream pretty prints the contents of the phone book
func PrettyPrintJSONStream(data any) (string, error) {
    buffer := new(bytes.Buffer)
    encoder := json.NewEncoder(buffer)
    encoder.SetIndent("", "\t")

    err := encoder.Encode(data)
    if err != nil {
        return "", err
    }

    return buffer.String(), nil
}

// Implement sort.Interface
func (s DFslice) Len() int {
    return len(s)
}

func (s DFslice) Less(i, j int) bool {
    if s[i].Mean == s[j].Mean {
        return s[i].StdDev < s[j].StdDev
    }

    return s[i].Mean < s[j].Mean
}

func (s DFslice) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
