package common

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestStructToJSONBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    []byte
		wantErr bool
	}{
		{
			name:    "integer",
			input:   1,
			want:    []byte("1"),
			wantErr: false,
		},
		{
			name:    "string",
			input:   "hello",
			want:    []byte(`"hello"`),
			wantErr: false,
		},
		{
			name:    "map",
			input:   map[string]string{"key": "value"},
			want:    []byte(`{"key":"value"}`),
			wantErr: false,
		},
		{
			name:    "slice",
			input:   []int{1, 2, 3},
			want:    []byte("[1,2,3]"),
			wantErr: false,
		},
		{
			name:    "struct",
			input:   struct{ Name string }{Name: "test"},
			want:    []byte(`{"Name":"test"}`),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StructToJSONBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("StructToJSONBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StructToJSONBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetFileByURL(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantErr    bool
	}{
		{
			name:       "invalid URL",
			url:        "://invalid",
			wantStatus: 0,
			wantErr:    true,
		},
		{
			name:       "connection refused",
			url:        "http://localhost:99999/nonexistent",
			wantStatus: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetFileByURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileByURL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetFileByURL_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test content"))
	}))
	defer server.Close()

	body, err := GetFileByURL(server.URL)
	if err != nil {
		t.Errorf("GetFileByURL() error = %v, wantErr false", err)
	}
	if string(body) != "test content" {
		t.Errorf("GetFileByURL() = %v, want 'test content'", string(body))
	}
}

func TestGetPubKey(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		want    string
		wantErr bool
	}{
		{
			name:    "invalid URL",
			url:     "://invalid",
			want:    "",
			wantErr: true,
		},
		{
			name:    "invalid JSON",
			url:     "http://localhost:99999/invalid",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPubKey(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPubKey() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("GetPubKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPubKey_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"value": "test-key-123"}`))
	}))
	defer server.Close()

	got, err := GetPubKey(server.URL)
	if err != nil {
		t.Errorf("GetPubKey() error = %v, wantErr false", err)
	}
	if got != "test-key-123" {
		t.Errorf("GetPubKey() = %v, want 'test-key-123'", got)
	}
}

func TestGetPubKey_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("not a json"))
	}))
	defer server.Close()

	_, err := GetPubKey(server.URL)
	if err == nil {
		t.Error("GetPubKey() want error for invalid JSON, got nil")
	}
}

func TestGetPubKey_MissingValueKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"other": "value"}`))
	}))
	defer server.Close()

	got, err := GetPubKey(server.URL)
	if err != nil {
		t.Errorf("GetPubKey() error = %v, wantErr false", err)
	}
	if got != "" {
		t.Errorf("GetPubKey() = %v, want empty string for missing key", got)
	}
}

func TestReadFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	content := []byte("test file content")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "existing file",
			path:    tmpFile.Name(),
			want:    content,
			wantErr: false,
		},
		{
			name:    "non-existing file",
			path:    "/nonexistent/path/file.txt",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
