package storage

import (
	"encoding/binary"
	"fmt"
	"io"
	"path/filepath"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

func WriteLengthPrefixedProtobufMessage(writer io.Writer, message protoreflect.ProtoMessage) error {
	payload, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize key value: %w", err)
	}

	err = binary.Write(writer, binary.LittleEndian, uint32(len(payload)))
	if err != nil {
		return fmt.Errorf("failed to write payload len: %w", err)
	}

	_, err = writer.Write(payload)
	if err != nil {
		return fmt.Errorf("failed to write payload: %w", err)
	}

	return nil
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

func GetPath(dbDir, filename string) string {
	if dbDir == "" {
		panic("dbDir is empty")
	}
	return filepath.Join(dbDir, filename)
}
