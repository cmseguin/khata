package colors

import "runtime"

var Reset = "\033[0;0m"

var Red = "\033[0;31m"
var Green = "\033[0;32m"
var Yellow = "\033[0;33m"
var Blue = "\033[0;34m"
var Purple = "\033[0;35m"
var Cyan = "\033[0;36m"
var Gray = "\033[0;37m"
var White = "\033[0;97m"
var Black = "\033[0;30m"

var BoldRed = "\033[1;31m"
var BoldGreen = "\033[1;32m"
var BoldYellow = "\033[1;33m"
var BoldBlue = "\033[1;34m"
var BoldPurple = "\033[1;35m"
var BoldCyan = "\033[1;36m"
var BoldGray = "\033[1;37m"
var BoldWhite = "\033[1;97m"
var BoldBlack = "\033[1;30m"

var UnderlineRed = "\033[4;31m"
var UnderlineGreen = "\033[4;32m"
var UnderlineYellow = "\033[4;33m"
var UnderlineBlue = "\033[4;34m"
var UnderlinePurple = "\033[4;35m"
var UnderlineCyan = "\033[4;36m"
var UnderlineGray = "\033[4;37m"
var UnderlineWhite = "\033[4;97m"
var UnderlineBlack = "\033[4;30m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}
