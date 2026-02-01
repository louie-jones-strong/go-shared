package filecache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupPaths(t *testing.T) (string, string) {
	t.Helper()
	tmp := t.TempDir()
	manifestPath := filepath.Join(tmp, "manifest.json")
	itemsDir := filepath.Join(tmp, "items")
	if err := os.MkdirAll(itemsDir, 0o755); err != nil {
		t.Fatalf("failed creating items dir: %v", err)
	}
	return manifestPath, itemsDir
}

func TestNilReceiverAndErrorPaths(t *testing.T) {
	// FileCache nil receiver for methods that return errors
	var fc *FileCache[string] = nil

	if _, err := fc.RemoveFiles("a"); err == nil {
		t.Fatalf("expected error calling RemoveFiles on nil FileCache")
	}

	if _, err := fc.CleanupExpiredItems(time.Hour); err == nil {
		t.Fatalf("expected error calling CleanupExpiredItems on nil FileCache")
	}

	// wLock/rLock panic on nil receiver
	var paniced bool
	func() {
		defer func() {
			if r := recover(); r != nil {
				paniced = true
			}
		}()
		_ = (*FileCache[string])(nil).wLock()
	}()
	if !paniced {
		t.Fatalf("expected panic from wLock on nil receiver")
	}

	paniced = false
	func() {
		defer func() {
			if r := recover(); r != nil {
				paniced = true
			}
		}()
		_ = (*FileCache[string])(nil).rLock()
	}()
	if !paniced {
		t.Fatalf("expected panic from rLock on nil receiver")
	}

	// getManifest and saveManifest panic on nil
	paniced = false
	func() {
		defer func() {
			if r := recover(); r != nil {
				paniced = true
			}
		}()
		_ = (*FileCache[string])(nil).getManifest()
	}()
	if !paniced {
		t.Fatalf("expected panic from getManifest on nil receiver")
	}

	paniced = false
	func() {
		defer func() {
			if r := recover(); r != nil {
				paniced = true
			}
		}()
		_ = (*FileCache[string])(nil).saveManifest()
	}()
	if !paniced {
		t.Fatalf("expected panic from saveManifest on nil receiver")
	}
}

func TestFileGroupInfo_NilAndCacheHolderErrors(t *testing.T) {
	var nilFi *FileGroupInfo = nil
	if err := nilFi.SaveFile("k", ".txt", []byte("x")); err == nil {
		t.Fatalf("expected error calling SaveFile on nil FileGroupInfo")
	}
	if _, err := nilFi.LoadFile("k"); err == nil {
		t.Fatalf("expected error calling LoadFile on nil FileGroupInfo")
	}
	if _, err := nilFi.LoadFiles(); err == nil {
		t.Fatalf("expected error calling LoadFiles on nil FileGroupInfo")
	}
	if err := nilFi.deleteFiles(); err == nil {
		t.Fatalf("expected error calling deleteFiles on nil FileGroupInfo")
	}

	// new FileGroupInfo without cacheHolder should return errors
	fi := newFGI()
	if err := fi.SaveFile("k", ".txt", []byte("x")); err == nil {
		t.Fatalf("expected error SaveFile with nil cacheHolder")
	}
	if _, err := fi.LoadFile("k"); err == nil {
		t.Fatalf("expected error LoadFile with nil cacheHolder")
	}
	if _, err := fi.LoadFiles(); err == nil {
		t.Fatalf("expected error LoadFiles with nil cacheHolder")
	}
	if err := fi.deleteFiles(); err == nil {
		t.Fatalf("expected error deleteFiles with nil cacheHolder")
	}
}

func TestSaveFileCreatesAndDeleteRemovesFileOnDisk(t *testing.T) {
	manifestPath, itemsDir := setupPaths(t)
	fc := New[string](manifestPath, itemsDir)

	fi := fc.GetOrCreateFileInfo("gg")
	if fi == nil {
		t.Fatalf("expected filegroupinfo")
	}

	// Save a file and ensure it exists on disk
	if err := fi.SaveFile(DefaultFileKey, ".bin", []byte("payload")); err != nil {
		t.Fatalf("SaveFile error: %v", err)
	}

	fileName := fi.Files[DefaultFileKey]
	if fileName == "" {
		t.Fatalf("expected filename set in Files map")
	}
	fp := filepath.Join(itemsDir, fileName)
	if _, err := os.Stat(fp); err != nil {
		t.Fatalf("expected saved file on disk: %v", err)
	}

	// deleteFiles should remove the physical file
	if err := fi.deleteFiles(); err != nil {
		t.Fatalf("deleteFiles error: %v", err)
	}
	if _, err := os.Stat(fp); !os.IsNotExist(err) {
		t.Fatalf("expected file removed, stat error: %v", err)
	}
}

func TestFileGroupInfo_IsValid_TableDriven(t *testing.T) {
	tests := []struct {
		name   string
		setup  func() *FileGroupInfo
		expire time.Duration
		want   bool
	}{
		{
			name:   "nil-fi",
			setup:  func() *FileGroupInfo { return nil },
			expire: time.Hour,
			want:   false},
		{
			name:   "no-files",
			setup:  func() *FileGroupInfo { return newFGI() },
			expire: time.Hour,
			want:   false,
		},
		{
			name:   "negative-expire-with-file",
			setup:  func() *FileGroupInfo { f := newFGI(); f.Files["a"] = "b"; return f },
			expire: -1,
			want:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fi := tc.setup()
			got := fi.IsValid(tc.expire)
			if got != tc.want {
				t.Fatalf("%s: expected %v, got %v", tc.name, tc.want, got)
			}
		})
	}
}
