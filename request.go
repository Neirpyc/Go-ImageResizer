package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	SizeSmall = iota
	SizeMedium
	sizeBig
)

type Request struct {
	lod    int
	file   string
	size   int
	format string
}

func parseReqUrl(url string) (Request, error) {
	var r Request

	if lod12Regex.MatchString(url) {
		groups := lod12Regex.FindStringSubmatch(url)

		//LOD
		lod, err := strconv.ParseInt(groups[1], 10, 32)
		r.lod = int(lod)
		if err != nil {
			return r, err
		}

		//file
		r.file = groups[2]

		//format
		r.format = groups[3]
	} else if lod3Regex.MatchString(url) {
		groups := lod3Regex.FindStringSubmatch(url)

		//LOD
		r.lod = 3

		//size
		switch groups[1] {
		case "small":
			r.size = SizeSmall
		case "medium":
			r.size = SizeMedium
		case "big":
			r.size = sizeBig
		default:
			return r, errors.New("Unvalid size:"+ groups[2])
		}

		//file
		r.file = groups[2]

		//format
		r.format = groups[3]
	}
	return r, nil
}

func rewriteUrlWithWebp(url string, header http.Header) string{
	strSplit := strings.Split(url, ".")
	ext := strSplit[len(strSplit)-1]
	if acceptsWebp(ext, header){
		strSplit[len(strSplit)-1] = "webp"
		return strings.Join(strSplit, ".")
	}else{
		strSplit[len(strSplit)-1] = "jpg"
		return strings.Join(strSplit, ".")
	}
}

func acceptsWebp(ext string, header http.Header) bool {
	suffixes := []string{"jpg", "jpeg", "JPG", "JPEG"}
	for _, s := range suffixes {
		if ext == s {
			for _, subElem := range strings.Split(header.Get("Accept"), ",") {
				if subElem == "image/webp" {
					return true
				}
			}
		}
	}
	return false
}
