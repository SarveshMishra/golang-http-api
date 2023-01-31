package main

import(
	"log"
	"net/http"
	"time"
)
type timeHandler struct {
	format string
}

func (th timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}
func main (){
	// use the http.NewServeMux() function to create an empty servemux.
	mux := http.NewServeMux()

	rh := http.RedirectHandler("https://sarvesh.xyz", 307)
	th := timeHandler{format: time.RFC1123}
	mux.Handle("/foo", rh)
	mux.Handle("/time", th)

	log.Print("listening...")
	http.ListenAndServe(":3000", mux)

}
