package configuration

import "github.com/logrusorgru/aurora/v4"

const (
	NAME    = "xurl"
	VERSION = "0.0.0"
)

var BANNER = func(au *aurora.Aurora) (banner string) {
	banner = au.Sprintf(
		au.BrightBlue(`
                 _ 
__  ___   _ _ __| |
\ \/ / | | | '__| |
 >  <| |_| | |  | |
/_/\_\\__,_|_|  |_|
             %s`).Bold(),
		au.BrightRed("v"+VERSION).Bold().Italic(),
	) + "\n\n"

	return
}
