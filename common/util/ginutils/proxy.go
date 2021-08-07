package ginutils

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func ReverseProxy(target string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//director := func(req *http.Request) {
		//	r := c.Request
		//	req = r
		//	req.URL.Scheme = "http"
		//	req.URL.Host = target
		//	req.Host = ""
		//	//req.Header["my-header"] = []string{r.Header.Get("my-header")}
		//	//// Golang camelcases headers
		//	//delete(req.Header, "My-Header")
		//}
		//
		//fmt.Println(target)
		//proxy := &httputil.ReverseProxy{Director: director}
		//proxy.ServeHTTP(c.Writer, c.Request)

		if strings.Index(target, "http") != 0 {
			target = "https://" + target
		}

		url, _ := url.Parse(target)
		//proxy := httputil.NewSingleHostReverseProxy(url)
		proxy := ReverseProxyCustomV2(url)
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func ReverseProxyCustomV2(target *url.URL) *httputil.ReverseProxy {
	targetQuery := target.RawQuery
	director := func(req *http.Request) {
		//req.URL.Scheme = target.Scheme
		//req.URL.Host = target.Host
		//req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		req.URL = target
		req.Host = target.Host
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
	}
	return &httputil.ReverseProxy{Director: director}
}

//func ReverseProxyCustomV1(target string) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		reverseProxyCustomV1(target, c.Writer, c.Request)
//	}
//}
func reverseProxyCustomV1(target string, w http.ResponseWriter, r *http.Request) {

	uri := target + r.RequestURI

	fmt.Println(r.Method + ": " + uri)

	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		fatal(err)
		fmt.Printf("Body: %v\n", string(body))
	}

	rr, err := http.NewRequest(r.Method, uri, r.Body)
	fatal(err)
	copyHeader(r.Header, &rr.Header)

	// Create a client and query the target
	var transport http.Transport
	resp, err := transport.RoundTrip(rr)
	fatal(err)

	fmt.Printf("Resp-Headers: %v\n", resp.Header)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fatal(err)

	dH := w.Header()
	copyHeader(resp.Header, &dH)
	dH.Add("Requested-Host", rr.Host)

	w.Write(body)
}
func fatal(err error) {
	if err != nil {
		logrus.Fatal(err)
		os.Exit(1)
	}
}
func copyHeader(source http.Header, dest *http.Header) {
	for n, v := range source {
		for _, vv := range v {
			dest.Add(n, vv)
		}
	}
}
