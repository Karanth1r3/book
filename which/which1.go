package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Provide an argument")
		return
	}

	file := arguments[1]

	path := os.Getenv("PATH")
	pathSplit := filepath.SplitList(path)
	for _, directory := range pathSplit {
		fullPath := filepath.Join(directory, file)
		// Does path exist?
		fileInfo, err := os.Stat(fullPath)
		if err == nil {
			mode := fileInfo.Mode()
			// Is file of common mode?
			if mode.IsRegular() {
				// Is file executable?
				if mode&0111 != 0 { // How this works?
					fmt.Println(fullPath)
					return
				}
			}

		}
	}
}
