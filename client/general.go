package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func (c *Client) Login(username, token string) error {
	body, err := json.Marshal(struct {
		Username string `json:"username"`
		Token    string `json:"api_token"`
	}{username, token})
	if err != nil {
		return err
	}

	resp, err := c.http.Post(c.urlStr+"/api/login", mimeJSON, bytes.NewBuffer(body))
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var data struct {
		Result        Status
		Session_token string
	}

	if err := json.Unmarshal(blob, &data); err != nil {
		return err
	}
	c.sessionToken = data.Session_token
	return data.Result.error()
}

func (c *Client) Logout() error {
	req, err := c.newRequest("GET", "/api/logout", nil)
	if err != nil {
		return err
	}
	resp, err := c.http.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var data struct {
		Result Status
	}
	if err := json.Unmarshal(blob, &data); err != nil {
		return err
	}
	return data.Result.error()
}

type Hello struct {
	Service    string
	ServerDate DateTime
	Version    string
	Language   string
	ClientIp   IPAddr
}

func (h *Hello) String() string {
	return fmt.Sprintf("{%s, %s, %s, %s, %s}",
		h.Service, h.ServerDate, h.Version, h.Language, h.ClientIp)
}

func (h *Hello) UnmarshalJSON(b []byte) error {
	var d struct {
		Result     *Status
		Service    string
		ServerDate DateTime `json:"server_date"`
		Version    string
		Language   string
		ClientIp   IPAddr `json:"client_ip"`
	}
	if err := json.Unmarshal(b, &d); err != nil {
		return err
	}
	*h = Hello{d.Service, d.ServerDate, d.Version, d.Language, d.ClientIp}
	return d.Result.error()
}

func (c *Client) Hello() (*Hello, error) {
	req, err := c.newRequest("GET", "/api/hello", nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data Hello
	if err := json.Unmarshal(blob, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
