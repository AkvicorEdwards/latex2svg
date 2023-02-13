package main

import (
	"flag"
	"fmt"
	"github.com/AkvicorEdwards/util"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	c := flag.String("c", "config.ini", "config file")
	k := flag.String("k", "random", "config file, 'random' and 'empty' are special values")
	flag.Parse()
	Load(*c, *k)
	// check bash path
	if util.FileStat(config.BashPath) != 2 {
		log.Fatalf("[%s] not found", config.BashPath)
	}
	// check temp dir
	if util.FileStat(config.Path.Temp) == 0 {
		_ = os.MkdirAll(config.Path.Temp, 0700)
	}
	if util.FileStat(config.Path.Temp) != 1 {
		log.Fatalf("[%s] is not dir", config.Path.Temp)
	}
	// check cache dir
	if util.FileStat(config.Path.Cache) == 0 {
		_ = os.MkdirAll(config.Path.Cache, 0700)
	}
	if util.FileStat(config.Path.Cache) != 1 {
		log.Fatalf("[%s] is not dir", config.Path.Cache)
	}
	// check script
	if util.FileStat(config.Path.Script) != 2 {
		f, err := os.OpenFile(config.Path.Script, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0700)
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
	dir, err := ioutil.ReadDir(config.Path.Temp)
	if err != nil {
		return
	}
	for _, d := range dir {
		_ = os.RemoveAll(path.Join(config.Path.Temp, d.Name()))
	}

	http.HandleFunc("/", index)
	http.HandleFunc("/clear/cache", clearCache)
	http.HandleFunc("/clear/temp", clearTemp)

	addr := fmt.Sprintf("%s:%d", config.Addr, config.Port)
	log.Printf("ListenAndServe %s. clear_key:[%s]", addr, config.ClearKey)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
