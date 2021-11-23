package excel

import (
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"io"
)

type ExFile struct {
	File *excelize.File
}

func Open(r io.Reader) (*ExFile, error) {
	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, err
	}
	return &ExFile{File: f}, nil
}

func OpenFile(path string) (*ExFile, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	return &ExFile{File: f}, nil
}

// 根据表名获取表中所有数据
func (ex ExFile) GetSheetData(sheet string) ([][]string, error) {
	rows, err := ex.File.GetRows(sheet)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
