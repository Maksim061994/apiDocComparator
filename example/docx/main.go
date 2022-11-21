package main

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/nguyenthenguyen/docx"
	"github.com/sergi/go-diff/diffmatchpatch"
)

type DocStruct struct {
	XMLName xml.Name `xml:"document"`
	Body    Body     `xml:"body"`
}

type Body struct {
	XMLName   xml.Name    `xml:"body"`
	Paragraph []Paragraph `xml:"p"`
}

type Paragraph struct {
	XMLName       xml.Name        `xml:"p"`
	TextParagraph []TextParagraph `xml:"r"`
}

type TextParagraph struct {
	XMLName xml.Name `xml:"r"`
	Text    string   `xml:"t"`
}

func processingOneDocumnet(path string) string {
	doc, err := docx.ReadDocxFile(path)
	if err != nil {
		panic(err)
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

func main() {
	paragraphs1 := processingOneDocumnet("example/doc2.docx")
	paragraphs2 := processingOneDocumnet("example/doc1.docx")
	dmp := diffmatchpatch.New()
	diffs2 := dmp.DiffMain(paragraphs1, paragraphs2, true)
	prettyHtml := dmp.DiffPrettyHtml(diffs2)
	fmt.Println(prettyHtml + "\n")
	prettyText := dmp.DiffPrettyText(diffs2)
	fmt.Println(prettyText)
}
