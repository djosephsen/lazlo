package lib

import (
	"net/http"
	"fmt"
)

type LinkCallback struct{
   ID           string
   Path         string
   Handler      func(res http.ResponseWriter, req *http.Request)
   Chan         chan *http.Request
}


func(b *Broker)StartHttp(){
	http.HandleFunc("/", httpHi)
	err := http.ListenAndServe(":"+b.Config.Port, nil)
	if err != nil {
		Logger.Error(err)
	}
}

func httpHi(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Hi. I am a Lazlo bot")
}

func newHttpPath(cb *LinkCallback) error{
	// assign a default handler if the user hasn't given us one
	if cb.Handler == nil{
		cb.Handler = func(res http.ResponseWriter, req *http.Request){
			cb.Chan <- req
			fmt.Fprintln(res, "Path: %s handled thanks!\n", cb.Path)
		}
	}
	//register with the the httpd 
	http.HandleFunc(cb.Path, cb.Handler)
	return nil
}

func (b *Broker) LinkCallback(path string, f ...func(http.ResponseWriter, *http.Request)) *LinkCallback{
   callback := &LinkCallback{
      ID:         fmt.Sprintf("link:%d",len(b.cbIndex[L])),
      Path:       path,
      Chan:       make(chan *http.Request),
   }

   //user-provided http handler function
   if f[0] != nil{ callback.Handler=f[0] }

   if err := newHttpPath(callback); err !=nil{
      Logger.Error("error registering callback ", callback.ID, ":: ",err)
      return nil
   }

   if err := b.RegisterCallback(callback); err != nil{
      Logger.Error("error registering callback ", callback.ID, ":: ",err)
      return nil
   }
	return callback
}
