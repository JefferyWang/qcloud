package cos

import (
	"time"
)

// Error 接口返回错误时的结构
type Error struct {
	Code      string `xml:"Code"`
	Message   string `xml:"Message"`
	Resource  string `xml:"Resource"`
	RequestID string `xml:"RequestId"`
	TraceID   string `xml:"TraceId"`
}

// BucketsResult 接口返回的bucket列表结构
type BucketsResult struct {
	Error
	Owner   Owner `xml:"Owner"`
	Buckets struct {
		List []Bucket `xml:"Bucket"`
	} `xml:"Buckets"`
}

// Owner 持有者信息
type Owner struct {
	ID          string `xml:"ID"`
	DisplayName string `xml:"DisplayName"`
}

// Bucket 存储桶信息
type Bucket struct {
	Name         string    `xml:"Name"`
	Location     string    `xml:"Location"`
	CreationDate time.Time `xml:"CreationDate"`
}
