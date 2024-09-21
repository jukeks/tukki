package marshalling

import (
	"encoding/binary"
	"fmt"
	"io"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func WriteLengthPrefixedProtobufMessage(writer io.Writer, message protoreflect.ProtoMessage) (uint32, error) {
	payload, err := proto.Marshal(message)
	if err != nil {
		return 0, fmt.Errorf("failed to serialize key value: %w", err)
	}

	payloadLen := uint32(len(payload))
	err = binary.Write(writer, binary.LittleEndian, uint32(len(payload)))
	if err != nil {
		return 0, fmt.Errorf("failed to write payload len: %w", err)
	}

	_, err = writer.Write(payload)
	if err != nil {
		return 0, fmt.Errorf("failed to write payload: %w", err)
	}

	return 4 + payloadLen, nil // 4 bytes for the length prefix
}

func ReadLengthPrefixedProtobufMessage(reader io.Reader, message protoreflect.ProtoMessage) error {
	var length uint32
	err := binary.Read(reader, binary.LittleEndian, &length)
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}

		return fmt.Errorf("failed to read payload len: %w", err)
	}

	payload := make([]byte, length)
	_, err = io.ReadFull(reader, payload)
	if err != nil {
		return fmt.Errorf("failed to read payload: %w", err)
	}

	err = proto.Unmarshal(payload, message)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return nil
}
