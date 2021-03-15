// Package porkbun implements a DNS record management client compatible
// with the libdns interfaces for porkbun.com.
package porkbun

import (
	"context"
	"fmt"
	"strings"

	"github.com/libdns/libdns"
)

// Provider facilitates DNS record manipulation with porkbun.com.
// For more info, https://porkbun.com/api/json/v3/documentation.
type Provider struct {
	APIKey       string `json:"api_key,omitempty"`
	SecretAPIKey string `json:"secret_api_key"`
}

// GetRecords lists all the records in the zone.
func (p *Provider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
	resp, err := p.doRetrieve(ctx, domain(zone), retrieveRecordsRequest{
		APIKey:       p.APIKey,
		SecretAPIKey: p.SecretAPIKey,
	})
	if err != nil {
		return nil, err
	}
	var ret []libdns.Record
	for _, r := range resp.Records {
		ret = append(ret, libdns.Record{
			ID:    r.ID,
			Type:  r.Type,
			Name:  r.Name,
			Value: r.Content,
			TTL:   r.TTL,
		})
	}
	return ret, nil
}

// AppendRecords adds records to the zone. It returns the records that were added.
func (p *Provider) AppendRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var ret []libdns.Record
	for _, r := range records {
		resp, err := p.doCreate(ctx, domain(zone), createRecordRequest{
			APIKey:       p.APIKey,
			SecretAPIKey: p.SecretAPIKey,
			Name:         r.Name,
			Type:         r.Type,
			Content:      r.Value,
			TTL:          r.TTL,
		})
		if err == nil {
			ret = append(ret, libdns.Record{
				ID:    fmt.Sprintf("%d", resp.ID),
				Type:  r.Type,
				Name:  r.Name,
				Value: r.Value,
				TTL:   r.TTL,
			})
		}
	}
	return ret, nil
}

// SetRecords sets the records in the zone, either by updating existing records or creating new ones.
// It returns the updated records.
func (p *Provider) SetRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var ret []libdns.Record
	for _, r := range records {
		if r.ID == "" {
			// New record, creating
			resp, err := p.doCreate(ctx, domain(zone), createRecordRequest{
				APIKey:       p.APIKey,
				SecretAPIKey: p.SecretAPIKey,
				Name:         r.Name,
				Type:         r.Type,
				Content:      r.Value,
				TTL:          r.TTL,
			})
			if err == nil {
				ret = append(ret, libdns.Record{
					ID:    fmt.Sprintf("%d", resp.ID),
					Type:  r.Type,
					Name:  r.Name,
					Value: r.Value,
					TTL:   r.TTL,
				})
			}
		} else {
			// Existing record, update it.
			_, err := p.doEdit(ctx, domain(zone), r.ID, editRecordRequest{
				APIKey:       p.APIKey,
				SecretAPIKey: p.SecretAPIKey,
				Name:         r.Name,
				Type:         r.Type,
				Content:      r.Value,
				TTL:          r.TTL,
			})
			if err == nil {
				ret = append(ret, r)
			}
		}

	}
	return ret, nil
}

// DeleteRecords deletes the records from the zone. It returns the records that were deleted.
// Requires Record.ID
func (p *Provider) DeleteRecords(ctx context.Context, zone string, records []libdns.Record) ([]libdns.Record, error) {
	var ret []libdns.Record
	for _, r := range records {
		if r.ID == "" {
			continue
		}
		if _, err := p.doDelete(ctx, domain(zone), r.ID, deleteRecordRequest{
			APIKey:       p.APIKey,
			SecretAPIKey: p.SecretAPIKey,
		}); err == nil {
			ret = append(ret, r)
		}
	}
	return ret, nil
}

func domain(zone string) string {
	return strings.TrimRight(zone, ".")
}

// Interface guards
var (
	_ libdns.RecordGetter   = (*Provider)(nil)
	_ libdns.RecordAppender = (*Provider)(nil)
	_ libdns.RecordSetter   = (*Provider)(nil)
	_ libdns.RecordDeleter  = (*Provider)(nil)
)
