package go_validator

import (
	"errors"
	goerr "pkg.tanyudii.me/go-pkg/go-err"
	"testing"
)

func TestResultNestedValidate(t *testing.T) {
	type testInfo struct {
		Email        string `validate:"required"`
		PhoneNumber  string `validate:"required" label:"Phone number"`
		ExampleField string `validate:"required" label:"Example Field" field:"example_field"`
		AvatarURL    string `validate:"required" label:"Avatar URL" field:"avatar_url"`
	}

	type testItem struct {
		Link            string `validate:"required"`
		DescriptionData string `validate:"required" field:"description_data"`
	}

	type testUserRequest struct {
		Name        string      `validate:"required"`
		SomeLink    string      `validate:"required" label:"SomeLink"`
		AvatarURL   string      `validate:"required" label:"Avatar URL"`
		ClientID    string      `validate:"required" label:"Client ID" field:"client_id"`
		ClientName  string      `validate:"required" label:"Client Name" field:"client_name"`
		FirstInfo   *testInfo   `validate:"required" label:"First info"`
		SecondInfo  testInfo    `validate:"required" label:"Second info"`
		FirstItems  []*testItem `validate:"required,dive"`
		SecondItems []testItem  `validate:"required,dive"`
	}

	testCases := []struct {
		name           string
		req            testUserRequest
		expectedFields goerr.ErrorField
	}{
		{
			name: "all field errors",
			req:  testUserRequest{},
			expectedFields: goerr.ErrorField{
				"name":                     "Name is a required field",
				"someLink":                 "SomeLink is a required field",
				"avatarUrl":                "Avatar URL is a required field",
				"client_id":                "Client ID is a required field",
				"client_name":              "Client Name is a required field",
				"firstInfo":                "First info is a required field",
				"firstItems":               "FirstItems is a required field",
				"secondItems":              "SecondItems is a required field",
				"secondInfo.email":         "Email is a required field",
				"secondInfo.phoneNumber":   "Phone number is a required field",
				"secondInfo.example_field": "Example Field is a required field",
				"secondInfo.avatar_url":    "Avatar URL is a required field",
			},
		},
		{
			name: "nested field errors",
			req: testUserRequest{
				Name:        "Name",
				SomeLink:    "Some Link",
				AvatarURL:   "Avatar URL",
				ClientID:    "Client ID",
				ClientName:  "Client Name",
				FirstInfo:   &testInfo{},
				SecondInfo:  testInfo{},
				FirstItems:  []*testItem{{}, {}},
				SecondItems: []testItem{{}, {}},
			},
			expectedFields: goerr.ErrorField{
				"firstInfo.email":                "Email is a required field",
				"firstInfo.phoneNumber":          "Phone number is a required field",
				"firstInfo.example_field":        "Example Field is a required field",
				"firstInfo.avatar_url":           "Avatar URL is a required field",
				"secondInfo.email":               "Email is a required field",
				"secondInfo.phoneNumber":         "Phone number is a required field",
				"secondInfo.example_field":       "Example Field is a required field",
				"secondInfo.avatar_url":          "Avatar URL is a required field",
				"firstItems.0.link":              "Link is a required field",
				"firstItems.0.description_data":  "DescriptionData is a required field",
				"firstItems.1.link":              "Link is a required field",
				"firstItems.1.description_data":  "DescriptionData is a required field",
				"secondItems.0.link":             "Link is a required field",
				"secondItems.0.description_data": "DescriptionData is a required field",
				"secondItems.1.link":             "Link is a required field",
				"secondItems.1.description_data": "DescriptionData is a required field",
			},
		},
	}

	v := NewValidator()

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Struct(tt.req)
			if err == nil {
				t.Errorf("Error should not be nil")
				return
			}

			var badReqErr *goerr.BadRequestError
			if ok := errors.As(err, &badReqErr); !ok {
				t.Errorf("Error should be BadRequestError, got %v", err)
				return
			}

			errFields := badReqErr.GetFields()
			for key, val := range tt.expectedFields {
				if _, ok := errFields[key]; !ok {
					t.Errorf("Field %s should be exist", key)
					continue
				}
				if errFields[key] != val {
					t.Errorf("Field %s should be '%s', got '%s'", key, val, errFields[key])
					continue
				}
			}

			for key := range errFields {
				if _, ok := tt.expectedFields[key]; !ok {
					t.Errorf("Field %s should not be exist", key)
					continue
				}
			}
		})
	}
}
