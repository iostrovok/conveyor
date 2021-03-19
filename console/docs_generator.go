package main

/*

	Just for fun - simple script which grabs documentation from "godoc" local server

	Simple script for testing go-log & go-log/event
*/

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	reducePathLen = 2
)

func filepathDir() string {
	dir := os.Getenv("CURDIR")
	if dir == "" {
		log.Fatal("empty CURDIR")
	}

	fmt.Printf("CURDIR: dir-dir: %s\n", dir)
	if err := os.Chdir(dir); err != nil {
		log.Fatal(err)
	}

	return dir
}

func main() {
	curDir := filepathDir()
	docFullDir := filepath.Join(curDir, "docs")
	goModFile := filepath.Join(curDir, "go.mod")
	tmpDir := filepath.Join(curDir, "127.0.0.1:6060")

	if err := os.RemoveAll(docFullDir); err != nil {
		log.Printf("rm error: %v\n", err)
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		log.Printf("rm error: %v\n", err)
	}

	file, err := os.Open(goModFile)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	module := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			tmp := strings.Split(line, " ")
			module = tmp[1]

			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("scanner error: %v\n", err)

		return
	}

	_, pkgName := filepath.Split(module)
	fmt.Printf("\nCreate documentation for: \x1b[35m %s \x1b[0m\n\n", module)

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		cmd := exec.CommandContext(ctx, "./bin/godoc", "-v", "-http=:6060", "-goroot=./")
		if err := cmd.Run(); err != nil {
			log.Printf("godoc error: %v\n", err)
		}
		wg.Done()
	}()

	time.Sleep(1 * time.Second)

	wgetOptions := []string{"-r", "-np", "-N", "-E", "-p", "-k", "http://127.0.0.1:6060/pkg/" + module}

	cmd := exec.Command("/opt/local/bin/wget", wgetOptions...)
	if err := cmd.Run(); err != nil {
		log.Printf("wget error: %v\n", err)
	}

	cancel()
	wg.Wait()

	if err := os.Mkdir(docFullDir, 0777); err != nil {
		log.Printf("mkdir error: %v\n", err)

		return
	}

	if err := os.Rename(filepath.Join(tmpDir, "pkg", module), filepath.Join(docFullDir, pkgName)); err != nil {
		log.Printf("mv error: %v\n", err)

		return
	}

	if err := os.Rename(filepath.Join(tmpDir, "pkg", module+".html"),
		filepath.Join(docFullDir, pkgName, "index.html")); err != nil {
		log.Printf("mv error: %v\n", err)

		return
	}

	if err := os.Rename(filepath.Join(tmpDir, "lib"), filepath.Join(docFullDir, "lib")); err != nil {
		log.Printf("mv error: %v\n", err)

		return
	}

	if err := os.RemoveAll(tmpDir); err != nil {
		log.Printf("rm error: %v\n", err)
	}

	cleanAllFiles(curDir, docFullDir, pkgName)
}

func cleanOneFile(path, pkgName string) error {
	// rg is the constant variable
	rg := regexp.MustCompile(`(href|src)="(../)+`)

	level := len(strings.Split(path, "/")) - reducePathLen

	file, err := os.Open(path)
	if err != nil {
		return err
	}

	out := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		res := rg.FindAllString(line, -1)
		if len(res) == 0 {
			out += line + "\n"

			continue
		}

		for _, oldSrc := range res {
			newSrc := strings.TrimPrefix(strings.TrimSuffix(oldSrc, "../"), "../") + strings.Repeat("../", level)
			line = strings.Replace(line, oldSrc, newSrc, 1)
		}
		out += line + "\n"
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		panic(err)
	}

	out = strings.ReplaceAll(out, ` href="`+pkgName+"/", ` href="`)
	out = strings.ReplaceAll(out, ` href="`+pkgName+`.html`, ` href="index.html`)

	return ioutil.WriteFile(path, []byte(out), 0644)
}

func cleanAllFiles(curDir, docFullDir, pkgName string) {
	f := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("prevent panic by handling failure accessing a path %q: %v\n", path, err)

			return err
		}

		if info.IsDir() || strings.LastIndex(info.Name(), ".html") == -1 {
			return nil
		}

		path = strings.Replace(path, curDir+"/", "", 1)
		log.Printf("visited file or dir: \x1b[35m%q\x1b[0m\n", path)

		return cleanOneFile(path, pkgName)
	}

	if err := filepath.Walk(docFullDir, f); err != nil {
		log.Fatalf("error walking the path %q: %v\n", docFullDir, err)
	}
}
