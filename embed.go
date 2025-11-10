package backend

import "embed"

//go:embed web/database/dist/*
var WebDataBaseFiles embed.FS

//go:embed web/admin/dist/*
var WebAdminFiles embed.FS

type EmbedsInfo struct {
	Fs   embed.FS
	Path string
}

var WebEmbeds = map[string]EmbedsInfo{
	"/admin": {Fs: WebAdminFiles, Path: "web/admin/dist"},
	"":       {Fs: WebDataBaseFiles, Path: "web/database/dist"},
}
