package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/dustin/go-humanize"
	sqlite "github.com/gwenn/gosqlite"
	"github.com/mitchellh/go-homedir"
)

var (
	// Our configuration info
	Conf TomlConfig

	// SQLite connection handle
	sdb *sqlite.Conn
)

func main() {
	// Read the disk cache directory name from the dbhub config file
	var err error
	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		userHome, err := homedir.Dir()
		if err != nil {
			log.Fatalf("User home directory couldn't be determined: %s", "\n")
		}
		configFile = filepath.Join(userHome, ".dbhub", "config.toml")
	}

	// Reads the server configuration from disk
	if _, err = toml.DecodeFile(configFile, &Conf); err != nil {
		log.Fatalf("Config file couldn't be parsed: %v\n", err)
	}

	// Create SQLite database in memory
	sdb, err = sqlite.Open(":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer sdb.Close()

	// Create SQLite database schema
	err = sdb.Exec("CREATE TABLE IF NOT EXISTS files (name TEXT, size INTEGER, lastmod TEXT, lastaccess TEXT); DELETE FROM files;")
	if err != nil {
		log.Fatal(err)
	}

	// Make a list of all the database files in the disk cache directory
	cacheDir := Conf.DiskCache.Directory
	entryList, err := ioutil.ReadDir(cacheDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, e := range entryList {
		if e.IsDir() {
			// We only want files in a sub directory off the main cache directory
			subDir := filepath.Join(cacheDir, e.Name())
			subDirList, err := ioutil.ReadDir(subDir)
			if err != nil {
				log.Fatal(err)
			}
			for _, z := range subDirList {
				// Add the file details to the SQLite database
				addFile(subDir, z)
			}
		}
	}

	// Determine which of the files hasn't been accessed/used in the greatest period of time (eg least needed)
	cleanTarget := 5 * 1024 * 1024 * 1024 // 5GB
	dbQuery := `
		SELECT name, size, lastaccess
		FROM "files"
		ORDER BY lastaccess ASC`
	var name string
	var size int
	var lastAccess time.Time
	var totalSize int
	err = sdb.Select(dbQuery, func(s *sqlite.Stmt) (err error) {
		if err = s.Scan(&name, &size, &lastAccess); err != nil {
			log.Fatal(err)
		}
		if totalSize <= cleanTarget {
			// Delete the least recently used files
			fmt.Printf("Deleting file: '%v', last accessed '%v', size: %v\n", name, lastAccess.Local(), humanize.Bytes(uint64(size)))
			err := os.Remove(name)
			if err != nil {
				log.Fatal(err)
			}
			totalSize += size
		}
		return
	})
	if err != nil {
		log.Fatal(err)
	}

	// Report on the deleted files
	fmt.Printf("\nTarget to clean up: %s\n", humanize.Bytes(uint64(cleanTarget)))
	fmt.Printf("Space recovered: %s\n", humanize.Bytes(uint64(totalSize)))
}

// Adds a given database file to the in memory SQLite database
func addFile(path string, z os.FileInfo) {
	// Get the last access file for the file
	aRaw := z.Sys().(*syscall.Stat_t).Atim
	aTime := time.Unix(aRaw.Sec, aRaw.Nsec)

	// Add file to the database
	_, err := sdb.ExecDml("INSERT INTO files (name, size, lastmod, lastaccess) VALUES (?, ?, ?, ?)",
		filepath.Join(path, z.Name()), z.Size(), z.ModTime(), aTime)
	if err != nil {
		log.Fatal(err)
	}
}
