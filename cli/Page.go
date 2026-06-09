package cli

import (
	"fast_ai_client/core"
	"fmt"
	"strings"
)

const (
	pageWidth    = 200
	messageWidth = 120
	messageLeftX = 2
	messageTopY  = 20
)

type Page struct {
	Title   string
	Body    strings.Builder
	NextRow int
	IsChat  bool
}

func NewPage(title string, body string, isChat bool) *Page {
	p := Page{Title: title, Body: strings.Builder{}, NextRow: messageTopY, IsChat: isChat}
	p.Body.WriteString(body)
	return &p
}

func (p *Page) AppendBody(body string) {
	p.Body.WriteString(body)
}

func (p *Page) AddMessage(msg string, user bool) {
	x := messageLeftX
	if user {
		x = pageWidth - messageWidth - 1
	}

	lines := wrapText(msg, messageWidth)
	var formatter core.Formatter
	for _, line := range lines {
		moveCursor(&p.Body, x, p.NextRow)
		if user {
			p.Body.WriteString(padLeft(line, messageWidth))
		} else {
			formatted := formatter.FormatChunk(line + "\n")
			p.Body.WriteString(strings.TrimSuffix(formatted, "\n"))
		}
		p.NextRow++
	}
	if !user {
		p.Body.WriteString(formatter.Close())
	}

	p.NextRow++
}

func (p *Page) Update() {
	ClearTerminal()
	fmt.Println(
		p.Title +
			"\n____________________________________________________________________\n" +
			p.Body.String(),
	)
}

func wrapText(text string, width int) []string {
	paragraphs := strings.Split(text, "\n")
	if len(paragraphs) == 0 {
		return []string{""}
	}

	var lines []string
	for _, paragraph := range paragraphs {
		wrapped := wrapLine(paragraph, width)
		lines = append(lines, wrapped...)
	}

	return lines
}

func wrapLine(text string, width int) []string {
	if strings.TrimSpace(text) == "" {
		return []string{""}
	}

	indent := leadingSpaces(text)
	if len(indent) >= width {
		indent = ""
	}

	words := strings.Fields(text)
	var lines []string
	var line strings.Builder
	wordWidth := width - len(indent)

	for _, word := range words {
		for len(word) > wordWidth {
			if line.Len() > 0 {
				lines = append(lines, line.String())
				line.Reset()
			}
			lines = append(lines, indent+word[:wordWidth])
			word = word[wordWidth:]
		}

		if line.Len() == 0 {
			line.WriteString(indent)
			line.WriteString(word)
			continue
		}

		if line.Len()+1+len(word) > width {
			lines = append(lines, line.String())
			line.Reset()
			line.WriteString(indent)
			line.WriteString(word)
			continue
		}

		line.WriteByte(' ')
		line.WriteString(word)
	}

	if line.Len() > 0 {
		lines = append(lines, line.String())
	}

	return lines
}

func leadingSpaces(text string) string {
	for i, r := range text {
		if r != ' ' && r != '\t' {
			return text[:i]
		}
	}
	return text
}

func padLeft(text string, width int) string {
	if len(text) >= width {
		return text
	}
	return strings.Repeat(" ", width-len(text)) + text
}
