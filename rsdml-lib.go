package rsdml

import "fmt"
import "log"
import "os"
import "path/filepath"
import "strings"
import "sort"

func UpdateDirectory(dir string) error {
	log.Printf("UpdateDirectory(%s)", dir)

	// Get current mtime of directory
	dirInfo, err := os.Lstat(dir)
	if err != nil {
		log.Print(err)
		return fmt.Errorf("Path cannot be accessed: %s", dir)
	}
	dir_mtime := dirInfo.ModTime()

	// Find newest mtime in directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Print(err)
		return fmt.Errorf("Path cannot be accessed: %s", dir)
	}

	newest_mtime := dir_mtime
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			log.Print(err)
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
			log.Print(err)
			return fmt.Errorf("Path cannot be touched: %s", dir)
		}
		log.Printf("Touching %s with date: %s", dir, newest_mtime)
	}

	return nil
}

func RecurseDirectory(root_dir string) error {
	// Input Sanitisation
	root_dir = filepath.Clean(root_dir)
	log.Printf("RecurseDirectory(%s)\n", root_dir)
	fileInfo, err := os.Stat(root_dir)
	if err != nil {
		log.Print(err)
		return fmt.Errorf("Path is not a directory or cannot be accessed: %s", root_dir)
	}

	if ! fileInfo.IsDir() {
		return fmt.Errorf("Path is not a directory: %s", root_dir)
	}

	// Collect list of directories in tree
	var dirs []string
	err = filepath.Walk(root_dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		if info.IsDir() {
			dirs = append(dirs, path)
		}

		return nil
	})

	if err != nil {
		log.Print(err)
		return fmt.Errorf("Cannot traverse path: %s", root_dir)
	}

	// Sort dirs from bottom to top (deepest first)
	sort.Slice(dirs, func(i, j int) bool {
		i_parts := strings.Split(dirs[i], string(os.PathSeparator))
		j_parts := strings.Split(dirs[j], string(os.PathSeparator))
		return len(i_parts) > len(j_parts)
	})

	// Adjust mtime values for each dir (bottom to top)
	failures := 0
	for _, subdir := range dirs {
		err = UpdateDirectory(subdir)
		if err != nil {
			log.Print(err)
			failures = 1
		}
	}
	if failures == 1 {
		return fmt.Errorf("Could not update entire tree of: %s", root_dir)
	}

	return nil
}
