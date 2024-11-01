package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	sat "github.com/tokopedia/golang-sat"
	"github.com/tokopedia/golang-sat/signature"
)

func main() {
	c := &http.Client{
		Timeout: 30 * time.Second,
	}

	cln, err := sat.NewClient(
		CLIENT_ID,
		CLIENT_SECRET,
		PRIVATE_KEY,
		sat.WithServerPublicKeyString(PUBLIC_KEY),
		sat.WithPaddingType(signature.PaddingTypePSS),
		sat.WithIsDebug(true),
		sat.WithHTTPClient(c),
	)
	if err != nil {
		panic(err)
	}

	ie := integrationExample{client: cln}
	clbe := &callbackExample{}

	http.HandleFunc("/test", ie.Handle)
	http.HandleFunc("/callback", cln.HandleCallback(clbe))
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", SERVER_PORT))
	if err == nil {
		fmt.Println("Listening on port ", SERVER_PORT)
	}

	http.Serve(l, nil)
}
