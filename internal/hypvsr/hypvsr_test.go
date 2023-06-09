package hypvsr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockKindManager struct{}

func (m *MockKindManager) Create() error                    { return nil }
func (m *MockKindManager) Start() error                     { return nil }
func (m *MockKindManager) Stop() error                      { return nil }
func (m *MockKindManager) ForceStop() error                 { return nil }
func (m *MockKindManager) Status() (string, error)          { return "", nil }
func (m *MockKindManager) Delete() error                    { return nil }
func (m *MockKindManager) CopyContent(string, string) error { return nil }
func (m *MockKindManager) Prepare() error                   { return nil }
func (m *MockKindManager) DownloadImage() error             { return nil }
func (m *MockKindManager) CreateDisks() error               { return nil }
func (m *MockKindManager) Wait() error                      { return nil }
func (m *MockKindManager) Shell() error                     { return nil }
func (m *MockKindManager) RunInitScripts() error            { return nil }
func (m *MockKindManager) GetVMs() []Machine                { return nil }
func (m *MockKindManager) CreateDir() error                 { return nil }

func TestConvertMemory(t *testing.T) {
	// Test cases for converting memory sizes
	testCases := []struct {
		memory   string
		expected string
	}{
		{memory: "1G", expected: "1024"},
		{memory: "2G", expected: "2048"},
		{memory: "512M", expected: "512"},
		{memory: "1024", expected: "1024"},
	}

	// Iterate over test cases
	for _, tc := range testCases {
		// Call the convertMemory function
		result, err := convertMemory(tc.memory)

		// Assert that the error is nil
		assert.Nil(t, err)

		// Assert that the converted memory value matches the expected value
		assert.Equal(t, tc.expected, result)
	}
}
