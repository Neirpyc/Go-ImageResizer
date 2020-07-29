package main

import (
	"bytes"
	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
	"image/jpeg"
	"log"
	"net/http"
	"os"
)

func serveImage(w http.ResponseWriter, r *http.Request) {
	content, err, code := getContent(r.RequestURI, r.Header)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(code)
	}
	w.Write(content)
}

func getContent(url string, header http.Header) ([]byte, error, int) {
	url = rewriteUrlWithWebp(url, header)
	content, urlHash, err := fetchFromCache(url)
	if err == nil {
		cache.Update(urlHash)
		return content, nil, http.StatusOK
	}

	if req, err := parseReqUrl(url); err == nil {
		if content, err, code := getImage(req); err == nil {
			f, err := os.Create(settings.CacheFolder + "/" + urlHash)
			if err != nil {
				return content, nil, http.StatusInternalServerError
			}
			_, err = f.Write(content)
			if err != nil {
				return content, nil, http.StatusInternalServerError
			}
			f.Close()
			fInfo, err := os.Stat(settings.CacheFolder + "/" + urlHash)
			if err != nil {
				return content, nil, http.StatusInternalServerError
			}
			cache.Insert(fInfo)
			return content, nil, http.StatusOK
		}else{
			return nil, err, code
		}
	}else {
		return nil, err, http.StatusBadRequest
	}
}

func getImage(r Request) ([]byte, error, int) {
	f, err := os.Open(settings.InputFolder + "/" + r.file + ".jpg")
	if err != nil {
		return nil, err, http.StatusNotFound
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	var width, height uint
	var quality int
	switch r.lod {
	case 1:
		width = Lod1Width
		quality = Lod1Quality
	case 2:
		width = Lod2Width
		quality = Lod2Quality
	case 3:
		switch r.size {
		case SizeSmall:
			width = Lod3SmallWidth
		case SizeMedium:
			width = Lod3MediumWidth
		case sizeBig:
			width = uint(img.Bounds().Dx())
		}

		switch r.format {
		case "jpg":
			quality = Lod3JpegQuality
		case "webp":
			quality = Lod3WebpQuality
		}
	}

	height = uint(float64(width) / float64(img.Bounds().Dx()) * float64(img.Bounds().Dy()))

	buf := new(bytes.Buffer)

	resized := resize.Resize(width, height, img, resize.Bicubic)

	switch r.format {
	case "jpg":
		err = jpeg.Encode(buf, resized, &jpeg.Options{Quality: quality})
	case "webp":
		err = webp.Encode(buf, resized, &webp.Options{
			Lossless: false,
			Quality:  float32(quality),
			Exact:    false,
		})
	}
	return buf.Bytes(), err, http.StatusOK
}
