package img

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"

	"github.com/Anacardo89/lenic/config"
)

type ImageManager struct {
	homeDir string
	cfg     *config.Img
}

func NewImgManager(cfg *config.Img, homeDir string) (*ImageManager, error) {
	if err := os.MkdirAll(cfg.Path, 0755); err != nil {
		return nil, err
	}
	return &ImageManager{
		homeDir: homeDir,
		cfg:     cfg,
	}, nil
}

func (m *ImageManager) SaveImg(file io.Reader, filename string) error {
	// Check if supported type
	ext := strings.ToLower(filepath.Ext(filename))
	if ext != ".jpg" && ext != ".jpeg" &&
		ext != ".png" &&
		ext != ".gif" {
		return fmt.Errorf("unsupported file type: %s", ext)
	}
	// Make path
	dir := filepath.Join(m.homeDir, m.cfg.Path, m.cfg.ImgDirs["originals"])
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	imgPath := filepath.Join(dir, filename)
	var (
		err     error
		imgData image.Image
		gifData *gif.GIF
	)
	// Decode
	if ext == ".gif" {
		data, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		gifData, err = gif.DecodeAll(bytes.NewReader(data))
		if err != nil {
			return err
		}
	} else {
		imgData, _, err = image.Decode(file)
		if err != nil {
			return err
		}
	}
	// Create output file
	outFile, err := os.Create(imgPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	// Encode
	if ext == ".gif" {
		if err := gif.EncodeAll(outFile, gifData); err != nil {
			return err
		}
	} else {
		switch ext {
		case ".jpg", ".jpeg":
			if err := jpeg.Encode(outFile, imgData, &jpeg.Options{Quality: m.cfg.JPEGQuality}); err != nil {
				return err
			}
		case ".png":
			if err := png.Encode(outFile, imgData); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *ImageManager) CreatePreview(filename string) error {
	// Get original
	originalDir := filepath.Join(m.homeDir, m.cfg.Path, m.cfg.ImgDirs["originals"])
	originalImgPath := filepath.Join(originalDir, filename)
	imgData, err := imaging.Open(originalImgPath)
	if err != nil {
		return err
	}
	// Make preview
	preview := imaging.Thumbnail(imgData, m.cfg.PreviewDims["width"], m.cfg.PreviewDims["height"], imaging.Lanczos)
	// Store preview
	previewDir := filepath.Join(m.homeDir, m.cfg.Path, m.cfg.ImgDirs["previews"])
	if err := os.MkdirAll(previewDir, 0755); err != nil {
		return err
	}
	previewImgPath := filepath.Join(previewDir, filename)
	if err := imaging.Save(preview, previewImgPath); err != nil {
		return err
	}

	return nil
}

func (m *ImageManager) GetImg(original bool, filename string) (*os.File, error) {
	var dir string
	if original {
		dir = filepath.Join(m.homeDir, m.cfg.Path, m.cfg.ImgDirs["originals"])
	} else {
		dir = filepath.Join(m.homeDir, m.cfg.Path, m.cfg.ImgDirs["previews"])
	}
	path := filepath.Join(dir, filename)
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %s", path)
		}
		return nil, err
	}
	return f, nil
}
