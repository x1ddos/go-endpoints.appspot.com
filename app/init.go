package app

import (
	"net/http"

	"github.com/crhym3/go-endpoints/endpoints"
	"github.com/crhym3/go-tictactoe/tictactoe"
)

func init() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/guestbook", guestbookHandler)

	registerGuestbookApi()
	tictactoe.RegisterService()

	endpoints.HandleHttp()
}
