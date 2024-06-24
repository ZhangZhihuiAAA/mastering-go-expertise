/*
Copyright © 2024 Zhang Zhihui <ZhangZhihuiAAA@126.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Entry struct {
    Filename string  `json:"filename"`
    Len      int     `json:"length"`
    Minimum  float64 `json:"minimum"`
    Maximum  float64 `json:"maximum"`
    Mean     float64 `json:"mean"`
    StdDev   float64 `json:"stddev"`
}

type DFslice []Entry

// JSONFILE resides in the current directory
const JSONFILE = "./data.json"

// Global variables
var data = DFslice{}
var enableLogging bool
var logger *slog.Logger
var index map[string]int

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
    Use:   "stats",
    Short: "Statistics application",
    Long:  `The statistics application`,
    PersistentPreRun: func(cmd *cobra.Command, args []string) {
        setDefaultLogger()
    },
    Run: func(cmd *cobra.Command, args []string) {},
}

func init() {
    rootCmd.PersistentFlags().BoolVarP(&enableLogging, "log", "l", true, "Logging information")

    err := readJSONFile(JSONFILE)
    if err != nil && strings.Contains(err.Error(), "no such file") {
        // Create the file if not exist.
        saveJSONFile(JSONFILE)
    // io.EOF is fine because it means the file is empty.
    } else if err != nil && err != io.EOF {
        fmt.Println(err)
        return
    }

    createIndex()
}

func setDefaultLogger() {
    logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
    if !enableLogging {
        logger = slog.New(slog.NewTextHandler(io.Discard, nil))
    }
    slog.SetDefault(logger)
}

func createIndex() {
    index = make(map[string]int)
    for i, k := range data {
        key := k.Filename
        index[key] = i
    }
}

func saveJSONFile(filepath string) error {
    f, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer f.Close()

    err = Serialize(&data, f)
    return err
}

func readJSONFile(filepath string) error {
    _, err := os.Stat(filepath)
    if err != nil {
        return err
    }

    f, err := os.Open(filepath)
    if err != nil {
        return err
    }
    defer f.Close()

    err = DeSerialize(&data, f)
    if err != nil {
        return err
    }

    return nil
}

// DeSerialize decodes a serialized slice with JSON records
func DeSerialize(s any, r io.Reader) error {
    decoder := json.NewDecoder(r)
    return decoder.Decode(s)
}

// Serialize serializes a slice with JSON records
func Serialize(s any, w io.Writer) error {
    encoder := json.NewEncoder(w)
    return encoder.Encode(s)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
    err := rootCmd.Execute()
    if err != nil {
        os.Exit(1)
    }
}
