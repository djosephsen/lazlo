package lib

import (
	"net/http"
	"fmt"
)

func(b *Broker)StartHttp(){
	http.HandleFunc("/", httpHi)
	err := http.ListenAndServe(":"+b.Config.Port, nil)
	if err != nil {
		Logger.Error(err)
	}
}

func httpHi(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(res, "Hi. I am Lazlo")
}
