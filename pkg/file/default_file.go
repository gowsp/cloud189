package file

import "sync"

var Sec ReadableFile = ReadableFile{Id: "-10", Name: "私密空间", IsFolder: true}
var Root ReadableFile = ReadableFile{Id: "-11", Name: "全部文件", IsFolder: true}
var Sync ReadableFile = ReadableFile{Id: "0", Name: "同步盘", IsFolder: true}
var Picture ReadableFile = ReadableFile{Id: "-12", Name: "我的图片", IsFolder: true}
var Vedio ReadableFile = ReadableFile{Id: "-13", Name: "我的视频", IsFolder: true}
var Music ReadableFile = ReadableFile{Id: "-14", Name: "我的音乐", IsFolder: true}
var Document ReadableFile = ReadableFile{Id: "-15", Name: "我的文档", IsFolder: true}
var App ReadableFile = ReadableFile{Id: "-16", Name: "我的应用", IsFolder: true}

var nameMapDir map[string]ReadableFile
var idMapDir map[string]ReadableFile
var defaultDirInstance sync.Once

var loader = func() {
	nameMapDir = map[string]ReadableFile{}
	nameMapDir["同步盘"] = Sync
	nameMapDir["私密空间"] = Sec
	nameMapDir["全部文件"] = Root
	nameMapDir["我的图片"] = Picture
	nameMapDir["我的视频"] = Vedio
	nameMapDir["我的音乐"] = Music
	nameMapDir["我的文档"] = Document
	nameMapDir["我的应用"] = App

	idMapDir = map[string]ReadableFile{}
	idMapDir["0"] = Sync
	idMapDir["-10"] = Sec
	idMapDir["-11"] = Root
	idMapDir["-12"] = Picture
	idMapDir["-13"] = Vedio
	idMapDir["-14"] = Music
	idMapDir["-15"] = Document
	idMapDir["-16"] = App
}

func DefaultIdDir() map[string]ReadableFile {
	defaultDirInstance.Do(loader)
	return idMapDir
}

func DefaultNameDir() map[string]ReadableFile {
	defaultDirInstance.Do(loader)
	return nameMapDir
}
