package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/mholt/archiver/v3"
	"gopkg.in/src-d/go-git.v4"
)

func splitStr(s string, delimeter string) []string {
	return strings.Split(s, delimeter)
}
func lastStr(s []string) string {
	return s[len(s)-1]
}

func gitclone(url string, out string) {
	_, e := git.PlainClone(out, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if e != nil {
		fmt.Println(e)
		return
	}
}

func taring(path []string,dest string) error{
	return archiver.Archive(path,dest)
}

func untar(path string, dest string) error {
	return archiver.Unarchive(path, dest)
}

func uriMatcher(uri string) (string, string) {
	gitMatch := regexp.MustCompile(`(git+:\/\/)+`)
	httpMatch := regexp.MustCompile(`(:\/\/)`)
	localMatch := regexp.MustCompile(`(.)+`)
	tarMatch := regexp.MustCompile(`(:\/\/)?(tar)+`)

	switch {
	case gitMatch.MatchString(uri):
		return uri, "git"
	case tarMatch.MatchString(uri):
		return uri, "tar"
	case httpMatch.MatchString(uri):
		return uri, "http"
	case localMatch.MatchString(uri):
		return uri, "local"
	}
	return uri, ""
}
