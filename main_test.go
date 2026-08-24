package main

import (
	"strings"
	"testing"
)

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
