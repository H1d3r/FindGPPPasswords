package ldap

import "testing"

func TestLDAPTLSConfigVerifiesServerIdentity(t *testing.T) {
	config := ldapTLSConfig("dc.example.com")

	if config.InsecureSkipVerify {
		t.Fatal("InsecureSkipVerify is enabled")
	}
	if config.ServerName != "dc.example.com" {
		t.Fatalf("ServerName = %q, want dc.example.com", config.ServerName)
	}
}
