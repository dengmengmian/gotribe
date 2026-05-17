package adminweb

import "embed"

// Dist contains the built Admin SPA assets.
//
//go:embed all:dist
var Dist embed.FS
