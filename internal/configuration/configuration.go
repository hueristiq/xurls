package configuration

import "github.com/logrusorgru/aurora/v4"

const (
	NAME    = "xurlunpack3r"
	VERSION = "0.0.0"
)

var BANNER = func(au *aurora.Aurora) (banner string) {
	banner = au.Sprintf(
		au.BrightBlue(`
                 _                              _    _____
__  ___   _ _ __| |_   _ _ __  _ __   __ _  ___| | _|___ / _ __
\ \/ / | | | '__| | | | | '_ \| '_ \ / _`+"`"+` |/ __| |/ / |_ \| '__|
 >  <| |_| | |  | | |_| | | | | |_) | (_| | (__|   < ___) | |
/_/\_\\__,_|_|  |_|\__,_|_| |_| .__/ \__,_|\___|_|\_\____/|_|
                              |_|                         %s`).Bold(),
		au.BrightRed("v"+VERSION).Bold().Italic(),
	) + "\n\n"

	return
}
