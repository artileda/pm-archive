package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"bufio"
	"os/exec"
	"strings"

	"github.com/BurntSushi/toml"
	"time"
)


// this package is reflection
// from dist.toml on each package
type Package struct {
	Name        string
	Version     string
	Depends     []string
	Sources     [][]string
	Buildscript string
	Prescript   *string
	Postscript  *string
}

func (p Package) Download() {
	path := getCachePath() + "/source/" + p.Name
	isSrcCached(p.Name)

	if !isExist(path) {
		os.Mkdir(path, 0755)
	}
	getHTTPRes := func(link string) {
		res, err := http.Get(link)
		if err != nil {
			fmt.Println(err)
		}
		defer res.Body.Close()

		out, err := os.Create(path + "/" + lastStr(splitStr(link, "/")))
		if err != nil {
			fmt.Println(err)
			//	return err
		}
		defer out.Close()

		_, err = io.Copy(out, res.Body)
		// // err
	}

	for _, el := range p.Sources {
		_, types := uriMatcher(el[0])
		if types != "" {
			fmt.Println(types)
		}
		// this for match
		// url type by regex
		switch {
		case types == "http":
			getHTTPRes(el[0])
			fmt.Println("http: TBD")
		case types == "tar":
			getHTTPRes(el[0])
			fmt.Println("tar: TBD")
		case types == "git":
			fmt.Println("git : TBD")
			gitclone(lastStr(splitStr(el[0], "+")), "")
		case types == "local":
			fmt.Println("local: TBD")
		}
	}
}

func (p Package) satisfy() ([]string, bool) {
	depends := []string{}
	satisfied := true
	// this for iterate 
	// dependencies list on
	// manifest file
	for _, item := range p.Depends {
		if !isInstalled(item) {
			depends = append(depends, item)
			satisfied = false
		}
	}
	return depends, satisfied
}

func (p Package) build() {

	// Checking for dependencies
	depends, satisfied := p.satisfy()
	if len(depends) != 0 && satisfied {
		fmt.Println("[!] Dependencies not satisfied...")
		for _, item := range depends {
			fmt.Println("- ", item)
		}
		os.Exit(1)
	}

	// Make temporary dir
	makeTempCacheDir()
	srcpath := getCachePath() + "/source-" + getPid() + "/" + p.Name
	os.Mkdir(srcpath,0755)
	if !isExist(srcpath) {
		fmt.Println("[!] No source unique cache dir created")
		os.Exit(1)
	} else {
		os.MkdirAll(srcpath, 0755)
		e := makeFile([]byte(p.Buildscript), srcpath+"/build.sh")
		if e != nil {

		}
	}

	// make binary path and manifestpath for included when
	// archived
	binpath := getCachePath() + "/binary-" + getPid() + "/" + p.Name
	manifestPath := binpath + "/var/db/kartini/installed/" + p.Name
	// make dummy system environment likes
	os.MkdirAll(binpath, 0755)
	os.MkdirAll(manifestPath, 0755)

	if p.Prescript != nil {
		e := makeFile([]byte(*p.Prescript), manifestPath+"/preinstall.sh")
		if e != nil {

		}
		exec.Command("sh", srcpath+"/preinstall.sh", binpath)
	}

	// run build script
	os.Chdir(srcpath)
	runCmd("sh", "./build.sh", binpath)

	// scan path for all builded files
	manifest := scanDir(binpath)

	manifestFile, e := os.Create(manifestPath + "/manifest")
	if e != nil {
		fmt.Println(e)
	}

	// this for index all builded file
	// to manifest list in file
	for _, item := range manifest {
		hash := hashFile(item)
		truePath := strings.Split(item, binpath)[1]
		manifestFile.Write([]byte(truePath + " " + hash))
	}
	manifest = append(manifest,manifestPath + "/manifest")
	manifestFile.Write([]byte(manifestPath + "/manifest"+ " " + hashFile(manifestPath + "/manifest")))

	// this will execute postscript
	// if available
	if p.Postscript != nil {
		e := makeFile([]byte(*p.Postscript), srcpath+"/postinstall.sh")
		if e != nil {

		}
		exec.Command("sh", srcpath+"/preinstall.sh", binpath)
	}

	for index,element := range manifest{
		manifest[index] = "."+splitStr(element,binpath)[1]
	}

	// Change dir to binpath and make archive
	os.Chdir(binpath)
	runCmd("tar",
	"-cvf",
	getCachePath() + "/binary/"+ p.Name + "%"+ p.Version + ".tar.xz",
	".")
	removeTempCacheDir()

}

func (p Package) install() {
	bin := getCachePath() + "/binary/" + p.Name + "%" + p.Version + ".tar.xz"
	untar(getFileSystem(), bin)
}

func (p Package) extract(path string) {
	var caches string = (getCachePath() + "/source/" + p.Name)
	var temp string = (getCachePath() + "/source-"+ getPid() + "/"+ p.Name)
	os.MkdirAll(temp,755)
	for _, item := range p.Sources {
		// target untar should be temporary caches
		runCmd("tar","-xf",caches+"/"+lastStr(splitStr(item[0], "/")),"--strip-components=1","-C",temp)
	}
}

func (p Package) install(){
	// archiving whole file and folder inside
	// binpath
	binpath := getCachePath()+"/binary/"+p.Name+"%"+p.Version+".tar.xz"
	runCmd("tar","-xvf",binpath,"-C",os.Getenv("KARTINI_ROOT"))
}
func (p Package) remove(){

	// Remove file based on manifest files 
	// each lines
	manifestPath := getManifestPath() +"/" + p.Name
	f,_ := os.Open(manifestPath + "/manifest")
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)

	for scan.Scan(){
		os.Remove(os.Getenv("KARTINI_ROOT") +"/" + splitStr(scan.Text()," ")[0])
	}
}

func (p Package) details() {
	fmt.Println("name : ", p.Name)
	fmt.Println("version : ", p.Version)
	for _, el := range p.Depends {
		fmt.Println("Depends on ", el)
	}
	fmt.Println(p.Buildscript)
}

func tomlToPackage(pathToml string) Package {
	p := new(Package)
	if _, e := toml.DecodeFile(pathToml, &p); e != nil {
		fmt.Println("[!] Invalid package descriptor:", e)
		os.Exit(1)
	}
	return *p
}

func main() {

	// this router for
	// connecting user terminal
	// input with desired one
	args := os.Args
	now := time.Now()
	if len(args) < 2 {
		fmt.Println("[?] Needs supplied by argument")
		os.Exit(1)
	}
	switch ""{
	case os.Getenv("KARTINI_ROOT"):
		fmt.Println("[!] KARTINIT_ROOT need be set")
		os.Exit(1)
	case os.Getenv("KARTINI_PATH"):
		fmt.Println("[!] KARTINI_PATH need be set")
		os.Exit(1)
	case os.Getenv("KARTINI_CACHE"):
		fmt.Println("[!] KARTINI_CACHE need be set")
		os.Exit(1)
	}

	switch args[1] {
	case "add":
		for _, item := range args[2:] {
			fmt.Println("Shall added : ", item)
			installPackage(item)
		}
	case "build":
		for _, item := range args[2:] {
			extractPackage(item)
		}
	case "gv":
		runCmd("go","version")
	case "del":
		for _, item := range args[2:] {
			fmt.Println("Shall deleted : ", item)
			removePackage(item)
		}
	case "get":
		for _, item := range args[2:] {
			fmt.Println("Shall searched: ", item)
			getPackage(item)
		}
	case "find":
		for _, item := range args[2:] {
			if p, found := findPackage(item); found {
				p.details()
			} else {
				fmt.Println("Pacakage Unavailable:", item)
			}
		}
	case "env":
		fmt.Println("Cache point: ", getCachePath())
		fmt.Println("Manifest point: ", getManifestPath())
		fmt.Println("Repo point: ", getLocalRepo())
	case "help":
		fmt.Println(`
usage: kartini <opts> [args,...]

opts:
	add	<package name>	install builded packages
	build	<package name>	build package
	env			check environment variables
	find	<package name>	find package by name 
	get	<package name>	download package resources
	help			this message

		`)
	case "id":
		fmt.Println("PID: ", os.Getpid())
		fmt.Println("PPID: ", os.Getppid())
	case "mother":
		fmt.Println("")
	}
	fmt.Println("["+ args[1] + "] done in " + time.Now().Sub(now).String()+ "")
}
