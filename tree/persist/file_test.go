package persist_test

import (
	"os"
	"testing"

	"github.com/ankored/tree-test/tree/persist"
)

func TestFileCreatesFilesIfDoNotExist(t *testing.T) {
	t.Cleanup(func() {
		os.Remove("./test_log")
		os.Remove("./test_cmp")
	})

	if _, err := persist.NewFile("./test_cmp", "./test_log"); err != nil {
		t.Fatalf("unexpected error: %#v", err)
	}

	// Should have made the files
	if _, err := os.Stat("./test_cmp"); err != nil {
		t.Fatalf("stats test_cmp: %s", err)
	}
	if _, err := os.Stat("./test_log"); err != nil {
		t.Fatalf("stats test_log: %s", err)
	}
}
