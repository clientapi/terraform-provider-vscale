package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const DefaultBaseURL = "https://api.vscale.io/v1"

type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

func NewClient(token string) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) newRequest(method, path string, body interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.BaseURL, path)
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Token", c.Token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if v != nil && resp.StatusCode != http.StatusNoContent {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(bodyBytes, v)
	}

	return nil
}

// Account structs
type AccountInfo struct {
	ActDate    string `json:"actdate"`
	Country    string `json:"country"`
	Email      string `json:"email"`
	FaceID     string `json:"face_id"`
	ID         string `json:"id"`
	Locale     string `json:"locale"`
	Middlename string `json:"middlename"`
	Mobile     string `json:"mobile"`
	Name       string `json:"name"`
	State      string `json:"state"`
	Surname    string `json:"surname"`
}

type AccountResponse struct {
	Info AccountInfo `json:"info"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
	Bonus   float64 `json:"bonus"`
}

func (c *Client) GetAccountInfo() (*AccountInfo, error) {
	req, err := c.newRequest("GET", "/account", nil)
	if err != nil {
		return nil, err
	}
	var resp AccountResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp.Info, nil
}

func (c *Client) GetBalance() (*BalanceResponse, error) {
	req, err := c.newRequest("GET", "/billing/balance", nil)
	if err != nil {
		return nil, err
	}
	var resp BalanceResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// Background / Info structs
type Image struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Active      bool     `json:"active"`
	Size        int      `json:"size"`
	Locations   []string `json:"locations"`
	RPlans      []string `json:"rplans"`
}

func (c *Client) GetImages() ([]Image, error) {
	req, err := c.newRequest("GET", "/images", nil)
	if err != nil {
		return nil, err
	}
	var images []Image
	if err := c.do(req, &images); err != nil {
		return nil, err
	}
	return images, nil
}

type Location struct {
	ID                string   `json:"id"`
	Description       string   `json:"description"`
	Active            bool     `json:"active"`
	PrivateNetworking bool     `json:"private_networking"`
	Templates         []string `json:"templates"`
	RPlans            []string `json:"rplans"`
}

func (c *Client) GetLocations() ([]Location, error) {
	req, err := c.newRequest("GET", "/locations", nil)
	if err != nil {
		return nil, err
	}
	var locations []Location
	if err := c.do(req, &locations); err != nil {
		return nil, err
	}
	return locations, nil
}

type RPlan struct {
	ID        string   `json:"id"`
	Memory    int      `json:"memory"`
	CPUs      int      `json:"cpus"`
	Disk      int      `json:"disk"`
	Addresses int      `json:"addresses"`
	Locations []string `json:"locations"`
	Templates []string `json:"templates"`
}

func (c *Client) GetRPlans() ([]RPlan, error) {
	req, err := c.newRequest("GET", "/rplans", nil)
	if err != nil {
		return nil, err
	}
	var rplans []RPlan
	if err := c.do(req, &rplans); err != nil {
		return nil, err
	}
	return rplans, nil
}

type PriceResponse struct {
	Default map[string]interface{} `json:"default"`
	Period  string                 `json:"period"`
}

func (c *Client) GetPrices() (*PriceResponse, error) {
	req, err := c.newRequest("GET", "/billing/prices", nil)
	if err != nil {
		return nil, err
	}
	var resp PriceResponse
	if err := c.do(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// SSH Key structs
type SSHKey struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
	Key  string `json:"key"`
}

func (c *Client) CreateSSHKey(name, key string) (*SSHKey, error) {
	reqBody := map[string]string{
		"name": name,
		"key":  key,
	}
	req, err := c.newRequest("POST", "/sshkeys", reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var sshKey SSHKey
	if err := json.Unmarshal(bodyBytes, &sshKey); err == nil && sshKey.ID != 0 {
		return &sshKey, nil
	}

	var sshKeys []SSHKey
	if err := json.Unmarshal(bodyBytes, &sshKeys); err == nil && len(sshKeys) > 0 {
		for _, k := range sshKeys {
			if k.Name == name {
				return &k, nil
			}
		}
		return &sshKeys[0], nil
	}

	return nil, fmt.Errorf("failed to parse SSH key response: %s", string(bodyBytes))
}

func (c *Client) GetSSHKey(id int) (*SSHKey, error) {
	req, err := c.newRequest("GET", "/sshkeys", nil)
	if err != nil {
		return nil, err
	}

	var keys []SSHKey
	if err := c.do(req, &keys); err != nil {
		return nil, err
	}

	for _, k := range keys {
		if k.ID == id {
			return &k, nil
		}
	}
	return nil, fmt.Errorf("SSH key with ID %d not found", id)
}

func (c *Client) DeleteSSHKey(id int) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/sshkeys/%d", id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// Scalet (Server) structs
type Address struct {
	Netmask string `json:"netmask"`
	Gateway string `json:"gateway"`
	Address string `json:"address"`
}

type ScaletKey struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type Scalet struct {
	CTID           int         `json:"ctid,omitempty"`
	Name           string      `json:"name"`
	MakeFrom       string      `json:"made_from"`
	RPlan          string      `json:"rplan"`
	Location       string      `json:"location"`
	DoStart        bool        `json:"do_start"`
	Password       string      `json:"password,omitempty"`
	Keys           []ScaletKey `json:"keys,omitempty"`
	Status         string      `json:"status,omitempty"`
	Active         bool        `json:"active,omitempty"`
	Locked         bool        `json:"locked,omitempty"`
	Hostname       string      `json:"hostname,omitempty"`
	PublicAddress  *Address    `json:"public_address,omitempty"`
	PrivateAddress *Address    `json:"private_address,omitempty"`
}

type CreateScaletRequest struct {
	MakeFrom string `json:"make_from"`
	RPlan    string `json:"rplan"`
	DoStart  bool   `json:"do_start"`
	Name     string `json:"name"`
	Keys     []int  `json:"keys,omitempty"`
	Password string `json:"password,omitempty"`
	Location string `json:"location"`
}

func (c *Client) CreateScalet(reqData *CreateScaletRequest) (*Scalet, error) {
	req, err := c.newRequest("POST", "/scalets", reqData)
	if err != nil {
		return nil, err
	}

	var scalet Scalet
	if err := c.do(req, &scalet); err != nil {
		return nil, err
	}
	return &scalet, nil
}

func (c *Client) GetScalet(ctid int) (*Scalet, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/scalets/%d", ctid), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var scalet Scalet
	if err := json.Unmarshal(bodyBytes, &scalet); err == nil && scalet.CTID != 0 {
		return &scalet, nil
	}

	var scalets []Scalet
	if err := json.Unmarshal(bodyBytes, &scalets); err == nil && len(scalets) > 0 {
		return &scalets[0], nil
	}

	return nil, fmt.Errorf("failed to parse scalet response: %s", string(bodyBytes))
}

func (c *Client) DeleteScalet(ctid int) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/scalets/%d", ctid), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

func (c *Client) UpgradeScalet(ctid int, rplan string) (*Scalet, error) {
	reqBody := map[string]string{
		"rplan": rplan,
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/scalets/%d/upgrade", ctid), reqBody)
	if err != nil {
		return nil, err
	}
	var scalet Scalet
	if err := c.do(req, &scalet); err != nil {
		return nil, err
	}
	return &scalet, nil
}

func (c *Client) RebuildScalet(ctid int, password string) (*Scalet, error) {
	reqBody := map[string]string{
		"password": password,
	}
	req, err := c.newRequest("PATCH", fmt.Sprintf("/scalets/%d/rebuild", ctid), reqBody)
	if err != nil {
		return nil, err
	}
	var scalet Scalet
	if err := c.do(req, &scalet); err != nil {
		return nil, err
	}
	return &scalet, nil
}

func (c *Client) UpdateScaletSSHKeys(ctid int, keys []int) (*Scalet, error) {
	reqBody := map[string][]int{
		"keys": keys,
	}
	req, err := c.newRequest("PATCH", fmt.Sprintf("/scalets/%d", ctid), reqBody)
	if err != nil {
		return nil, err
	}
	var scalet Scalet
	if err := c.do(req, &scalet); err != nil {
		return nil, err
	}
	return &scalet, nil
}

// Domain structs
type Domain struct {
	ID         int    `json:"id,omitempty"`
	Name       string `json:"name"`
	BindZone   string `json:"bind_zone,omitempty"`
	CreateDate int64  `json:"create_date,omitempty"`
	ChangeDate int64  `json:"change_date,omitempty"`
	UserID     int    `json:"user_id,omitempty"`
}

func (c *Client) CreateDomain(name string, bindZone string) (*Domain, error) {
	reqBody := map[string]string{
		"name": name,
	}
	if bindZone != "" {
		reqBody["bind_zone"] = bindZone
	}
	req, err := c.newRequest("POST", "/domains", reqBody)
	if err != nil {
		return nil, err
	}

	var domain Domain
	if err := c.do(req, &domain); err != nil {
		return nil, err
	}
	return &domain, nil
}

func (c *Client) GetDomain(id int) (*Domain, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/domains/%d", id), nil)
	if err != nil {
		return nil, err
	}

	var domain Domain
	if err := c.do(req, &domain); err != nil {
		return nil, err
	}
	return &domain, nil
}

func (c *Client) DeleteDomain(id int) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/domains/%d", id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// Domain Record structs
type DomainRecord struct {
	ID       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	TTL      int    `json:"ttl"`
	Content  string `json:"content,omitempty"`
	Priority *int   `json:"priority,omitempty"`
	Weight   *int   `json:"weight,omitempty"`
	Port     *int   `json:"port,omitempty"`
	Target   string `json:"target,omitempty"`
	Email    string `json:"email,omitempty"`
}

func (c *Client) CreateDomainRecord(domainID int, record *DomainRecord) (*DomainRecord, error) {
	req, err := c.newRequest("POST", fmt.Sprintf("/domains/%d/records", domainID), record)
	if err != nil {
		return nil, err
	}

	var created DomainRecord
	if err := c.do(req, &created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (c *Client) GetDomainRecord(domainID, recordID int) (*DomainRecord, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/domains/%d/records/%d", domainID, recordID), nil)
	if err != nil {
		return nil, err
	}

	var record DomainRecord
	if err := c.do(req, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (c *Client) UpdateDomainRecord(domainID, recordID int, record *DomainRecord) (*DomainRecord, error) {
	req, err := c.newRequest("PUT", fmt.Sprintf("/domains/%d/records/%d", domainID, recordID), record)
	if err != nil {
		return nil, err
	}

	var updated DomainRecord
	if err := c.do(req, &updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (c *Client) DeleteDomainRecord(domainID, recordID int) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/domains/%d/records/%d", domainID, recordID), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// Domain Tag structs
type DomainTag struct {
	ID      int      `json:"id,omitempty"`
	Name    string   `json:"name"`
	Domains []string `json:"domains"`
}

func (c *Client) CreateDomainTag(name string, domains []string) (*DomainTag, error) {
	reqBody := map[string]interface{}{
		"name":    name,
		"domains": domains,
	}
	req, err := c.newRequest("POST", "/domains/tags", reqBody)
	if err != nil {
		return nil, err
	}
	var tag DomainTag
	if err := c.do(req, &tag); err != nil {
		return nil, err
	}
	return &tag, nil
}

func (c *Client) GetDomainTag(id int) (*DomainTag, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/domains/tags/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var tag DomainTag
	if err := c.do(req, &tag); err != nil {
		return nil, err
	}
	return &tag, nil
}

func (c *Client) UpdateDomainTag(id int, name string, domains []string) (*DomainTag, error) {
	reqBody := map[string]interface{}{
		"name":    name,
		"domains": domains,
	}
	req, err := c.newRequest("PUT", fmt.Sprintf("/domains/tags/%d", id), reqBody)
	if err != nil {
		return nil, err
	}
	var tag DomainTag
	if err := c.do(req, &tag); err != nil {
		return nil, err
	}
	return &tag, nil
}

func (c *Client) DeleteDomainTag(id int) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/domains/tags/%d", id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// PTR Record structs
type PTRRecord struct {
	ID      int    `json:"id,omitempty"`
	IP      string `json:"ip"`
	Content string `json:"content"`
	UserID  int    `json:"user_id,omitempty"`
}

func (c *Client) CreatePTRRecord(ip, content string) (*PTRRecord, error) {
	reqBody := map[string]string{
		"ip":      ip,
		"content": content,
	}
	req, err := c.newRequest("POST", "/domains/ptr", reqBody)
	if err != nil {
		return nil, err
	}
	var record PTRRecord
	if err := c.do(req, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (c *Client) GetPTRRecord(id int) (*PTRRecord, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/domains/ptr/%d", id), nil)
	if err != nil {
		return nil, err
	}
	var record PTRRecord
	if err := c.do(req, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (c *Client) UpdatePTRRecord(id int, ip, content string) (*PTRRecord, error) {
	reqBody := map[string]string{
		"ip":      ip,
		"content": content,
	}
	req, err := c.newRequest("PUT", fmt.Sprintf("/domains/ptr/%d", id), reqBody)
	if err != nil {
		return nil, err
	}
	var record PTRRecord
	if err := c.do(req, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

func (c *Client) DeletePTRRecord(id int) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/domains/ptr/%d", id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// Backup structs
type Backup struct {
	ID       string `json:"id"`
	Template string `json:"template"`
	Active   bool   `json:"active"`
	Name     string `json:"name"`
	ScaletID int    `json:"scalet"`
	Status   string `json:"status"`
	Size     int    `json:"size"`
	Locked   bool   `json:"locked"`
	Location string `json:"location"`
	Created  string `json:"created"`
}

func (c *Client) CreateBackup(scaletID int, name string) (*Backup, error) {
	reqBody := map[string]string{
		"name": name,
	}
	req, err := c.newRequest("POST", fmt.Sprintf("/scalets/%d/backup", scaletID), reqBody)
	if err != nil {
		return nil, err
	}
	var backup Backup
	if err := c.do(req, &backup); err != nil {
		return nil, err
	}
	return &backup, nil
}

func (c *Client) GetBackup(id string) (*Backup, error) {
	req, err := c.newRequest("GET", fmt.Sprintf("/backups/%s", id), nil)
	if err != nil {
		return nil, err
	}
	var backup Backup
	if err := c.do(req, &backup); err != nil {
		return nil, err
	}
	return &backup, nil
}

func (c *Client) GetBackups() ([]Backup, error) {
	req, err := c.newRequest("GET", "/backups", nil)
	if err != nil {
		return nil, err
	}
	var backups []Backup
	if err := c.do(req, &backups); err != nil {
		return nil, err
	}
	return backups, nil
}

func (c *Client) DeleteBackup(id string) error {
	req, err := c.newRequest("DELETE", fmt.Sprintf("/backups/%s", id), nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
