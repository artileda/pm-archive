package main

import (
	"fmt"
	"os"
)

func getManifestPath() string {
	return os.Getenv("KARTINI_ROOT") + "/var/db/kartini/installed"
}

func getLocalRepo() []string {
	path := os.Getenv("KARTINI_PATH")
	return splitStr(path, ":")
}

func getCachePath() string {
	return os.Getenv("KARTINI_CACHE")
}

func isExist(path string) bool {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return true
	} else {
		return false
	}
}

func isInstalled(name string) bool {
	os.Chdir(getManifestPath())
	if isExist(getManifestPath() + "/" + name) {
		return true
	} else {
		return false
	}
}

func isCached(name string) bool {
	path, found := findPackage(name)
	if !found {
		errPackNotExist()
	}
	p := tomlToPackage(path + "/dist.toml")
	path = getCachePath() + "/" + p.Name
	for _, item := range p.Sources {
		if !isExist(path + "/" + lastStr(splitStr(item[0], "/"))) {
			return false
		}
	}
	return true
}
func makePersistCacheDir() {
	if getCachePath() == "" {
		fmt.Println("KARTINI_CACHE variable need be set up")
		os.Exit(1)
		return
	} else if !(isExist(getCachePath()+"/bin") &&
		isExist(getCachePath()+"/source")) {
		fmt.Println("[!] Make Cache ")
		os.Mkdir(getCachePath()+"/bin", 0755)
		os.Mkdir(getCachePath()+"/source", 0755)
	}
}

func isSrcCached(name string) bool {
	makePersistCacheDir()
	return isExist(getCachePath() + "/source/" + name)
}

func isBinCached(name string) bool {
	makePersistCacheDir()
	return isExist(getCachePath() + "/bin/" + name)
}

func makeTempCacheDir() []string {
	if getCachePath() == "" {
		fmt.Println("KARTINI_CACHE variable need be set up")
		return []string{}
	} else {
		var dir = []string{
			getCachePath() + "/bin-" + string(os.Getpid()),
			getCachePath() + "/source-" + string(os.Getpid()),
		}
		for _, d := range dir {
			os.Create(d)
		}
		return dir
	}
}

func errPackNotExist() {
	fmt.Println("[!] Package not exist on local repos")
	os.Exit(1)
}

// func makeManifest(tomlpath string) bool {

// }
