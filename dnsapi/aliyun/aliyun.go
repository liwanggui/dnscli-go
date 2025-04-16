package aliyun

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/liwanggui/dnscli-go/dnsapi"
	"strings"
)

var regionID = "cn-hangzhou"

// Client 实现阿里云DNS服务提供商
type Client struct {
	api *alidns.Client
}

// NewClient 创建新的阿里云DNS提供商实例
func NewClient(accessKeyID, accessKeySecret string) (*Client, error) {
	client, err := alidns.NewClientWithAccessKey(
		regionID,
		accessKeyID,
		accessKeySecret,
	)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云DNS客户端失败: %v", err)
	}

	return &Client{api: client}, nil
}

// ListRecords 获取指定域名的所有解析记录
func (client *Client) ListRecords(param *dnsapi.Parameter) ([]dnsapi.Record, error) {
	if param.Domain == "" {
		return nil, fmt.Errorf("域名不能为空")
	}

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.DomainName = param.Domain
	request.RRKeyWord = param.Name
	request.Type = param.Type
	request.Line = param.Line
	request.ValueKeyWord = param.Value
	request.PageSize = requests.NewInteger(100)

	pageNumber := 1
	records := make([]dnsapi.Record, 0, 20)

	for {
		request.PageNumber = requests.NewInteger(pageNumber)
		response, err := client.api.DescribeDomainRecords(request)
		if err != nil {
			// 尝试提供更详细的错误信息
			if strings.Contains(err.Error(), "InvalidDomainName") {
				return nil, fmt.Errorf("无效的域名格式: %s", param.Domain)
			} else if strings.Contains(err.Error(), "DomainNotFound") {
				return nil, fmt.Errorf("域名不存在或未在阿里云DNS中管理: %s", param.Domain)
			} else if strings.Contains(err.Error(), "InvalidAccessKeyId.NotFound") {
				return nil, fmt.Errorf("无效的AccessKey ID或Secret")
			}
			return nil, fmt.Errorf("获取域名记录失败: %v", err)
		}

		for _, r := range response.DomainRecords.Record {
			records = append(records, dnsapi.Record{
				ID:       r.RecordId,
				Domain:   param.Domain,
				Name:     r.RR,
				Type:     r.Type,
				Value:    r.Value,
				TTL:      int(r.TTL),
				Line:     r.Line,
				Priority: int(r.Priority),
			})
		}

		if int64(len(records)) >= response.TotalCount {
			break
		}

		pageNumber++
	}

	return records, nil
}

// GetRecord 获取特定记录的详情
func (client *Client) GetRecord(param *dnsapi.Parameter) (*dnsapi.Record, error) {
	request := alidns.CreateDescribeDomainRecordInfoRequest()
	if param.ID == "" {
		return nil, fmt.Errorf("记录 ID 不能为空")
	}
	request.RecordId = param.ID
	response, err := client.api.DescribeDomainRecordInfo(request)
	if err != nil {
		return nil, fmt.Errorf("获取记录详情失败: %v", err)
	}

	return &dnsapi.Record{
		ID:       response.RecordId,
		Domain:   param.Domain,
		Name:     response.RR,
		Type:     response.Type,
		Value:    response.Value,
		TTL:      int(response.TTL),
		Priority: int(response.Priority),
	}, nil
}

// AddRecord 添加新的解析记录
func (client *Client) AddRecord(param *dnsapi.Parameter) error {
	request := alidns.CreateAddDomainRecordRequest()
	request.DomainName = param.Domain
	request.RR = param.Name
	request.Type = param.Type
	request.Value = param.Value
	if param.TTL != 0 {
		request.TTL = requests.NewInteger(param.TTL)
	}
	if param.Priority > 0 {
		request.Priority = requests.NewInteger(param.Priority)
	}

	_, err := client.api.AddDomainRecord(request)
	if err != nil {
		return fmt.Errorf("添加解析记录失败: %v", err)
	}

	return nil
}

// UpdateRecord 更新现有解析记录
func (client *Client) UpdateRecord(param *dnsapi.Parameter) error {
	request := alidns.CreateUpdateDomainRecordRequest()
	request.RecordId = param.ID
	request.RR = param.Name
	request.Type = param.Type
	request.Value = param.Value
	if param.TTL != 0 {
		request.TTL = requests.NewInteger(param.TTL)
	}
	if param.Priority > 0 {
		request.Priority = requests.NewInteger(param.Priority)
	}

	_, err := client.api.UpdateDomainRecord(request)
	if err != nil {
		return fmt.Errorf("更新解析记录失败: %v", err)
	}

	return nil
}

// DeleteRecord 删除解析记录
func (client *Client) DeleteRecord(param *dnsapi.Parameter) error {
	request := alidns.CreateDeleteDomainRecordRequest()
	if param.ID == "" {
		return fmt.Errorf("记录 ID 不能为空")
	}
	request.RecordId = param.ID
	_, err := client.api.DeleteDomainRecord(request)
	if err != nil {
		if err, ok := err.(errors.Error); ok {
			return fmt.Errorf("删除解析记录失败: %v", err.Message())
		} else {
			return fmt.Errorf("删除解析记录失败: %v", err)
		}
	}
	return nil
}

// ListDomains 列出账号下所有域名
func (client *Client) ListDomains() ([]string, error) {
	request := alidns.CreateDescribeDomainsRequest()

	response, err := client.api.DescribeDomains(request)
	if err != nil {
		return nil, fmt.Errorf("获取域名列表失败: %v", err)
	}

	domains := make([]string, 0, len(response.Domains.Domain))
	for _, d := range response.Domains.Domain {
		domains = append(domains, d.DomainName)
		//fmt.Println(d.DomainName, !d.InstanceExpired, time.UnixMilli(d.CreateTimestamp).Format("2006-01-02 15:04:05"))
	}

	return domains, nil
}
