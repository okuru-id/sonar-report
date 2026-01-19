package report

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Storage manages report file storage and history
type Storage struct {
	basePath string
	mu       sync.RWMutex
	records  []ReportRecord
}

// NewStorage creates a new storage manager
func NewStorage(basePath string) (*Storage, error) {
	// Ensure directory exists
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	s := &Storage{
		basePath: basePath,
		records:  []ReportRecord{},
	}

	// Load existing records
	if err := s.loadRecords(); err != nil {
		// Ignore error on first run
		s.records = []ReportRecord{}
	}

	return s, nil
}

// Save saves a report and returns the record
func (s *Storage) Save(data *ReportData, content []byte, format string) (*ReportRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()[:8]
	timestamp := time.Now().Format("20060102-150405")

	var ext string
	switch format {
	case "pdf":
		ext = "pdf"
	default:
		ext = "md"
	}

	fileName := fmt.Sprintf("%s-%s-%s.%s", data.ProjectKey, timestamp, id, ext)
	filePath := filepath.Join(s.basePath, fileName)

	// Write file
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		return nil, fmt.Errorf("failed to write report file: %w", err)
	}

	// Create record
	record := ReportRecord{
		ID:          id,
		ProjectKey:  data.ProjectKey,
		ProjectName: data.ProjectName,
		Branch:      data.Branch,
		Format:      format,
		FileName:    fileName,
		FilePath:    filePath,
		FileSize:    int64(len(content)),
		GeneratedAt: time.Now(),
	}

	// Add to records
	s.records = append([]ReportRecord{record}, s.records...)

	// Save records
	if err := s.saveRecords(); err != nil {
		return nil, fmt.Errorf("failed to save records: %w", err)
	}

	return &record, nil
}

// GetHistory returns all report records
func (s *Storage) GetHistory() []ReportRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Return copy
	result := make([]ReportRecord, len(s.records))
	copy(result, s.records)
	return result
}

// GetRecord returns a specific report record
func (s *Storage) GetRecord(id string) (*ReportRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, r := range s.records {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("report not found: %s", id)
}

// GetFilePath returns the file path for a report
func (s *Storage) GetFilePath(id string) (string, error) {
	record, err := s.GetRecord(id)
	if err != nil {
		return "", err
	}
	return record.FilePath, nil
}

// Delete removes a report
func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, r := range s.records {
		if r.ID == id {
			// Delete file
			os.Remove(r.FilePath)

			// Remove from records
			s.records = append(s.records[:i], s.records[i+1:]...)

			// Save records
			return s.saveRecords()
		}
	}
	return fmt.Errorf("report not found: %s", id)
}

// ClearAll removes all reports
func (s *Storage) ClearAll() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, r := range s.records {
		os.Remove(r.FilePath)
	}

	s.records = []ReportRecord{}
	return s.saveRecords()
}

// CleanOld removes reports older than the specified number of days
func (s *Storage) CleanOld(days int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	cutoff := time.Now().AddDate(0, 0, -days)
	var newRecords []ReportRecord

	for _, r := range s.records {
		if r.GeneratedAt.After(cutoff) {
			newRecords = append(newRecords, r)
		} else {
			os.Remove(r.FilePath)
		}
	}

	s.records = newRecords
	return s.saveRecords()
}

func (s *Storage) loadRecords() error {
	indexPath := filepath.Join(s.basePath, "index.json")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &s.records)
}

func (s *Storage) saveRecords() error {
	// Sort by generated time descending
	sort.Slice(s.records, func(i, j int) bool {
		return s.records[i].GeneratedAt.After(s.records[j].GeneratedAt)
	})

	data, err := json.MarshalIndent(s.records, "", "  ")
	if err != nil {
		return err
	}

	indexPath := filepath.Join(s.basePath, "index.json")
	return os.WriteFile(indexPath, data, 0644)
}
