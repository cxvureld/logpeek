package parser

import (
	"encoding/json"
	"fmt"
	"time"
)

type JSONLogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Source    string `json:"source"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	LatencyMs int    `json:"latency_ms"`
	ClientIP  string `json:"client_ip"`
	Bytes     int    `json:"bytes"`
}

type JSONParser struct{}

func NewJSONParser() *JSONParser {
	return &JSONParser{}
}

func (p *JSONParser) ParseLine(line string) (*LogEntry, error) {
	var jsonEntry JSONLogEntry

	if err := json.Unmarshal([]byte(line), &jsonEntry); err != nil {
		return nil, fmt.Errorf("failed to parse JSON log: %w", err)
	}

	var t time.Time
	var err error

	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
	}

	for _, format := range formats {
		t, err = time.Parse(format, jsonEntry.Timestamp)
		if err == nil {
			break
		}
	}

	if err != nil {
		t = time.Now()
	}

	request := fmt.Sprintf("%s %s", jsonEntry.Method, jsonEntry.Path)
	if jsonEntry.Method == "" {
		request = "-"
	}

	return &LogEntry{
		RemoteAddr: jsonEntry.ClientIP,
		TimeLocal:  t,
		Request:    request,
		Status:     jsonEntry.Status,
		BodyBytes:  jsonEntry.Bytes,
		HTTPRefer:  jsonEntry.Source,
		UserAgent:  jsonEntry.Level,
	}, nil
}
