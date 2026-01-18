package main

import (
	"os"
	"testing"
)

// Helper to clean up test files
func setupTestFile(t *testing.T) {
	os.Remove("test_tasks.json") // Ensure clean state
}

func TestAddTask(t *testing.T) {
	setupTestFile(t)
	defer os.Remove("test_tasks.json")

	tests := []struct {
		name        string
		description string
		wantErr     bool
	}{
		{"Valid task", "Buy coffee", false},
		{"Empty task (UC-4)", "", true},
		{"Whitespace only (UC-4)", "   ", true},
		{"Too long task", string(make([]byte, 301)), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := addTask(12345, tt.description)
			if (err != nil) != tt.wantErr {
				t.Errorf("addTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	setupTestFile(t)
	defer os.Remove("test_tasks.json")

	// Add a dummy task first
	addTask(12345, "Test Task")

	tests := []struct {
		name    string
		number  int
		wantErr bool
	}{
		{"Delete valid index", 1, false},
		{"Delete zero (Invalid)", 0, true},
		{"Delete out of bounds", 99, true},
		{"Delete negative", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := deleteTask(12345, tt.number)
			if (err != nil) != tt.wantErr {
				t.Errorf("deleteTask() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
