package embedder

import _ "embed"

//go:embed Warning.README.md
var WarningReadMe []byte

const WarningReadMeName = "Warning.README.md"
