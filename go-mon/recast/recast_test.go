package recast_test

import (
	"pkg.tanyudii.me/go-pkg/go-mon/pointer"
	"pkg.tanyudii.me/go-pkg/go-mon/recast"
	"testing"
)

type Data struct {
	Name     string
	Email    string
	Property *map[string]string
}

type Request struct {
	Name     string
	Email    string
	Property map[string]string
}

func TestRecast(t *testing.T) {
	testCases := map[string]struct {
		data Data
		req  Request
	}{
		"test update field & new attribute in map": {
			data: Data{
				Name:  "John Doe",
				Email: "john.doe@gmail.com",
				Property: pointer.Val(map[string]string{
					"FirstValue":  "First",
					"SecondValue": "Second",
				}),
			},
			req: Request{
				Name:  "John Doe",
				Email: "john.doe+new@gmail.com",
				Property: map[string]string{
					"FirstValue":  "First",
					"SecondValue": "Second",
					"ThirdValue":  "Third",
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			_ = recast.Recast(&tc.req, &tc.data)

			if tc.req.Name != tc.data.Name {
				t.Errorf("expected name %s got %s", tc.req.Name, tc.data.Name)
			}

			if tc.req.Email != tc.data.Email {
				t.Errorf("expected email %s got %s", tc.req.Email, tc.data.Email)
			}

			if len(tc.req.Property) != len(pointer.Extract(tc.data.Property)) {
				t.Errorf("expected property length %d got %d", len(tc.req.Property), len(pointer.Extract(tc.data.Property)))
			}

			for k, v := range tc.req.Property {
				if pointer.Extract(tc.data.Property)[k] != v {
					t.Errorf("expected property %s value %s got %s", k, v, pointer.Extract(tc.data.Property)[k])
				}
			}
		})
	}
}
