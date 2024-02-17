package v1

import "embed"

//go:embed api.yml redoc.html
var Content embed.FS
