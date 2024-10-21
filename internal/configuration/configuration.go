package configuration

import "github.com/logrusorgru/aurora/v3"

const (
	NAME    = "xurlbits"
	VERSION = "0.1.0"
)

// BANNER is this project's CLI display banner.
var BANNER = aurora.Sprintf(
	aurora.BrightBlue(`
                 _ _     _ _
__  ___   _ _ __| | |__ (_) |_ ___
\ \/ / | | | '__| | '_ \| | __/ __|
 >  <| |_| | |  | | |_) | | |_\__ \
/_/\_\\__,_|_|  |_|_.__/|_|\__|___/
                             %s`).Bold(),
	aurora.BrightRed("v"+VERSION).Bold(),
)
