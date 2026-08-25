package ldap

import (
	"net"
	"testing"

	ldapv3 "github.com/go-ldap/ldap/v3"
)

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
