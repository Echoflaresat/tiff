# `tiff` – Streamable TIFF Decoder for Go

This package provides a memory-efficient, standards-compatible TIFF decoder for Go, with support for **on-demand access** to striped and tiled images. It is a drop-in replacement for the standard `image/tiff` package, but optimized for large or remote TIFFs.

## ✨ Features

- ✅ Supports **TIFF Strips** and **Tiled TIFFs**
- ✅ Compatible with the standard `image.Image` interface
- ✅ Integrates with Go’s `image.RegisterFormat` (supports `image.Decode`)
- ✅ Uses `io.ReaderAt` for efficient random access
- ✅ Lazy-loading + LRU tile cache for tiled images
- ✅ Gracefully falls back to `golang.org/x/image/tiff` for unsupported formats

## 🔧 Supported TIFF Tags

| Tag                        | Support |
|---------------------------|---------|
| Compression: `None`, `Deflate` | ✅      |
| Photometric: `RGB`, `BlackIsZero` | ✅      |
| PlanarConfig: `Contig` only | ✅      |
| StripOffsets / StripByteCounts | ✅      |
| TileOffsets / TileByteCounts   | ✅      |

## 🚀 Usage

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
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	// img.Bounds(), img.At(x, y), etc.
}
```

### Optional: DecodeConfig

```go
cfg, err := image.DecodeConfig(f)
// cfg.Width, cfg.Height, cfg.ColorModel
```

## 🧠 Internals

- If the input is an `io.ReaderAt`, it uses direct random access.
- If the input is an `io.ReadSeeker`, it wraps it to support `ReadAt`.
- For striped/tiled TIFFs with supported layout, it avoids reading the full image.
- Uses [Hashicorp LRU](https://github.com/hashicorp/golang-lru) for caching decompressed tiles.

## 📦 Compatibility

The package conforms to Go’s standard `image` interface:

- `ColorModel()`
- `Bounds()`
- `At(x, y)`

It can be used anywhere `image.Image` is expected, including PNG or JPEG pipelines.

## 🛠 Future Roadmap

- [ ] Support for `PlanarConfig = Separate`
- [ ] Support `Predictor = Horizontal`
- [ ] Add support for more photometric interpretations (e.g., CMYK, YCbCr)

## 🧪 Tests & Real World Samples

This package is tested with both synthetic and real TIFFs, including high-res satellite and scientific imagery.

## License

MIT – see [LICENSE](./LICENSE)
