package main

import (
	"fmt"
	"github.com/AkvicorEdwards/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

func index(w http.ResponseWriter, r *http.Request) {
	tx := r.URL.Query().Get("latex")
	if len(tx) == 0 {
		_, _ = w.Write([]byte(latexErrorReturn))
		return
	}
	id := fmt.Sprintf("%x", util.SHA256String(tx))
	tPath := path.Join(tempPath, id)
	cPath := path.Join(cachePath, id)
	if util.FileStat(cPath) == 2 {
		res, err := os.ReadFile(cPath)
		if err == nil {
			_, _ = w.Write(res)
			return
		}
		log.Printf("read error %v", err)
		_ = os.Remove(cPath)
	}
	f, e := os.Create(tPath)
	if e != nil {
		_, _ = w.Write([]byte(latexErrorReturn))
		log.Printf("read error %v", e)
		return
	}
	_, _ = f.WriteString(fmt.Sprintf(latexBase, tx))
	_ = f.Close()
	cmd := exec.Command(bashPath, "-c", fmt.Sprintf("%s %s", scriptPath, id))
	fin := make(chan bool, 1)
	err := make(chan bool, 1)
	go func() {
		e := cmd.Run()
		if e != nil {
			err <- true
			log.Printf("read error %v", e)
			return
		}
		fin <- true
	}()
	select {
	case <-time.After(runWait * time.Second):
		_, _ = w.Write([]byte(latexErrorReturn))
		return
	case <-err:
		_, _ = w.Write([]byte(latexErrorReturn))
		return
	case <-fin:
	}
	_ = os.Rename(tPath, cPath)
	res, e := os.ReadFile(cPath)
	if e == nil {
		_, _ = w.Write(res)
		return
	}
	_, _ = w.Write([]byte(latexErrorReturn))
	log.Printf("read error %v", e)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if util.FileStat(bashPath) != 2 {
		log.Fatalf("[%s] not found", bashPath)
	}
	// temp dir
	if util.FileStat(tempPath) == 0 {
		_ = os.MkdirAll(tempPath, 0700)
	}
	if util.FileStat(tempPath) != 1 {
		log.Fatalf("[%s] is not dir", tempPath)
	}
	// cache dir
	if util.FileStat(cachePath) == 0 {
		_ = os.MkdirAll(cachePath, 0700)
	}
	if util.FileStat(cachePath) != 1 {
		log.Fatalf("[%s] is not dir", cachePath)
	}
	// script
	if util.FileStat(scriptPath) != 2 {
		f, err := os.OpenFile(scriptPath, os.O_CREATE|os.O_WRONLY, 0700)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.WriteString(scriptBase())
		if err != nil {
			log.Fatal(err)
		}
		_ = f.Close()
	}
	// clear temp
	dir, err := ioutil.ReadDir(tempPath)
	if err != nil {
		return
	}
	for _, d := range dir {
		_ = os.RemoveAll(path.Join(tempPath, d.Name()))
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/clear", clear)

	log.Printf("ListenAndServe %s", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func clear(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("key") != clearKey {
		return
	}
	dir, err := ioutil.ReadDir(cachePath)
	if err != nil {
		_, _ = w.Write([]byte("clear cache failed"))
		return
	}
	for _, d := range dir {
		_ = os.RemoveAll(path.Join(cachePath, d.Name()))
	}
	_, _ = w.Write([]byte("clear cache successful"))
}
