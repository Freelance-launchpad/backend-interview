package registry

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

const magicByte = 0x00

// Encode prefixes the JSON payload with the Confluent wire format header:
//
//	[0x00][int32 schemaID big-endian][json payload]
func Encode(schemaID int32, payload []byte) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := buf.WriteByte(magicByte); err != nil {
		return nil, err
	}
	if err := binary.Write(buf, binary.BigEndian, schemaID); err != nil {
		return nil, err
	}
	if _, err := buf.Write(payload); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func Decode(payload []byte) (int32, []byte, error) {
	const headerSize = 5 // 1 magic byte + 4 bytes schema ID

	if len(payload) < headerSize {
		return 0, nil, fmt.Errorf("invalid payload: too short")
	}
	if payload[0] != magicByte {
		return 0, nil, fmt.Errorf("invalid magic byte: 0x%02x", payload[0])
	}

	schemaID := int32(binary.BigEndian.Uint32(payload[1:headerSize]))
	return schemaID, payload[headerSize:], nil
}
