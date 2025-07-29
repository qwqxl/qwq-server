package util

import (
	"fmt"
)

type Banner struct{}

func (u *Banner) Run() {
	fmt.Printf("\n%s\n", u.StringThreePoint())
	return
}

func (u *Banner) String3DAscii() string {
	banner := `
 ________   ___       __    ________      
|\   __  \ |\  \     |\  \ |\   __  \     
\ \  \|\  \\ \  \    \ \  \\ \  \|\  \    
 \ \  \\\  \\ \  \  __\ \  \\ \  \\\  \   
  \ \  \\\  \\ \  \|\__\_\  \\ \  \\\  \  
   \ \_____  \\ \____________\\ \_____  \ 
    \|___| \__\\|____________| \|___| \__\
          \|__|                      \|__|
`
	return fmt.Sprintf(banner)
}

func (u *Banner) StringBloody() string {
	banner := `
  █████   █     █░  █████  
▒██▓  ██▒▓█░ █ ░█░▒██▓  ██▒
▒██▒  ██░▒█░ █ ░█ ▒██▒  ██░
░██  █▀ ░░█░ █ ░█ ░██  █▀ ░
░▒███▒█▄ ░░██▒██▓ ░▒███▒█▄ 
░░ ▒▒░ ▒ ░ ▓░▒ ▒  ░░ ▒▒░ ▒ 
 ░ ▒░  ░   ▒ ░ ░   ░ ▒░  ░ 
   ░   ░   ░   ░     ░   ░ 
    ░        ░        ░    
`
	return fmt.Sprintf(banner)
}

func (u *Banner) StringBell() string {
	banner := "                        \n   ___.  ,  _  /   ___. \n .'   `  |  |  | .'   ` \n |    |  `  ^  ' |    | \n  `---|.  \\/ \\/   `---|.\n      |/              |/"
	return fmt.Sprintf(banner)
}

func (u *Banner) StringChunky() string {
	banner := `
.-----..--.--.--..-----.
|  _  ||  |  |  ||  _  |
|__   ||________||__   |
   |__|             |__|
`
	return fmt.Sprintf(banner)
}

func (u *Banner) StringThreePoint() string {
	banner := `
 _    _  
(_|VV(_| 
  |/   |/
`
	return fmt.Sprintf(banner)
}

func (u *Banner) StringRegular() string {
	banner := `
 ██████  ██     ██  ██████  
██    ██ ██     ██ ██    ██ 
██    ██ ██  █  ██ ██    ██ 
██ ▄▄ ██ ██ ███ ██ ██ ▄▄ ██ 
 ██████   ███ ███   ██████  
    ▀▀                 ▀▀   
`
	return fmt.Sprintf(banner)
}

func (u *Banner) StringShadow() string {
	banner := `
   ╔════════════════════════════════╗
   ║   ██████╗ ██╗    ██╗ ██████╗   ║
   ║  ██╔═══██╗██║    ██║██╔═══██╗  ║
   ║  ██║   ██║██║ █╗ ██║██║   ██║  ║
   ║  ██║▄▄ ██║██║███╗██║██║▄▄ ██║  ║
   ║  ╚██████╔╝╚███╔███╔╝╚██████╔╝  ║
   ║   ╚══▀▀═╝  ╚══╝╚══╝  ╚══▀▀═╝   ║
   ╚════════════════════════════════╝
     ▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀▄▀▄▄▀▄▀▄▀
 ══════════════════════════════════════
`
	return fmt.Sprintf(banner)
}
