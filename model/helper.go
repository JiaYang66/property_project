package model

import (
	"github.com/jinzhu/gorm"
)

type PaginationQ struct {
	//每页显示的数量
	PageSize uint `json:"pageSize"`
	//当前页码
	Page uint `json:"page"`
	//分页的数据内容
	Data interface{} `json:"data"`
	//全部的页码数量
	Total uint `json:"total"`
}

// 分页查询
func (p *PaginationQ) SearchAll(queryTx *gorm.DB) (data *PaginationQ, err error) {
	err = queryTx.Count(&p.Total).Error
	if err != nil {
		return p, err
	}
	if p.PageSize == 9999 {
		err = queryTx.Scan(p.Data).Error
		return p, err
	}
	if p.PageSize < 1 {
		p.PageSize = 10
	}
	if p.Page < 1 {
		p.Page = 1
	}
	offset := p.PageSize * (p.Page - 1)
	err = queryTx.Limit(p.PageSize).Offset(offset).Scan(p.Data).Error
	return p, err
}
