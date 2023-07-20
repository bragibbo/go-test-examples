package hexifier

import (
	"fmt"
	"math"

	"github.com/uber/h3-go/v4"
)

var cellDiameterMiles = map[string]float64{
	"H0":  1588.75745,
	"H1":  598.99048,
	"H2":  226.31607,
	"H3":  85.53424,
	"H4":  32.32898,
	"H5":  12.21907,
	"H6":  4.61842,
	"H7":  1.74403,
	"H8":  0.65895,
	"H9":  0.24897,
	"H10": 0.09407,
	"H11": 0.03554,
	"H12": 0.01343,
	"H13": 0.00507,
	"H14": 0.00192,
	"H15": 0.00072,
}

// LatLngToCell converts a latitude, longitudem, and a given resolution
// to the h3 cell equivalent. It returns the cell id as an int64 value.
func LatLngToCell(lat, lng float64, resolution int) int64 {
	latLng := h3.NewLatLng(lat, lng)
	cell := h3.LatLngToCell(latLng, resolution)
	return int64(cell)
}

// SurroundingCells returns the ids for the cells surrounding the origin
// cell. Given the origin cell id and a distance (number of rings) from the
// origin, a slice of ids is returned. The lookup respects the original
// resolution of the origin cell. Includes the origin cell in the result.
func SurroundingCells(cellID int64, dist int) []int64 {
	var origin = h3.Cell(cellID)
	cells := h3.GridDisk(origin, dist)
	cellIDs := make([]int64, len(cells))
	for i, cell := range cells {
		cellIDs[i] = int64(cell)
	}
	return cellIDs
}

// MilesToCellRing returns the disk ring number at a given resolution for a number
// of miles. This is a proxy for a search radius. Using this disk ring number, the
// cells that fall within a given radius can be found. It is important to note that
// these numbers are averages since cell diameters vary based on the cell's proximity
// to a pentagon.
//
// Returns 0 if the resolution is invalid.
func MilesToCellRing(miles int, resolution int) int {
	key := fmt.Sprintf("H%v", resolution)
	dia, exists := cellDiameterMiles[key]
	if !exists {
		return 0
	}
	// handle radius of origin cell to not overstate ring size
	distRemaining := float64(miles) - (dia / 2)
	return int(math.Round(distRemaining / dia))
}