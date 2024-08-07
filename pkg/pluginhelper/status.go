package pluginhelper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloudnative-pg/cnpg-i/pkg/operator"
)

// ErrNilObject is used when a nill object is passed to the builder.
var ErrNilObject = errors.New("nil object passed, use NoOpResponse")

// NotAnObjectError is used when the passed value cannot be represented
// as a JSON object.
type NotAnObjectError struct {
	representation []byte
}

func (err NotAnObjectError) Error() string {
	return fmt.Sprintf(
		"the passed variable cannot be serialized as a JSON object: %s",
		err.representation,
	)
}

// SetStatusResponseBuilder a SetStatus response builder.
type SetStatusResponseBuilder struct{}

// NewSetStatusResponseBuilder is an helper that creates the SetStatus endpoint responses.
func NewSetStatusResponseBuilder() *SetStatusResponseBuilder {
	return &SetStatusResponseBuilder{}
}

// NoOpResponse this response will ensure that no changes will be done to the plugin status.
func (s SetStatusResponseBuilder) NoOpResponse() *operator.SetClusterStatusResponse {
	return &operator.SetClusterStatusResponse{JsonStatus: nil}
}

// SetEmptyStatusResponse will set the plugin status to an empty object '{}'.
func (s SetStatusResponseBuilder) SetEmptyStatusResponse() *operator.SetClusterStatusResponse {
	b, err := json.Marshal(map[string]string{})
	if err != nil {
		panic("JSON mashaller failed for empty map")
	}

	return &operator.SetClusterStatusResponse{JsonStatus: b}
}

// JSONStatusResponse requires a struct or map that can be translated to a JSON object,
// will set the status to the passed object.
func (s SetStatusResponseBuilder) JSONStatusResponse(obj any) (*operator.SetClusterStatusResponse, error) {
	if obj == nil {
		return nil, ErrNilObject
	}

	jsonObject, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("while marshalling resource definition: %w", err)
	}

	var js map[string]interface{}
	if err := json.Unmarshal(jsonObject, &js); err != nil {
		return nil, NotAnObjectError{representation: jsonObject}
	}

	return &operator.SetClusterStatusResponse{
		JsonStatus: jsonObject,
	}, nil
}
