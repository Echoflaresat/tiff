// Package tiff provides a memory-efficient, standards-compatible TIFF decoder for Go,
// with support for on-demand access to striped and tiled TIFF images.
//
// When working with supported TIFF formats, this decoder avoids loading the entire image
// into memory. Instead, it reads only the required pixel data on demand using io.ReaderAt.
//
//	⚠️ When a supported format is detected, the caller must keep the underlying reader
//	(typically a file) open for as long as the image.Image is in use.
//
// If the format is unsupported, the decoder gracefully falls back to golang.org/x/image/tiff,
// in which case the full image is decoded eagerly and no special reader lifetime is required.
//
// Supported features in random access mode:
//
//   - Striped and Tiled TIFF decoding
//   - Compression: None, Deflate (zlib)
//   - Photometric: RGB, BlackIsZero (grayscale)
//   - PlanarConfig: Contig (interleaved samples only)
//
// Example usage:
//
//	import (
//	    "image"
//	    _ "github.com/echoflaresat/tiff"
//	)
//
//	func main() {
//	    f, _ := os.Open("image.tif") // Must remain open when using the image
//	    defer f.Close()
//
//	    img, _, err := image.Decode(f)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    // Use img.At(x, y), img.Bounds(), etc.
//	}
//
// For full details and source, visit: https://pkg.go.dev/github.com/echoflaresat/tiff
package tiff

import (
	"image"
	"io"

	"github.com/echoflaresat/tiff/impl"
	stdtiff "golang.org/x/image/tiff"
)

const (
	// littleEndianHeader is the TIFF header for little-endian byte order.
	littleEndianHeader = "II\x2A\x00"
	// bigEndianHeader is the TIFF header for big-endian byte order.
	bigEndianHeader = "MM\x00\x2A"
)

// DecodeConfig returns the color model and dimensions of a TIFF image without decoding the entire image.
// It uses the standard library's TIFF decoder for configuration extraction.
func DecodeConfig(r io.Reader) (image.Config, error) {
	return stdtiff.DecodeConfig(r)
}

// Decode reads a TIFF image from r and returns it as an image.Image.
// It first attempts to decode using custom striped and tiled TIFF loaders,
// falling back to the standard library's TIFF decoder if those fail.
func Decode(r io.Reader) (image.Image, error) {
	var readerAt io.ReaderAt

	if ra, ok := r.(io.ReaderAt); ok {
		readerAt = ra
	} else if rs, ok := r.(io.ReadSeeker); ok {
		readerAt = &readerAtFromSeeker{rs: rs}
	}

	if readerAt != nil {
		if img, err := impl.LoadStripedTiff(readerAt); err == nil {
			return img, nil
		}
		if img, err := impl.LoadTiledTiff(readerAt); err == nil {
			return img, nil
		}
	}

	// Fallback to standard decoder
	return stdtiff.Decode(r)
}

// readerAtFromSeeker adapts an io.ReadSeeker to io.ReaderAt.
type readerAtFromSeeker struct {
	rs io.ReadSeeker
}

// ReadAt implements the io.ReaderAt interface for readerAtFromSeeker.
// It seeks to the specified offset and reads into p.
func (r *readerAtFromSeeker) ReadAt(p []byte, off int64) (int, error) {
	if _, err := r.rs.Seek(off, io.SeekStart); err != nil {
		return 0, err
	}
	return r.rs.Read(p)
}

// init registers the TIFF format with the image package, supporting both little-endian and big-endian headers.
func init() {
	image.RegisterFormat("tiff", littleEndianHeader, Decode, DecodeConfig)
	image.RegisterFormat("tiff", bigEndianHeader, Decode, DecodeConfig)
}
