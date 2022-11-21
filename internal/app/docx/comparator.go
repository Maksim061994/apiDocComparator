package docx

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/nguyenthenguyen/docx"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/sirupsen/logrus"
)

func readOneDocumnet(docInput io.ReaderAt, size int64) string {
	doc, err := docx.ReadDocxFromMemory(docInput, size)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	docEdit := doc.Editable()
	byteValue := []byte(docEdit.GetContent())
	var docStruct DocStruct
	xml.Unmarshal(byteValue, &docStruct)
	paragraphs := docStruct.Body.Paragraph
	output := make([]string, 0)
	for _, paragraph := range paragraphs {
		fullText := ""
		for _, text := range paragraph.TextParagraph {
			fullText += text.Text
		}
		output = append(output, fullText)
	}
	result := strings.Join(output, "\n")
	return result
}

func diffFormatHtml(paragraphs1 string, paragraphs2 string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(paragraphs1, paragraphs2, true)
	prettyFormat := dmp.DiffPrettyHtml(diffs)
	return prettyFormat
}

func diffFormatText(paragraphs1 string, paragraphs2 string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(paragraphs1, paragraphs2, true)
	prettyFormat := dmp.DiffPrettyText(diffs)
	fmt.Println(prettyFormat)
	return prettyFormat
}

func Comparator(docx1 bytes.Buffer, docx2 bytes.Buffer, typeResponse string) string {
	response := ""
	if typeResponse == "" {
		typeResponse = "text"
	}
	readerDocx1 := bytes.NewReader(docx1.Bytes())
	readerDocx2 := bytes.NewReader(docx2.Bytes())
	paragraphs1 := readOneDocumnet(readerDocx1, int64(readerDocx1.Size()))
	paragraphs2 := readOneDocumnet(readerDocx2, int64(readerDocx2.Size()))
	if typeResponse == "html" {
		response = diffFormatHtml(paragraphs1, paragraphs2)
	} else if typeResponse == "text" {
		response = diffFormatText(paragraphs1, paragraphs2)
	}
	return response
}
