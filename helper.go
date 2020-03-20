package main

import (
	"strings"
	"github.com/mholt/archiver/v3"
)

func splitStr(s string, delimeter string) []string {
	return strings.Split(s, delimeter)
}
func lastStr(s []string) string {
	return s[len(s)-1]
}
func untar(path string, dest string)error{
	return archiver.Unarchive(path,dest)
}