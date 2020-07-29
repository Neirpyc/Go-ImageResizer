package main

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
)

var (
	lod12Regex     *regexp.Regexp
	lod3Regex      *regexp.Regexp
	cacheFilesList []os.FileInfo
)

const (
	Lod1Width       = 40
	Lod1Quality     = 40
	Lod2Width       = 250
	Lod2Quality     = 65
	Lod3JpegQuality = 70
	Lod3WebpQuality = 80
	Lod3SmallWidth  = 450
	Lod3MediumWidth = 820
)

var settings Settings
var cache ByATime

func init() {
	f, err := os.Open("settings.toml")
	if err != nil {
		exec, err := os.Executable()
		if err != nil {
			panic(err)
		}
		f, err = os.Open(path.Dir(exec) + "/settings.toml")
		if err != nil {
			panic(err)
		}
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	var sT SettingsToml
	if _, err = toml.Decode(string(data), &sT); err != nil {
		panic(err)
	}
	settings = sT.parse()
	lod12Regex = regexp.MustCompile("/images/lod/([12])/(.*)\\.(jpg|webp)")
	lod3Regex = regexp.MustCompile("/images/lod/3/(small|medium|big)/(.*)\\.(jpg|webp)")
	cache = NewByATimeFromFileList(createCacheFileList(), settings.CacheFolder)
	if settings.PreCache{
		go preCache()
	}
}

func main() {
	defer settings.LogFile.Close()

	log.SetOutput(settings.LogFile)
	log.SetFlags(log.Ldate | log.Llongfile | log.Lmicroseconds | log.Ltime)
	log.Println("Starting image resizer")

	http.HandleFunc("/", serveImage)
	log.Println(http.ListenAndServe(":"+settings.Port, nil))
}
