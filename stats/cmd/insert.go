/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"

	"github.com/spf13/cobra"
)

var file string

// insertCmd represents the insert command
var insertCmd = &cobra.Command{
    Use:   "insert",
    Short: "insert command",
    Long: `A longer description of the insert command.`,
    Run: func(cmd *cobra.Command, args []string) {
        if file == "" {
            logger.Info("Need a file to read!")
            return
        }

        _, ok := index[file]
        if ok {
            fmt.Println("Found key:", file)
            delete(index, file)

            // Now, delete it from data
            for i, k := range data {
				if k.Filename == file {
					data = slices.Delete(data, i, i + 1)
					break
				}
			}
        }
    },
}

func init() {
    rootCmd.AddCommand(insertCmd)

    // Here you will define your flags and configuration settings.

    // Cobra supports Persistent Flags which will work for this command
    // and all subcommands, e.g.:
    // insertCmd.PersistentFlags().String("foo", "", "A help for foo")

    // Cobra supports local flags which will only run when this command
    // is called directly, e.g.:
    // insertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func readFile(filepath string) ([]float64, error) {
	_, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return nil ,err
	}
	defer f.Close()

	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}

	values := make([]float64, 0, len(lines))
	for _, line := range lines {
		value, err := strconv.ParseFloat(line[0], 64)
		if err != nil {
			fmt.Println("Error reading:", line[0], err)
		}
		values = append(values, value)
	}

	return values, nil
}

func stdDev(x []float64) (float64, float64) {
	// Mean value
	var sum float64
	for _, val := range x {
		sum += val
	}
	meanValue := sum / float64(len(x))

	// Standard deviation
	var squared float64
	for i := 0; i < len(x); i++ {
		squared += math.Pow(x[i] - meanValue, 2)
	}
	standardDeviation := math.Sqrt(squared / float64(len(x)))

	return meanValue, standardDeviation
}

func ProcessFile(file string) error {
	currentFile := Entry{}
	currentFile.Filename = file

	values, err := readFile(file)
	if err != nil {
		return nil
	}

	currentFile.Len = len(values)
	currentFile.Minimum = slices.Min(values)
	currentFile.Maximum = slices.Max(values)
	currentFile.Mean, currentFile.StdDev = stdDev(values)

	data = append(data, currentFile)

	return nil
}