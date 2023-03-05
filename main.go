package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	// "io"
	// "io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type timeHandler struct {
	format string
}
type bhriguData struct {
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

const (
	CATEGORY        string = "cat"
	SOURCE          string = "src"
	ACTION          string = "act"
	LABEL           string = "lbl"
	PAGEINFO        string = "pi"
	QUERYSTRING     string = "qs"
	REFERRER        string = "ref"
	COOKIEID        string = "cid"
	SESSIONID       string = "sid"
	EVENTTIMESTAMP  string = "ts"
	BHRIGUCOOKIES   string = "bhrigu_cookies" // this should be shorter in length
	RANDOM_STRING   string = "RS"
	BHRIGUCOOKIESV2 string = "bco"
	COOKIE          string = "cookie"
	// below are compulsory params
	PAGETYPE      string = "pt"
	APPLICATIONID string = "aid"
	PLATFORMID    string = "pid"
)

func (th timeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(th.format)
	w.Write([]byte("The time is: " + tm))
}
func main() {
	// use the http.NewServeMux() function to create an empty servemux.
	mux := http.NewServeMux()
	port := ":8089"
	rh := http.RedirectHandler("https://sarvesh.xyz", 307)
	// thf := http.HandlerFunc(timeHandlerFunc)
	mux.Handle("/foo", rh)
	mux.Handle("/time", timeHandler{format: time.RFC1123})
	mux.HandleFunc("/timeFunc", timeHandlerFunc)
	mux.HandleFunc("/beacon", postDataHandler)
	mux.HandleFunc("/get", getRequestHandler)
	mux.HandleFunc("/bhrigu.gif", getOldEventFromNewEvent)

	log.Printf("listening... on %v", port)
	http.ListenAndServe(port, mux)

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

	var t bhriguData
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

// getOldEventFromNewEvent handles post request without payload
func getOldEventFromNewEvent(w http.ResponseWriter, r *http.Request) {
	// log.Println(r)
	var err error
	var mandatoryData map[string]string
	var trackingData map[string]string
switch r.Method {
	case "OPTION":
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	case "GET":
		mandatoryData, trackingData, err = getDataFromQS(r)
	case "POST":
		mandatoryData, trackingData, err = getDataFromBody(r)
	}
	// r.url.Values.Add(r.URL.Query(), "test", "new")
	var bhrigu_cookies = ""
	params := make(url.Values)
	for key, value := range trackingData {
		if(key == "CWC" || key == "_cwv"){
			bhrigu_cookies += key + "=" + value + "; "
			continue
		}
		params.Add(key, value)
	}
	params.Add("bhrigu_cookies", bhrigu_cookies)
	r.URL.Path = "https://bhrigu.carwale.com/bhrigu.gif?" + params.Encode()
	log.Println(mandatoryData)
	log.Println(trackingData)
	log.Println(r.URL.Path)
	log.Println(err)
}

// getDataFromQS return tracking data from query string url for GET
func getDataFromQS(r *http.Request) (map[string]string,map[string]string, error) {
	dataMap := make(map[string]string)
	mandatoryData := make(map[string]string)
	queries := r.URL.Query()
	for key, value := range queries {
		if key == PAGETYPE || key == APPLICATIONID || key == PLATFORMID {
			mandatoryData[key] = value[len(value)-1]
			continue
		}
		dataMap[key] = value[len(value)-1]
	}
	return mandatoryData, dataMap, nil
}

// getDataFromBody return tracking data from POST request body
func getDataFromBody(r *http.Request) (map[string]string, map[string]string, error) {
	var dataMap = map[string]string{}
	var mandatoryData = map[string]string{}
	if r.Body != nil {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return mandatoryData, dataMap, fmt.Errorf("[getDataFromBody] reading request body failed, request object is and error is : %v", err)
		}
		err = json.Unmarshal(body, &dataMap)
		if err != nil {
			return mandatoryData, dataMap, fmt.Errorf("[getDataFromBody] unmarshalling to json failed, request object is  and error is : %v", err)
		}
		// extracting validation data from datamap to avoid duplication in tracking data
		for key, value := range dataMap {
			if key == PAGETYPE || key == APPLICATIONID || key == PLATFORMID {
				mandatoryData[key] = value
				delete(dataMap, key)
			}
		}
		return mandatoryData, dataMap, nil
	}
	err := fmt.Errorf("request does not contain payload in body")
	return mandatoryData, dataMap, err
}
