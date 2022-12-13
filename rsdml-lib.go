package rsdml

import "fmt"
import "os"
import "path/filepath"
import "strings"
import "sort"

func UpdateDirectory(dir string) error {
	// Get current mtime of directory
	dirInfo, err := os.Lstat(dir)
	if err != nil {
		return fmt.Errorf("Path cannot be accessed: %s: %s", dir, err)
	}
	dir_mtime := dirInfo.ModTime()

	// Find newest mtime in directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Path cannot be accessed: %s: %s", dir, err)
	}

	newest_mtime := dir_mtime
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		entry_mtime := info.ModTime()
		if entry_mtime.After(newest_mtime) {
			newest_mtime = entry_mtime
		}
	}

	// Update directory mtime with newest mtime found
	if newest_mtime != dir_mtime {
		err := os.Chtimes(dir, newest_mtime, newest_mtime)
		if err != nil {
			return fmt.Errorf("Path cannot be touched: %s: %s", dir, err)
		}
	}

	return nil
}

func RecurseDirectory(root_dir string) ([]string, error) {
	// Input Sanitisation
	root_dir = filepath.Clean(root_dir)
	fileInfo, err := os.Stat(root_dir)
	if err != nil {
		return nil, fmt.Errorf("Path is not a directory or cannot be accessed: %s: %s", root_dir, err)
	}

	if !fileInfo.IsDir() {
		return nil, fmt.Errorf("Path is not a directory: %s", root_dir)
	}

	// Collect list of directories in tree
	var dirs []string
	err = filepath.Walk(root_dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			dirs = append(dirs, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Cannot traverse path: %s: %s", root_dir, err)
	}

	// Sort dirs from bottom to top (deepest first)
	sort.Slice(dirs, func(i, j int) bool {
		i_parts := strings.Split(dirs[i], string(os.PathSeparator))
		j_parts := strings.Split(dirs[j], string(os.PathSeparator))
		return len(i_parts) > len(j_parts)
	})

	// Adjust mtime values for each dir (bottom to top)
	failures := 0
	directories_touched := []string{}
	for _, dir := range dirs {
		err = UpdateDirectory(dir)
		if err != nil {
			failures = 1
		} else {
			directories_touched = append(directories_touched, dir)
		}
	}
	if failures == 1 {
		return directories_touched, fmt.Errorf("Could not update entire tree of: %s", root_dir)
	}

	return directories_touched, nil
}
