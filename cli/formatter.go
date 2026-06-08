package cli

import "strings"

type Formatter struct {
	green       bool
	pendingStar bool
	blue        bool
	pendingHash bool
}

func (f *Formatter) FormatChunk(s string) string {
	var b strings.Builder

	for _, r := range s {
		if f.pendingHash {
			if r == '#' {
				if !f.blue {
					b.WriteString("\033[34m")
					f.blue = true
				}
			} else {
				b.WriteRune('#')
				b.WriteRune(r)
			}

			f.pendingHash = false
			continue
		}

		if f.pendingStar {
			if r == '*' {
				if f.green {
					b.WriteString("\033[0m")
					if f.blue {
						b.WriteString("\033[34m")
					}
				} else {
					b.WriteString("\033[32m")
				}
				f.green = !f.green
			} else {
				b.WriteRune('*')
				b.WriteRune(r)
			}

			f.pendingStar = false
			continue
		}

		if r == '\n' {
			if f.green {
				b.WriteString("\033[0m")
				f.green = false
			}
			if f.blue {
				b.WriteString("\033[0m")
				f.blue = false
			}
			b.WriteRune(r)
			continue
		}

		if r == '#' {
			f.pendingHash = true
			continue
		}

		if r == '*' {
			f.pendingStar = true
			continue
		}

		b.WriteRune(r)
	}

	return b.String()
}

func (f *Formatter) Close() string {
	var b strings.Builder

	if f.pendingStar {
		b.WriteRune('*')
		f.pendingStar = false
	}

	if f.pendingHash {
		b.WriteRune('#')
		f.pendingHash = false
	}

	if f.green {
		b.WriteString("\033[0m")
		f.green = false
	}

	if f.blue {
		b.WriteString("\033[0m")
		f.blue = false
	}

	return b.String()
}

func FormatText(s string) string {
	var f Formatter
	return f.FormatChunk(s) + f.Close()
}
