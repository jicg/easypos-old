package excelize

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
)

// NewFile provides function to create new file by default template. For
// example:
//
//    xlsx := NewFile()
//
func NewFile() *File {
	file := make(map[string]string)
	file["_rels/.rels"] = XMLHeader + templateRels
	file["docProps/app.xml"] = XMLHeader + templateDocpropsApp
	file["docProps/core.xml"] = XMLHeader + templateDocpropsCore
	file["xl/_rels/workbook.xml.rels"] = XMLHeader + templateWorkbookRels
	file["xl/theme/theme1.xml"] = XMLHeader + templateTheme
	file["xl/worksheets/sheet1.xml"] = XMLHeader + templateSheet
	file["xl/styles.xml"] = XMLHeader + templateStyles
	file["xl/workbook.xml"] = XMLHeader + templateWorkbook
	file["[Content_Types].xml"] = XMLHeader + templateContentTypes
	return &File{
		Sheet: make(map[string]*xlsxWorksheet),
		XLSX:  file,
	}
}

// Save provides function to override the xlsx file with origin path.
func (f *File) Save() error {
	if f.Path == "" {
		return fmt.Errorf("No path defined for file, consider File.WriteTo or File.Write")
	}
	return f.SaveAs(f.Path)
}

// SaveAs provides function to create or update to an xlsx file at the provided
// path.
func (f *File) SaveAs(name string) error {
	file, err := os.OpenFile(name, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	return f.Write(file)
}

// Write provides function to write to an io.Writer.
func (f *File) Write(w io.Writer) error {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	f.contentTypesWriter()
	f.workbookWriter()
	f.workbookRelsWriter()
	f.worksheetWriter()
	f.styleSheetWriter()
	for path, content := range f.XLSX {
		fi, err := zw.Create(path)
		if err != nil {
			return err
		}
		_, err = fi.Write([]byte(content))
		if err != nil {
			return err
		}
	}
	err := zw.Close()
	if err != nil {
		return err
	}

	if _, err := buf.WriteTo(w); err != nil {
		return err
	}

	return nil
}
