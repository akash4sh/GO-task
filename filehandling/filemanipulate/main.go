package main

import (
	"fmt"
	"path"
	"path/filepath"
)

func ManipulatePath() {
	dirs := []string{"home", "mano", "go-projects", "..", "src"}
	path := path.Join(dirs...)
	fmt.Printf("Path after join: %s\n", path)

	fmt.Printf("Path after split: ")
	splitted := filepath.SplitList(path)
	for _, d := range splitted {
		fmt.Printf("%s%c", d, filepath.Separator)
	}
}

func main() {
	ManipulatePath()
}
