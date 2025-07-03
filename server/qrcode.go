package server

import (
	"fmt"

	qrcodeTerminal "github.com/Baozisoftware/qrcode-terminal-go"
)

func PrintQRCode(addr string) {
	qr := qrcodeTerminal.New()
	qr.Get("http://" + addr).Print()
	fmt.Printf("Starting server at http://%s\n", addr)
}
