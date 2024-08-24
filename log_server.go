package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Record struct {
	Value  []byte `json:"value"`
	Offset uint64 `json:"offset"` // Offset in the log
}

type Log struct {
	mu      sync.Mutex
	records []Record // Slice of records
}

// Append a record to the log
func (l *Log) Append(record Record) uint64 {
	l.mu.Lock()
	defer l.mu.Unlock()

	record.Offset = uint64(len(l.records))
	l.records = append(l.records, record)

	return record.Offset
}

// Read a record from the log
func (l *Log) Read(offset uint64) (Record, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if offset >= uint64(len(l.records)) {
		return Record{}, fmt.Errorf("offset out of range")
	}
	return l.records[offset], nil
}

func (l *Log) handleProduce(w http.ResponseWriter, r *http.Request) {
	var record Record

	err := json.NewDecoder(r.Body).Decode(&record)
	if err != nil {
		log.Printf("Error decoding record: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Record: %+v", record)
	offset := l.Append(record)
	json.NewEncoder(w).Encode(map[string]uint64{"offset": offset})
}

func (l *Log) handleConsume(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Offset uint64 `json:"offset"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	record, err := l.Read(req.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(record)
}

// ServeHTTP method for the Log type
func (l *Log) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		l.handleProduce(w, r)
	case "GET":
		l.handleConsume(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	log := &Log{}
	http.ListenAndServe(":8080", log)
}
