package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

var lock = sync.RWMutex{}

const baseMath = "math"
const baseCJK = "cjk"
const baseDoc = "doc"
const baseEmpty = "empty"
const baseError = baseCJK
const latexError = "ERROR"

func index(w http.ResponseWriter, r *http.Request) {
	if len(config.ClearKey) != 0 && r.URL.Query().Get("key") != config.ClearKey {
		return
	}
	latex := strings.TrimSpace(r.URL.Query().Get("latex"))
	if len(latex) == 0 {
		return
	}
	crop := strings.TrimSpace(r.URL.Query().Get("crop"))
	filetype := strings.TrimSpace(r.URL.Query().Get("type"))
	transp := strings.TrimSpace(r.URL.Query().Get("transp"))
	base := strings.TrimSpace(r.URL.Query().Get("base"))
	ok := conv(w, base, latex, filetype, crop, transp)
	if !ok {
		ok = conv(w, baseError, latexError, filetype, crop, transp)
	}
	if !ok {
		log.Println("response failed")
	}
}

func clearCache(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("key") != config.ClearKey {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	dir, err := ioutil.ReadDir(config.Path.Cache)
	if err != nil {
		_, _ = w.Write([]byte("clear cache failed"))
		return
	}
	for _, d := range dir {
		_ = os.RemoveAll(config.Path.JoinCache(d.Name()))
	}
	_, _ = w.Write([]byte("clear cache successful"))
}

func clearTemp(w http.ResponseWriter, r *http.Request) {
	if r.URL.Query().Get("key") != config.ClearKey {
		return
	}
	lock.Lock()
	defer lock.Unlock()
	dir, err := ioutil.ReadDir(config.Path.Temp)
	if err != nil {
		_, _ = w.Write([]byte("clear temp failed"))
		return
	}
	for _, d := range dir {
		_ = os.RemoveAll(config.Path.JoinTemp(d.Name()))
	}
	_, _ = w.Write([]byte("clear temp successful"))
}
