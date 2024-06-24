/*
Copyright © 2024 Zhang Zhihui <ZhangZhihuiAAA@126.com>
*/
package cmd

import (
	"encoding/csv"
	"errors"
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
    Long:  `A longer description of the insert command.`,
    Run: func(cmd *cobra.Command, args []string) {
        insertRun()
    },
}

func init() {
    rootCmd.AddCommand(insertCmd)
    insertCmd.Flags().StringVarP(&file, "file", "f", "", "Filename to process")
    insertCmd.MarkFlagRequired("file")
}

func insertRun() {
    if file == "" {
        fmt.Println("Need a file to read!")
        return
    }

    _, ok := index[file]
    if ok {
        fmt.Println("Found key:", file)
        delete(index, file)

        // Now, delete it from data
        for i, k := range data {
            if k.Filename == file {
                data = slices.Delete(data, i, i+1)
                break
            }
        }
    }

    err := ProcessFile(file)
    if err != nil {
        fmt.Println(err)
        return
    }

    err = saveJSONFile(JSONFILE)
    if err != nil {
        fmt.Printf("Error saving data: %s", err)
    }
}

func readFile(filepath string) (values []float64, err error) {
    _, err = os.Stat(filepath)
    if err != nil {
        return nil, err
    }

    f, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    lines, err := csv.NewReader(f).ReadAll()
    if err != nil {
        return nil, err
    }

    values = make([]float64, 0, len(lines))
    for i, line := range lines {
        value, err4 := strconv.ParseFloat(line[0], 64)
        if err4 != nil {
            logger.Error(fmt.Sprintln("Invalid value", line[0], "in line", i, err4))

            if err == nil {
                err = errors.New("failed to read at least one value")
            }
        }
        values = append(values, value)
    }

    return
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
        squared += math.Pow(x[i]-meanValue, 2)
    }
    standardDeviation := math.Sqrt(squared / float64(len(x)))

    return meanValue, standardDeviation
}

func ProcessFile(file string) error {
    currentFile := Entry{}
    currentFile.Filename = file

    values, err := readFile(file)
    if err != nil {
        return err
    }

    currentFile.Len = len(values)
    currentFile.Minimum = slices.Min(values)
    currentFile.Maximum = slices.Max(values)
    currentFile.Mean, currentFile.StdDev = stdDev(values)

    data = append(data, currentFile)

    return nil
}
