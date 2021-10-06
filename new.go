package doxyproxy

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

//Proxy is a doxyproxy config
type Proxy struct {
	API string
	App string
	Key string

	CacheTTL time.Duration
	cache    map[string]IPEntry
	mutex    sync.RWMutex
	Client   *http.Client
}

//IPEntry holds the true IP address of a client
type IPEntry struct {
	ID string

	IP     string
	Expire int64

	proxy *Proxy
}

//New creates a new DoxyProxy
func New(api string, app string, key string) *Proxy {

	return &Proxy{
		API: fmt.Sprintf("%s/%s/", api, url.PathEscape(app)),
		App: app,
		Key: key,

		cache:    make(map[string]IPEntry),
		CacheTTL: time.Minute * 5,
		Client:   &http.Client{},
	}
}
