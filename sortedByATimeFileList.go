package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"syscall"
)

type ByATime struct {
	Files []os.FileInfo
	Mutex sync.RWMutex
	Size  int64
	Directory string
}

func createCacheFileList() []os.FileInfo {
	cacheFilesList = []os.FileInfo{}
	filepath.Walk(settings.CacheFolder, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			cacheFilesList = append(cacheFilesList, info)
		}
		return err
	})
	return cacheFilesList
}

func NewByATimeFromFileList(list []os.FileInfo, directory string) ByATime {
	var size int64
	for _, f := range list {
		size += f.Size()
	}
	byATime := ByATime{
		Files: list,
		Mutex: sync.RWMutex{},
		Size:  size,
		Directory: directory,
	}
	byATime.Sort()
	byATime.ValidateSize()
	return byATime
}

func (f ByATime) Sort() {
	sort.Sort(f)
}

func (f ByATime) ValidateSize(){
	f.Mutex.RLock()
	if f.Size > settings.MaxCacheSize{
		f.Mutex.RUnlock()

		f.Mutex.Lock()
		for f.Size > settings.MaxCacheSize{
			f.Size -= f.Files[0].Size()
			if err := os.Remove(f.Directory + "/" + f.Files[0].Name()); err != nil{
				log.Println(err.Error())
			}
			f.Files = f.Files[1:]
		}
		f.Mutex.Unlock()
	}else{
		f.Mutex.RUnlock()
	}
}

func (f ByATime) Find(name string) int{
	f.Mutex.RLock()
	defer f.Mutex.RUnlock()
	for i := 0; i < f.Len(); i++{
		if f.Files[i].Name() == name{
			return i
		}
	}
	return -1
}

func (f ByATime) Insert(info os.FileInfo) {
	f.Mutex.Lock()
	f.Files = append(f.Files, info)
	f.Size += info.Size()
	f.Mutex.Unlock()
	for i := f.Len() - 1; i > 1; i-- {
		if f.Less(i, i-1) {
			f.Swap(i, i-1)
		} else {
			break
		}
	}
	f.ValidateSize()
}

func (f ByATime) Update(name string) {
	i := f.Find(name)
	if i >= f.Len()-1 || i < 0{
		return
	}
	for j := i; j+1 < f.Len()-1; j++ {
		if f.Less(j+1, j) {
			f.Swap(j+1, j)
		} else {
			break
		}
	}
}

func (f ByATime) Sorted() bool {
	for i := f.Len() - 1; i > 1; i-- {
		if !f.Less(i, i-1) {
			return false
		}
	}
	return true
}

func (f ByATime) Len() int {
	f.Mutex.RLock()
	defer f.Mutex.RUnlock()
	return len(f.Files)
}

func (f ByATime) Swap(i, j int) {
	f.Mutex.Lock()
	f.Files[i], f.Files[j] = f.Files[j], f.Files[i]
	f.Mutex.Unlock()
}

func (f ByATime) Less(i, j int) bool {
	f.Mutex.RLock()
	statI := f.Files[i].Sys().(*syscall.Stat_t)
	statJ := f.Files[j].Sys().(*syscall.Stat_t)
	f.Mutex.RUnlock()
	return statI.Atim.Sec < statJ.Atim.Sec
}
