package hexifier_test

import (
	hexifier "go-test-examples/no-interfaces-needed"
	"reflect"
	"sort"
	"testing"
)

func TestLatLngToCell(t *testing.T) {
	type args struct {
		lat        float64
		lng        float64
		resolution int
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "Returns correct cell id for given lat, long, at resolution 5",
			args: args{
				lat:        40.50818200904383,
				lng:        -111.9784200763932,
				resolution: 5,
			},
			want: 599657577537601535,
		},
		{
			name: "Returns correct cell id for given lat, long, at resolution 8",
			args: args{
				lat:        40.50818200904383,
				lng:        -111.9784200763932,
				resolution: 8,
			},
			want: 613168375937368063,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hexifier.LatLngToCell(tt.args.lat, tt.args.lng, tt.args.resolution); got != tt.want {
				t.Errorf("LatLngToCell() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSurroundingCells(t *testing.T) {
	type args struct {
		cellID int64
		dist   int
	}
	tests := []struct {
		name string
		args args
		want []int64
	}{
		{
			name: "Returns correct cell ids for the given ring number",
			args: args{
				cellID: 613168375937368063,
				dist:   2,
			},
			want: []int64{
				// ring 0
				613168375937368063,
				// ring 1
				613168375945756671,
				613168376132403199,
				613168376124014591,
				613168375440343039,
				613168375941562367,
				613168375935270911,
				// ring 2
				613168375943659519,
				613168375905910783,
				613168375908007935,
				613168376128208895,
				613168376119820287,
				613168376126111743,
				613168375442440191,
				613168375431954431,
				613168375444537343,
				613168375425662975,
				613168375939465215,
				613168375947853823,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := hexifier.SurroundingCells(tt.args.cellID, tt.args.dist)
			sort.Slice(tt.want, func(i, j int) bool { return tt.want[i] < tt.want[j] })
			sort.Slice(got, func(i, j int) bool { return got[i] < got[j] })
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SurroundingCells() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMilesToCellRing(t *testing.T) {
	type args struct {
		miles      int
		resolution int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Returns the correct size of ring when miles is smaller than diameter",
			args: args{
				miles:      2,
				resolution: 2,
			},
			want: 0,
		},
		{
			name: "Returns the correct size of ring when miles is greate than diameter",
			args: args{
				miles:      2,
				resolution: 8,
			},
			want: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hexifier.MilesToCellRing(tt.args.miles, tt.args.resolution); got != tt.want {
				t.Errorf("MilesToCellRing() = %v, want %v", got, tt.want)
			}
		})
	}
}
