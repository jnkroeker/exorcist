package telemetry

import (
	"encoding/binary"
	"errors"
)

type GPSFix struct {
	F uint32
}

func (gpsf *GPSFix) Parse(bytes []byte) error {
	if 4 != len(bytes) {
		return errors.New("invalid length GPSF packet")
	}

	gpsf.F = binary.BigEndian.Uint32(bytes[0:4])

	return nil
}
