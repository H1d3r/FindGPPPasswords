package main

import (
	"strings"
	"testing"
)

func TestDefaultLDAPPort(t *testing.T) {
	tests := []struct {
		name     string
		port     int
		useLdaps bool
		want     int
	}{
		{name: "ldap default", want: 389},
		{name: "ldaps default", useLdaps: true, want: 636},
		{name: "explicit ldap port", port: 1389, want: 1389},
		{name: "explicit ldaps port", port: 1636, useLdaps: true, want: 1636},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := defaultLDAPPort(tt.port, tt.useLdaps); got != tt.want {
				t.Fatalf("defaultLDAPPort(%d, %t) = %d, want %d", tt.port, tt.useLdaps, got, tt.want)
			}
		})
	}
}

func TestCredentialTestMessageNoColors(t *testing.T) {
	message := credentialTestMessage(true, "DOMAIN", "user", "password", true)
	if strings.Contains(message, "\x1b[") {
		t.Fatalf("message contains ANSI escape: %q", message)
	}
	if message != `   [+] DOMAIN\user : password` {
		t.Fatalf("message = %q", message)
	}
}

func TestCredentialTestMessageColors(t *testing.T) {
	message := credentialTestMessage(false, "", "user", "password", false)
	if !strings.HasPrefix(message, "\x1b[91m") || !strings.HasSuffix(message, "\x1b[0m") {
		t.Fatalf("message does not contain expected ANSI wrapper: %q", message)
	}
}
