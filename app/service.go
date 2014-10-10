package app

import (
	"net/http"
	"time"

	"appengine/datastore"

	"github.com/crhym3/go-endpoints/endpoints"
)

const clientId = "YOUR-CLIENT-ID"

var (
	scopes = []string{
		endpoints.EmailScope,
		"https://www.googleapis.com/auth/userinfo.profile",
	}
	clientIds = []string{clientId, endpoints.APIExplorerClientID}
	audiences = []string{clientId}
)

// Greeting is a datastore entity that represents a single greeting.
// It also serves as (part of) a response of GreetingService.
type Greeting struct {
	Id      string    `json:"id,omitempty" datastore:"-"`
	Author  string    `json:"author"`
	Content string    `json:"content" datastore:",noindex"`
	Date    time.Time `json:"date"`
}

// GreetingService can sign the guesbook, list all greetings and delete
// a greeting from the guestbook.
type GreetingService struct {
}

// GreetingsList is a response type of GreetingService.List method
type GreetingsList struct {
	Items []*Greeting `json:"items"`
}

// List responds with a list of all greetings ordered by Date field.
// Most recent greets come first.
func (gs *GreetingService) List(
	r *http.Request, _ *endpoints.VoidMessage, resp *GreetingsList) error {

	c := endpoints.NewContext(r)
	q := datastore.NewQuery("Greeting").Order("-Date").Limit(10)
	greets := make([]*Greeting, 0, 10)
	keys, err := q.GetAll(c, &greets)
	if err != nil {
		return err
	}
	for i, k := range keys {
		greets[i].Id = k.Encode()
	}
	resp.Items = greets
	return nil
}

// NewGreet is the expected data structure for signing the guestbook.
type NewGreet struct {
	Message string `json:"message" endpoints:"req"`
}

// Sign creates a new Greeting based on provided NewGreet.
// It stores newly created Greeting with Content being that of NewGreet.Message.
// Author field will be either currently signed in user or "Anonymous User".
func (gs *GreetingService) Sign(
	r *http.Request, req *NewGreet, greet *Greeting) error {

	c := endpoints.NewContext(r)

	greet.Content = req.Message
	greet.Date = time.Now()

	user, err := endpoints.CurrentUser(c, scopes, audiences, clientIds)
	if err == nil {
		greet.Author = user.String()
	} else {
		greet.Author = "Anonymous User (" + err.Error() + ")"
	}

	key, err := datastore.Put(
		c, datastore.NewIncompleteKey(c, "Greeting", nil), greet)
	if err != nil {
		return err
	}

	greet.Id = key.Encode()
	return nil
}

// GreetIdReq serves as a data structure for identifying a single Greeting.
type GreetIdReq struct {
	Id string `json:"id" endpoints:"req"`
}

// Delete deletes a single greeting from the guesbook and replies with nothing.
func (gs *GreetingService) Delete(
	r *http.Request, req *GreetIdReq, _ *endpoints.VoidMessage) error {

	c := endpoints.NewContext(r)
	key, err := datastore.DecodeKey(req.Id)
	if err != nil {
		return err
	}
	return datastore.Delete(c, key)
}

// TestMessageGet is here just to test various field types
type TestMessageGet struct {
	A int   `endpoints:"req"`
	B int32 `endopints:"100"`                               // default
	C int64 `json:",string" endpoints:",This is a C field"` // description
	D float32
	E float64
	F string
	G bool
	// TODO: add enum
}

// TestMessagePost is here just to test various field types
type TestMessagePost struct {
	A int   `endpoints:"req"`
	B int32 `endpoints:"d=100"`              // default
	C int64 `endpoints:",This is a C field"` // description
	D float32
	E float64
	F string
	G bool
	H time.Time
	I []byte
	J *TestMessageGet
	K []*TestMessagePost
	// TODO: add enum
}

// EchoGet is a method to test different message field types.
func (gs *GreetingService) EchoGet(
	r *http.Request, req *TestMessageGet, res *TestMessageGet) error {

	*res = *req
	return nil
}

// QueryEchoGet is the same as EchoGet, only that it's API path template
// does not contain message fields so they all should go in the query part
// of a request URL.
func (gs *GreetingService) QueryEchoGet(
	r *http.Request, req *TestMessageGet, res *TestMessageGet) error {

	*res = *req
	return nil
}

// EchoPost is a method to test different message field types.
func (gs *GreetingService) EchoPost(
	r *http.Request, req *TestMessagePost, res *TestMessagePost) error {

	*res = *req
	return nil
}

func registerGuestbookApi() (*endpoints.RPCService, error) {
	greetService := &GreetingService{}
	rpcService, err := endpoints.RegisterServiceWithDefaults(greetService)
	if err != nil {
		return nil, err
	}
	rpcService.Info().Name = "greeting"

	info := rpcService.MethodByName("List").Info()
	info.Name, info.HTTPMethod, info.Path, info.Desc =
		"greets.list", "GET", "greetings", "List most recent greetings."

	info = rpcService.MethodByName("Sign").Info()
	info.Name, info.HTTPMethod, info.Path, info.Desc =
		"greets.sign", "POST", "greetings", "Sign the guestbook."
	info.Scopes = scopes
	info.Audiences = audiences
	info.ClientIds = clientIds

	info = rpcService.MethodByName("Delete").Info()
	info.Name, info.HTTPMethod, info.Path, info.Desc =
		"greets.delete", "DELETE", "greetings/{id}", "Delete a single Greeting."

	info = rpcService.MethodByName("EchoGet").Info()
	info.Name, info.HTTPMethod, info.Path =
		// These should match TestMessageGet field names
		"tests.echoGet", "GET", "tests/echo/{A}/{B}/{C}/{D}/{E}/{F}/{G}"

	info = rpcService.MethodByName("QueryEchoGet").Info()
	info.Name, info.HTTPMethod, info.Path =
		"tests.queryEchoGet", "GET", "tests/query"

	info = rpcService.MethodByName("EchoPost").Info()
	info.Name, info.HTTPMethod, info.Path =
		"tests.echoPost", "POST", "tests/echo/{B}"

	return rpcService, nil
}
