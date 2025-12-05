package fileconfig

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"zatrano/configs/envconfig"
	"zatrano/configs/logconfig"
)

type FileConfig struct {
	BasePath      string
	AllowedExtMap map[string][]string
	mu            sync.Mutex
}

var Config *FileConfig

func InitFileConfig() {
	// .env'den oku, yoksa default belirle
	basePath := envconfig.String("FILE_BASE_PATH", "")
	if basePath == "" {
		if envconfig.IsProd() {
			basePath = "./uploads"
		} else {
			basePath = "./uploads"
		}
	}

	// dizin oluşturulmamışsa oluştur
	if err := os.MkdirAll(basePath, 0755); err != nil {
		// 🔧 Fix: Sugared logger kullan
		logconfig.SLog.Fatalw("Upload klasörü oluşturulamadı",
			"path", basePath, "error", err)
	}

	Config = &FileConfig{
		BasePath:      basePath,
		AllowedExtMap: make(map[string][]string),
	}

	logconfig.SLog.Infow("FileConfig initialized", "base_path", basePath)
}

func (fc *FileConfig) GetPath(contentType string) string {
	contentType = sanitize(contentType)
	return filepath.Join(fc.BasePath, contentType)
}

func (fc *FileConfig) GetAllowedExtensions(contentType string) []string {
	contentType = sanitize(contentType)
	fc.mu.Lock()
	defer fc.mu.Unlock()
	return fc.AllowedExtMap[contentType]
}

func (fc *FileConfig) SetAllowedExtensions(contentType string, extensions []string) {
	contentType = sanitize(contentType)
	fc.mu.Lock()
	defer fc.mu.Unlock()
	fc.AllowedExtMap[contentType] = extensions

	dir := fc.GetPath(contentType)
	if err := os.MkdirAll(dir, 0755); err != nil {
		// 🔧 Fix: Sugared logger kullan
		logconfig.SLog.Fatalw("Klasör oluşturulamadı",
			"dir", dir, "error", err)
	}
}

func (fc *FileConfig) IsExtensionAllowed(contentType, ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	for _, allowed := range fc.GetAllowedExtensions(contentType) {
		if allowed == ext {
			return true
		}
	}
	return false
}

func sanitize(str string) string {
	str = strings.ToLower(str)
	str = strings.TrimSpace(str)
	str = strings.ReplaceAll(str, " ", "_")
	return str
}
