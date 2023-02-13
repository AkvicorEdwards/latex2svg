package main

import (
	"github.com/AkvicorEdwards/util"
	"github.com/go-ini/ini"
	"log"
	"path"
)

var cfg *ini.File
var config *Model

type Model struct {
	Addr          string    `ini:"http_addr"`
	Port          int       `ini:"http_port"`
	BashPath      string    `ini:"bash_path"`
	Workspace     string    `ini:"workspace"`
	ClearKey      string    `ini:"clear_key"`
	ScriptTimeout int       `ini:"script_timeout"`
	Path          PathModel `ini:"-"`
}

type PathModel struct {
	Script string
	Temp   string
	Cache  string
}

func (p *PathModel) JoinTemp(s string) string {
	return path.Join(p.Temp, s)
}

func (p *PathModel) JoinCache(s string) string {
	return path.Join(p.Cache, s)
}

func Load(c string, key string) {
	var err error
	config = new(Model)
	cfg, err = ini.Load(c)
	if err != nil {
		cfg = ini.Empty()
		config.Addr = "0.0.0.0"
		config.Port = 8080
		config.BashPath = "/usr/bin/sh"
		config.Workspace = "./workspace"
		config.ScriptTimeout = 15
		if key == "random" {
			config.ClearKey = util.RandomString(7)
		} else if key == "empty" {
			config.ClearKey = ""
		} else {
			config.ClearKey = key
		}
		err = ini.ReflectFrom(cfg, config)
		if err != nil {
			log.Fatalln("write config failed", err)
			return
		}
		err = cfg.SaveTo(c)
		if err != nil {
			log.Fatalln("write config failed", err)
			return
		}
	} else {
		err = cfg.MapTo(config)
		if err != nil {
			log.Fatalln("read config failed", err)
			return
		}
	}
	config.Path.Script = path.Join(config.Workspace, "conv.sh")
	config.Path.Temp = path.Join(config.Workspace, "temp")
	config.Path.Cache = path.Join(config.Workspace, "cache")
}
