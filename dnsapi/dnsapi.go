package dnsapi

import "fmt"

// RecordTypes 表示DNS记录类型
var RecordTypes = []string{"A", "AAAA", "CNAME", "MX", "TXT", "NS", "SRV", "CAA"}

func ValidRecordType(rType string) error {
	for _, r := range RecordTypes {
		if r == rType {
			return nil
		}
	}
	return fmt.Errorf("无效的 DNS 记录类型：%s", rType)
}

type Domain struct {
	DomainName string
	Status     bool   // ali: InstanceExpired tencent: Status
	CreateTime string // ali: CreateTime      tencent: CreatedOn
}

// Record 表示一条DNS记录
type Record struct {
	ID       string `json:"id" table:"记录ID"`
	Domain   string `json:"domain" table:"域名"`
	Name     string `json:"name" table:"主机记录"`
	Type     string `json:"type" table:"记录类型"`
	Value    string `json:"value" table:"记录值"`
	TTL      int    `json:"ttl"`
	Line     string `json:"line,omitempty" table:"线路名"`
	Priority int    `json:"priority,omitempty" table:"优先级"` // 用于MX和SRV记录
	Proxied  bool   `json:"proxied,omitempty"`              // 适用于 Cloudflare
	Updated  string `json:"updated,omitempty" table:"更新时间"`
}

// Parameter 解析请求参数
type Parameter struct {
	// ID 记录的唯一标识符
	ID string
	// Domain 域名
	Domain string
	// Name 主机记录（子域名）
	Name string
	// Type 记录类型（A、AAAA、CNAME、TXT、MX等）
	Type string
	// Value 记录值
	Value string
	// TTL 生存时间，免费版大多最低只支持设置为600s
	TTL int
	// Line DNS线路类型
	Line string
	// Priority 优先级，用于MX和SRV记录
	Priority int
	// Proxied 是否启用Cloudflare代理
	Proxied bool
	// Status 记录状态
	Status string
	// Remark 备注信息
	Remark string
}

func CreateParameter(domain string) *Parameter {
	return &Parameter{Domain: domain}
}

// DNSAPI 定义 DNS API 接口
type DNSAPI interface {
	// ListRecords 获取指定域名的所有解析记录
	ListRecords(param *Parameter) ([]Record, error)

	// GetRecord 获取特定记录的详情
	GetRecord(param *Parameter) (*Record, error)

	// AddRecord 添加新的解析记录
	AddRecord(param *Parameter) error

	// UpdateRecord 更新现有解析记录
	UpdateRecord(param *Parameter) error

	// DeleteRecord 删除解析记录
	DeleteRecord(param *Parameter) error

	// ListDomains 列出账号下所有域名
	ListDomains() ([]string, error)
}
