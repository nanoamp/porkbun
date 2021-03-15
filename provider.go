// Package porkbun implements a DNS record management client compatible
// with the libdns interfaces for porkbun.com.
package porkbun

import (
	"context"
	"strings"
	"time"

	"github.com/libdns/libdns"
)

// TODO: Providers must not require additional provisioning steps by the callers; it
// should work simply by populating a struct and calling methods on it. If your DNS
// service requires long-lived state or some extra provisioning step, do it implicitly
// when methods are called; sync.Once can help with this, and/or you can use a
// sync.(RW)Mutex in your Provider struct to synchronize implicit provisioning.

// Provider facilitates DNS record manipulation with <TODO: PROVIDER NAME>.
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
			TTL:   time.Second * time.Duration(r.TTL),
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
			TTL:          r.TTL.Seconds(),
		})
		if err == nil {
			ret = append(ret, libdns.Record{
				ID:    resp.ID,
				Type:  r.ID,
				Name:  r.Name,
				Value: r.Type,
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
				TTL:          r.TTL.Seconds(),
			})
			if err == nil {
				ret = append(ret, libdns.Record{
					ID:    resp.ID,
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
				TTL:          r.TTL.Seconds(),
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

func (p *Provider) getRecords(ctx context.Context, zone string) (idToRecord map[string]libdns.Record, nameToID map[string]string, err error) {
	nameToID = make(map[string]string)
	idToRecord = make(map[string]libdns.Record)

	var rec []libdns.Record
	if rec, err = p.GetRecords(ctx, zone); err == nil {
		for _, rec := range rec {
			nameToID[rec.Name] = rec.ID
			idToRecord[rec.ID] = rec
		}
	}
	return
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
