package iql

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/influxdata/influx-stress/internal/engine"
	parse "github.com/influxdata/influx-stress/internal/iql/parse"
	"github.com/influxdata/influxql"
)

// Token represents a lexical token.
type Token int

// These are the lexical tokens used by the file parser
const (
	ILLEGAL Token = iota
	EOF
	STATEMENT
	BREAK
)

var eof = rune(0)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func isNewline(r rune) bool {
	return r == '\n'
}

// Scanner scans the file and tokenizes the raw text
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a Scanner
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() { _ = s.r.UnreadRune() }

func (s *Scanner) peek() rune {
	ch := s.read()
	s.unread()
	return ch
}

// Scan moves the Scanner forward one character
func (s *Scanner) Scan() (tok Token, lit string) {
	ch := s.read()

	if isNewline(ch) {
		s.unread()
		return s.scanNewlines()
	} else if ch == eof {
		return EOF, ""
	} else {
		s.unread()
		return s.scanStatements()
	}
	// golint marks as unreachable code
	// return ILLEGAL, string(ch)
}

func (s *Scanner) scanNewlines() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isNewline(ch) {
			s.unread()
			break
		} else {
			buf.WriteRune(ch)
		}
	}

	return BREAK, buf.String()
}

func (s *Scanner) scanStatements() (tok Token, lit string) {
	var buf bytes.Buffer
	buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if isNewline(ch) && isNewline(s.peek()) {
			s.unread()
			break
		} else if isNewline(ch) {
			s.unread()
			buf.WriteRune(ch)
		} else {
			buf.WriteRune(ch)
		}
	}

	return STATEMENT, buf.String()
}

// ParseStatements takes a configFile and returns a slice of Statements
func ParseStatements(file string) ([]engine.Statement, error) {
	seq := []engine.Statement{}

	f, err := os.Open(file)
	check(err)

	s := NewScanner(f)

	for {
		t, l := s.Scan()

		if t == EOF {
			break
		}
		_, err := influxql.ParseStatement(l)
		if err == nil {

			seq = append(seq, &engine.InfluxqlStatement{Query: l, StatementID: parse.RandStr(10)})
		} else if t == BREAK {
			continue
		} else {
			f := strings.NewReader(l)
			p := parse.NewParser(f)
			parsed, err := p.Parse()
			if err != nil {
				return nil, err
			}
			s, err := buildStatement(parsed)
			if err != nil {
				return nil, err
			}
			seq = append(seq, s)

		}
	}

	f.Close()

	return seq, nil

}

func buildStatement(stmt parse.Statement) (engine.Statement, error) {
	switch s := stmt.(type) {
	case *parse.QueryStatement:
		return &engine.QueryStatement{
			StatementID:    s.StatementID,
			Name:           s.Name,
			TemplateString: s.TemplateString,
			Args:           s.Args,
			Count:          s.Count,
		}, nil
	case *parse.InsertStatement:
		return &engine.InsertStatement{
			StatementID:    s.StatementID,
			Name:           s.Name,
			TemplateString: s.TemplateString,
			TagCount:       s.TagCount,
			Timestamp:      s.Timestamp,
			Templates:      s.Templates,
		}, nil
	case *parse.ExecStatement:
		return &engine.ExecStatement{
			StatementID: s.StatementID,
			Script:      s.Script,
		}, nil
	case *parse.SetStatement:
		return &engine.SetStatement{
			StatementID: s.StatementID,
			Var:         s.Var,
			Value:       s.Value,
		}, nil
	case *parse.WaitStatement:
		return &engine.WaitStatement{
			StatementID: s.StatementID,
		}, nil
	case *parse.GoStatement:
		body, err := buildStatement(s.Statement)
		if err != nil {
			return nil, err
		}
		return &engine.GoStatement{
			Statement:   body,
			StatementID: s.StatementID,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported IQL statement type %T", stmt)
	}
}
