package main

import (
	"fmt"
	"os"
)

// Refactor : Return Package Object instead string and bool
func findPackage(name string) (string, bool) {
	repo := getLocalRepo()
	found := false
	fpath := ""
	for _, path := range repo {
		found = isExist(path+"/"+name) && isExist(path+"/"+name+"/dist.toml")
		if found {
			fpath = (path + "/" + name)
			return fpath, found
		}
	}
	fmt.Println("Package not found: ", name)
	os.Exit(1)
	return fpath, found
}

func getPackage(name string) {
	path, _ := findPackage(name)
	p := tomlToPackage(path + "/dist.toml")
	fmt.Println("[*] Fetching resource ...")
	p.Download()
}

func installPackage(name string) {
	path, found := findPackage(name)
	if !found {
		os.Exit(1)
	}
	p := tomlToPackage(path)
	p.Download()
}

func extractPackage(name string) {
	path, found := findPackage(name)
	if !found {
		os.Exit(1)
	}
	p := tomlToPackage(path + "/dist.toml")
	fmt.Println("[", name, "] Extarcting resources...")
	p.extract("")
}
func buildPackage(name string) {

}
func removePackage(name string) {
}
