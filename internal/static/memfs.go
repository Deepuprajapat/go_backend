package static

import (
	"mime"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// FileInfo holds metadata about a file stored in memory
type FileInfo struct {
	Content     []byte
	ContentType string
	ModTime     time.Time
	Size        int64
}

// MemoryFileSystem stores files in memory for fast serving
type MemoryFileSystem struct {
	files map[string]*FileInfo
	mutex sync.RWMutex
}

// NewMemoryFileSystem creates a new in-memory filesystem
func NewMemoryFileSystem() *MemoryFileSystem {
	return &MemoryFileSystem{
		files: make(map[string]*FileInfo),
	}
}

// StoreFile stores a file in memory with automatic MIME type detection
func (mfs *MemoryFileSystem) StoreFile(path string, content []byte) {
	mfs.mutex.Lock()
	defer mfs.mutex.Unlock()

	// Normalize path to always start with /
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	// Detect MIME type from file extension
	contentType := mime.TypeByExtension(filepath.Ext(path))
	if contentType == "" {
		// Default fallbacks for common web files
		switch strings.ToLower(filepath.Ext(path)) {
		case ".js":
			contentType = "application/javascript"
		case ".css":
			contentType = "text/css"
		case ".html", ".htm":
			contentType = "text/html; charset=utf-8"
		case ".json":
			contentType = "application/json"
		case ".svg":
			contentType = "image/svg+xml"
		case ".woff", ".woff2":
			contentType = "font/woff"
		case ".ttf":
			contentType = "font/ttf"
		case ".eot":
			contentType = "application/vnd.ms-fontobject"
		default:
			contentType = "application/octet-stream"
		}
	}

	mfs.files[path] = &FileInfo{
		Content:     content,
		ContentType: contentType,
		ModTime:     time.Now(),
		Size:        int64(len(content)),
	}
}

// GetFile retrieves a file from memory
func (mfs *MemoryFileSystem) GetFile(path string) (*FileInfo, bool) {
	mfs.mutex.RLock()
	defer mfs.mutex.RUnlock()

	// Normalize path
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	file, exists := mfs.files[path]
	return file, exists
}

// HasFile checks if a file exists in memory
func (mfs *MemoryFileSystem) HasFile(path string) bool {
	mfs.mutex.RLock()
	defer mfs.mutex.RUnlock()

	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	_, exists := mfs.files[path]
	return exists
}

// ListFiles returns all file paths stored in memory
func (mfs *MemoryFileSystem) ListFiles() []string {
	mfs.mutex.RLock()
	defer mfs.mutex.RUnlock()

	paths := make([]string, 0, len(mfs.files))
	for path := range mfs.files {
		paths = append(paths, path)
	}
	return paths
}

// Clear removes all files from memory
func (mfs *MemoryFileSystem) Clear() {
	mfs.mutex.Lock()
	defer mfs.mutex.Unlock()

	mfs.files = make(map[string]*FileInfo)
}

// GetCacheControl returns appropriate cache-control header based on file type
func GetCacheControl(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	filename := strings.ToLower(filepath.Base(path))

	// Never cache index.html and manifest files
	if filename == "index.html" || filename == "manifest.json" || filename == "robots.txt" || filename == "sitemap.xml" {
		return "no-cache, no-store, must-revalidate"
	}

	// Long cache for versioned assets (contain hash in filename)
	if strings.Contains(filename, ".") && (ext == ".js" || ext == ".css") {
		parts := strings.Split(filename, ".")
		if len(parts) > 2 {
			// Likely has version hash like main.abc123.js
			return "public, max-age=31536000, immutable" // 1 year
		}
	}

	// Medium cache for other static assets
	switch ext {
	case ".js", ".css":
		return "public, max-age=86400" // 1 day
	case ".png", ".jpg", ".jpeg", ".gif", ".webp", ".svg", ".ico":
		return "public, max-age=604800" // 1 week
	case ".woff", ".woff2", ".ttf", ".eot":
		return "public, max-age=2592000" // 30 days
	case ".pdf":
		return "public, max-age=86400" // 1 day
	default:
		return "public, max-age=3600" // 1 hour
	}
}