package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"tailscale.com/tsnet"
)

var (
	tailnetName   string
	reverseServer string
)

func init() {
	flag.StringVar(&tailnetName, "n", "", "tail net name")
	flag.StringVar(&reverseServer, "r", "", "reverse server address")
	flag.Parse()
}

// NewProxy takes target host and creates a reverse proxy
// NewProxy 拿到 targetHost 后，创建一个反向代理
func NewProxy(targetHost string) (*httputil.ReverseProxy, error) {
	url, err := url.Parse(targetHost)
	if err != nil {
		return nil, err
	}

	return httputil.NewSingleHostReverseProxy(url), nil
}

// ProxyRequestHandler handles the http request using proxy
// ProxyRequestHandler 使用 proxy 处理请求
func ProxyRequestHandler(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}
}

func main() {
	s := &tsnet.Server{Hostname: tailnetName}
	defer s.Close()

	ln, err := s.ListenFunnel("tcp", ":443") // does TLS
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	proxy, err := NewProxy("http://my-api-server.com")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", ProxyRequestHandler(proxy))

	//h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintln(w, "<html>Hello from Funnel!")
	//})

	log.Fatal(http.Serve(ln, nil))
}
