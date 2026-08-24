package gpp

import (
	"bytes"
	"testing"
)

const encryptedPodalirius = "bdajdgpjZqolVYI3h2O2mp+JpxDuZd0xoi2M86z7JuI="

func TestExtractCPasswordsFromGroupsXML(t *testing.T) {
	buffer := bytes.NewBufferString(`<Groups><User><Properties userName="alice" newName="renamed" cpassword="` + encryptedPodalirius + `"/></User></Groups>`)

	entries, err := ExtractCPasswordsFromRawXML(buffer)
	if err != nil {
		t.Fatalf("ExtractCPasswordsFromRawXML() error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("ExtractCPasswordsFromRawXML() returned %d entries, want 1", len(entries))
	}
	if entries[0].UserName != "alice" || entries[0].NewName != "renamed" || entries[0].Password != "Podalirius" {
		t.Fatalf("ExtractCPasswordsFromRawXML() entry = %#v", entries[0])
	}
}

func TestExtractCPasswordsFromScheduledTasksXML(t *testing.T) {
	buffer := bytes.NewBufferString(`<ScheduledTasks><Task name="task"><Properties runAs="DOMAIN\alice" cpassword="` + encryptedPodalirius + `"/></Task></ScheduledTasks>`)

	entries, err := ExtractCPasswordsFromRawXML(buffer)
	if err != nil {
		t.Fatalf("ExtractCPasswordsFromRawXML() error = %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("ExtractCPasswordsFromRawXML() returned %d entries, want 1", len(entries))
	}
	if entries[0].RunAs != `DOMAIN\alice` || entries[0].Password != "Podalirius" {
		t.Fatalf("ExtractCPasswordsFromRawXML() entry = %#v", entries[0])
	}
}

func TestExtractCPasswordsRejectsMalformedXML(t *testing.T) {
	buffer := bytes.NewBufferString(`<Groups><User></Groups>`)

	if _, err := ExtractCPasswordsFromRawXML(buffer); err == nil {
		t.Fatal("ExtractCPasswordsFromRawXML() error = nil, want malformed XML error")
	}
}

func TestExtractCPasswordsRejectsInvalidCPassword(t *testing.T) {
	buffer := bytes.NewBufferString(`<Groups><User><Properties userName="alice" cpassword="invalid"/></User></Groups>`)

	if _, err := ExtractCPasswordsFromRawXML(buffer); err == nil {
		t.Fatal("ExtractCPasswordsFromRawXML() error = nil, want decryption error")
	}
}
