package infra

import (
	"log"
	"os"
)

func listFolders(path string) []string {
	folders := []string{}

	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			folders = append(folders, entry.Name())
		}
	}

	return folders
}
