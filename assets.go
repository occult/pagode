package assets

import "embed"

//go:embed public/build/assets/*
var StaticFS embed.FS
