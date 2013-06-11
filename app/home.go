package app

import (
	"html/template"
	"net/http"
	"path/filepath"
	"io/ioutil"
	"time"

	"appengine"
	"appengine/memcache"
	"appengine/urlfetch"

	"github.com/russross/blackfriday"
)

// where to fetch README.md from
const readmeUrl = "https://raw.github.com/crhym3/go-endpoints/master/README.md"

// template custom functions / filters
var funcs = template.FuncMap{
	"safe": func(html string) template.HTML {
		return template.HTML(html)
	},
}

func getReadme(c appengine.Context) ([]byte, error) {
	item, err := memcache.Get(c, readmeUrl)
	if err == nil {
		return item.Value, nil
	}

	c.Debugf("Fetching readme from %s", readmeUrl)
	resp, err := urlfetch.Client(c).Get(readmeUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	readme, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	readme = blackfriday.MarkdownCommon(readme)
	item = &memcache.Item{
		Key: readmeUrl,
		Value: readme,
		Expiration: 24 * time.Hour,
	}
	memcache.Set(c, item)
	return readme, nil
}

func getTemplate(name string) (*template.Template, error) {
	t, err := template.New(name).Funcs(funcs).ParseFiles(
		filepath.Join("templates", name))
	if err != nil {
		return nil, err
	}
	return t, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	t, err := getTemplate("home.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	c := appengine.NewContext(r)
	readme, err := getReadme(c)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	data := &struct {
		Readme template.HTML
	}{
		Readme: template.HTML(readme),
	}
	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
