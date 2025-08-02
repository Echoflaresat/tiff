// Package tifftag defines known TIFF tag identifiers used in image metadata.
// These tag constants correspond to the TIFF 6.0 specification and supplements,
// including common fields such as ImageWidth, Compression, and TileOffsets.
//
// For reference, see:
// https://www.loc.gov/preservation/digital/formats/content/tiff_tags.shtml
package tifftag

import "fmt"

// Tag represents a TIFF field tag number.
// These are used to identify metadata entries in TIFF image headers.
type Tag uint16

const (
	// ImageWidth specifies the number of columns (pixels) in the image.
	ImageWidth Tag = 256

	// ImageLength specifies the number of rows (pixels) in the image.
	ImageLength Tag = 257

	// BitsPerSample defines the number of bits per image component.
	BitsPerSample Tag = 258

	// Compression defines the compression scheme used on the image data.
	Compression Tag = 259

	// PhotometricInterpretation defines how pixel values should be interpreted.
	PhotometricInterpretation Tag = 262

	// StripOffsets contains the offsets to image data strips.
	StripOffsets Tag = 273

	// SamplesPerPixel defines the number of components per pixel.
	SamplesPerPixel Tag = 277

	// RowsPerStrip specifies how many rows are in each strip.
	RowsPerStrip Tag = 278

	// StripByteCounts contains the byte size of each strip.
	StripByteCounts Tag = 279

	// PlanarConfiguration specifies whether components are stored together or separately.
	PlanarConfiguration Tag = 284

	// TileWidth defines the width of a tile in pixels.
	TileWidth Tag = 322

	// TileLength defines the height of a tile in pixels.
	TileLength Tag = 323

	// TileOffsets contains the offsets to each tile.
	TileOffsets Tag = 324

	// TileByteCounts contains the byte size of each tile.
	TileByteCounts Tag = 325
)

// String returns a human-readable name for the TIFF tag.
// If the tag is unknown, it returns a formatted numeric identifier.
func (t Tag) String() string {
	switch t {
	case ImageWidth:
		return "ImageWidth"
	case ImageLength:
		return "ImageLength"
	case BitsPerSample:
		return "BitsPerSample"
	case Compression:
		return "Compression"
	case PhotometricInterpretation:
		return "PhotometricInterpretation"
	case StripOffsets:
		return "StripOffsets"
	case SamplesPerPixel:
		return "SamplesPerPixel"
	case RowsPerStrip:
		return "RowsPerStrip"
	case StripByteCounts:
		return "StripByteCounts"
	case PlanarConfiguration:
		return "PlanarConfiguration"
	case TileWidth:
		return "TileWidth"
	case TileLength:
		return "TileLength"
	case TileOffsets:
		return "TileOffsets"
	case TileByteCounts:
		return "TileByteCounts"
	default:
		return fmt.Sprintf("Tag(%d)", t)
	}
}
