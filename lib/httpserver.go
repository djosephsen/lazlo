package lib

import (
	"fmt"
	"github.com/bmizerany/pat"
	"net/http"
)

//This map is redundant with the brokers cbIndex but there's
//no way to get a reference to the broker into metaHandler so kludge
var httpRoutes = make(map[string]*LinkCallback)

type LinkCallback struct {
	p       string // the httpRoutes index value
	ID      string
	Path    string // the computed URL
	URL     string
	Handler func(res http.ResponseWriter, req *http.Request)
	Chan    chan *http.Request
}

func (b *Broker) StartHttp() {
	m := pat.New()
	m.Get("/", http.HandlerFunc(metaHandler))
	m.Get("/linkcb/:name", http.HandlerFunc(metaHandler))
	http.Handle("/", m)
	err := http.ListenAndServe(":"+b.Config.Port, nil)
	if err != nil {
		Logger.Error(err)
	}
}

func metaHandler(res http.ResponseWriter, req *http.Request) {
	Logger.Debug("entered metaHandler")
	path := req.URL.Query().Get(":name")
	if path == `` {
		Logger.Debug("path is /")
		fmt.Fprintln(res, "Hi. I am a Lazlo bot")
	} else if cb, ok := httpRoutes[path]; ok {
		Logger.Debug("path is known")
		if cb.Handler == nil {
			go func(cb *LinkCallback) {
				cb.Chan <- req
				fmt.Fprintln(res, "Path: %s handled. Thanks!", path)
			}(cb)
		} else {
			cb.Handler(res, req)
		}
	} else {
		Logger.Debug("path is unknown (", path, ")")
		fmt.Fprintf(res, "sorry, no modules have registered to handle %s\n", path)
	}
}

func (b *Broker) LinkCallback(p string, f ...func(http.ResponseWriter, *http.Request)) *LinkCallback {
	path := fmt.Sprintf("linkcb/%s", p)
	callback := &LinkCallback{
		p:    p,
		ID:   fmt.Sprintf("link:%d", len(b.cbIndex[L])),
		Path: path,
		URL:  fmt.Sprintf("%s:%s/%s", b.Config.URL, b.Config.Port, path),
		Chan: make(chan *http.Request),
	}

	//user-provided http handler function
	if f != nil {
		callback.Handler = f[0]
	}

	//append the path to the list of routes used by metaHandler()
	httpRoutes[p] = callback

	if err := b.RegisterCallback(callback); err != nil {
		Logger.Error("error registering callback ", callback.ID, ":: ", err)
		return nil
	}
	return callback
}

func (l *LinkCallback) Delete() {
	delete(httpRoutes, l.p)
}
