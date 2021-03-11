package validator

import (
	"fmt"
	"io"
	"strings"
)

type ValidateTag struct {
	Op   string
	Args []interface{}
}

type LookaheadReader struct {
	cur, next rune
	reader    *strings.Reader
}

func NewLookaheadReader(s string) *LookaheadReader {
	r := &LookaheadReader{
		reader: strings.NewReader(s),
	}
	r.Next()
	return r
}

func (r *LookaheadReader) Next() rune {
	r.cur = r.next
	next, _, err := r.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			next = 0
		} else {
			panic(err.Error())
		}
	}
	r.next = next
	return r.cur
}

func (r *LookaheadReader) Peek() rune {
	return r.next
}

func (r *LookaheadReader) Match(want rune) bool {
	if r.next == want {
		r.Next()
		return true
	}
	return false
}

func (r *LookaheadReader) Read(want rune) error {
	if r.next == want {
		r.Next()
		return nil
	}
	return fmt.Errorf("unexpected character: %c", r.next)
}

func (r *LookaheadReader) HasNext() bool {
	return r.next != 0
}

func parseValidateTags(tag string) []ValidateTag {
	tags := []ValidateTag{}
	r := NewLookaheadReader(tag)
	for {
		eatWhitespace(r)
		if !r.HasNext() {
			break
		}
		op := readLiteral(r)
		args := []interface{}{}
		eatWhitespace(r)
		if r.Match('(') {
			//We got an argument list
			for {
				eatWhitespace(r)
				if r.Match(')') {
					break
				}
				arg := readLiteral(r)
				args = append(args, arg)
				eatWhitespace(r)
				if r.Match(',') {
					continue
				}
			}
		}
		tags = append(tags, ValidateTag{
			Op:   op,
			Args: args,
		})
		eatWhitespace(r)
		if !r.Match(',') {
			break
		}
	}

	return tags
}

func eatWhitespace(r *LookaheadReader) {
	for {
		if !r.Match(' ') {
			break
		}
	}
}

func read(r *strings.Reader, exp rune) {
	ch, _, err := r.ReadRune()
	if err != nil {
		panic(fmt.Sprintf("Unexpected string reader error: %s", err))
	}
	if ch != exp {
		panic(fmt.Sprintf("Unexpected rune: got: %c, want: %c", ch, exp))
	}
}

func isAlphanum(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		(ch >= '0' && ch <= '9')
}

func isNumericSep(ch rune) bool {
	return ch == '.' ||
		ch == '-' ||
		ch == 'e' ||
		ch == 'E' ||
		ch == '+'
}

func isUtilChar(ch rune) bool {
	return ch == '_'
}

func readLiteral(r *LookaheadReader) string {
	var res strings.Builder
	for r.HasNext() {
		if ch := r.Peek(); !(isAlphanum(ch) || isNumericSep(ch) || isUtilChar(ch)) {
			goto Res
		}
		res.WriteRune(r.Next())
	}
Res:
	return res.String()
}
