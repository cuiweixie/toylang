package syntax

import "fmt"

type Pos struct {
	fileName string
	line, col int
}

func(pos *Pos) String() string {
	return fmt.Sprintf("%s:%d:%d", pos.fileName, pos.line, pos.col)
}

type Scanner struct {
	content []byte
	index int
	Pos
	literal string
	tToken TokenType
	err ScannerError
	Prec Prec
	isBinaryOp bool
}

type Prec int
const (
	NONE Prec = iota
	LOGICPREC
	ADDPREC
	MULPREC
)

type ScannerError struct {
	msg string
	pos Pos
}

func NewScanner(fileName string, content []byte) *Scanner {
	s := &Scanner{content: content}
	s.Pos.fileName = fileName
	s.Pos.line = 1
	s.Pos.col = 1
	return s
}

func (s *Scanner) Next() {
	s.isBinaryOp = false
	for {
		ch, ok := s.nextCh()
		if !ok {
			s.tToken = EOF
			return
		}
		switch ch {
		case '\r':
		case '\n':
			fallthrough
		case ';':
			s.line++
			s.col = 1
			s.tToken = SEMICOLON
			return
		case '"':
			s.tToken = STRING
			s.col ++
			str := ""
			for {
				ch, ok := s.nextCh()
				if !ok {
					s.err = ScannerError{
						msg: "",
						pos: s.Pos,
					}
					return
				}
				s.col ++
				if ch == '\\' {
					ch, ok := s.nextCh()
					if !ok {
						s.err = ScannerError{
							msg: "",
							pos: s.Pos,
						}
						return
					}
					s.col ++
					if ch == 'n' {
						str += "\n"
					}
					if ch == 't' {
						str += "\t"
					}
				} else {
					if ch == '"' {
						s.literal = str
						return
					} else {
						str += string(ch)
					}
				}
			}
		case '{':
			s.tToken = LEFTBRACE
			s.col++
			return
		case '}':
			s.tToken = RIGHTBRACE
			s.col++
			return
		case ',':
			s.tToken = COMMA
			s.col ++
			return
		case ' ':
			s.col ++
		case '+':
			s.col ++
			s.tToken = PLUS
			s.isBinaryOp = true
			s.Prec = ADDPREC
			return
		case '-':
			s.col ++
			s.tToken = MINUS
			s.isBinaryOp = true
			s.Prec = ADDPREC
			return
		case '*':
			s.col ++
			s.tToken = MUL
			s.isBinaryOp = true
			s.Prec = MULPREC
			return
		case '/':
			s.col ++
			s.tToken = DIV
			s.isBinaryOp = true
			s.Prec = MULPREC
			return
		case '=':
			ch, ok := s.nextCh()
			if !ok {
				s.col++
				s.tToken = ASSIGN
				return
			}
			if ch == '=' {
				s.tToken = EQUAL
				s.col += 2
				s.isBinaryOp = true
				s.Prec = LOGICPREC
				return
			}
			s.unGetCh()
			s.col++
			s.tToken = ASSIGN
			return
		case '(':
			s.tToken = LEFTPAREN
			s.col ++
			return
		case ')':
			s.tToken = RIGHTPAREN
			s.col ++
			return
		case '<':
			s.isBinaryOp = true
			s.Prec = LOGICPREC
			ch, ok := s.nextCh()
			if !ok {
				s.col++
				s.tToken = LT
				return
			}
			if ch == '=' {
				s.tToken = LEQ
				s.col += 2
				return
			}
			s.unGetCh()
			s.col++
			s.tToken = LT
			return
		case '>':
			s.isBinaryOp = true
			s.Prec = LOGICPREC
			ch, ok := s.nextCh()
			if !ok {
				s.col++
				s.tToken = GT
				return
			}
			if ch == '=' {
				s.tToken = GEQ
				s.col += 2
				return
			}
			s.unGetCh()
			s.col++
			s.tToken = GT
			return
		default:
			if isDigit(ch) {
				s.tToken = NUM
				str := string(ch)
				s.col++
				for {
					ch, ok := s.nextCh()
					if !ok {
						s.literal = str
						return
					}
					if isDigit(ch) {
						str += string(ch)
						s.col ++
					} else {
						s.unGetCh()
						s.literal = str
						return
					}
				}
			} else {
				if isLegalIdent(ch, true) {
					str := string(ch)
					s.col ++
					for {
						ch, ok := s.nextCh()
						if !ok {
							s.Ident(str)
							return
						}
						if isLegalIdent(ch, false) {
							s.col++
							str += string(ch)
						} else {
							s.unGetCh()
							s.Ident(str)
							return
						}
					}
				}
			}
		}
	}
}

func (s *Scanner) Ident(str string) {
	s.literal = str
	switch str {
	case "if":
		s.tToken = _KIF
	case "var":
		s.tToken = _KVAR
	case "else":
		s.tToken = _KELSE
	case "for":
		s.tToken = _KFOR
	case "func":
		s.tToken = _KFUNC
	case "break":
		s.tToken = _KBREAK
	case "continue":
		s.tToken = _KCONTINUE
	case "return":
		s.tToken = _KRETURN
	default:
		s.tToken = IDENT
	}
}
func isLegalIdent(ch byte, first bool) bool {
	if ch <= 'z' && ch >= 'a' {
		return true
	}
	if ch <= 'Z' && ch >= 'A' {
		return true
	}
	if ch == '_' {
		return true
	}
	if !first {
		return isDigit(ch)
	}
	return false
}


func isDigit(ch byte) bool {
	a := int(ch - '0')
	return a >=0 && a <= 9
}

func (s *Scanner) nextCh() (byte, bool) {
	if s.index >= len(s.content) {
		return 0, false
	}
	ch := s.content[s.index]
	s.index++
	return ch, true
}

func (s *Scanner) unGetCh() {
	s.index--
}

//go:generate stringer -type TokenType -linecomment scanner.go
type TokenType int

const (
	_ TokenType = iota
	IDENT
	_KVAR
	_KFUNC
	_KIF
	_KELSE
	_KFOR
	_KBREAK
	_KCONTINUE
	_KRETURN
	NUM
	STRING
	EOF
	MINUS
	PLUS
	MUL
	DIV
	LT
	LEQ
	GT
	GEQ
	LEFTPAREN
	RIGHTPAREN
	LEFTBRACE
	RIGHTBRACE
	ASSIGN
	EQUAL
	SEMICOLON
	COMMA
)


