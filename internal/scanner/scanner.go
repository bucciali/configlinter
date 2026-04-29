package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"configlinter/internal/domain"
	"configlinter/internal/engine"
	"configlinter/internal/parser"
)

var supportedExtensions = map[string]bool{
	".yaml": true,
	".yml":  true,
	".json": true,
	".toml": true,
}

type FileResult struct {
	Path     string
	Findings []domain.Finding
	Err      error
}

func ScanDir(dir string, reg *parser.Registry, eng *engine.Engine) ([]FileResult, error) {
	info, err := os.Stat(dir)
	if err != nil {
		return nil, fmt.Errorf("cannot access %s: %w", dir, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", dir)
	}

	var files []string
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if supportedExtensions[ext] {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk error: %w", err)
	}

	results := make([]FileResult, len(files))
	var wg sync.WaitGroup

	for i, path := range files {
		wg.Add(1)
		go func(i int, path string) {
			defer wg.Done()
			results[i] = lintFile(path, reg, eng)
		}(i, path)
	}

	wg.Wait()
	return results, nil
}

func lintFile(path string, reg *parser.Registry, eng *engine.Engine) FileResult {
	data, err := os.ReadFile(path)
	if err != nil {
		return FileResult{Path: path, Err: fmt.Errorf("read: %w", err)}
	}

	p, err := reg.GetByFilename(path)
	if err != nil {
		return FileResult{Path: path, Err: fmt.Errorf("parser: %w", err)}
	}

	parsed, err := p.Parse(data)
	if err != nil {
		return FileResult{Path: path, Err: fmt.Errorf("parse: %w", err)}
	}

	findings := eng.Analyze(parsed)
	return FileResult{Path: path, Findings: findings}
}
