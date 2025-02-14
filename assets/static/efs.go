package static

import (
	"embed"
)

// Embed static files and subdirectories under "static" directory
// Adjust the path to point to the correct directory
// Here, it assumes `assets/efs.go` is in the `assets` directory and static files are in the sibling directory

//go:embed *
var StaticFiles embed.FS
