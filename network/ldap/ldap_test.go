package ldap

import (
	"net"
	"testing"

	ldapv3 "github.com/go-ldap/ldap/v3"
)

func TestLDAPTLSConfigVerifiesServerIdentity(t *testing.T) {
	config := ldapTLSConfig("dc.example.com")

	if config.InsecureSkipVerify {
		t.Fatal("InsecureSkipVerify is enabled")
	}
	if config.ServerName != "dc.example.com" {
		t.Fatalf("ServerName = %q, want dc.example.com", config.ServerName)
	}
}

func TestSessionCloseIsIdempotent(t *testing.T) {
	client, server := net.Pipe()
	defer server.Close()

	session := Session{connection: ldapv3.NewConn(client, false)}
	session.connection.Start()
	session.Close()
	if session.connection != nil {
		t.Fatal("connection was not cleared")
	}

	session.Close()
}
