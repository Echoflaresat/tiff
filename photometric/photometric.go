// Package photometric defines TIFF PhotometricInterpretation tag values.
// These specify how pixel values are to be interpreted and displayed.
//
// This corresponds to TIFF tag 262.
// Reference: https://www.awaresystems.be/imaging/tiff/tifftags/photometricinterpretation.html
package photometric

import "fmt"

// Interpretation represents a TIFF PhotometricInterpretation value.
// It defines the color space and pixel layout of the image.
type Interpretation int

const (
	// Unknown indicates an undefined or unrecognized photometric interpretation.
	Unknown Interpretation = -1

	// WhiteIsZero (0) means 0 is interpreted as white, and maximum value is black.
	WhiteIsZero Interpretation = 0

	// BlackIsZero (1) means 0 is black, and maximum value is white.
	BlackIsZero Interpretation = 1

	// RGB (2) means image pixels are stored as RGB triplets.
	RGB Interpretation = 2

	// Paletted (3) means pixel values are indexes into a color lookup table.
	Paletted Interpretation = 3

	// TransMask (4) is a transparency mask—black pixels are transparent.
	TransMask Interpretation = 4

	// CMYK (5) means image is stored in the CMYK color model.
	CMYK Interpretation = 5

	// YCbCr (6) means image uses YCbCr color encoding (common in JPEGs).
	YCbCr Interpretation = 6

	// CIELab (8) means image uses the CIE L*a*b* color space.
	CIELab Interpretation = 8
)

// String returns the symbolic name of the photometric interpretation.
func (p Interpretation) String() string {
	switch p {
	case WhiteIsZero:
		return "WhiteIsZero"
	case BlackIsZero:
		return "BlackIsZero"
	case RGB:
		return "RGB"
	case Paletted:
		return "Paletted"
	case TransMask:
		return "TransMask"
	case CMYK:
		return "CMYK"
	case YCbCr:
		return "YCbCr"
	case CIELab:
		return "CIELab"
	case Unknown:
		return "Unknown"
	default:
		return fmt.Sprintf("Interpretation(%d)", int(p))
	}
}
