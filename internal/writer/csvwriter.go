package writer

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"sync"

	"github.com/ahuangg/gh-crawler/internal/models"
)

type CSVWriter struct {
    outputDir string
    mutex     sync.Mutex
    file      *os.File
    writer    *csv.Writer
}

func NewCSVWriter(outputDir, location string) (*CSVWriter, error) {
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return nil, err
    }

    file, err := os.Create(outputDir + "/" + location + "_users.csv")
    if err != nil {
        return nil, err
    }

    writer := csv.NewWriter(file)
    return &CSVWriter{outputDir: outputDir, file: file, writer: writer}, nil
}

func (w *CSVWriter) WriteUser(user *models.User) error {
    w.mutex.Lock()
    defer w.mutex.Unlock()

    langJSON, err := json.Marshal(user.LanguageStats)
    if err != nil {
        return err
    }

    record := []string{user.Username, user.Location, string(langJSON)}
    if err := w.writer.Write(record); err != nil {
        return err
    }
    w.writer.Flush()
    return nil
}

func (w *CSVWriter) Close() error {
    return w.file.Close()
}