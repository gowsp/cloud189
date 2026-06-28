package cmd

import (
	"io/fs"
	"testing"
	"time"
)

type testFileInfo struct {
	name    string
	size    int64
	dir     bool
	id      string
	pid     string
	md5     string
	modTime time.Time
}

func (f testFileInfo) Name() string       { return f.name }
func (f testFileInfo) Size() int64        { return f.size }
func (f testFileInfo) Mode() fs.FileMode  { return 0644 }
func (f testFileInfo) ModTime() time.Time { return f.modTime }
func (f testFileInfo) IsDir() bool        { return f.dir }
func (f testFileInfo) Sys() any           { return nil }
func (f testFileInfo) Id() string         { return f.id }
func (f testFileInfo) PId() string        { return f.pid }

type testChecksumInfo struct {
	testFileInfo
	Md5 string
}

func TestFileToJSONEntry(t *testing.T) {
	modTime := time.Date(2026, 6, 28, 12, 0, 0, 0, time.FixedZone("CST", 8*60*60))
	info := testChecksumInfo{
		testFileInfo: testFileInfo{
			name:    "a.txt",
			size:    123,
			id:      "123",
			pid:     "456",
			modTime: modTime,
		},
		Md5: "abc",
	}
	entry := fileToJSONEntry("/target/a.txt", info)
	if entry.Path != "/target/a.txt" {
		t.Fatalf("Path = %q", entry.Path)
	}
	if entry.RelativePath != "a.txt" {
		t.Fatalf("RelativePath = %q", entry.RelativePath)
	}
	if entry.Name != "a.txt" || entry.Size != 123 || entry.IsDir {
		t.Fatalf("unexpected basic entry: %#v", entry)
	}
	if entry.FileID != "123" || entry.ParentFileID != "456" {
		t.Fatalf("unexpected ids: %#v", entry)
	}
	if entry.ModifiedAt != "2026-06-28T04:00:00Z" {
		t.Fatalf("ModifiedAt = %q", entry.ModifiedAt)
	}
	if entry.Checksum != "abc" {
		t.Fatalf("Checksum = %q", entry.Checksum)
	}
}

func TestRelativeCloudPath(t *testing.T) {
	cases := map[string]struct {
		base string
		name string
		want string
	}{
		"child":     {base: "/目标", name: "/目标/a.txt", want: "a.txt"},
		"nested":    {base: "/目标", name: "/目标/dir/a.txt", want: "dir/a.txt"},
		"root base": {base: "/", name: "/目标/a.txt", want: "目标/a.txt"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := relativeCloudPath(tc.base, tc.name); got != tc.want {
				t.Fatalf("relativeCloudPath(%q, %q) = %q, want %q", tc.base, tc.name, got, tc.want)
			}
		})
	}
}
