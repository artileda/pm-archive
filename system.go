package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	//"github.com/src-d/go-git/utils/ioutil"
	"golang.org/x/crypto/sha3"
)

func hashThese(b []byte) string {
	hash := make([]byte, 64)
	sha3.ShakeSum256(hash, b)
	return string(hash)
}

func hashFile(p string) string {
	file, e := ioutil.ReadFile(p)
	if e != nil {
		return e.Error()
	}
	return hashThese(file)
}

func getPid() string {
	return string(os.Getpid())
}
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

func scanDir(path string) []string {
	files := []string{}
	_ = filepath.Walk(path, func(path string, info os.FileInfo, e error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files
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
	p, found := findPackage(name)
	if !found {
		errPackNotExist()
	}
	for _, item := range p.Sources {
		if !isExist(getCachePath() + "/sources/" + p.Name + "/" + lastStr(splitStr(item[0], "/"))) {
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

func makeFile(buf []byte, pathname string) error {
	out, e := os.Create(pathname)
	if e != nil {
		return e
		os.Exit(1)
	}
	defer out.Close()
	out.Write([]byte(buf))
	out.Sync()
	return nil
}

func makeTempCacheDir() []string {
	if getCachePath() == "" {
		fmt.Println("KARTINI_CACHE variable need be set up")
		return []string{}
	} else {
		var dir = []string{
			getCachePath() + "/bin-" + string(getPid()),
			getCachePath() + "/source-" + string(getPid()),
		}
		for _, d := range dir {
			os.Mkdir(d, 0755)
		}
		return dir
	}
}

func removeTempCacheDir() {

}

func cleanStep() {

}

func errPackNotExist() {
	fmt.Println("[!] Package not exist on local repos")
	os.Exit(1)
}

// func makeManifest(tomlpath string) bool {

// }
