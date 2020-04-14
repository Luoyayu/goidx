package api

import (
	"github/luoyayu/goidx/config"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Do(method, opName, path, query, auth string) *http.Response {
	if auth == "" {
		auth = config.RootBasicAuth
	}

	url_ := (&url.URL{
		Scheme:   "https",
		Host:     config.WorkerHost,
		Path:     path,
		RawQuery: query,
	}).String()

	//log.Println(method, url_)

	req, _ := http.NewRequest(strings.ToUpper(method), url_, nil)
	req.Header.Set("authorization", auth)
	resp, err := (&http.Client{Timeout: time.Second * 10}).Do(req)

	if err != nil {
		log.Fatalln(opName, err)
	} else if resp.StatusCode != http.StatusOK {
		log.Fatalln(opName, resp.StatusCode, http.StatusText(resp.StatusCode), err)
	}
	//log.Println(opName, resp.StatusCode, http.StatusText(resp.StatusCode), err)
	return resp
}
