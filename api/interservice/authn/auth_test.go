package authn_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/chef/automate/api/interservice/authn"
	"github.com/chef/automate/lib/grpc/grpctest"
)

func TestGeneratedProtobufUpToDate(t *testing.T) {
	grpctest.AssertCompiledInUpToDate(t, "api/interservice/authn/")
}

func TestValidationCreateTokenID(t *testing.T) {
	negativeCases := map[string]*authn.CreateTokenReq{
		"with uppercase characters": {
			Id:       "TestID",
			Projects: []string{"project1", "project2"},
		},
		"with non-dash characters": {
			Id:       "test_underscore",
			Projects: []string{"project1", "project2"},
		},
		"with spaces": {
			Id:       "test space",
			Projects: []string{"project1", "project2"},
		},
		"whitespace projects list": {
			Id:       "valid-id",
			Projects: []string{"     ", "test"},
		},
		"repeated projects in list": {
			Id:       "valid-id",
			Projects: []string{"repeat", "repeat"},
		},
		"project has invalid characters": {
			Id:       "valid-id",
			Projects: []string{"valid", "wrong~"},
		},
		"project has spaces": {
			Id:       "valid-id",
			Projects: []string{"valid", "wrong space"},
		},
		"project is too long": {
			Id:       "valid-id",
			Projects: []string{"much-too-long-longest-word-in-english-pneumonoultramicroscopicsilicovolcanoconiosis", "valid"},
		},
	}
	positiveCases := map[string]*authn.CreateTokenReq{
		"projects are missing": {
			Id: "valid-id",
		},
		"with no characters": {
			Id:       "",
			Projects: []string{"project1", "project2"},
		},
		"with ID all lowercase": {
			Id:       "test",
			Projects: []string{"project1", "project2"},
		},
		"with ID that has dashes": {
			Id:       "-test-with-dashes-",
			Projects: []string{"project1", "project2"},
		},
		"with ID that has dashes and numbers": {
			Id:       "1-test-with-1-and-dashes-0",
			Projects: []string{"project1", "project2"},
		},
		"with ID that has only numbers": {
			Id:       "1235",
			Projects: []string{"project1", "project2"},
		},
		"with a single project": {
			Id:       "1235",
			Projects: []string{"project1"},
		},
		"projects are empty": {
			Id:       "valid-id",
			Projects: []string{},
		},
	}

	classes := map[bool]map[string]*authn.CreateTokenReq{
		true:  positiveCases,
		false: negativeCases,
	}

	for expectedSuccess, cases := range classes {
		for name, tc := range cases {
			t.Run(name, func(t *testing.T) {
				err := tc.Validate()
				if expectedSuccess {
					assert.NoError(t, err)
				} else {
					assert.Error(t, err)
				}
			})
		}
	}
}
