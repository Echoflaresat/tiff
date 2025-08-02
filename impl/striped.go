// Package impl contains internal TIFF image decoding implementations.
// This file implements support for stripped TIFF images using on-demand access.
package impl

import (
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/echoflaresat/tiff/compression"
	"github.com/echoflaresat/tiff/photometric"
)

// stripedTiff represents a memory-efficient view of a TIFF image using strips.
//
// This implementation accesses pixel data lazily by reading only the necessary
// strip from the underlying io.ReaderAt when At(x, y) is called.
type stripedTiff struct {
	header TiffHeader
	reader io.ReaderAt
}

// LoadStripedTiff attempts to parse and load a TIFF image using a striped layout.
// It returns an image.Image implementation that lazily accesses pixel data as needed.
//
// Supported format constraints:
//   - Compression: None
//   - PhotometricInterpretation: RGB or BlackIsZero
//   - BitsPerSample: 8-bit per channel
//
// Note: The returned image.Image requires that the `reader` remains open for future reads.
func LoadStripedTiff(reader io.ReaderAt) (image.Image, error) {
	header, err := parseTiffHeader(reader)
	if err != nil {
		return nil, err
	}

	if header.Compression != compression.None {
		return nil, fmt.Errorf("unsupported compression: %d", header.Compression)
	}
	if header.Photometric != photometric.RGB && header.Photometric != photometric.BlackIsZero {
		return nil, fmt.Errorf("expected RGB or grayscale photometric interpretation, got %d", header.Photometric)
	}

	switch header.Photometric {
	case photometric.BlackIsZero:
		if header.SamplesPerPixel != 1 || header.BitsPerSample[0] != 8 {
			return nil, fmt.Errorf("unsupported grayscale format")
		}
	case photometric.RGB:
		if header.SamplesPerPixel != 3 || header.BitsPerSample[0] != 8 {
			return nil, fmt.Errorf("unsupported RGB format")
		}
	default:
		return nil, fmt.Errorf("unsupported photometric: %d", header.Photometric)
	}

	if len(header.StripOffsets) == 0 || len(header.StripOffsets) != len(header.StripByteCounts) {
		return nil, fmt.Errorf("invalid strip offset/length")
	}

	return &stripedTiff{header: header, reader: reader}, nil
}

// ColorModel returns the color model used by the TIFF image.
// Always returns color.RGBAModel for simplicity.
func (t *stripedTiff) ColorModel() color.Model {
	return color.RGBAModel
}

// Bounds returns the image rectangle.
func (t *stripedTiff) Bounds() image.Rectangle {
	return image.Rect(0, 0, t.header.Width, t.header.Height)
}

// At returns the color of the pixel at (x, y).
// This function reads the relevant bytes from the correct strip using t.reader.
func (t *stripedTiff) At(x, y int) color.Color {
	h := t.header

	strip := y / h.RowsPerStrip
	localY := y % h.RowsPerStrip
	bytesPerPixel := h.SamplesPerPixel
	idx := h.StripOffsets[strip] + (localY*h.Width+x)*bytesPerPixel

	switch h.Photometric {
	case photometric.RGB:
		var buf [3]byte
		_, err := t.reader.ReadAt(buf[:], int64(idx))
		if err != nil {
			panic(fmt.Sprintf("could not read RGB pixel at (%d,%d): %v", x, y, err))
		}
		return color.RGBA{R: buf[0], G: buf[1], B: buf[2], A: 255}

	case photometric.BlackIsZero:
		var b [1]byte
		_, err := t.reader.ReadAt(b[:], int64(idx))
		if err != nil {
			panic(fmt.Sprintf("could not read grayscale pixel at (%d,%d): %v", x, y, err))
		}
		return color.RGBA{R: b[0], G: b[0], B: b[0], A: 255}

	default:
		panic(fmt.Sprintf("unsupported PhotometricInterpretation: %d", h.Photometric))
	}
}
