package telemetry

import (
	"encoding/binary"
	"errors"
)

type GPSPrecision struct {
	Accuracy uint16
}

func (gpsp *GPSPrecision) Parse(bytes []byte) error {
	if 2 != len(bytes) {
		return errors.New("invalid length GPSP packet")
	}

	gpsp.Accuracy = binary.BigEndian.Uint16(bytes[0:2])

	return nil
}
