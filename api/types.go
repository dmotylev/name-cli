package api

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	mimeJSON   = "application/json"
	timeLayout = "2006-01-02 15:04:05"
)

var NameComLocation = mustLoad(time.LoadLocation("US/Mountain"))

func mustLoad(l *time.Location, err error) *time.Location {
	if err != nil {
		panic(err)
	}
	return l
}

type Currency float32

func (c *Currency) UnmarshalJSON(data []byte) error {
	f, err := strconv.ParseFloat(strings.Trim(string(data), "\""), 32)
	if err != nil {
		return err
	}
	*c = Currency(float32(f))
	return nil
}

type IPAddr net.IP

func (a *IPAddr) UnmarshalJSON(data []byte) error {
	*a = IPAddr(net.ParseIP(strings.Trim(string(data), "\"")))
	return nil
}

func (a IPAddr) String() string {
	return net.IP(a).String()
}

type DateTime time.Time

func (n *DateTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse(`"`+timeLayout+`"`, string(data))
	// workaround: name.com does not provide TZ info, but they work in US/Mountain
	*n = DateTime(time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond(),
		NameComLocation).Local())
	return err
}

func (d DateTime) String() string {
	return time.Time(d).String()
}

func (d DateTime) Format(layout string) string {
	return time.Time(d).Format(layout)
}

// 100 command successful
// 203 required parameter missing
// 204 parameter value error
// 211 invalid command url
// 221 authorization error
// 240 command failed
// 250 unexpected error (exception)
// 251 authentication error
// 260 insufficient funds
// 261 unable to authorize funds
type Status struct {
	Code    int
	Message string
}

func (s *Status) Error() string {
	return fmt.Sprintf("%d : %s", s.Code, s.Message)
}

func (s *Status) error() error {
	if s.Code == 100 {
		return nil
	}
	return s
}

type EndPoint struct {
	http         *http.Client
	urlStr       string
	sessionToken string
}

func NewEndPoint(url string) *EndPoint {
	return &EndPoint{
		urlStr: url,
		http: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		},
	}
}

// Creates a new http request object and set value on session header
func (c *EndPoint) newRequest(method, uriStr string, body io.Reader) (*http.Request, error) {
	r, err := http.NewRequest(method, c.urlStr+uriStr, body)
	if err != nil {
		return nil, err
	}
	if c.sessionToken == "" {
		return r, nil
	}
	r.Header.Set("Api-Session-Token", c.sessionToken)
	return r, nil
}
