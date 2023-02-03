package main

import (
	"encoding/json"
	// "io"
	// "io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type timeHandler struct {
	format string
}
type test_struct struct {
	Link string
}
type Event struct {
	LogDateTime time.Time
	CookieID    string
	SessionID   string
	Category    string
	Action      string
	Label       string
	Pageinfo    string
	Query       string
	Cookie      string
	IP          string
	UserAgent   string
	Referrer    string
	Source      string
}

// Keys in query string
const (
	CATEGORY       string = "cat"
	SOURCE         string = "src"
	ACTION         string = "act"
	LABEL          string = "lbl"
	PAGEINFO       string = "pi"
	QUERYSTRING    string = "qs"
	REFERRER       string = "ref"
	COOKIEID       string = "cid"
	SESSIONID      string = "sid"
	EVENTTIMESTAMP string = "ts"
	BHRIGUCOOKIES  string = "bhrigu_cookies" // this should be shorter in length
	RANDOM_STRING  string = "RS"
)

func (th timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}
func main() {
	// use the http.NewServeMux() function to create an empty servemux.
	mux := http.NewServeMux()

	rh := http.RedirectHandler("https://sarvesh.xyz", 307)
	// thf := http.HandlerFunc(timeHandlerFunc)
	mux.Handle("/foo", rh)
	mux.Handle("/time", timeHandler{format: time.RFC1123})
	mux.HandleFunc("/timeFunc", timeHandlerFunc)
	mux.HandleFunc("/beacon", postDataHandler)
	mux.HandleFunc("/get", getRequestHandler)

	log.Print("listening...")
	http.ListenAndServe(":3001", mux)

}
func timeHandlerFunc(w http.ResponseWriter, r *http.Request) {
	log.Println(strings.Split(r.Host, "."))
	qs := r.URL.Query()
	// fmt.Print(r.URL.RawQuery)
	data := Event{
		LogDateTime: time.Now().UTC().Round(time.Millisecond),
		IP:          strings.TrimSpace(getIP(r)),
		UserAgent:   strings.TrimSpace(r.UserAgent()),
		SessionID:   strings.TrimSpace(qs.Get(SESSIONID)),
	}
	log.Println(data)
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("The time is: " + tm))
}
func getIP(r *http.Request) string {
	if ipProxy := r.Header.Get("X-Forwarded-For"); len(ipProxy) > 0 {
		return ipProxy
	} else if ipProxy := r.Header.Get("Client-IP"); len(ipProxy) > 0 {
		return ipProxy
	} else if ipProxy := r.Header.Get("X-Original-Forwarded-For"); len(ipProxy) > 0 {
		return ipProxy
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// postDataHandler handles json body body request in POST request
func postDataHandler(w http.ResponseWriter, r *http.Request) {
	// maxBytesSize := 10
	// r.Body = http.MaxBytesReader(w, r.Body, maxBytesSize)

	// check if request body is not too large
	// bodySize, err := ioutil.ReadAll(r.Body)
	// log.Printf("Request body size: %v\n",len(bodySize))
	// if err != nil {
	// 	log.Println(err)
	// }
	// if len(data) >= maxBytesSize {
	//      //exceeded
	// }
	// some other error
	for _, c := range r.Cookies() {
		log.Println(c)
	}
	decoder := json.NewDecoder(r.Body)

	var t test_struct
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	log.Println(t)
	var newUserSessionCookie = &http.Cookie{
		Name:     "TestCookie",
		Value:    "kljfajfkljklafjejfkjka",
		Path:     "/",
		Expires:  time.Now().Add(5 * 365 * 24 * time.Hour),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Domain:   "localhost",
	}

	http.SetCookie(w, newUserSessionCookie)

}

// getRequestHandler handles get request in long query url
func getRequestHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r)
	for _, c := range r.Cookies() {
		log.Println(c)
	}
	var newUserSessionCookie = &http.Cookie{
		Name:     "TestCookie",
		Value:    "kljfajfkljklafjejfkjka",
		Path:     "/",
		Expires:  time.Now().Add(5 * 365 * 24 * time.Hour),
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		Domain:   "localhost",
	}

	http.SetCookie(w, newUserSessionCookie)
}

// postRequestHandlerWithoutPayload handles post request without payload
func postRequestHandlerWithoutPayload(w http.ResponseWriter, r *http.Request) {

}
