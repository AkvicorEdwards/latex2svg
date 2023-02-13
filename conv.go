package main

import (
	"fmt"
	"github.com/AkvicorEdwards/util"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

// filetype: pdf/svg/png/jpg
// crop: 1/0
// transparent: 1/0
func conv(w io.Writer, base, latex, filetype, crop, transparent string) bool {
	lock.RLock()
	defer lock.RUnlock()
	if filetype != "pdf" && filetype != "svg" && filetype != "png" && filetype != "jpg" {
		filetype = "svg"
	}
	if filetype != "png" || transparent != "1" {
		transparent = "0"
	}
	if crop != "1" {
		crop = "0"
	}
	if base != baseMath && base != baseCJK && base != baseDoc && base != baseEmpty {
		base = baseEmpty
	}
	// generate file id
	id := fmt.Sprintf("%x", util.SHA256String(fmt.Sprintf("%s,%s,%s,%s,%s", base, filetype, crop, transparent, latex)))
	cacheFile := config.Path.JoinCache(id)
	if util.FileStat(cacheFile) == 2 {
		res, err := os.ReadFile(cacheFile)
		if err == nil {
			_, _ = w.Write(res)
			return true
		}
		log.Printf("read error %v", err)
	}
	tempFile := config.Path.JoinTemp(id)
	ok := false
	switch base {
	case baseMath:
		ok = writeFile(tempFile, latexMath(latex))
	case baseCJK:
		ok = writeFile(tempFile, latexCJK(latex))
	case baseDoc:
		ok = writeFile(tempFile, latexDoc(latex))
	case baseEmpty:
		ok = writeFile(tempFile, latex)
	}
	if !ok {
		return false
	}

	cmd := exec.Command(config.BashPath, "-c", fmt.Sprintf("%s %s %s %s %s", config.Path.Script, id, crop, filetype, transparent))
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
	case <-time.After(time.Duration(config.ScriptTimeout) * time.Second):
		return false
	case <-err:
		return false
	case <-fin:
	}

	e := os.Rename(tempFile, cacheFile)
	if e != nil {
		_ = os.Remove(tempFile)
		log.Printf("move error %v", e)
		return false
	}
	res, e := os.ReadFile(cacheFile)
	if e != nil {
		log.Printf("read error %v", e)
		return false
	}
	_, _ = w.Write(res)
	return true
}

func writeFile(name, content string) bool {
	f, e := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0700)
	if e != nil {
		log.Printf("write file [%s] error %v", name, e)
		return false
	}
	defer func() {
		_ = f.Close()
	}()
	_, _ = f.WriteString(content)
	return true
}
