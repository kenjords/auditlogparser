package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// LogEntry represents a single log entry from the audit.log file
type LogEntry struct {
	LogLevel  string          `json:"log_level"`
	Stage     string          `json:"stage"`
	DateTime  string          `json:"date_time"`
	Class     string          `json:"class"`
	Operation OperationRecord `json:"operation"`
}

// OperationRecord represents the operation record from a log entry
type OperationRecord struct {
	Host             string `json:"host"`
	Source           string `json:"source"`
	User             string `json:"user"`
	Authenticated    string `json:"authenticated"`
	Timestamp        string `json:"timestamp"`
	Category         string `json:"category"`
	Type             string `json:"type"`
	Batch            string `json:"batch"`
	KS               string `json:"ks"`
	CF               string `json:"cf"`
	OperationMessage string `json:"operation_message"`
	Operation        string `json:"operation"`
	ConsistencyLevel string `json:"consistency_level"`
}

func main() {
	var auditLogPath string
	flag.StringVar(&auditLogPath, "file", "audit.log", "Path to the audit.log file")
	flag.Parse()

	file, err := os.Open(auditLogPath)
	if err != nil {
		log.Printf("Error opening file: %v", err)
		os.Exit(1)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
			os.Exit(2)
		}
	}()

	scanner := bufio.NewScanner(file)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 10*1024*1024)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		logEntry := LogEntry{}
		line := scanner.Text()
		spaceFields := strings.Fields(line)

		logEntry.LogLevel = spaceFields[0]
		logEntry.Stage = spaceFields[1]
		logEntry.DateTime = fmt.Sprintf("%s %s", spaceFields[2], strings.TrimRight(spaceFields[3], ","))
		logEntry.Class = spaceFields[4]

		pipeFields := strings.Split(line, "|")
		for i, field := range pipeFields {
			if strings.Contains(field, " - host:/") {
				fieldParts := strings.Split(field, " - ")
				field = fieldParts[1]
				pipeFields[i] = field
			}
			keyValue := strings.SplitN(field, ":", 2)
			if len(keyValue) < 2 {
				continue
			}
			key := keyValue[0]
			value := keyValue[1]

			switch key {
			case "host":
				logEntry.Operation.Host = strings.TrimPrefix(value, "/")
			case "source":
				logEntry.Operation.Source = strings.TrimPrefix(value, "/")
			case "user":
				logEntry.Operation.User = value
			case "authenticated":
				logEntry.Operation.Authenticated = value
			case "timestamp":
				logEntry.Operation.Timestamp = value
			case "category":
				logEntry.Operation.Category = value
			case "type":
				logEntry.Operation.Type = value
			case "batch":
				logEntry.Operation.Batch = value
			case "ks":
				logEntry.Operation.KS = value
			case "cf":
				logEntry.Operation.CF = value
			case "operation":
				opMsgIdx := strings.Index(value, ".")
				logEntry.Operation.OperationMessage = strings.TrimRight(value[:opMsgIdx+1], ".")
				logEntry.Operation.Operation = strings.TrimLeft(value[opMsgIdx+2:], ".")
			case "consistency level":
				logEntry.Operation.ConsistencyLevel = value
			}
		}

		jsonData, err := json.MarshalIndent(logEntry, "", "\t")
		if err != nil {
			log.Printf("Error marshalling JSON: %v", err)
			os.Exit(3)
		}
		fmt.Println(string(jsonData))
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error scanning file: %v", err)
		os.Exit(4)
	}
}
