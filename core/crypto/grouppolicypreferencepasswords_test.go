package crypto

import (
	"bytes"
	"testing"
)

func TestExtractCPasswordsFromSupportedPreferenceRoots(t *testing.T) {
	tests := []struct {
		name     string
		xml      string
		username string
		runAs    string
	}{
		{name: "groups", xml: `<Groups><User><Properties userName="local-user" cpassword="abc"/></User></Groups>`, username: "local-user"},
		{name: "scheduled tasks", xml: `<ScheduledTasks><Task><Properties runAs="DOMAIN\task-user" cpassword="abc"/></Task></ScheduledTasks>`, runAs: `DOMAIN\task-user`},
		{name: "services", xml: `<NTServices><NTService><Properties accountName="DOMAIN\service-user" cpassword="abc"/></NTService></NTServices>`, runAs: `DOMAIN\service-user`},
		{name: "drives", xml: `<Drives><Drive><Properties userName="drive-user" cpassword="abc"/></Drive></Drives>`, username: "drive-user"},
		{name: "data sources", xml: `<DataSources><DataSource><Properties userName="database-user" cpassword="abc"/></DataSource></DataSources>`, username: "database-user"},
		{name: "printers", xml: `<Printers><SharedPrinter><Properties userName="printer-user" cpassword="abc"/></SharedPrinter></Printers>`, username: "printer-user"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries := ExtractCPasswordsFromRawXML(bytes.NewBufferString(tt.xml))
			if len(entries) != 1 {
				t.Fatalf("got %d entries, want 1", len(entries))
			}
			if entries[0].UserName != tt.username {
				t.Errorf("username = %q, want %q", entries[0].UserName, tt.username)
			}
			if entries[0].RunAs != tt.runAs {
				t.Errorf("runAs = %q, want %q", entries[0].RunAs, tt.runAs)
			}
			if entries[0].CPassword != "abc" {
				t.Errorf("cpassword = %q, want abc", entries[0].CPassword)
			}
		})
	}
}
