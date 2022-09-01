package entities

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseValidSSHHostKeys(t *testing.T) {
	givenSSHHostKeysContent := `ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzd root@ip-10-0-0-179
	ssh-ed25519 AAAAC3NzaC1lZD root@ip-10-0-0-179
	ssh-rsa AAAAB3NzaC root@ip-10-0-0-179
`
	expectedSSHHostKeys := []EnvSSHHostKey{
		{
			Algorithm:   "ecdsa-sha2-nistp256",
			Fingerprint: "AAAAE2VjZHNhLXNoYTItbmlzd",
		},

		{
			Algorithm:   "ssh-ed25519",
			Fingerprint: "AAAAC3NzaC1lZD",
		},

		{
			Algorithm:   "ssh-rsa",
			Fingerprint: "AAAAB3NzaC",
		},
	}

	returnedSSHHostKeys, err := ParseSSHHostKeys(givenSSHHostKeysContent)

	if err != nil {
		t.Fatalf(
			"expected no error, got %s",
			err,
		)
	}

	if !reflect.DeepEqual(expectedSSHHostKeys, returnedSSHHostKeys) {
		t.Fatalf(
			"expected SSH host keys to equal '%+v', got '%+v'",
			expectedSSHHostKeys,
			returnedSSHHostKeys,
		)
	}
}

func TestParseInvalidSSHHostKeys(t *testing.T) {
	giventSSHHostKeysContent := "host_keys_content"
	returnedSSHHostKeys, err := ParseSSHHostKeys(giventSSHHostKeysContent)

	if err == nil {
		t.Fatalf("expected error, got nothing")
	}

	if returnedSSHHostKeys != nil {
		t.Fatalf(
			"expected no SSH host keys, got '%+v'",
			returnedSSHHostKeys,
		)
	}
}

func TestCheckPortValidityWithValidPorts(t *testing.T) {
	testCases := []struct {
		test          string
		port          string
		reservedPorts []string
	}{
		{
			test:          "with valid port",
			port:          "8080",
			reservedPorts: []string{"2200"},
		},

		{
			test:          "with minimum port",
			port:          "1",
			reservedPorts: []string{},
		},

		{
			test:          "with maximum port",
			port:          "65535",
			reservedPorts: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := CheckPortValidity(
				tc.port,
				tc.reservedPorts,
			)

			if err != nil {
				t.Fatalf("expected no error, got '%+v'", err)
			}
		})
	}
}

func TestCheckPortValidityWithInvalidPorts(t *testing.T) {
	testCases := []struct {
		test          string
		port          string
		reservedPorts []string
	}{
		{
			test:          "with invalid port",
			port:          "invalid_port",
			reservedPorts: []string{"2200"},
		},

		{
			test:          "with less than minimum port",
			port:          "0",
			reservedPorts: []string{},
		},

		{
			test:          "with more than maximum port",
			port:          "65536",
			reservedPorts: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := CheckPortValidity(
				tc.port,
				tc.reservedPorts,
			)

			if err == nil {
				t.Fatalf("expected error, got nothing")
			}

			if !errors.As(err, &ErrInvalidPort{}) {
				t.Fatalf("expected invalid port error, got '%+v'", err)
			}
		})
	}
}

func TestCheckPortValidityWithReservedPorts(t *testing.T) {
	testCases := []struct {
		test          string
		port          string
		reservedPorts []string
	}{
		{
			test:          "with reserved port",
			port:          "2200",
			reservedPorts: []string{"2200", "8", "100"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.test, func(t *testing.T) {
			err := CheckPortValidity(
				tc.port,
				tc.reservedPorts,
			)

			if err == nil {
				t.Fatalf("expected error, got nothing")
			}

			if !errors.As(err, &ErrReservedPort{}) {
				t.Fatalf("expected reserved port error, got '%+v'", err)
			}
		})
	}
}
