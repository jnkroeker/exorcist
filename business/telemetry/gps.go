package telemetry

import (
	"encoding/binary"
	"errors"
)

type GPS struct {
	Latitute  float64 `json:"lat"`    // degrees latitude
	Longitude float64 `json:"lon"`    // degrees longitude
	Altitude  float64 `json:"alt"`    // meters above wgs84 ellipsoid?
	Speed     float64 `json:"spd"`    // m/s
	Speed3D   float64 `json:"spd_3d"` // m/s, standard error?
	Timestamp int64   `json:"utc"`
}

func (gps *GPS) Parse(bytes []byte, scale *Scale) error {
	if 20 != len(bytes) {
		return errors.New("invalid length GPS5 packet")
	}

	gps.Latitute = float64(int32(binary.BigEndian.Uint32(bytes[0:4]))) / float64(scale.Values[0])
	gps.Longitude = float64(int32(binary.BigEndian.Uint32(bytes[4:8]))) / float64(scale.Values[1])

	// convert from mm
	gps.Altitude = float64(int32(binary.BigEndian.Uint32(bytes[8:12]))) / float64(scale.Values[2])

	// convert from mm/s
	gps.Speed = float64(int32(binary.BigEndian.Uint32(bytes[12:16]))) / float64(scale.Values[3])

	// convert from mm/s
	gps.Speed3D = float64(int32(binary.BigEndian.Uint32(bytes[16:20]))) / float64(scale.Values[4])

	return nil
}
