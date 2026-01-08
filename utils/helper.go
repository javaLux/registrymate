package utils

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func WriteFile(path string, content []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("YAML-File [ %v ]- failed to close file: %v", path, err)
		}
	}()

	writer := bufio.NewWriter(file)

	_, err = writer.Write(content)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func EnsureYAMLExt(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	if ext == ".yaml" || ext == ".yml" {
		return path
	}

	return strings.TrimSuffix(path, ext) + ".yaml"
}
