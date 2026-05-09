package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"
)

type LogEntry struct {
	RemoteAddr string
	TimeLocal  time.Time
	Request    string
	Status     int
	BodyBytes  int
	HTTPRefer  string
	UserAgent  string
}

type NginxParser struct {
	re *regexp.Regexp
}

func NewNginxParser() (*NginxParser, error) {
	re, err := regexp.Compile(`^(\S+) - - \[([^\]]+)\] "([^"]*)" (\d{3}) (\d+) "([^"]*)" "([^"]*)"$`)
	if err != nil {
		return nil, fmt.Errorf("failed to compile nginx regexp: %w", err)
	}
	return &NginxParser{re: re}, nil
}

func (p *NginxParser) ParseLine(line string) (*LogEntry, error) {
	matches := p.re.FindStringSubmatch(line)
	if matches == nil {
		return nil, fmt.Errorf("line does not match nginx format: %s", line)
	}

	t, err := time.Parse("02/Jan/2006:15:04:05 -0700", matches[2])
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %w", err)
	}

	status, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("failed to parse status: %w", err)
	}

	bodyBytes, err := strconv.Atoi(matches[5])
	if err != nil {
		bodyBytes = 0
	}

	return &LogEntry{
		RemoteAddr: matches[1],
		TimeLocal:  t,
		Request:    matches[3],
		Status:     status,
		BodyBytes:  bodyBytes,
		HTTPRefer:  matches[6],
		UserAgent:  matches[7],
	}, nil
}

type Parser interface {
	ParseLine(line string) (*LogEntry, error)
}

func ReadAndParse(r io.Reader, p Parser, out chan<- *LogEntry, errCh chan<- error) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		entry, err := p.ParseLine(scanner.Text())
		if err != nil {
			errCh <- err
			continue
		}
		out <- entry
	}
	close(out)
}
