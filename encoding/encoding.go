package encoding

import "encoding/json"

// UnmarshalInterface encodes an interface using json.Marshal
// and then parses the data with json.Unmarshal and stores the result
// in the value pointed to by v.
func UnmarshalInterface(data, v interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, v)
}
