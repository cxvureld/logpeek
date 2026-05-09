package main

import (
	"flag"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/cxvureld/logpeek/internal/analyzer"
	"github.com/cxvureld/logpeek/internal/parser"
	"github.com/cxvureld/logpeek/internal/ui"
)

func main() {
	var filePath string
	var format string
	flag.StringVar(&filePath, "file", "", "Path to log file (leave empty for stdin)")
	flag.StringVar(&format, "format", "nginx", "Log format: nginx, json")
	flag.Parse()

	var source *os.File
	var err error
	if filePath != "" {
		source, err = os.Open(filePath)
		if err != nil {
			log.Fatalf("Error opening file: %v", err)
		}
		defer source.Close()
	} else {
		source = os.Stdin
	}

	var p parser.Parser
	switch format {
	case "json":
		p = parser.NewJSONParser()
	case "nginx":
		nginxParser, err := parser.NewNginxParser()
		if err != nil {
			log.Fatalf("Error creating nginx parser: %v", err)
		}
		p = nginxParser
	default:
		log.Fatalf("Unknown format: %s. Supported: nginx, json", format)
	}

	entriesCh := make(chan *parser.LogEntry, 100)
	errCh := make(chan error, 10)

	go parser.ReadAndParse(source, p, entriesCh, errCh)

	stats := analyzer.NewStats()
	go analyzer.RunProcessor(entriesCh, stats, 0)

	go func() {
		for e := range errCh {
			log.Printf("Parse error: %v\n", e)
		}
	}()

	model := ui.NewModel(stats)
	program := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Printf("Error running TUI: %v", err)
		os.Exit(1)
	}
}
