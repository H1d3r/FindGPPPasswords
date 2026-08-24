package core

import (
	"fmt"
	"sync"
	"testing"

	"FindGPPPasswords/core/crypto"
)

func TestGPPResultCollectorConcurrentMerge(t *testing.T) {
	const workers = 100

	combined := crypto.GroupPolicyPreferencePasswordsFound{
		Entries: make(map[string][]*crypto.CPasswordEntry),
	}
	collector := gppResultCollector{results: &combined}

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			path := fmt.Sprintf(`\\dc-%d\SYSVOL\Groups.xml`, id)
			workerResults := crypto.GroupPolicyPreferencePasswordsFound{
				Entries: map[string][]*crypto.CPasswordEntry{
					path:                         {{UserName: fmt.Sprintf("user-%d", id)}},
					`\\shared\SYSVOL\Groups.xml`: {{UserName: fmt.Sprintf("shared-user-%d", id)}},
				},
			}
			collector.merge(&workerResults)
		}(i)
	}
	wg.Wait()

	if len(combined.Entries) != workers+1 {
		t.Fatalf("got %d result paths, want %d", len(combined.Entries), workers+1)
	}
	for path, entries := range combined.Entries {
		want := 1
		if path == `\\shared\SYSVOL\Groups.xml` {
			want = workers
		}
		if len(entries) != want {
			t.Errorf("path %q has %d entries, want %d", path, len(entries), want)
		}
	}
}
