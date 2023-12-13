package server

import (
	_ "embed"

	"sunnyvaleserv.org/portal/ui"
)

//go:embed zxcvbn.js.gz
var zxcvbn []byte

func init() {
	ui.RegisterAsset("zxcvbn.js", "text/javascript", zxcvbn, 0x62c8cb55, true)
}
