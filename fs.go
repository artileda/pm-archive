package main

import (
	"path/filepath"
	"os"
	"fmt"
)


func scanDir(path string) []string {
	files := []string{}
	e := filepath.Walk(path, func(path string, info os.FileInfo, e error) error {
		fmt.Println(path)
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if e!=nil{
	  fmt.Println(e)
	}
	return files
}

func main(){
	scanDir("/root/kartini/container")
}
