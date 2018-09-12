package cvm

import (
	"fmt"
	"strconv"

	"github.com/JefferyWang/qcloud/util"
)

// SetInstanceIDs 设置实例id列表
func (params *Params) SetInstanceIDs(instanceIDs []string) {
	for i, v := range instanceIDs {
		key := fmt.Sprintf("InstanceIds.%v", i)
		(*params)[key] = v
	}
}

// SetFilters 设置实例id列表
func (params *Params) SetFilters(filters map[string][]string) {
	i := 0
	for k, v := range filters {
		key := fmt.Sprintf("Filters.%v", i)
		(*params)[key+".Name"] = k
		for j, v1 := range v {
			key1 := fmt.Sprintf("%v.Values.%v", key, j)
			(*params)[key1] = v1
		}
		i++
	}
}

// SetOffset 设置偏移量
func (params *Params) SetOffset(offset int) {
	(*params)["Offset"] = strconv.Itoa(offset)
}

// SetLimit 设置返回数量
func (params *Params) SetLimit(limit int) {
	(*params)["Limit"] = strconv.Itoa(limit)
}

// DescribeInstances 查看实例列表
func (params *Params) DescribeInstances() {
	(*params)["Action"] = "DescribeInstances"
	params.GetSign("GET", REQUEST_URL)
	resp, err := util.DoGet(REQUEST_URL, *params, nil)
	if err != nil {
		return
	}
	fmt.Println(string(resp))
}
