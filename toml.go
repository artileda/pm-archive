package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"bufio"
	//os/exec"
	"strings"

	"github.com/BurntSushi/toml"
)

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
		// if len(el) == 2 {
		// 	fmt.Println("URL : ", el[0], ", Extraction point :", el[1])
		// } else if len(el) == 1 {
		// 	fmt.Println("URI: ", el[0])
		// }
	}
}

func (p Package) satisfy() ([]string, bool) {
	depends := []string{}
	satisfied := true
	for _, item := range p.Depends {
		if !isInstalled(item) {
			depends = append(depends, item)
			satisfied = false
		}
	}
	return depends, satisfied
}

func (p Package) build() {

	depends, satisfied := p.satisfy()
	if len(depends) != 0 && satisfied {
		fmt.Println("[!] Dependencies not satisfied...")
		for _, item := range depends {
			fmt.Println("- ", item)
		}
		os.Exit(1)
	}
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

	binpath := getCachePath() + "/binary-" + getPid() + "/" + p.Name
	manifestPath := binpath + "/var/db/kartini/installed/" + p.Name
	// make dummy system environment likes
	os.MkdirAll(binpath, 0755)
	os.MkdirAll(manifestPath, 0755)

	if p.Prescript != nil {
		e := makeFile([]byte(*p.Prescript), manifestPath+"/preinstall.sh")
		if e != nil {

		}
		//exec.Command("sh", srcpath+"/preinstall.sh", binpath)
	}

	// run build script
	runCmd("sh", srcpath+"/build.sh", binpath)

	//fmt.Println(string(out))
	manifest := scanDir(binpath)


//	fmt.Println(manifest)
//	fmt.Println(binpath)

	// create manifest subfolder

	manifestFile, e := os.Create(manifestPath + "/manifest")
	if e != nil {
		fmt.Println(e)
	}
		// for making fingerprint	// make manifest

	for _, item := range manifest {
		hash := hashFile(item)
		truePath := strings.Split(item, binpath)[1]
		manifestFile.Write([]byte(truePath + " " + hash))
	}
	manifest = append(manifest,manifestPath + "/manifest")
	manifestFile.Write([]byte(manifestPath + "/manifest"+ " " + hashFile(manifestPath + "/manifest")))
	
	if p.Postscript != nil {
		e := makeFile([]byte(*p.Postscript), srcpath+"/postinstall.sh")
		if e != nil {

		}
		//exec.Command("sh", srcpath+"/preinstall.sh", binpath)
	}

	for index,element := range manifest{
		manifest[index] = "."+splitStr(element,binpath)[1]
	}

	//fmt.Println(manifest)
	//Archive this
	os.Chdir(binpath)
	s,_ := os.Getwd()
	fmt.Println(s)
//	fmt.Println(taring(
//		manifest,
//		getCachePath()+"/binary/"+p.Name+"%"+p.Version+".tar.xz"))
	runCmd("tar",
	"-cvf",
	getCachePath() + "/binary/"+ p.Name + "%"+ p.Version + ".tar.xz",
	".")

}

func (p Package) extract(path string) {
	var caches string = (getCachePath() + "/source/" + p.Name)
	for _, item := range p.Sources {
		// target untar should be temporary caches
		untar(caches+"/"+lastStr(splitStr(item[0], "/")), caches)
	}
}

func (p Package) install(){
	binpath := getCachePath()+"/binary/"+p.Name+"%"+p.Version+".tar.xz"
//	fmt.Println(binpath)
//	fmt.Println(os.Getenv("KARTINI_ROOT"))
	runCmd("tar","-xvf",binpath,"-C",os.Getenv("KARTINI_ROOT"))
}
func (p Package) remove(){
	manifestPath := getManifestPath() +"/" + p.Name
	fmt.Println(manifestPath)
	f,_ := os.Open(manifestPath + "/manifest")
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)
	//var manifestList []string

	for scan.Scan(){
		//fmt.Println(scan.Text())
		//nifestList = append(manifestList,
		os.Remove(os.Getenv("KARTINI_ROOT") +"/" + splitStr(scan.Text()," ")[0])
	}
	// this still not remove a manifest folder
	//fmt.Println(manifestList)
}

func (p Package) details() {
	fmt.Println("Nama : ", p.Name)
	fmt.Println("Versi : ", p.Version)
	//p.Download()
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
	// p := new(Package)
	// _, e := toml.DecodeFile("./package.toml", &p)
	// if e != nil {
	// 	fmt.Println(e)
	// }

	args := os.Args
	if len(args) < 2 {
		fmt.Println("[?] Needs supplied by argument")
		os.Exit(1)
	}

	//fmt.Println(args[1])
	switch args[1] {
	case "add":
		fmt.Println("add's subcommand summoned")
		fmt.Println("opts supplied: ", args[2:])
		for _, item := range args[2:] {
			fmt.Println("Shall added : ", item)
			installPackage(item)
		}
	case "build":
		fmt.Println("build's subcommand summoned")
		fmt.Println("opts supplie: ", args[2:])
		for _, item := range args[2:] {
			extractPackage(item)
		}
	case "gv":
		runCmd("go","version")
	case "del":
		fmt.Println("del's subcommand summoned")
		fmt.Println("opts supplied: ", args[2:])
		for _, item := range args[2:] {
			fmt.Println("Shall deleted : ", item)
			removePackage(item)
		}
	case "get":
		fmt.Println("get's subcommand summoned")
		fmt.Println("opts supplied: ", args[2:])
		for _, item := range args[2:] {
			fmt.Println("Shall searched: ", item)
			getPackage(item)
		}
	case "find":
		fmt.Println("find's subcommand summoned")
		fmt.Println(args[2:])
		for _, item := range args[2:] {
			if p, found := findPackage(item); found {
				fmt.Println("Pacakage available:", item)
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
}
