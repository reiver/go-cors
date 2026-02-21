package cors

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	// Let you chose which HTTP methods are in the value of the HTTP "Access-Control-Allow-Methods" response header.
	//
	// If left empty, defaults to:
	//
	//	[]string{"GET", "DELETE", "HEAD", "OPTIONS", "PATCH", "POST", "PUT", "TRACE"}
	//
	// To include WebDAV HTTP request methods, use:
	//
	//	[]string{"COPY", "GET", "DELETE", "HEAD","LOCK", "MKCOL", "MOVE", "OPTIONS", "PATCH", "POST", "PROPFIND", "PROPPATCH", "PUT", "TRACE", "UNLOCK"}
	//
	// Remember that HTTP request methods are case-sensitive.
	// I.e., "GET" and "get" are, according to the HTTP specifications, different methods.
	// (Most HTTP request methods tends to be upper-case.)
	Methods   []string

	// Where logs gets written to.
	LogWriter io.Writer
}

func (receiver *ProxyHandler) logProxyRequest(r *http.Request) {
	if nil == receiver {
/////////////// RETURN
		return
	}

	var w io.Writer = receiver.LogWriter
	if nil == w {
/////////////// RETURN
		return
	}

	fmt.Fprintf(w, "PROXY REQUEST:  %s http://%s%s\n", r.Method, r.Host, r.URL.RequestURI())
}

func (receiver *ProxyHandler) logProxyResponse(resp *http.Response) {
	if nil == receiver {
/////////////// RETURN
		return
	}

	var w io.Writer = receiver.LogWriter
	if nil == w {
/////////////// RETURN
		return
	}

	fmt.Fprintf(w, "PROXY RESPONSE: %s %s http://%s%s\n", resp.Status, resp.Request.Method, resp.Request.Host, resp.Request.URL.RequestURI())
}

func (receiver *ProxyHandler) logRequest(r *http.Request) {
	if nil == receiver {
/////////////// RETURN
		return
	}

	var w io.Writer = receiver.LogWriter
	if nil == w {
/////////////// RETURN
		return
	}

	fmt.Fprintf(w, "CLIENT REQUEST: %s http://%s%s\n", r.Method, r.Host, r.URL.RequestURI())
}

// ServeHTTP makes [ProxyHandler] an [http.Handler].
func (receiver *ProxyHandler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if nil == receiver {
/////////////// RETURN
		return
	}

	if nil == rw {
/////////////// RETURN
		return
	}
	if nil == r {
		http.Error(rw, "WTF?!?!? — Internal Server Error", http.StatusInternalServerError)
/////////////// RETURN
		return
	}

	receiver.logRequest(r)

	switch r.Method {
	case http.MethodOptions:
		receiver.serveOptions(rw, r)
	default:
		receiver.serveHTTP(rw, r)
	}
}

func (receiver *ProxyHandler) serveHTTP(rw http.ResponseWriter, r *http.Request) {
	if nil == receiver {
		return
	}

	var proxiedMethod string = r.Method

	var proxiedURL string
	{
		proxiedURL = r.URL.RequestURI()
		if 1 <= len(proxiedURL) && '/' == proxiedURL[0] {
			proxiedURL = proxiedURL[1:]
		}
	}

	var proxiedRequest *http.Request
	{
		var err error

		proxiedRequest, err = http.NewRequest(proxiedMethod, proxiedURL, r.Body)
		if nil != err {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
/////////////////////// RETURN
			return
		}
	}
	receiver.logProxyRequest(proxiedRequest)

	var proxiedResponse *http.Response
	{
		var client http.Client

		var err error

		proxiedResponse, err = client.Do(proxiedRequest)
		if nil != err {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
/////////////////////// RETURN
			return
		}
	}
	receiver.logProxyResponse(proxiedResponse)

	var response *http.Response
	{
		var respBuffer strings.Builder
		err := proxiedResponse.Write(&respBuffer)
		if nil != err {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
/////////////////////// RETURN
			return
		}

		var reader io.Reader = strings.NewReader(respBuffer.String())
		if nil == reader {
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
/////////////////////// RETURN
			return
		}

		var bufReader *bufio.Reader = bufio.NewReader(reader)
		if nil == bufReader {
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
/////////////////////// RETURN
			return
		}

		response, err = http.ReadResponse(bufReader, proxiedRequest)
		if nil != err {
			http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
/////////////////////// RETURN
			return
		}
	}

	addCORSHeaders(response.Header)

	for name, values := range response.Header {

		for _, value := range values {

			rw.Header().Add(name, value)
		}
	}

	rw.WriteHeader(response.StatusCode)
	io.Copy(rw, response.Body)
	response.Body.Close()
}

// serveOptions serves the HTTP response to an "OPTIONS" request.
func (receiver *ProxyHandler) serveOptions(rw http.ResponseWriter, r *http.Request) {
	if nil == receiver {
		return
	}

	addCORSHeaders(rw.Header(), receiver.Methods...)
	rw.WriteHeader(http.StatusNoContent)
}

// addCORSHeaders ads CORS headers to an [http.Header].
//
// If the HTTP methods specified are empty, it defaults to:
//
//	"GET, DELETE, HEAD, OPTIONS, PATCH, POST, PUT, TRACE"
//
// See also: [allowedMethods]
func addCORSHeaders(header http.Header, methods ...string) {
	header.Add("Access-Control-Allow-Origin", "*")
	header.Add("Access-Control-Allow-Methods", allowedMethods(methods...))
}

// allowedMethods creates the value for the "Access-Control-Allow-Methods" HTTP response header.
func allowedMethods(methods ...string) string {
	if len(methods) <= 0 {
		return  "GET, DELETE, HEAD, OPTIONS, PATCH, POST, PUT, TRACE"
	}

	return strings.Join(methods, ", ")
}
