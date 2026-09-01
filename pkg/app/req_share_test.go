package app

import (
	"context"
	"fmt"
	"testing"

	"github.com/gowsp/cloud189/pkg"
	"github.com/gowsp/cloud189/pkg/invoker"
)

func TestShareSaveSync(t *testing.T) {
	var api pkg.DriveApi = New(invoker.DefaultPath())
	//api := New(invoker.DefaultPath())
	f, accessCode, shareMode, err := api.GetShareInfo("YZfQ3ujUN77b（访问码：v5oy）") // QF7R7vBnQVR3（访问码：8qo9） Nz2iaqqMnUBn
	if err != nil {
		t.Logf("share info: %v, access code: %s", f, accessCode)
		t.Error(err)
	}

	shareId := f.PId()
	ret, err := api.ListShareDir(f.Id(), shareId, accessCode, shareMode, f.IsDir(), true)
	if err != nil {
		t.Log(ret)
		t.Error(err)
	}
	//t.Error(ret[0])

	_, err = api.ShareSaveSync(
		context.Background(), int(DealWayIgnore),
		"-11", shareId,
		ret[0],
	)
	if err != nil {
		t.Log(ret)
		t.Error(err)
	}
}

func TestShareDirInfo_Tree(t *testing.T) {
	// Construct a dummy tree
	// root
	// ├── file1.txt
	// ├── folder1
	// │   ├── subfile1.txt
	// │   └── subfolder1
	// │       └── deepfile.txt
	// └── folder2
	//     └── file2.txt

	deepFile := &shareFile{
		fileInfo: fileInfo{FileName: "deepfile.txt"},
	}
	subFolder1 := &shareFolder{
		folder: folder{DirName: "subfolder1"},
		ShareFileList: &ShareFileList{
			FileList: []*shareFile{deepFile},
		},
	}

	subFile1 := &shareFile{
		fileInfo: fileInfo{FileName: "subfile1.txt"},
	}
	folder1 := &shareFolder{
		folder: folder{DirName: "folder1"},
		ShareFileList: &ShareFileList{
			FileList:   []*shareFile{subFile1},
			FolderList: []*shareFolder{subFolder1},
		},
	}

	file2 := &shareFile{
		fileInfo: fileInfo{FileName: "file2.txt"},
	}
	folder2 := &shareFolder{
		folder: folder{DirName: "folder2"},
		ShareFileList: &ShareFileList{
			FileList: []*shareFile{file2},
		},
	}

	file1 := &shareFile{
		fileInfo: fileInfo{FileName: "file1.txt"},
	}
	root := &ShareFileList{
		FileList:   []*shareFile{file1},
		FolderList: []*shareFolder{folder1, folder2},
	}

	expected := `├── file1.txt
├── folder1
│   ├── subfile1.txt
│   └── subfolder1
│       └── deepfile.txt
└── folder2
    └── file2.txt
`

	got := root.Tree()
	if got != expected {
		t.Errorf("Tree() mismatch:\nGot:\n%s\nExpected:\n%s", got, expected)
	} else {
		fmt.Println("Tree() output matches expected:")
		fmt.Print(got)
	}
}
