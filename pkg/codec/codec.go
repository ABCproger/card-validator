package codec

import (
	"encoding/json"
	"fmt"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// JSON is a Connect codec that always emits all fields, including proto3 zero values.
// This ensures `valid: false` is present in error responses, matching the API contract.
type JSON struct{}

func (JSON) Name() string { return "json" }

func (JSON) Marshal(message any) ([]byte, error) {
	if msg, ok := message.(proto.Message); ok {
		return protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseProtoNames:   true,
		}.Marshal(msg)
	}
	return json.Marshal(message)
}

func (JSON) Unmarshal(data []byte, message any) error {
	if msg, ok := message.(proto.Message); ok {
		return protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(data, msg)
	}
	if err := json.Unmarshal(data, message); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	return nil
}
