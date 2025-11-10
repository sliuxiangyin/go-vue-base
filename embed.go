//go:build !dev
// +build !dev

package databaseAi

import "embed"

//go:embed web/learn_en/dist/*
var WebLearnEnFiles embed.FS

//go:embed web/admin/dist/*
var WebAdminFiles embed.FS

type EmbedsInfo struct {
	Fs   embed.FS
	Path string
}

var WebEmbeds = map[string]EmbedsInfo{
	"/admin": {Fs: WebAdminFiles, Path: "web/admin/dist"},
	"":       {Fs: WebLearnEnFiles, Path: "web/learn_en/dist"},
}
