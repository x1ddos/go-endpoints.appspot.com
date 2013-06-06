package app

import (
	"net/http"

	"github.com/crhym3/go-endpoints/endpoints"
	"github.com/crhym3/go-tictactoe/tictactoe"
)

func init() {
	http.HandleFunc("/", homeHandler)

	registerGuestbookApi()
	tictactoe.RegisterService()

	endpoints.HandleHttp()
}
