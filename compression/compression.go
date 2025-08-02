// Package compression defines constants for TIFF compression methods
// as used in TIFF tag 259 (Compression). These values map directly to
// standard TIFF specification compression types.
//
// This package is used to interpret the Compression tag in TIFF image headers.
// It supports both modern and legacy methods like Deflate, PackBits, and Fax encodings.
package compression

import "fmt"

// Type represents a TIFF compression method as defined by the TIFF spec.
type Type int

const (
	// Unknown represents an undefined or unsupported compression type.
	Unknown Type = -1

	// None means no compression.
	None Type = 1

	// CCITT is CCITT 1D Group 3 compression.
	CCITT Type = 2

	// G3 is Group 3 Fax compression (2D encoding).
	G3 Type = 3

	// G4 is Group 4 Fax compression.
	G4 Type = 4

	// LZW is Lempel-Ziv-Welch compression.
	LZW Type = 5

	// JPEGOld is the original JPEG compression (deprecated).
	JPEGOld Type = 6

	// JPEG is modern JPEG compression.
	JPEG Type = 7

	// Deflate is zlib-style Deflate compression (RFC 1951).
	Deflate Type = 8

	// PackBits is a simple byte-oriented run-length encoding.
	PackBits Type = 32773

	// DeflateOld is an older value used for Deflate, superseded by Deflate.
	DeflateOld Type = 32946
)

// String returns a readable name for the compression type.
func (c Type) String() string {
	switch c {
	case Unknown:
		return "Unknown"
	case None:
		return "None"
	case CCITT:
		return "CCITT"
	case G3:
		return "G3Fax"
	case G4:
		return "G4Fax"
	case LZW:
		return "LZW"
	case JPEGOld:
		return "JPEGOld"
	case JPEG:
		return "JPEG"
	case Deflate:
		return "Deflate"
	case PackBits:
		return "PackBits"
	case DeflateOld:
		return "DeflateOld"
	default:
		return fmt.Sprintf("CompressionType(%d)", int(c))
	}
}
