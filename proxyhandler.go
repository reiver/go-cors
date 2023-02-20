package cors

import (
	"fmt"
	"io"
	"net/http"
)

// ProxyHandler is used to create a CORS proxy.
//
// If ProxyHandler is being used on the Internet domain ‘example.com’, then a request to:
//
//	http://example.com/http://something.tld/blog/feed.atom
//
// Would make a server-side request to:
//
//	http://something.tld/blog/feed.atom
//
// And request the response with the addition of two headers:
//
//	Access-Control-Allow-Origin: *
//	Access-Control-Allow-Methods: GET, DELETE, HEAD, OPTIONS, PATCH, POST, PUT, TRACE
//
// Example usage:
//
//	var handler http.Handler = &cors.ProxyHandler{
//		LogWriter: os.Stdout,
//	}
//	
//	err := http.ListenAndServe(addr, handler)
type ProxyHandler struct {
	LogWriter io.Writer
}

func (receiver *ProxyHandler) logProxyRequest(r *http.Request) {
	var w io.Writer = receiver.LogWriter
	if nil == w {
		return
	}

	fmt.Fprintf(w, "PROXY REQUEST:  %s http://%s%s\n", r.Method, r.Host, r.URL.RequestURI())
}

func (receiver *ProxyHandler) logProxyResponse(resp *http.Response) {
	var w io.Writer = receiver.LogWriter
	if nil == w {
		return
	}

	fmt.Fprintf(w, "PROXY RESPONSE: %s %s http://%s%s\n", resp.Status, resp.Request.Method, resp.Request.Host, resp.Request.URL.RequestURI())
}

func (receiver *ProxyHandler) logRequest(r *http.Request) {
	var w io.Writer = receiver.LogWriter
	if nil == w {
		return
	}

	fmt.Fprintf(w, "CLIENT REQUEST: %s http://%s%s\n", r.Method, r.Host, r.URL.RequestURI())
}

func (receiver *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if nil == w {
/////////////// RETURN
		return
	}
	if nil == r {
		http.Error(w, "WTF?!?!? — Internal Server Error", http.StatusInternalServerError)
/////////////// RETURN
		return
	}

	receiver.logRequest(r)

	switch r.Method {
	case http.MethodOptions:
		receiver.serveOptions(w, r)
	default:
		receiver.serveHTTP(w, r)
	}
}

func (receiver *ProxyHandler) serveHTTP(w http.ResponseWriter, r *http.Request) {

	var method string = r.Method

	var url string
	{
		url = r.URL.RequestURI()
		if 1 <= len(url) && '/' == url[0] {
			url = url[1:]
		}
	}

	var req *http.Request
	{
		var err error

		req, err = http.NewRequest(method, url, nil)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
/////////////////////// RETURN
			return
		}
	}
	receiver.logProxyRequest(req)

	var resp *http.Response
	{
		var client http.Client

		var err error

		resp, err = client.Do(req)
		if nil != err {
			http.Error(w, err.Error(), http.StatusInternalServerError)
/////////////////////// RETURN
			return
		}
	}
	receiver.logProxyResponse(resp)

	addCORSHeaders(resp.Header)
	resp.Write(w)
}

func (*ProxyHandler) serveOptions(w http.ResponseWriter, r *http.Request) {
	addCORSHeaders(w.Header())
	w.WriteHeader(http.StatusNoContent)
}

func addCORSHeaders(header http.Header) {
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", "GET, DELETE, HEAD, OPTIONS, PATCH, POST, PUT, TRACE")
}
