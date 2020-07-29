package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func fetchFromCache(url string) ([]byte, string, error){
	hash := sha512Str(url)

	f, err := os.Open(settings.CacheFolder + "/" + hash)
	if err == nil {
		cache.Update(hash)
		content, err := ioutil.ReadAll(f)
			return content, hash, err
	}
	return nil, hash, err
}

func preCache(){
	err := filepath.Walk(settings.InputFolder,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir(){
				urls := generateUrlsForFile(strings.TrimRight(strings.TrimLeft(path, settings.InputFolder+"/"), ".jpg"))
				for _, url := range(urls){
					if _, urlHash, err := fetchFromCache(url); err != nil{
						if req, err := parseReqUrl(url); err == nil {
							if content, err, _ := getImage(req); err == nil {
								f, err := os.Create(settings.CacheFolder + "/" + urlHash)
								if err != nil {
									return err
								}
								_, err = f.Write(content)
								if err != nil {
									return err
								}
								f.Close()
								fInfo, err := os.Stat(settings.CacheFolder + "/" + urlHash)
								if err != nil {
									return err
								}
								cache.Insert(fInfo)
							}else{
								return err
							}
						}else {
							return err
						}
					}
				}
			}
			return nil
		})
	if err != nil {
		fmt.Println(err)
	}
}

func generateUrlsForFile(file string)[]string{
	exts := []string{".jpg", ".webp"}
	lods := []string{"1", "2", "3"}
	sizes := []string{"small", "medium", "big"}
	var urls []string
	for _, ext := range(exts){
		for _, lod := range(lods){
			if lod == "3"{
				for _, size := range(sizes){
					urls = append(urls, generateUrl(file, lod, size, ext))
				}
			}else{
				urls = append(urls, generateUrl(file, lod, "", ext))
			}
		}
	}
	return urls
}

func generateUrl(file string, lod string, size string, ext string) string{
	if lod == "3"{
		return "/images/lod/3/"+size+"/"+file+ext
	}
	return "/images/lod/"+lod+"/"+file+ext
}
