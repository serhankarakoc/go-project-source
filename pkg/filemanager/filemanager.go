package filemanager

import (
	"crypto/rand"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"zatrano/configs/fileconfig"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrFileNotProvided        = errors.New("dosya sağlanmadı")
	ErrInvalidFileType        = errors.New("geçersiz dosya türü veya uzantısı")
	ErrFileTooLarge           = errors.New("dosya boyutu çok büyük")
	ErrImageCouldNotBeDecoded = errors.New("resim dosyası çözümlenemedi, format desteklenmiyor olabilir")
)

const (
	DefaultMaxFileSize     = 2 * 1024 * 1024
	JpegProcessingQuality  = 75
	MimeSniffingBufferSize = 512
)

func UploadFile(c *fiber.Ctx, formFieldName, contentType string) (string, error) {
	fileHeader, err := c.FormFile(formFieldName)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return "", ErrFileNotProvided
		}
		return "", err
	}

	if err := validateFile(fileHeader, contentType); err != nil {
		return "", err
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("yüklenen dosya açılamadı: %w", err)
	}
	defer file.Close()

	buffer := make([]byte, MimeSniffingBufferSize)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("dosya türü okunurken hata: %w", err)
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("dosya okuma imleci sıfırlanamadı: %w", err)
	}

	detectedContentType := http.DetectContentType(buffer[:n])

	if detectedContentType == "image/jpeg" || detectedContentType == "image/png" {
		return processAndSaveImage(file, fileHeader.Filename, contentType)
	}

	return saveOriginalFile(c, fileHeader, contentType)
}

func processAndSaveImage(file multipart.File, originalFilename, contentType string) (string, error) {
	img, format, err := image.Decode(file)
	if err != nil {
		return "", ErrImageCouldNotBeDecoded
	}

	newFileName, err := generateUniqueFileName(originalFilename)
	if err != nil {
		return "", fmt.Errorf("benzersiz dosya adı oluşturulamadı: %w", err)
	}

	destination := filepath.Join(fileconfig.Config.GetPath(contentType), newFileName)
	destFile, err := os.Create(destination)
	if err != nil {
		return "", fmt.Errorf("hedef dosya oluşturulamadı: %w", err)
	}
	defer destFile.Close()

	switch format {
	case "jpeg":
		options := &jpeg.Options{Quality: JpegProcessingQuality}
		if err := jpeg.Encode(destFile, img, options); err != nil {
			os.Remove(destination)
			return "", fmt.Errorf("jpeg dosyası kodlanamadı: %w", err)
		}
	case "png":
		if err := png.Encode(destFile, img); err != nil {
			os.Remove(destination)
			return "", fmt.Errorf("png dosyası kodlanamadı: %w", err)
		}
	}

	return newFileName, nil
}

func saveOriginalFile(c *fiber.Ctx, fileHeader *multipart.FileHeader, contentType string) (string, error) {
	newFileName, err := generateUniqueFileName(fileHeader.Filename)
	if err != nil {
		return "", fmt.Errorf("benzersiz dosya adı oluşturulamadı: %w", err)
	}
	destination := filepath.Join(fileconfig.Config.GetPath(contentType), newFileName)
	if err := c.SaveFile(fileHeader, destination); err != nil {
		return "", err
	}
	return newFileName, nil
}

func validateFile(file *multipart.FileHeader, contentType string) error {
	if file.Size > DefaultMaxFileSize {
		return ErrFileTooLarge
	}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !fileconfig.Config.IsExtensionAllowed(contentType, ext) {
		return ErrInvalidFileType
	}
	return nil
}

func generateUniqueFileName(originalName string) (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	randomStr := fmt.Sprintf("%x", b)
	ext := filepath.Ext(originalName)
	safeBaseName := regexp.MustCompile(`[^a-zA-Z0-9_-]+`).ReplaceAllString(strings.TrimSuffix(originalName, ext), "")
	if safeBaseName == "" {
		safeBaseName = "file"
	}
	return fmt.Sprintf("%s-%s%s", randomStr, safeBaseName, ext), nil
}

func DeleteFile(contentType, fileName string) {
	if fileName == "" || contentType == "" {
		return
	}
	go func() {
		const maxRetries = 5
		const retryDelay = 1 * time.Second
		absolutePath, err := filepath.Abs(filepath.Join(fileconfig.Config.GetPath(contentType), fileName))
		if err != nil {
			return
		}
		for i := 0; i < maxRetries; i++ {
			err = os.Remove(absolutePath)
			if err == nil || os.IsNotExist(err) {
				return
			}
			time.Sleep(retryDelay)
		}
	}()
}
