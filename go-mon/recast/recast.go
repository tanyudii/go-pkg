package recast

import "encoding/json"

// Recast a to b. Not working for struct has field map
//
// Example:
//
//	type Data struct {
//		Name     string
//		Email    string
//		Property map[string]string
//	}
//
//	type Request struct {
//		Name     string
//		Email    string
//		Property map[string]string
//	}
//
// if some key in property in `Request` is not configured, it still presents in `Data`
func Recast(a, b interface{}) error {
	js, err := json.Marshal(a)
	if err != nil {
		return err
	}
	return json.Unmarshal(js, b)
}
