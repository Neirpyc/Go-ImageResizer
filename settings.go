package main

import (
	"errors"
	"os"
	"path"
	"strconv"
	"strings"
)

func (sT SettingsToml) parse() Settings {
	s := Settings{}

	s.Port = sT.Port
	s.PreCache = sT.PreCache

	var err error

	if f, err := os.Open(sT.CacheFolder); err != nil {
		exec, err := os.Executable()
		if err != nil {
			panic(err)
		}
		f, err = os.Open(path.Dir(exec) + "/" + sT.CacheFolder)
		if err != nil {
			panic(err)
		}
		s.CacheFolder = path.Dir(exec) + "/" + sT.CacheFolder
		f.Close()
	} else {
		f.Close()
		s.CacheFolder = sT.CacheFolder
	}

	if f, err := os.Open(sT.InputFolder); err != nil {
		exec, err := os.Executable()
		if err != nil {
			panic(err)
		}
		f, err = os.Open(path.Dir(exec) + "/" + sT.InputFolder)
		if err != nil {
			panic(err)
		}
		s.InputFolder = path.Dir(exec) + "/" + sT.InputFolder
		f.Close()
	} else {
		f.Close()
		s.InputFolder = sT.InputFolder
	}

	exec, err := os.Executable()
	if err != nil {
		panic(err)
	}
	if s.LogFile, err = os.OpenFile(exec+"/"+sT.LogFile, os.O_APPEND|os.O_WRONLY, 0600); err != nil {
		if s.LogFile, err = os.Create(sT.LogFile); err != nil {
			panic(err)
		}
	}

	var suffix string
	if s.MaxCacheSize, err = strconv.ParseInt(sT.MaxCacheSize[:len(sT.MaxCacheSize)-2], 10, 64); err != nil {
		if s.MaxCacheSize, err = strconv.ParseInt(sT.MaxCacheSize[:len(sT.MaxCacheSize)-3], 10, 64); err != nil {
			panic(err)
		} else {
			suffix = sT.MaxCacheSize[len(sT.MaxCacheSize)-3:]
		}
	} else {
		suffix = sT.MaxCacheSize[len(sT.MaxCacheSize)-2:]
	}
	switch len(suffix) {
	case 1:
		switch suffix {
		case "b":
		case "B":
			s.MaxCacheSize *= 8
		default:
			panic(errors.New("Invalid size suffix " + suffix))
		}
	case 2:
		switch suffix[1] {
		case 'b':
		case 'B':
			s.MaxCacheSize *= 8
		default:
			panic(errors.New("Invalid size suffix " + suffix))
		}

		switch strings.ToLower(suffix)[0] {
		case 'k':
			s.MaxCacheSize *= 1024
		case 'm':
			s.MaxCacheSize *= 1024 * 1024
		case 'g':
			s.MaxCacheSize *= 1024 * 1024 * 1024
		case 't':
			s.MaxCacheSize *= 1024 * 1024 * 1024 * 1024
		default:
			panic(errors.New("Invalid size suffix " + suffix))
		}
	default:
		panic(errors.New("Invalid size suffix " + suffix))
	}
	s.MaxCacheSize /= 8

	return s
}

type Settings struct {
	Port         string
	CacheFolder  string
	MaxCacheSize int64 //in byte
	InputFolder  string
	LogFile      *os.File
	PreCache     bool
}

type SettingsToml struct {
	Port         string `toml:"port"`
	CacheFolder  string `toml:"cacheFolder"`
	MaxCacheSize string `toml:"maxCacheSize"`
	InputFolder  string `toml:"inputFolder"`
	LogFile      string `toml:"logFile"`
	PreCache     bool   `toml:"preCache"`
}
