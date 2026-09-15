package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Remover handles file deletion and marker-based patching.
type Remover struct {
	root string
}

// NewRemover creates a Remover rooted at the given project directory.
func NewRemover(root string) *Remover {
	return &Remover{root: root}
}

// DeleteFiles removes standalone files listed in the module manifest.
// Missing files are silently skipped.
func (r *Remover) DeleteFiles(files []string) (int, error) {
	count := 0
	for _, f := range files {
		p := filepath.Join(r.root, f)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			fmt.Printf("  skip (not found): %s\n", f)
			continue
		}
		if err := os.Remove(p); err != nil {
			return count, fmt.Errorf("failed to delete %s: %w", f, err)
		}
		count++
		fmt.Printf("  deleted: %s\n", f)
	}
	return count, nil
}

// DeleteDirs removes directories recursively.
// Missing directories are silently skipped.
func (r *Remover) DeleteDirs(dirs []string) (int, error) {
	count := 0
	for _, d := range dirs {
		p := filepath.Join(r.root, d)
		if _, err := os.Stat(p); os.IsNotExist(err) {
			fmt.Printf("  skip dir (not found): %s\n", d)
			continue
		}
		if err := os.RemoveAll(p); err != nil {
			return count, fmt.Errorf("failed to delete dir %s: %w", d, err)
		}
		count++
		fmt.Printf("  deleted dir: %s\n", d)
	}
	return count, nil
}

// markerStyle pairs start and end regex patterns for a single comment syntax.
type markerStyle struct {
	start string
	end   string
}

// markerStyles returns paired start/end regex patterns for all supported
// comment syntaxes for the given feature name.
func markerStyles(featureName string) []markerStyle {
	q := regexp.QuoteMeta(featureName)
	return []markerStyle{
		{ // Go / JS / TS single-line comments
			start: fmt.Sprintf(`(?m)^[^\S\n]*//\s*\[feature:%s\]\s*start\s*\n`, q),
			end:   fmt.Sprintf(`(?m)^[^\S\n]*//\s*\[feature:%s\]\s*end\s*\n?`, q),
		},
		{ // YAML / Makefile hash comments
			start: fmt.Sprintf(`(?m)^[^\S\n]*#\s*\[feature:%s\]\s*start\s*\n`, q),
			end:   fmt.Sprintf(`(?m)^[^\S\n]*#\s*\[feature:%s\]\s*end\s*\n?`, q),
		},
		{ // JSX / TSX block comments
			start: fmt.Sprintf(`(?m)^[^\S\n]*\{/\*\s*\[feature:%s\]\s*start\s*\*/\}\s*\n`, q),
			end:   fmt.Sprintf(`(?m)^[^\S\n]*\{/\*\s*\[feature:%s\]\s*end\s*\*/\}\s*\n?`, q),
		},
	}
}

// PatchFeatureMarkers removes all [feature:name] start/end blocks from all
// files in the project that contain them. Uses a single directory walk
// and a single read/write per file regardless of how many features are removed.
func (r *Remover) PatchFeatureMarkers(featureNames []string) (int, error) {
	fileFeatures, err := r.findAllFilesWithMarkers(featureNames)
	if err != nil {
		return 0, err
	}

	filesPatched := 0
	for _, ff := range fileFeatures {
		patched, err := r.patchFile(ff.path, ff.features)
		if err != nil {
			return filesPatched, err
		}
		if patched {
			filesPatched++
		}
	}

	return filesPatched, nil
}

// fileWithFeatures associates a file path with the feature markers found in it.
type fileWithFeatures struct {
	path     string
	features []string
}

// findAllFilesWithMarkers performs a single directory walk and returns files
// that contain any of the given feature markers, along with which markers
// each file contains. Results are in walk order.
func (r *Remover) findAllFilesWithMarkers(featureNames []string) ([]fileWithFeatures, error) {
	markers := make(map[string]string, len(featureNames))
	for _, name := range featureNames {
		markers[name] = fmt.Sprintf("[feature:%s]", name)
	}

	var result []fileWithFeatures

	err := filepath.Walk(r.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			base := info.Name()
			if base == "node_modules" || base == ".git" || base == "vendor" {
				return filepath.SkipDir
			}
			if base == "setup" && filepath.Dir(path) == filepath.Join(r.root, "cmd") {
				return filepath.SkipDir
			}
			return nil
		}

		if !isTextFile(info.Name()) {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		content := string(data)

		var found []string
		for _, name := range featureNames {
			if strings.Contains(content, markers[name]) {
				found = append(found, name)
			}
		}
		if len(found) > 0 {
			result = append(result, fileWithFeatures{path: path, features: found})
		}
		return nil
	})

	return result, err
}

// isTextFile returns true if the file should be scanned for feature markers.
func isTextFile(name string) bool {
	switch filepath.Ext(name) {
	case ".go", ".tsx", ".ts", ".yaml", ".yml", ".jsx", ".js":
		return true
	}
	return name == "Makefile"
}

// patchFile reads the file once, strips marker blocks for all given features,
// runs artifact cleanup, and writes the result once.
func (r *Remover) patchFile(filePath string, featureNames []string) (bool, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}

	original := string(data)
	result := original

	for _, name := range featureNames {
		result = stripMarkerBlocks(result, relPath(r.root, filePath), name)
	}

	if result == original {
		return false, nil
	}

	result = cleanupArtifacts(result, filePath)

	info, _ := os.Stat(filePath)
	if err := os.WriteFile(filePath, []byte(result), info.Mode()); err != nil {
		return false, err
	}

	fmt.Printf("  patched: %s (removed [feature:%s] blocks)\n",
		relPath(r.root, filePath), strings.Join(featureNames, ", "))

	return true, nil
}

// stripMarkerBlocks removes all start/end delimited blocks for the given
// feature from content, returning the modified string.
func stripMarkerBlocks(content, displayPath, featureName string) string {
	for _, style := range markerStyles(featureName) {
		startRe := regexp.MustCompile(style.start)
		endRe := regexp.MustCompile(style.end)

		for {
			startLoc := startRe.FindStringIndex(content)
			if startLoc == nil {
				break
			}

			remainder := content[startLoc[1]:]
			endLoc := endRe.FindStringIndex(remainder)
			if endLoc == nil {
				fmt.Printf("  warning: no matching end marker for [feature:%s] in %s\n", featureName, displayPath)
				break
			}

			content = content[:startLoc[0]] + content[startLoc[1]+endLoc[1]:]
		}
	}
	return content
}

// cleanupArtifacts handles post-removal cleanup:
// - Collapse consecutive blank lines to one
// - Remove empty import groups: import ()
// - Remove dangling commas before closing braces in imports
func cleanupArtifacts(content, filePath string) string {
	multiBlank := regexp.MustCompile(`\n{3,}`)
	content = multiBlank.ReplaceAllString(content, "\n\n")

	ext := filepath.Ext(filePath)

	if ext == ".go" {
		emptyImport := regexp.MustCompile(`(?m)^import\s*\(\s*\)\s*\n?`)
		content = emptyImport.ReplaceAllString(content, "")
	}

	if ext == ".tsx" || ext == ".ts" || ext == ".jsx" || ext == ".js" {
		danglingComma := regexp.MustCompile(`,\s*\n(\s*}\s*from\s)`)
		content = danglingComma.ReplaceAllString(content, "\n$1")
	}

	return content
}

func relPath(root, abs string) string {
	rel, err := filepath.Rel(root, abs)
	if err != nil {
		return abs
	}
	return rel
}
