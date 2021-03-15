package porkbun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/resty.v1"
)

var errAPINotSuccess = errors.New("api status is not success")

const SUCCESS = "SUCCESS"

type retrieveRecordsRequest struct {
	APIKey       string `json:"api_key,omitempty"`
	SecretAPIKey string `json:"secret_api_key,omitempty"`
}

type retrieveRecordsResponse struct {
	Status  string   `json:"status,omitempty"`
	Records []record `json:"records,omitempty"`
}

type record struct {
	ID      string  `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Type    string  `json:"type,omitempty"`
	Content string  `json:"content,omitempty"`
	TTL     float64 `json:"ttl,omitempty"`
	Prio    string  `json:"prio,omitempty"`
	Notes   string  `json:"notes,omitempty"`
}

func (p *Provider) doRetrieve(ctx context.Context, domain string, req retrieveRecordsRequest) (*retrieveRecordsResponse, error) {
	resp, err := resty.R().SetContext(ctx).SetBody(req).Get(fmt.Sprintf("https://porkbun.com/api/json/v3/dns/retrieve/%s", domain))
	if err != nil {
		return nil, err
	}
	r := &retrieveRecordsResponse{}
	if err := json.Unmarshal(resp.Body(), r); err != nil {
		return nil, err
	}
	if r.Status != SUCCESS {
		return nil, errAPINotSuccess
	}
	return r, nil
}

type deleteRecordRequest struct {
	APIKey       string `json:"api_key,omitempty"`
	SecretAPIKey string `json:"secret_api_key,omitempty"`
}

type deleteRecordResponse struct {
	Status string `json:"status,omitempty"`
}

func (p *Provider) doDelete(ctx context.Context, domain, id string, req deleteRecordRequest) (*deleteRecordResponse, error) {
	resp, err := resty.R().SetContext(ctx).SetBody(req).Get(fmt.Sprintf("https://porkbun.com/api/json/v3/dns/delete/%s/%s", domain, id))
	if err != nil {
		return nil, err
	}
	r := &deleteRecordResponse{}
	if err := json.Unmarshal(resp.Body(), r); err != nil {
		return nil, err
	}
	if r.Status != SUCCESS {
		return nil, errAPINotSuccess
	}
	return r, nil
}

type createRecordRequest struct {
	APIKey       string  `json:"api_key,omitempty"`
	SecretAPIKey string  `json:"secret_api_key,omitempty"`
	Name         string  `json:"name,omitempty"`
	Type         string  `json:"type,omitempty"`
	Content      string  `json:"content,omitempty"`
	TTL          float64 `json:"ttl,omitempty"`
	Prio         string  `json:"prio,omitempty"`
}

type createRecordResponse struct {
	Status string `json:"status,omitempty"`
	ID     string `json:"id,omitempty"`
}

func (p *Provider) doCreate(ctx context.Context, domain string, req createRecordRequest) (*createRecordResponse, error) {
	resp, err := resty.R().SetContext(ctx).SetBody(req).Get(fmt.Sprintf("https://porkbun.com/api/json/v3/dns/create/%s", domain))
	if err != nil {
		return nil, err
	}
	r := &createRecordResponse{}
	if err := json.Unmarshal(resp.Body(), r); err != nil {
		return nil, err
	}
	if r.Status != SUCCESS {
		return nil, errAPINotSuccess
	}
	return r, nil
}

type editRecordRequest struct {
	APIKey       string  `json:"api_key,omitempty"`
	SecretAPIKey string  `json:"secret_api_key,omitempty"`
	Name         string  `json:"name,omitempty"`
	Type         string  `json:"type,omitempty"`
	Content      string  `json:"content,omitempty"`
	TTL          float64 `json:"ttl,omitempty"`
	Prio         string  `json:"prio,omitempty"`
}
type editRecordResponse struct {
	Status string `json:"status,omitempty"`
}

func (p *Provider) doEdit(ctx context.Context, domain, id string, req editRecordRequest) (*editRecordResponse, error) {
	resp, err := resty.R().SetContext(ctx).SetBody(req).Get(fmt.Sprintf("https://porkbun.com/api/json/v3/dns/edit/%s/%s", domain, id))
	if err != nil {
		return nil, err
	}
	r := &editRecordResponse{}
	if err := json.Unmarshal(resp.Body(), r); err != nil {
		return nil, err
	}
	if r.Status != SUCCESS {
		return nil, errAPINotSuccess
	}
	return r, nil
}
