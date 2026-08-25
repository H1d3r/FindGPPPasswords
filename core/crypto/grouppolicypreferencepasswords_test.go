package crypto

import (
	"bytes"
	"testing"
)

func TestExtractCPasswordsFromRawXMLReturnsParseErrors(t *testing.T) {
	tests := []struct {
		name string
		xml  string
	}{
		{name: "groups", xml: `<Groups><User></Groups>`},
		{name: "scheduled tasks", xml: `<ScheduledTasks><Task></ScheduledTasks>`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entries, err := ExtractCPasswordsFromRawXML(bytes.NewBufferString(tt.xml))
			if err == nil {
				t.Fatal("expected parse error, got nil")
			}
			if entries != nil {
				t.Fatalf("entries = %#v, want nil", entries)
			}
		})
	}
}

func TestExtractCPasswordsFromRawXMLValidDocument(t *testing.T) {
	entries, err := ExtractCPasswordsFromRawXML(bytes.NewBufferString(`<Groups><User><Properties userName="user" cpassword="abc"/></User></Groups>`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}
}
