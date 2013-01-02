package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type WhoisPrivacyStatus struct {
	Enabled bool
	Expired DateTime `json:"expire_date"`
}

type Addon struct {
	Price Currency
}

func (a Addon) String() string {
	return fmt.Sprintf("%.2f", a.Price)
}

type Domain struct {
	TLD          string
	Created      DateTime           `json:"create_date"`
	Expired      DateTime           `json:"expire_date"`
	WhoisPrivacy WhoisPrivacyStatus `json:"whois_privacy"`
	Addons       map[string]Addon
}

func (d *Domain) String() string {
	return fmt.Sprintf("{%s, %s, %s, %v, %s}",
		d.TLD, d.Created, d.Expired, d.WhoisPrivacy, d.Addons)
}

type domains map[string]Domain

func (d *domains) UnmarshalJSON(data []byte) error {
	// workaround: name.com returns empty array when no domains on account
	if len(data) == 2 && data[0] == '[' && data[1] == ']' {
		return nil
	}
	var m map[string]Domain
	err := json.Unmarshal(data, &m)
	*d = m
	return err
}

func (c *EndPoint) ListDomains() (map[string]Domain, error) {
	req, err := c.newRequest("GET", "/api/domain/list", nil)
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
	var data struct {
		Result  *Status
		Domains domains
	}
	if err := json.Unmarshal(blob, &data); err != nil {
		return nil, err
	}
	return data.Domains, data.Result.error()
}
