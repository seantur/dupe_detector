package main

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const errMsg = "Usage: dupe_detector <dir>"

func sha1Sum(path string) (string, error) {
	f, err := os.Open(path)

	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha1.New()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func walk(path string) error {
	size := make(map[int64][]string)
	dupeSize := make(map[int64]bool)

	links := make(map[string][]string)
	dupeLinks := make(map[string]bool)

	// get size duplicates
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		infoSize := info.Size()

		size[infoSize] = append(size[infoSize], path)

		if len(size[infoSize]) > 1 {
			dupeSize[infoSize] = true
		}

		return nil
	})

	if err != nil {
		return err
	}

	// for each size duplicate, check for hash duplicates
	for s := range dupeSize {
		for _, path := range size[s] {
			hash, err := sha1Sum(path)

			if err != nil {
				return err
			}

			links[hash] = append(links[hash], path)

			if len(links[hash]) > 1 {
				dupeLinks[hash] = true
			}
		}
	}

	if len(dupeLinks) != 0 {
		fmt.Println("duplicates:")

		for hash := range dupeLinks {
			fmt.Printf("%s: %v\n", hash, links[hash])
		}
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal(errMsg)
	}

	err := walk(os.Args[1])

	if err != nil {
		fmt.Println(err)
	}
}
