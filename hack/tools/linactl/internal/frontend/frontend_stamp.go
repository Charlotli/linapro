// This file tracks the official-plugin frontend source stamp so plugin page
// changes can force a turbo rebuild. The Vite plugin-pages virtual module
// globs apps/lina-plugins outside the frontend workspace, so its inputs are
// invisible to turbo's cache key; linactl persists a content hash and lets the
// build translate changes into a turbo --force pass.

package frontend

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// stampFilePath stores the last built plugin frontend content hash under the
// repository temp directory owned by linactl.
func stampFilePath(root string) string {
	return filepath.Join(root, "temp", "frontend-plugin-pages.stamp")
}

// pluginFrontendStampChanged reports whether the current plugin frontend and
// manifest content differs from the last recorded build stamp. A missing stamp
// forces one full rebuild so the baseline is always trustworthy.
func PluginFrontendStampChanged(root string) (bool, error) {
	current, err := computePluginFrontendStamp(root)
	if err != nil {
		return false, err
	}
	previous, err := os.ReadFile(stampFilePath(root))
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	return strings.TrimSpace(string(previous)) != current, nil
}

// SavePluginFrontendStamp persists the current plugin frontend content hash
// after one successful frontend build.
func SavePluginFrontendStamp(root string) error {
	current, err := computePluginFrontendStamp(root)
	if err != nil {
		return err
	}
	stampPath := stampFilePath(root)
	if err = os.MkdirAll(filepath.Dir(stampPath), 0o755); err != nil {
		return fmt.Errorf("create frontend stamp directory: %w", err)
	}
	if err = os.WriteFile(stampPath, []byte(current), 0o644); err != nil {
		return fmt.Errorf("write frontend stamp: %w", err)
	}
	return nil
}

// computePluginFrontendStamp hashes every official plugin's frontend sources
// and manifest files with stable path-qualified ordering so any page, asset,
// or governance declaration change invalidates the stamp.
func computePluginFrontendStamp(root string) (string, error) {
	pluginsRoot := filepath.Join(root, "apps", "lina-plugins")
	if _, err := os.Stat(pluginsRoot); err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", err
	}

	hasher := sha256.New()
	walkErr := filepath.WalkDir(pluginsRoot, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			switch entry.Name() {
			case "node_modules", "dist", ".git":
				return filepath.SkipDir
			}
			return nil
		}
		relative, relErr := filepath.Rel(pluginsRoot, path)
		if relErr != nil {
			return relErr
		}
		relativeSlash := filepath.ToSlash(relative)
		if !strings.HasPrefix(relativeSlash, "frontend/") &&
			filepath.Base(relativeSlash) != "plugin.yaml" {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		fmt.Fprintf(hasher, "%s\x00%x\x00", relativeSlash, sha256.Sum256(content))
		return nil
	})
	if walkErr != nil {
		return "", fmt.Errorf("hash official plugin frontend sources: %w", walkErr)
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
