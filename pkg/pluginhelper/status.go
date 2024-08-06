package pluginhelper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cloudnative-pg/cnpg-i/pkg/operator"
)

// SetStatusResponseBuilder a SetStatus response builder
type SetStatusResponseBuilder struct{}

// NewSetStatusResponseBuilder is an helper that creates the SetStatus endpoint responses
func NewSetStatusResponseBuilder() *SetStatusResponseBuilder {
	return &SetStatusResponseBuilder{}
}

// NoOpResponse this response will ensure that no changes will be done to the plugin status
func (s SetStatusResponseBuilder) NoOpResponse() *operator.SetClusterStatusResponse {
	return &operator.SetClusterStatusResponse{JsonStatus: nil}
}

// SetEmptyStatusResponse will set the plugin status to an empty object '{}'
func (s SetStatusResponseBuilder) SetEmptyStatusResponse() *operator.SetClusterStatusResponse {
	b, _ := json.Marshal(map[string]string{})
	return &operator.SetClusterStatusResponse{JsonStatus: b}
}

// JsonStatusResponse requires a struct or map that can be translated to a JSON object, will set the status to the passed
// object
func (s SetStatusResponseBuilder) JsonStatusResponse(obj any) (*operator.SetClusterStatusResponse, error) {
	if obj == nil {
		return nil, errors.New("nil object passed, use NoOpResponse")
	}
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	var js map[string]interface{}
	if err := json.Unmarshal(b, &js); err != nil {
		return nil, fmt.Errorf("invalid json: not an object: '%s'", string(b))
	}
	return &operator.SetClusterStatusResponse{
		JsonStatus: b,
	}, nil
}
