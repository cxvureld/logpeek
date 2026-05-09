package analyzer

import (
	"sync"
	"time"

	"github.com/cxvureld/logpeek/internal/parser"
)

type Stats struct {
	mu            sync.RWMutex
	TotalRequests int64
	StatusCodes   map[int]int64
	TopIPs        map[string]int64
	TopEndpoints  map[string]int64
	BytesSent     int64
}

func NewStats() *Stats {
	return &Stats{
		StatusCodes:  make(map[int]int64),
		TopIPs:       make(map[string]int64),
		TopEndpoints: make(map[string]int64),
	}
}

func (s *Stats) Add(entry *parser.LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.TotalRequests++
	s.StatusCodes[entry.Status]++
	s.TopIPs[entry.RemoteAddr]++
	s.BytesSent += int64(entry.BodyBytes)
	s.TopEndpoints[extractEndpoint(entry.Request)]++
}

func (s *Stats) Snapshot() *Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ss := &Stats{
		TotalRequests: s.TotalRequests,
		StatusCodes:   make(map[int]int64),
		TopIPs:        make(map[string]int64),
		TopEndpoints:  make(map[string]int64),
		BytesSent:     s.BytesSent,
	}
	for k, v := range s.StatusCodes {
		ss.StatusCodes[k] = v
	}
	for k, v := range s.TopIPs {
		ss.TopIPs[k] = v
	}
	for k, v := range s.TopEndpoints {
		ss.TopEndpoints[k] = v
	}
	return ss
}

func RunProcessor(entries <-chan *parser.LogEntry, stats *Stats, tickInterval time.Duration) {
	for entry := range entries {
		stats.Add(entry)
	}

}

func extractEndpoint(request string) string {
	parts := splitRequest(request)
	if len(parts) >= 2 {
		return parts[1]
	}
	return request
}

func splitRequest(request string) []string {
	var parts []string
	current := ""
	for _, ch := range request {
		if ch == ' ' {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}
