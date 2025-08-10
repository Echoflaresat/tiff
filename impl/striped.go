// Package impl contains internal TIFF image decoding implementations.
// This file implements support for stripped TIFF images using on-demand access.
package impl

import (
	"fmt"
	"image"
	"image/color"
	"io"
	"sync"

	lru "github.com/hashicorp/golang-lru"

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
	cache  *lru.Cache // maps tileIndex -> []byte
	mutex  *sync.Mutex
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

	cache, err := lru.New(256)
	if err != nil {
		return nil, fmt.Errorf("could not create cache; %w", err)
	}

	return &stripedTiff{
		header: header,
		reader: reader,
		cache:  cache,
		mutex:  &sync.Mutex{},
	}, nil
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
	row := t.getRow(strip, localY, bytesPerPixel)

	base := x * bytesPerPixel

	switch h.Photometric {
	case photometric.RGB:
		return color.RGBA{R: row[base+0], G: row[base+1], B: row[base+2], A: 255}
	case photometric.BlackIsZero:
		v := row[base]
		return color.RGBA{R: v, G: v, B: v, A: 255}
	default:
		panic(fmt.Sprintf("unsupported PhotometricInterpretation: %d", h.Photometric))
	}
}

// getRow returns a full row of raw bytes for (strip, rowInStrip).
// Fast path: no lock on reader; RLock+Get on cache.
// On miss: Lock, double-check, then single-threaded ReadAt and cache.
func (t *stripedTiff) getRow(strip, rowInStrip, bpp int) []byte {
	key := (uint64(strip) << 32) | uint64(uint32(rowInStrip))

	// Try cache under read lock.
	if row, ok := t.cache.Get(key); ok {
		return row.([]byte)
	}

	h := t.header
	rowSize := h.Width * bpp
	offset := int64(h.StripOffsets[strip] + (rowInStrip*h.Width)*bpp)

	row := make([]byte, rowSize)
	t.mutex.Lock()
	n, err := t.reader.ReadAt(row, offset)
	defer t.mutex.Unlock()

	if err != nil || n != len(row) {
		panic(fmt.Sprintf("could not read row strip=%d row=%d: read %d/%d bytes, err=%v",
			strip, rowInStrip, n, len(row), err))
	}

	t.cache.Add(key, row)
	return row
}
