package usrutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUserData(t *testing.T) {
	// Config
	cfg := &CloudConfig{
		Hostname: "example.com",
		Username: "john",
		Password: "p@ssw0rd",
		Groups:   []string{"users", "admins"},
	}

	//
	want := &UserData{
		Hostname: "example.com",
		Users: []User{
			{
				Name:     "john",
				Password: "p@ssw0rd",
				Groups:   "users, admins",
				Home:     "/home/john",
				Shell:    "/bin/bash",
			},
		},
		ChPassword: ChPassword{
			List: "john:john",
		},
	}

	got, err := NewUserData(cfg)
	assert.NoError(t, err)

	assert.Equal(t, got.Hostname, want.Hostname)
	assert.Equal(t, got.Users[0].Name, want.Users[0].Name)
	assert.Equal(t, got.Users[0].Password, want.Users[0].Password)
	assert.Equal(t, got.Users[0].Groups, want.Users[0].Groups)
	assert.Equal(t, got.Users[0].Home, want.Users[0].Home)
	assert.Equal(t, got.Users[0].Shell, want.Users[0].Shell)
}
