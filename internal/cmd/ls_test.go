package cmd

import (
	"io/fs"
	"sort"
	"testing"
)

type testDirEntry struct {
	info testFileInfo
}

func (d testDirEntry) Name() string               { return d.info.Name() }
func (d testDirEntry) IsDir() bool                { return d.info.IsDir() }
func (d testDirEntry) Type() fs.FileMode          { return d.info.Mode().Type() }
func (d testDirEntry) Info() (fs.FileInfo, error) { return d.info, nil }

type testDirReader map[string][]fs.DirEntry

func (r testDirReader) ReadDir(name string) ([]fs.DirEntry, error) {
	return r[name], nil
}

func TestListJSONEntriesRecursiveFilesOnly(t *testing.T) {
	reader := testDirReader{
		"/我的文档": {
			testDirEntry{info: testFileInfo{name: "b.txt", size: 2, id: "2", pid: "root"}},
			testDirEntry{info: testFileInfo{name: "dir", dir: true, id: "dir", pid: "root"}},
		},
		"/我的文档/dir": {
			testDirEntry{info: testFileInfo{name: "a.txt", size: 1, id: "1", pid: "dir"}},
		},
	}
	entries, err := listJSONEntries(reader, "/我的文档", true)
	if err != nil {
		t.Fatal(err)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].RelativePath < entries[j].RelativePath
	})
	if len(entries) != 2 {
		t.Fatalf("len(entries) = %d, want 2: %#v", len(entries), entries)
	}
	if entries[0].RelativePath != "b.txt" || entries[0].Path != "/我的文档/b.txt" {
		t.Fatalf("entry[0] = %#v", entries[0])
	}
	if entries[1].RelativePath != "dir/a.txt" || entries[1].Path != "/我的文档/dir/a.txt" {
		t.Fatalf("entry[1] = %#v", entries[1])
	}
	for _, entry := range entries {
		if entry.IsDir {
			t.Fatalf("recursive listing should omit directories: %#v", entries)
		}
	}
}

func TestListJSONEntriesNonRecursiveIncludesDirectories(t *testing.T) {
	reader := testDirReader{
		"/我的文档": {
			testDirEntry{info: testFileInfo{name: "dir", dir: true, id: "dir", pid: "root"}},
		},
	}
	entries, err := listJSONEntries(reader, "/我的文档", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("len(entries) = %d, want 1", len(entries))
	}
	if !entries[0].IsDir || entries[0].RelativePath != "dir" {
		t.Fatalf("entry = %#v", entries[0])
	}
}
