package public

import "embed"

//go:embed css/*.css
//go:embed js/*.js
//go:embed img/*
var AssetsFS embed.FS
