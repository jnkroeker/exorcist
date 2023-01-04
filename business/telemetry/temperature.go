package telemetry

import (
	"encoding/binary"
	"errors"
	"math"
)

// Degrees Celsius
type Temperature struct {
	Temp float32
}

func (temp *Temperature) Parse(bytes []byte) error {
	if 4 != len(bytes) {
		return errors.New("Invalid length TMPC packet")
	}

	bits := binary.BigEndian.Uint32(bytes[0:4])

	temp.Temp = math.Float32frombits(bits)

	return nil
}
