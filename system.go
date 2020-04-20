package main

import (
	"fmt"
//	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
//	"strings"
	//"time"
	"os/exec"
//	"bufio"


	//"github.com/src-d/go-git/utils/ioutil"
	"golang.org/x/crypto/sha3"
//	"github.com/go-cmd/cmd"
)

func runCmd(sh string,args ...string){
	fmt.Println(args)
	cmd,e := exec.Command(sh,args...).Output()
	if e != nil{
	  fmt.Println(e)
	}
	fmt.Println(string(cmd))

/*	buf := bufio.NewReader(stdout)
	num := 1
	for{
	   line,_,e := buf.ReadLine()
	   if e == io.EOF {
	      break
	   }
	   if e != nil{
		fmt.Println(e)
		break
	   }
	   num += 1
	   fmt.Println(string(line))
	}
*/
}

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
	fmt.Println(strconv.Itoa(os.Getpid()))
	return strconv.Itoa(os.Getpid())
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
	fmt.Println(path)
	files := []string{}
	e := filepath.Walk(path, func(path string, info os.FileInfo, e error) error {
		fmt.Println(path)
		if !info.IsDir() {
			files = append(files, path)
		}
		if e!= nil{
		  fmt.Println(e)
		}
		return nil
	})
	if e!=nil{
	  fmt.Println(e)
	}
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
			getCachePath() + "/bin-" + getPid(),
			getCachePath() + "/source-" + getPid(),
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
