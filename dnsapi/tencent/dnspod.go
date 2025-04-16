package tencent

import (
	"fmt"
	"github.com/liwanggui/dnscli-go/dnsapi"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	dnspod "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dnspod/v20210323"
	"strconv"
)

type Client struct {
	client *dnspod.Client
}

func NewClient(secretID, secretKey string) (*Client, error) {
	credential := common.NewCredential(secretID, secretKey)
	cpf := profile.NewClientProfile()
	client, err := dnspod.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return nil, fmt.Errorf("创建DNSPod客户端失败: %v", err)
	}
	return &Client{client: client}, nil
}

func (p *Client) ListRecords(param *dnsapi.Parameter) ([]dnsapi.Record, error) {
	request := dnspod.NewDescribeRecordListRequest()
	request.Domain = &param.Domain
	request.Subdomain = &param.Name
	request.RecordType = common.StringPtr(param.Type)
	request.Keyword = common.StringPtr(param.Value)
	request.Limit = common.Uint64Ptr(100)
	records := make([]dnsapi.Record, 0, 20)

	var offset uint64 = 0
	// 循环获取所有数据
	for {
		request.Offset = common.Uint64Ptr(offset)
		response, err := p.client.DescribeRecordList(request)
		if err != nil {
			return nil, fmt.Errorf("获取域名记录失败: %v", err)
		}

		for _, r := range response.Response.RecordList {
			ttl := int(*r.TTL)
			priority := 0
			if r.MX != nil {
				priority = int(*r.MX)
			}

			records = append(records, dnsapi.Record{
				ID:       strconv.FormatUint(*r.RecordId, 10),
				Domain:   param.Domain,
				Name:     *r.Name,
				Type:     *r.Type,
				Value:    *r.Value,
				TTL:      ttl,
				Line:     *r.Line,
				Priority: priority,
				Updated:  *r.UpdatedOn,
			})
		}

		offset += *response.Response.RecordCountInfo.ListCount
		// 当偏移量等于记录总数时，表示后面没有数据了
		if offset >= *response.Response.RecordCountInfo.TotalCount {
			break
		}
	}
	return records, nil
}

// GetRecord 获取特定记录的详情
func (p *Client) GetRecord(param *dnsapi.Parameter) (*dnsapi.Record, error) {
	id, err := strconv.ParseUint(param.ID, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("无效的记录ID: %v", err)
	}

	request := dnspod.NewDescribeRecordRequest()
	request.Domain = &param.Domain
	request.RecordId = &id

	response, err := p.client.DescribeRecord(request)
	if err != nil {
		return nil, fmt.Errorf("获取记录详情失败: %v", err)
	}

	r := response.Response.RecordInfo
	ttl := int(*r.TTL)
	priority := 0
	if r.MX != nil {
		priority = int(*r.MX)
	}

	return &dnsapi.Record{
		ID:       param.ID,
		Domain:   param.Domain,
		Name:     *r.SubDomain,
		Type:     *r.RecordType,
		Value:    *r.Value,
		TTL:      ttl,
		Priority: priority,
		Updated:  *r.UpdatedOn,
	}, nil
}

// AddRecord 添加新的解析记录
func (p *Client) AddRecord(param *dnsapi.Parameter) error {
	request := dnspod.NewCreateRecordRequest()
	request.Domain = &param.Domain
	request.SubDomain = &param.Name
	request.RecordType = common.StringPtr(param.Type)
	request.RecordLine = common.StringPtr("默认")
	request.Value = &param.Value
	request.TTL = common.Uint64Ptr(uint64(param.TTL))

	if param.Priority > 0 {
		mx := uint64(param.Priority)
		request.MX = &mx
	}

	_, err := p.client.CreateRecord(request)
	if err != nil {
		return fmt.Errorf("添加解析记录失败: %v", err)
	}

	return nil
}

// UpdateRecord 更新现有解析记录
func (p *Client) UpdateRecord(param *dnsapi.Parameter) error {
	id, err := strconv.ParseUint(param.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("无效的记录ID: %v", err)
	}

	request := dnspod.NewModifyRecordRequest()
	request.Domain = &param.Domain
	request.RecordId = &id
	request.SubDomain = &param.Name
	request.RecordType = common.StringPtr(param.Type)
	request.RecordLine = common.StringPtr("默认")
	request.Value = &param.Value
	ttl := uint64(param.TTL)
	request.TTL = &ttl

	if param.Priority > 0 {
		mx := uint64(param.Priority)
		request.MX = &mx
	}

	_, err = p.client.ModifyRecord(request)
	if err != nil {
		return fmt.Errorf("更新解析记录失败: %v", err)
	}

	return nil
}

// DeleteRecord 删除解析记录
func (p *Client) DeleteRecord(param *dnsapi.Parameter) error {
	id, err := strconv.ParseUint(param.ID, 10, 64)
	if err != nil {
		return fmt.Errorf("无效的记录ID: %v", err)
	}

	request := dnspod.NewDeleteRecordRequest()
	request.Domain = &param.Domain
	request.RecordId = &id

	_, err = p.client.DeleteRecord(request)
	if err != nil {
		return fmt.Errorf("删除解析记录失败: %v", err)
	}

	return nil
}

// ListDomains 列出账号下所有域名
func (p *Client) ListDomains() ([]string, error) {
	request := dnspod.NewDescribeDomainListRequest()
	response, err := p.client.DescribeDomainList(request)
	if err != nil {
		return nil, fmt.Errorf("获取域名列表失败: %v", err)
	}

	domains := make([]string, 0, len(response.Response.DomainList))
	for _, d := range response.Response.DomainList {
		domains = append(domains, *d.Name)
		fmt.Println(*d.Name, *d.Status, *d.CreatedOn)
	}
	return domains, nil
}
