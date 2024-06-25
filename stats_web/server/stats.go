package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"
)

type Entry struct {
    Name    string
    Len     int
    Minimum float64
    Maximum float64
    Mean    float64
    StdDev  float64
}

func process(file string, values []float64) Entry {
    currentEntry := Entry{}
    currentEntry.Name = file
    currentEntry.Len = len(values)
    currentEntry.Minimum = slices.Min(values)
    currentEntry.Maximum = slices.Max(values)
    currentEntry.Mean, currentEntry.StdDev = stdDev(values)

    return currentEntry
}

func stdDev(x []float64) (float64, float64) {
    var sum float64
    for _, val := range x {
        sum += val
    }
    meanValue := sum / float64(len(x))

    // Standard deviation
    var squared float64
    for i := 0; i < len(x); i++ {
        squared += math.Pow((x[i] - meanValue), 2)
    }
    standardDeviation := math.Sqrt(squared / float64(len(x)))

    return meanValue, standardDeviation
}

// JSONFILE resides in the current directory
const JSONFILE = "./data.json"

var data []*Entry
var index map[string]int

func createIndex() {
    index = make(map[string]int)
    for i, k := range data {
        key := k.Name
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
        if os.IsNotExist(err) {
            f, err1 := os.OpenFile(filepath, os.O_RDONLY | os.O_CREATE, 0644)
            if err1 != nil {
                return err1
            }
            defer f.Close()
            return nil
        }
        return err
    }

    f, err := os.Open(filepath)
    if err != nil {
        return err
    }
    defer f.Close()

    err = DeSerialize(&data, f)
    return err
}

func insert(e *Entry) error {
    // If it already exists, do not add it
    _, ok := index[e.Name]
    if ok {
        return fmt.Errorf("%s already exists", e.Name)
    }

    data = append(data, e)
    // Update the index
    createIndex()

    err := saveJSONFile(JSONFILE)
    return err
}

func deleteEntry(key string) error {
    i, ok := index[key]
    if !ok {
        return fmt.Errorf("%s cannot be found", key)
    }
    data = append(data[:i], data[i+1:]...)
    // Update the index
    delete(index, key)

    err := saveJSONFile(JSONFILE)
    return err
}

func search(key string) *Entry {
    i, ok := index[key]
    if !ok {
        return nil
    }

    return data[i]
}

func list() string {
    sb := strings.Builder{}
    for _, k := range data {
        sb.WriteString(fmt.Sprintf("%s\t%d\t%f\t%f\n", k.Name, k.Len, k.Mean, k.StdDev))
    }
    return sb.String()
}

func main() {
    err := readJSONFile(JSONFILE)
    if err != nil && err != io.EOF {
        fmt.Println("Error:", err)
        return
    }
    createIndex()

    mux := http.NewServeMux()
    server := &http.Server{
        Addr: PORT,
        Handler: mux,
        IdleTimeout: 10 * time.Second,
        ReadTimeout: time.Second,
        WriteTimeout: time.Second,
    }

    mux.Handle("/list", http.HandlerFunc(listHandler))
    mux.Handle("/insert/", http.HandlerFunc(insertHandler))
    mux.Handle("/insert", http.HandlerFunc(insertHandler))
    mux.Handle("/delete/", http.HandlerFunc(deleteHandler))
    mux.Handle("/delete", http.HandlerFunc(deleteHandler))
    mux.Handle("/search/", http.HandlerFunc(searchHandler))
    mux.Handle("/search", http.HandlerFunc(searchHandler))
    mux.Handle("/status", http.HandlerFunc(statusHandler))
    mux.Handle("/", http.HandlerFunc(defaultHandler))

    fmt.Println("Ready to serve at", PORT)
    err = server.ListenAndServe()
    if err != nil {
        fmt.Println(err)
        return
    }
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

