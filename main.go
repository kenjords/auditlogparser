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

// LogEntry represents a single log entry
type LogEntry struct {
	LogLevel         string `json:"log_level"`
	Stage            string `json:"stage"`
	DateTime         string `json:"date_time"`
	Class            string `json:"class"`
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
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
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
		for _, field := range pipeFields {
			keyValue := strings.SplitN(field, ":", 2)
			if len(keyValue) < 2 {
				continue
			}
			key := keyValue[0]
			value := keyValue[1]

			switch key {
			case "host":
				logEntry.Host = strings.TrimPrefix(value, "/")
			case "source":
				logEntry.Source = strings.TrimLeft(value, "/")
			case "user":
				logEntry.User = value
			case "authenticated":
				logEntry.Authenticated = value
			case "timestamp":
				logEntry.Timestamp = value
			case "category":
				logEntry.Category = value
			case "type":
				logEntry.Type = value
			case "batch":
				logEntry.Batch = value
			case "ks":
				logEntry.KS = value
			case "cf":
				logEntry.CF = value
			case "operation":
				opMsgIdx := strings.Index(value, ".")
				logEntry.OperationMessage = strings.TrimRight(value[:opMsgIdx+1], ".")
				logEntry.Operation = strings.TrimLeft(value[opMsgIdx+2:], ".")
			case "consistency level":
				logEntry.ConsistencyLevel = value
			}
		}

		jsonData, err := json.MarshalIndent(logEntry, "", "\t")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(jsonData))
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
