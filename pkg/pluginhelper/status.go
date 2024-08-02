package pluginhelper

import (
	"encoding/json"
	"fmt"
	"github.com/cloudnative-pg/cnpg-i/pkg/operator"
)

// BuildSetStatusResponse accepts as values:
// a) any struct containing values that can be converted into a valid JSON object. This will be set as the plugin
// status
// b) nil value. This is translated as a 'noop' (no changes to the plugin status will be done) in the CNPG main
// reconciliation loop
// c) any empty struct this will be translated as a cleanup of the currently registered plugin status.
// Examples:
// Lets suppose this function is invoked by a plugin named 'test'
// a) Body{Uptime: 100} will result into .status.plugins["test"] = '{"Uptime": 100}'
// b) will leave the status unchanged, so it will stay as the currently set value, example '{"Uptime": 100}'
// c) Body{} will result into .status.plugins["test"] = '{}'
func BuildSetStatusResponse(obj any) (*operator.SetClusterStatusResponse, error) {
	if obj == nil {
		return &operator.SetClusterStatusResponse{JsonStatus: nil}, nil
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
