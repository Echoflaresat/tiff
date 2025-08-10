# `tiff` – Random Access TIFF Decoder for Go

This package provides a memory-efficient, standards-compatible TIFF decoder for Go, with support for **on-demand access** to striped and tiled images with graceful fallback to `https://pkg.go.dev/golang.org/x/image/tiff` for unsupported formats.

For random-access decoding to work, the TIFF must conform to the following constraints:

| Feature        | Support 
|----------------|---------
| Compression    | `None`, `Deflate`   
| Photometric    | `RGB`, `BlackIsZero`    
| PlanarConfig   | `Contig`  only     

## Usage


```go
import (
	"image"
	"github.com/echoflaresat/tiff"
)

func main() {
	f, err := os.Open("large_image.tif")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // keep it open until you are done with the image

	
	img, _, err := tiff.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// img.Bounds(), img.At(x, y), etc.
}
```

## License

MIT – see [LICENSE](./LICENSE)
