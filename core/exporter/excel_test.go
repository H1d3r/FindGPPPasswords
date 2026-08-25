package exporter

import (
	"path/filepath"
	"testing"

	"FindGPPPasswords/core/config"
	"FindGPPPasswords/core/crypto"

	"github.com/xuri/excelize/v2"
)

func TestExcelOutputPath(t *testing.T) {
	outputDir := filepath.Join(string(filepath.Separator), "tmp", "scan")
	absolute := filepath.Join(string(filepath.Separator), "exports", "results.xlsx")

	if got := excelOutputPath(outputDir, "reports/results.xlsx"); got != filepath.Join(outputDir, "reports", "results.xlsx") {
		t.Fatalf("relative path = %q", got)
	}
	if got := excelOutputPath(outputDir, absolute); got != absolute {
		t.Fatalf("absolute path = %q, want %q", got, absolute)
	}
}

func TestGenerateExcelWritesOneCredentialPerRow(t *testing.T) {
	outputDir := t.TempDir()
	settings := config.Config{OutputDir: outputDir}
	settings.Credentials.Domain = "example"
	results := crypto.GroupPolicyPreferencePasswordsFound{
		Entries: map[string][]*crypto.CPasswordEntry{
			`\\dc\SYSVOL\Groups.xml`: {
				{RunAs: `EXAMPLE\task-user`, Password: "task-password"},
				{UserName: "local-user", NewName: "renamed-user", Password: "local-password"},
			},
		},
	}

	GenerateExcel(results, settings, filepath.Join("reports", "results.xlsx"))
	workbookPath := filepath.Join(outputDir, "reports", "results.xlsx")
	workbook, err := excelize.OpenFile(workbookPath)
	if err != nil {
		t.Fatalf("open workbook: %v", err)
	}
	defer workbook.Close()

	rows, err := workbook.GetRows("EXAMPLE")
	if err != nil {
		t.Fatalf("get rows: %v", err)
	}
	if len(rows) != 3 {
		t.Fatalf("got %d rows, want 3", len(rows))
	}
	wantHeader := []string{"RunAs", "Username", "NewName", "Password", "Path"}
	for column, want := range wantHeader {
		if rows[0][column] != want {
			t.Errorf("header column %d = %q, want %q", column, rows[0][column], want)
		}
	}
	if rows[1][0] != `EXAMPLE\task-user` || rows[1][3] != "task-password" {
		t.Errorf("scheduled-task row = %#v", rows[1])
	}
	if rows[2][1] != "local-user" || rows[2][2] != "renamed-user" || rows[2][3] != "local-password" {
		t.Errorf("local-user row = %#v", rows[2])
	}
}
