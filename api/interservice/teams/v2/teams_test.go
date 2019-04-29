package v2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	v2 "github.com/chef/automate/api/interservice/teams/v2"
)

func TestGetTeamReq(t *testing.T) {
	negativeCases := map[string]*v2.GetTeamReq{
		"empty ID": {
			Id: "",
		},
		"whitespace ID": {
			Id: "      ",
		},
		"missing ID": {},
	}
	positiveCases := map[string]*v2.GetTeamReq{
		"with spaces until GA": {
			Id: "1111valid ID for now~~~",
		},
		"alphanumeric-post-GA": {
			Id: "asdf-123-fun",
		},
	}

	classes := map[bool]map[string]*v2.GetTeamReq{
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

func TestDeleteTeamReq(t *testing.T) {
	negativeCases := map[string]*v2.DeleteTeamReq{
		"empty ID": {
			Id: "",
		},
		"whitespace ID": {
			Id: "       ",
		},
		"missing ID": {},
	}
	positiveCases := map[string]*v2.DeleteTeamReq{
		"with spaces until GA": {
			Id: "1111valid ID for now~~~",
		},
		"alphanumeric-post-GA": {
			Id: "asdf-123-fun",
		},
	}

	classes := map[bool]map[string]*v2.DeleteTeamReq{
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

func TestUpdateTeamReq(t *testing.T) {
	negativeCases := map[string]*v2.UpdateTeamReq{
		"empty ID": {
			Id:       "",
			Name:     "name of my team",
			Projects: []string{"test"},
		},
		"whitespace ID": {
			Id:       "      ",
			Name:     "name of my team",
			Projects: []string{"test"},
		},
		"missing ID": {
			Name:     "name of my team",
			Projects: []string{"test"},
		},
		"empty Name": {
			Id:       "test",
			Name:     "",
			Projects: []string{"test"},
		},
		"whitespace Name": {
			Id:       "test",
			Name:     "          ",
			Projects: []string{"test"},
		},
		"missing Name": {
			Id:       "test",
			Projects: []string{"test"},
		},
		"missing Name, ID, and Projects": {},
		"whitespace projects list": {
			Id:       "this-is-valid-1",
			Name:     "name of my team ~ fun characters 1 %",
			Projects: []string{"     ", "test"},
		},
		"repeated projects in list": {
			Id:       "this-is-valid-1",
			Name:     "name of my team ~ fun characters 1 %",
			Projects: []string{"repeat", "repeat"},
		},
		"Project has invalid characters": {
			Id:       "this-is-valid-1",
			Name:     "name of my team ~ fun characters 1 %",
			Projects: []string{"valid", "wrong~"},
		},
		"Project has spaces": {
			Id:       "this-is-valid-1",
			Name:     "name of my team ~ fun characters 1 %",
			Projects: []string{"valid", "wrong space"},
		},
	}

	positiveCases := map[string]*v2.UpdateTeamReq{
		"with spaces until GA": {
			Id:       "1111valid ID for now~~~",
			Name:     "this is valid ~ fun characters",
			Projects: []string{"test", "test-2"},
		},
		"alphanumeric-post-GA": {
			Id:       "asdf-123-fun",
			Name:     "this is valid ~ fun characters",
			Projects: []string{"test", "test-2"},
		},
		"missing projects list": {
			Id:   "this-is-valid-1",
			Name: "name of my team ~ fun characters 1 %",
		},
		"empty projects list": {
			Id:       "this-is-valid-1",
			Name:     "name of my team ~ fun characters 1 %",
			Projects: []string{},
		},
	}

	classes := map[bool]map[string]*v2.UpdateTeamReq{
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

func TestCreateTeamReq(t *testing.T) {
	negativeCases := map[string]*v2.CreateTeamReq{
		"empty ID": {
			Id:       "",
			Name:     "name of my team",
			Projects: []string{"test"},
		},
		"whitespace ID": {
			Id:       "           ",
			Name:     "name of my team",
			Projects: []string{"test"},
		},
		"missing ID": {
			Name:     "name of my team",
			Projects: []string{"test"},
		},
		"empty Name": {
			Id:       "test",
			Name:     "",
			Projects: []string{"test"},
		},
		"whitespace Name": {
			Id:       "test",
			Name:     "          ",
			Projects: []string{"test"},
		},
		"missing Name": {
			Id:       "test",
			Projects: []string{"test"},
		},
		"missing Name, ID, and projects": {},
		"ID has invalid characters": {
			Id:       "wrong~",
			Name:     "name of my team",
			Projects: []string{"test"},
		},
		"ID has spaces": {
			Id:       "wrong space",
			Name:     "name of my team",
			Projects: []string{"test"},
		},
		"ID has valid characters but is too long": {
			Id:       "super-duper-long-1232-pneumonoultramicroscopicsilicovolcanoconiosis",
			Name:     "name of my team",
			Projects: []string{"test"},
		},
		"whitespace projects list": {
			Id:       "this-is-valid-1",
			Name:     "name of my team ~ fun characters 1 %",
			Projects: []string{"     ", "test"},
		},
		"repeated projects in list": {
			Id:       "this-is-valid-1",
			Name:     "name of my team ~ fun characters 1 %",
			Projects: []string{"repeat", "repeat"},
		},
		"Project has invalid characters": {
			Id:       "this-is-valid-1",
			Name:     "name of my team ~ fun characters 1 %",
			Projects: []string{"valid", "wrong~"},
		},
		"Project has spaces": {
			Id:       "this-is-valid-1",
			Name:     "name of my team ~ fun characters 1 %",
			Projects: []string{"valid", "wrong space"},
		},
	}

	positiveCases := map[string]*v2.CreateTeamReq{
		"ID present with a name": {
			Id:       "this-is-valid-1",
			Name:     "name of my team ~ fun characters 1 %",
			Projects: []string{"test", "test-2"},
		},
		"empty projects list": {
			Id:       "this-is-valid-1",
			Name:     "name of my team ~ fun characters 1 %",
			Projects: []string{},
		},
		"missing projects list": {
			Id:   "this-is-valid-1",
			Name: "name of my team ~ fun characters 1 %",
		},
	}

	classes := map[bool]map[string]*v2.CreateTeamReq{
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

func TestAddTeamMembersReq(t *testing.T) {
	negativeCases := map[string]*v2.AddTeamMembersReq{
		"empty ID": {
			Id:      "",
			UserIds: []string{"valid"},
		},
		"whitespace ID": {
			Id:      "    ",
			UserIds: []string{"valid"},
		},
		"missing ID": {
			UserIds: []string{"valid"},
		},
		"empty user list": {
			Id:      "valid",
			UserIds: []string{},
		},
		"missing user list": {
			Id: "valid",
		},
		"whitespace user list": {
			Id:      "valid",
			UserIds: []string{"     ", "test"},
		},
		"repeated users in list": {
			Id:      "valid",
			UserIds: []string{"repeat", "repeat"},
		},
		"User ID has invalid characters": {
			Id:      "valid",
			UserIds: []string{"valid", "wrong~"},
		},
		"User ID has spaces": {
			Id:      "valid",
			UserIds: []string{"valid", "wrong space"},
		},
	}

	positiveCases := map[string]*v2.AddTeamMembersReq{
		"single user ID present with crazy ID until GA": {
			Id:      "1111valid ID for now~~~",
			UserIds: []string{"this-is-valid-1@gmail+chef.com"},
		},
		"multiple unique User IDs present with crazy ID until GA (including a GUID)": {
			Id:      "1111valid ID for now~~~",
			UserIds: []string{"this-is-valid-1", "this-is-valid-2@gmail+chef.com", "6413bd61-d532-47d2-b842-0a2c9f3a8ce1"},
		},
	}

	classes := map[bool]map[string]*v2.AddTeamMembersReq{
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

func TestRemoveTeamMembersReq(t *testing.T) {
	negativeCases := map[string]*v2.RemoveTeamMembersReq{
		"empty ID": {
			Id:      "",
			UserIds: []string{"valid"},
		},
		"whitespace ID": {
			Id:      "      ",
			UserIds: []string{"valid"},
		},
		"missing ID": {
			UserIds: []string{"valid"},
		},
		"empty user list": {
			Id:      "valid",
			UserIds: []string{},
		},
		"missing user list": {
			Id: "valid",
		},
		"whitespace user list": {
			Id:      "valid",
			UserIds: []string{"   ", "test"},
		},
		"repeated users in list": {
			Id:      "valid",
			UserIds: []string{"repeat", "repeat"},
		},
		"User ID has invalid characters": {
			Id:      "valid",
			UserIds: []string{"valid", "wrong~"},
		},
		"User ID has spaces": {
			Id:      "valid",
			UserIds: []string{"valid", "wrong space"},
		},
	}

	positiveCases := map[string]*v2.RemoveTeamMembersReq{
		"single user ID present with crazy ID until GA": {
			Id:      "1111valid ID for now~~~",
			UserIds: []string{"this-is-valid-1@gmail+chef.com"},
		},
		"multiple unique User IDs present with crazy ID until GA (including a GUID)": {
			Id:      "1111valid ID for now~~~",
			UserIds: []string{"this-is-valid-1", "this-is-valid-2@gmail+chef.com", "6413bd61-d532-47d2-b842-0a2c9f3a8ce1"},
		},
	}

	classes := map[bool]map[string]*v2.RemoveTeamMembersReq{
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

func TestGetTeamsForMemberReq(t *testing.T) {
	negativeCases := map[string]*v2.GetTeamsForMemberReq{
		"empty ID": {
			UserId: "",
		},
		"whitespace ID": {
			UserId: "        ",
		},
		"missing ID": {},
		"ID has invalid characters": {
			UserId: "wrong~",
		},
		"ID has spaces": {
			UserId: "wrong space",
		},
	}

	positiveCases := map[string]*v2.GetTeamsForMemberReq{
		"when the ID is an email": {
			UserId: "this-is-valid-1@gmail+chef.com",
		},
		"when the ID is a GUID": {
			UserId: "6413bd61-d532-47d2-b842-0a2c9f3a8ce1",
		},
	}

	classes := map[bool]map[string]*v2.GetTeamsForMemberReq{
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

func TestGetTeamMembershipReq(t *testing.T) {
	negativeCases := map[string]*v2.GetTeamMembershipReq{
		"empty ID": {
			Id: "",
		},
		"whitespace ID": {
			Id: "     ",
		},
		"missing ID": {},
	}
	positiveCases := map[string]*v2.GetTeamMembershipReq{
		"with spaces until GA": {
			Id: "1111valid ID for now~~~",
		},
		"alphanumeric-post-GA": {
			Id: "asdf-123-fun",
		},
	}

	classes := map[bool]map[string]*v2.GetTeamMembershipReq{
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
