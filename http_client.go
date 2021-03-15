package porkbun

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"gopkg.in/resty.v1"
)

const SUCCESS = "SUCCESS"

type retrieveRecordsRequest struct {
	APIKey       string `json:"apikey,omitempty"`
	SecretAPIKey string `json:"secretapikey,omitempty"`
}

type retrieveRecordsResponse struct {
	Status  string    `json:"status,omitempty"`
	Records []*record `json:"records,omitempty"`
	Message string    `json:"message,omitempty"`
}

type record struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Type    string `json:"type,omitempty"`
	Content string `json:"content,omitempty"`
	TTL     time.Duration
	TTLStr  string `json:"ttl,omitempty"`
	Prio    string `json:"prio,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

func (p *Provider) doRetrieve(ctx context.Context, domain string, req retrieveRecordsRequest) (*retrieveRecordsResponse, error) {
	resp, err := resty.R().SetContext(ctx).SetBody(req).Post(fmt.Sprintf("https://porkbun.com/api/json/v3/dns/retrieve/%s", domain))
	if err != nil {
		return nil, err
	}
	r := &retrieveRecordsResponse{}
	if err := json.Unmarshal(resp.Body(), r); err != nil {
		return nil, err
	}
	if r.Status != SUCCESS {
		return nil, fmt.Errorf("API Status: %s", r.Status)
	}
	for _, rec := range r.Records {
		t, err := strconv.Atoi(rec.TTLStr)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse TTL, got: %q, err: %v", rec.TTLStr, err)
		}
		rec.TTL = time.Second * time.Duration(t)
	}
	return r, nil
}

type deleteRecordRequest struct {
	APIKey       string `json:"apikey,omitempty"`
	SecretAPIKey string `json:"secretapikey,omitempty"`
}

type deleteRecordResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func (p *Provider) doDelete(ctx context.Context, domain, id string, req deleteRecordRequest) (*deleteRecordResponse, error) {
	resp, err := resty.R().SetContext(ctx).SetBody(req).Post(fmt.Sprintf("https://porkbun.com/api/json/v3/dns/delete/%s/%s", domain, id))
	if err != nil {
		return nil, err
	}
	r := &deleteRecordResponse{}
	if err := json.Unmarshal(resp.Body(), r); err != nil {
		return nil, err
	}
	if r.Status != SUCCESS {
		return nil, fmt.Errorf("API call failed: %v", r.Message)
	}
	return r, nil
}

type createRecordRequest struct {
	APIKey       string `json:"apikey,omitempty"`
	SecretAPIKey string `json:"secretapikey,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	Content      string `json:"content,omitempty"`
	TTL          time.Duration
	TTLStr       string `json:"ttl,omitempty"`
	Prio         string `json:"prio,omitempty"`
}

type createRecordResponse struct {
	Status string `json:"status,omitempty"`
	// Weirdly this field is a number type in JSON but not the record object in Retrieve.
	ID      int    `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

func (p *Provider) doCreate(ctx context.Context, domain string, req createRecordRequest) (*createRecordResponse, error) {
	req.TTLStr = strconv.Itoa(int(req.TTL.Seconds()))
	resp, err := resty.R().SetContext(ctx).SetBody(req).Post(fmt.Sprintf("https://porkbun.com/api/json/v3/dns/create/%s", domain))
	if err != nil {
		return nil, err
	}
	r := &createRecordResponse{}
	if err := json.Unmarshal(resp.Body(), r); err != nil {
		return nil, err
	}
	if r.Status != SUCCESS {
		return nil, fmt.Errorf("API call failed: %v", r.Message)
	}
	return r, nil
}

type editRecordRequest struct {
	APIKey       string `json:"apikey,omitempty"`
	SecretAPIKey string `json:"secretapikey,omitempty"`
	Name         string `json:"name,omitempty"`
	Type         string `json:"type,omitempty"`
	Content      string `json:"content,omitempty"`
	TTL          time.Duration
	TTLStr       string `json:"ttl,omitempty"`
	Prio         string `json:"prio,omitempty"`
}
type editRecordResponse struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func (p *Provider) doEdit(ctx context.Context, domain, id string, req editRecordRequest) (*editRecordResponse, error) {
	req.TTLStr = strconv.Itoa(int(req.TTL.Seconds()))
	resp, err := resty.R().SetContext(ctx).SetBody(req).Post(fmt.Sprintf("https://porkbun.com/api/json/v3/dns/edit/%s/%s", domain, id))
	if err != nil {
		return nil, err
	}
	r := &editRecordResponse{}
	if err := json.Unmarshal(resp.Body(), r); err != nil {
		return nil, err
	}
	if r.Status != SUCCESS {
		return nil, fmt.Errorf("API call failed: %v", r.Message)
	}
	return r, nil
}
