package docx

import (
	"encoding/xml"
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
