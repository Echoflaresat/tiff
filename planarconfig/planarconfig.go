// Package planarconfig defines the TIFF PlanarConfiguration tag values,
// which specify how image samples (e.g., R, G, B channels) are stored within the file.
//
// This corresponds to TIFF tag 284:
// https://www.awaresystems.be/imaging/tiff/tifftags/planarconfiguration.html
package planarconfig

import "fmt"

// Type represents the TIFF PlanarConfiguration field (tag 284).
// It indicates whether image samples are stored in chunky or planar format.
type Type int

const (
	// Unknown indicates an unrecognized or missing planar configuration.
	Unknown Type = -1

	// Contig (1) means samples are stored interleaved per pixel (e.g., RGB RGB RGB).
	Contig Type = 1

	// Separate (2) means each sample is stored in a separate plane
	// (e.g., RRR... GGG... BBB...).
	Separate Type = 2
)

// String returns a human-readable name for the planar configuration type.
// If the value is unknown, it returns a formatted fallback string.
func (c Type) String() string {
	switch c {
	case Unknown:
		return "Unknown"
	case Contig:
		return "Contig"
	case Separate:
		return "Separate"
	default:
		return fmt.Sprintf("Config(%d)", int(c))
	}
}
