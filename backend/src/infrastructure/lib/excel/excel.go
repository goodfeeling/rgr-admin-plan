// file: /src/infrastructure/lib/excel/excel.go
package excel

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

type ExcelHandler struct{}

func NewExcelHandler() *ExcelHandler {
	return &ExcelHandler{}
}

// ExcelData 表示Excel数据结构
type ExcelData struct {
	Headers []string
	Rows    [][]string
}

// CreateExcel 创建Excel文件
func (e *ExcelHandler) CreateExcel(sheetName string, data *ExcelData) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	// 创建或获取工作表
	if sheetName != "" {
		f.SetSheetName("Sheet1", sheetName)
	} else {
		sheetName = "Sheet1"
	}

	// 写入表头
	for i, header := range data.Headers {
		cell, err := excelize.CoordinatesToCellName(i+1, 1)
		if err != nil {
			return nil, fmt.Errorf("failed to convert coordinates for header %s: %w", header, err)
		}
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据行
	for rowIndex, row := range data.Rows {
		for colIndex, cellValue := range row {
			cell, err := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			if err != nil {
				return nil, fmt.Errorf("failed to convert coordinates at row %d, col %d: %w", rowIndex, colIndex, err)
			}
			f.SetCellValue(sheetName, cell, cellValue)
		}
	}

	// 将文件写入缓冲区
	buffer := &bytes.Buffer{}
	if err := f.Write(buffer); err != nil {
		return nil, fmt.Errorf("failed to write file to buffer: %w", err)
	}

	// 关闭文件以确保所有数据都被写入
	if err := f.Close(); err != nil {
		return nil, fmt.Errorf("failed to close file: %w", err)
	}

	return buffer, nil
}

// CreateCSV 创建CSV文件
func (e *ExcelHandler) CreateCSV(data *ExcelData) (*bytes.Buffer, error) {
	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)

	// 写入表头
	if err := writer.Write(data.Headers); err != nil {
		return nil, err
	}

	// 写入数据行
	for _, row := range data.Rows {
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	// 确保所有数据都写入缓冲区
	writer.Flush()

	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buffer, nil
}

// ReadExcel 读取Excel文件
func (e *ExcelHandler) ReadExcel(file multipart.File, sheetName string) (*ExcelData, error) {
	// 读取Excel文件
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 获取工作表名称
	if sheetName == "" {
		sheetName = f.GetSheetName(0) // 获取第一个工作表
	}

	// 读取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &ExcelData{}, nil
	}

	// 第一行作为表头
	headers := rows[0]

	// 其余行作为数据
	var dataRows [][]string
	if len(rows) > 1 {
		dataRows = rows[1:]
	}

	return &ExcelData{
		Headers: headers,
		Rows:    dataRows,
	}, nil
}

// ReadCSV 读取CSV文件
func (e *ExcelHandler) ReadCSV(file multipart.File) (*ExcelData, error) {
	// 读取CSV文件
	reader := csv.NewReader(file)

	// 读取所有记录
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return &ExcelData{}, nil
	}

	// 第一行作为表头
	headers := records[0]

	// 其余行作为数据
	var dataRows [][]string
	if len(records) > 1 {
		dataRows = records[1:]
	}

	return &ExcelData{
		Headers: headers,
		Rows:    dataRows,
	}, nil
}

// CreateTemplate 创建API模板
func (e *ExcelHandler) CreateApiTemplate(headers []string, sheetName string) (*bytes.Buffer, error) {
	templateData := &ExcelData{
		Headers: headers,
		Rows:    [][]string{},
	}

	return e.CreateExcel(sheetName, templateData)
}

// CreateApiCSVTemplate 创建API CSV模板
func (e *ExcelHandler) CreateApiCSVTemplate(headers []string) (*bytes.Buffer, error) {
	templateData := &ExcelData{
		Headers: headers,
		Rows:    [][]string{},
	}

	return e.CreateCSV(templateData)
}
