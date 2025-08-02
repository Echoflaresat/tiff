# `tiff` – Random Access TIFF Decoder for Go

This package provides a memory-efficient, standards-compatible TIFF decoder for Go, with support for **on-demand access** to striped and tiled images. 
It is a drop-in replacement for the standard `https://pkg.go.dev/golang.org/x/image/tiff` package, optimized for large TIFFs.

- Compatible with the standard `image.Image` interface
- Integrates with Go’s `image.RegisterFormat` (supports `image.Decode`)
- Gracefully falls back to `https://pkg.go.dev/golang.org/x/image/tiff` for unsupported formats.

**⚠️ Important**: When a TIFF is successfully parsed using this package’s `striped` or `tiled` backend, the decoded image lazily fetches pixel data. 
  Therefore, the original input must remain open and readable for the lifetime of the returned `image.Image`.

For random-access decoding to work, the TIFF must conform to the following constraints:

| Feature        | Support 
|----------------|---------
| Compression    | `None`, `Deflate`   
| Photometric    | `RGB`, `BlackIsZero`    
| PlanarConfig   | `Contig`  only     

## Usage

You can use this package exactly like any Go image decoder:

```go
import (
	"image"
	_ "github.com/echoflaresat/tiff"
)

func main() {
	f, err := os.Open("large_image.tif")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close() // keep it open until you are done with the image

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// img.Bounds(), img.At(x, y), etc.
}
```

## License

MIT – see [LICENSE](./LICENSE)
