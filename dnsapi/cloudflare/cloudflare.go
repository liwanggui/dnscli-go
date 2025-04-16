package cloudflare

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/liwanggui/dnscli-go/dnsapi"
	"strings"
	"time"
)

type Client struct {
	client *cloudflare.Client
}

func NewClient(apiToken, apiEmail, apiKey string) (*Client, error) {
	if apiToken != "" {
		return &Client{client: cloudflare.NewClient(option.WithAPIToken(apiToken))}, nil
	}
	return &Client{client: cloudflare.NewClient(option.WithAPIKey(apiKey), option.WithAPIEmail(apiEmail))}, nil
}

func (this *Client) getZoneID(domainName string) (zoneID string, err error) {
	page, err := this.client.Zones.List(context.TODO(), zones.ZoneListParams{Name: cloudflare.F(domainName)})
	if err != nil {
		panic(err.Error())
	}
	for _, v := range page.Result {
		if v.Name == domainName {
			return v.ID, nil
		}
	}
	return "", fmt.Errorf("zone %s not found", domainName)
}

func (this *Client) ListRecords(param *dnsapi.Parameter) ([]dnsapi.Record, error) {
	zoneID, err := this.getZoneID(param.Domain)
	if err != nil {
		return nil, err
	}
	recordListParams := dns.RecordListParams{}
	recordListParams.ZoneID = cloudflare.F(zoneID)
	recordListParams.Name = cloudflare.F(dns.RecordListParamsName{Startswith: cloudflare.F(param.Name)})
	recordListParams.Type = cloudflare.F(dns.RecordListParamsType(param.Type))
	recordListParams.Content = cloudflare.F(dns.RecordListParamsContent{Startswith: cloudflare.F(param.Value)})
	recordListParams.PerPage = cloudflare.F(2.0)
	records := make([]dnsapi.Record, 0)

	page, err := this.client.DNS.Records.List(context.TODO(), recordListParams)
	if err != nil {
		return nil, err
	}
	for {
		for _, v := range page.Result {
			records = append(records, dnsapi.Record{
				ID:       v.ID,
				Domain:   param.Domain,
				Name:     strings.Replace(v.Name, "."+param.Domain, "", -1),
				Value:    v.Content,
				Type:     string(v.Type),
				TTL:      int(v.TTL),
				Priority: int(v.Priority),
				Proxied:  v.Proxied,
				Updated:  v.ModifiedOn.Format(time.DateTime),
			})
		}
		page, _ = page.GetNextPage()
		if len(page.Result) == 0 {
			break
		}
	}
	return records, nil
}

// GetRecord 获取特定记录的详情
func (this *Client) GetRecord(param *dnsapi.Parameter) (*dnsapi.Record, error) {
	zoneID, err := this.getZoneID(param.Domain)
	if err != nil {
		return nil, err
	}
	page, err := this.client.DNS.Records.Get(context.TODO(), param.ID, dns.RecordGetParams{ZoneID: cloudflare.F(zoneID)})
	if err != nil {
		return nil, err
	}
	return &dnsapi.Record{
		ID:       page.ID,
		Domain:   param.Domain,
		Name:     page.Name,
		Value:    page.Content,
		Type:     string(page.Type),
		TTL:      int(page.TTL),
		Priority: int(page.Priority),
		Proxied:  page.Proxied,
		Updated:  page.ModifiedOn.Format(time.DateTime),
	}, nil
}

// AddRecord 添加新的解析记录
// https://developers.cloudflare.com/api/resources/dns/subresources/records/methods/create/
func (this *Client) AddRecord(param *dnsapi.Parameter) error {
	zoneID, err := this.getZoneID(param.Domain)
	if err != nil {
		return err
	}
	recordParam := getRecordUnionParam(param)
	_, err = this.client.DNS.Records.New(context.TODO(), dns.RecordNewParams{
		ZoneID: cloudflare.F(zoneID),
		Record: *recordParam,
	})

	return err
}

// UpdateRecord 更新现有解析记录
func (this *Client) UpdateRecord(param *dnsapi.Parameter) error {
	zoneID, err := this.getZoneID(param.Domain)
	if err != nil {
		return err
	}
	recordParam := getRecordUnionParam(param)
	recordUpdateParams := dns.RecordUpdateParams{
		ZoneID: cloudflare.F(zoneID),
		Record: *recordParam,
	}
	_, err = this.client.DNS.Records.Update(context.TODO(), param.ID, recordUpdateParams)
	return err
}

// DeleteRecord 删除解析记录
func (this *Client) DeleteRecord(param *dnsapi.Parameter) error {
	zoneID, err := this.getZoneID(param.Domain)
	if err != nil {
		return err
	}
	_, err = this.client.DNS.Records.Delete(context.TODO(), param.ID,
		dns.RecordDeleteParams{ZoneID: cloudflare.F(zoneID)})
	return err
}

// ListDomains 列出账号下所有域名
func (this *Client) ListDomains() ([]string, error) {
	page, err := this.client.Zones.List(context.TODO(), zones.ZoneListParams{})
	if err != nil {
		return nil, err
	}
	var zoneList = make([]string, 0)
	for _, v := range page.Result {
		zoneList = append(zoneList, v.Name)
	}
	return zoneList, nil
}

func getRecordUnionParam(param *dnsapi.Parameter) *dns.RecordUnionParam {
	var recordUnionParam dns.RecordUnionParam
	switch param.Type {
	case "A":
		recordUnionParam = dns.ARecordParam{
			Name:    cloudflare.F(fmt.Sprintf("%s.%s", param.Name, param.Domain)),
			Type:    cloudflare.F(dns.ARecordTypeA),
			Content: cloudflare.F(param.Value),
			TTL:     cloudflare.F(dns.TTL(param.TTL)),
			Proxied: cloudflare.F(param.Proxied),
		}
	case "AAAA":
		recordUnionParam = dns.AAAARecordParam{
			Name:    cloudflare.F(fmt.Sprintf("%s.%s", param.Name, param.Domain)),
			Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
			Content: cloudflare.F(param.Value),
			TTL:     cloudflare.F(dns.TTL(param.TTL)),
			Proxied: cloudflare.F(param.Proxied),
		}
	case "CNAME":
		recordUnionParam = dns.CNAMERecordParam{
			Name:    cloudflare.F(fmt.Sprintf("%s.%s", param.Name, param.Domain)),
			Type:    cloudflare.F(dns.CNAMERecordTypeCNAME),
			Content: cloudflare.F(param.Value),
			TTL:     cloudflare.F(dns.TTL(param.TTL)),
			Proxied: cloudflare.F(param.Proxied),
		}
	case "MX":
		recordUnionParam = dns.MXRecordParam{
			Name:     cloudflare.F(fmt.Sprintf("%s.%s", param.Name, param.Domain)),
			Type:     cloudflare.F(dns.MXRecordTypeMX),
			Content:  cloudflare.F(param.Value),
			TTL:      cloudflare.F(dns.TTL(param.TTL)),
			Priority: cloudflare.F(float64(param.Priority)),
			Proxied:  cloudflare.F(param.Proxied),
		}
	case "TXT":
		recordUnionParam = dns.TXTRecordParam{
			Name:    cloudflare.F(fmt.Sprintf("%s.%s", param.Name, param.Domain)),
			Type:    cloudflare.F(dns.TXTRecordTypeTXT),
			Content: cloudflare.F(param.Value),
			TTL:     cloudflare.F(dns.TTL(param.TTL)),
			Proxied: cloudflare.F(param.Proxied),
		}
	}
	return &recordUnionParam
}
