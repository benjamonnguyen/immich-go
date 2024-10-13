package metadata

import (
	"reflect"
	"testing"
	"time"

	"github.com/simulot/immich-go/internal/tzone"
)

func TestExifTool_ReadMetaData(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		want     *Metadata
		wantErr  bool
	}{
		{
			name:     "read JPG",
			fileName: "DATA/PXL_20231006_063000139.jpg",
			want: &Metadata{
				DateTaken: time.Date(2023, 10, 6, 8, 30, 0, 0, time.Local),
				Latitude:  +48.8583736,
				Longitude: +2.2919010,
			},
			wantErr: false,
		},
		{
			name:     "read mp4",
			fileName: "DATA/PXL_20220724_210650210.NIGHT.mp4",
			want: &Metadata{
				DateTaken: time.Date(2022, 7, 24, 21, 10, 56, 0, time.Local),
				Latitude:  47.538300,
				Longitude: -2.891900,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			et, err := NewExifTool(&ExifToolFlags{
				Timezone: tzone.Timezone{TZ: time.Local},
			})
			if err != nil {
				t.Error(err)
				return
			}
			got, err := et.ReadMetaData(tt.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExifTool.ReadMetaData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExifTool.ReadMetaData() = %v, want %v", got, tt.want)
			}
		})
	}
}
