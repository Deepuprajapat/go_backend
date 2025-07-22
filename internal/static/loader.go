package static

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/VI-IM/im_backend_go/shared/logger"
)

// StaticLoader handles downloading and loading static assets into memory
type StaticLoader struct {
	memfs *MemoryFileSystem
	url   string
}

// NewStaticLoader creates a new static assets loader
func NewStaticLoader(url string) *StaticLoader {
	return &StaticLoader{
		memfs: NewMemoryFileSystem(),
		url:   url,
	}
}

// LoadFromURL downloads a zip file from URL and extracts it to memory
func (sl *StaticLoader) LoadFromURL() error {
	if sl.url == "" {
		logger.Get().Info().Msg("No static assets URL configured, skipping static assets loading")
		return nil
	}

	logger.Get().Info().Msgf("Downloading static assets from: %s", sl.url)
	
	// Download zip file
	resp, err := http.Get(sl.url)
	if err != nil {
		return fmt.Errorf("failed to download zip file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download zip file: HTTP %d", resp.StatusCode)
	}

	// Read the entire response body into memory
	zipData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read zip data: %w", err)
	}

	logger.Get().Info().Msgf("Downloaded %d bytes, extracting to memory...", len(zipData))

	// Extract zip directly from memory
	if err := sl.extractZipFromMemory(zipData); err != nil {
		return fmt.Errorf("failed to extract zip: %w", err)
	}

	files := sl.memfs.ListFiles()
	logger.Get().Info().Msgf("Successfully loaded %d files into memory", len(files))
	
	// Log some example files for debugging
	for i, file := range files {
		if i < 5 {
			logger.Get().Debug().Msgf("  - %s", file)
		}
	}
	if len(files) > 5 {
		logger.Get().Debug().Msgf("  ... and %d more files", len(files)-5)
	}

	return nil
}

// extractZipFromMemory extracts zip data from memory to memory filesystem
func (sl *StaticLoader) extractZipFromMemory(zipData []byte) error {
	// Create a reader from the zip data
	zipReader := bytes.NewReader(zipData)
	
	// Open zip reader
	zipFile, err := zip.NewReader(zipReader, int64(len(zipData)))
	if err != nil {
		return fmt.Errorf("failed to create zip reader: %w", err)
	}

	// Clear any existing files
	sl.memfs.Clear()

	// Extract each file
	for _, file := range zipFile.File {
		if err := sl.extractFile(file); err != nil {
			logger.Get().Warn().Err(err).Msgf("Failed to extract file: %s", file.Name)
			continue
		}
	}

	return nil
}

// extractFile extracts a single file from the zip archive to memory
func (sl *StaticLoader) extractFile(file *zip.File) error {
	// Skip directories
	if file.FileInfo().IsDir() {
		return nil
	}

	// Open the file in the zip
	reader, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", file.Name, err)
	}
	defer reader.Close()

	// Read file content
	content, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", file.Name, err)
	}

	// Clean up the path - remove any leading directories if needed
	path := strings.TrimPrefix(file.Name, "./")
	path = strings.TrimPrefix(path, "build/")
	
	// Ensure path starts with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Store in memory filesystem
	sl.memfs.StoreFile(path, content)
	
	return nil
}

// GetMemoryFileSystem returns the memory filesystem instance
func (sl *StaticLoader) GetMemoryFileSystem() *MemoryFileSystem {
	return sl.memfs
}

// LoadFromDirectory loads files from local directory (for development)
func (sl *StaticLoader) LoadFromDirectory(dir string) error {
	logger.Get().Info().Msgf("Loading static assets from directory: %s", dir)
	
	// Clear existing files
	sl.memfs.Clear()
	
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		// Read file content
		content, err := os.ReadFile(path)
		if err != nil {
			logger.Get().Warn().Err(err).Msgf("Failed to read file: %s", path)
			return nil // Continue with other files
		}
		
		// Convert to web path
		relativePath := strings.TrimPrefix(path, dir)
		relativePath = strings.TrimPrefix(relativePath, string(filepath.Separator))
		webPath := "/" + strings.ReplaceAll(relativePath, string(filepath.Separator), "/")
		
		// Store in memory
		sl.memfs.StoreFile(webPath, content)
		
		return nil
	})
}

// Reload reloads the static assets from the configured URL
func (sl *StaticLoader) Reload() error {
	return sl.LoadFromURL()
}

// Global static loader instance
var globalLoader *StaticLoader

// InitializeStaticLoader initializes the global static loader
func InitializeStaticLoader(url string) error {
	globalLoader = NewStaticLoader(url)
	return globalLoader.LoadFromURL()
}

// GetGlobalMemoryFileSystem returns the global memory filesystem
func GetGlobalMemoryFileSystem() *MemoryFileSystem {
	if globalLoader == nil {
		// Return an empty filesystem if not initialized
		return NewMemoryFileSystem()
	}
	return globalLoader.GetMemoryFileSystem()
}