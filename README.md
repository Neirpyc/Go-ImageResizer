# Go-ImageResizer

# Introduction
This program is meant to load a directory of JPEG images, 
and serve them as either JPEG or WEBP, with three different quality settings.
This is mostly mean to be used on a web server to deliver smaller images to 
mobile devices, and also quickly load light previews of the images. 

# Installation
To install, run:
```
got get github.com/chai2010/webp
go get github.com/nfnt/resize
go get github.com/BurntSushi/toml
git clone https://github.com/Neirpyc/Go-ImageResizer
cd Go-ImageResizer
go build *.go -o imageResizer
./imageResizer
```

# Configuration
Setting are provided in the file `settings.toml`.
Here is an example:
```toml
port="8000"
cacheFolder="/tmp/cache_imageResizer"
inputFolder="/var/www/images/"
maxCacheSize="500MB"
logFile="log.log"
preCache=true
```

- **port**: The local port on which the program will listen.
- **cacheFolder**: The folder in which served images will be cached.
- **inputFolder**: The folder in which served images will be fetched.
- **MaxCacheSiwe**: Should be in `b`, `B`, `Kb`, `KB`, `Mb`, `MB`, `Gb` or `GB`. It represents the maximal size the *cache* folder will use.
- **logFile**: Path of a file in which logs will be stored.
- **preCache**: should be *true* or *false*
   - *true* means the program will on launch rescale the all images in the input folder and pre-cache them.
   - *false* means nothing will be done on launch.

# Examples
Examples assume you use the previously shown config file.

The following request: `localhost:8000/images/lod/2/holidays/sunset.webp` 
will seek for `/var/www/images/holidays/sunset.jpg` and return it, converted to webp, 
and scaled according to the constants defined in `main.go:19`, 
which are by default: 250px width, encoding quality of 65, and height 
scaled appropriately to keep aspect ratio.
```go
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
```

If the file has already been served by the server, it is likely it'll be in the cache,
 and will be found and sent without rescaling the image another time.

Here are a few request examples:
- `localhost:8000/images/lod/1/holidays/sunset.webp` 
- `localhost:8000/images/lod/2/holidays/sunset.webp` 
- `localhost:8000/images/lod/2/holidays/sunset.jpg` 
- `localhost:8000/images/lod/3/small/holidays/sunset.webp` 
- `localhost:8000/images/lod/3/medium/holidays/sunset.jpg` 
- `localhost:8000/images/lod/3/big/holidays/sunset.jpg` 

# Demo
A live demo can be found [here](https://beggiora.neirpyc.ovh/master/gallery).

# license
This program is licensed under the GNU GPL v3 license. Please see the LICENSE file for details.
