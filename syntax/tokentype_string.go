// Code generated by "stringer -type TokenType -linecomment scanner.go"; DO NOT EDIT.

package syntax

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[IDENT-1]
	_ = x[_KVAR-2]
	_ = x[_KFUNC-3]
	_ = x[_KIF-4]
	_ = x[_KELSE-5]
	_ = x[_KFOR-6]
	_ = x[_KBREAK-7]
	_ = x[_KCONTINUE-8]
	_ = x[_KRETURN-9]
	_ = x[NUM-10]
	_ = x[STRING-11]
	_ = x[EOF-12]
	_ = x[MINUS-13]
	_ = x[PLUS-14]
	_ = x[MUL-15]
	_ = x[DIV-16]
	_ = x[LT-17]
	_ = x[LEQ-18]
	_ = x[GT-19]
	_ = x[GEQ-20]
	_ = x[LEFTPAREN-21]
	_ = x[RIGHTPAREN-22]
	_ = x[LEFTBRACE-23]
	_ = x[RIGHTBRACE-24]
	_ = x[ASSIGN-25]
	_ = x[EQUAL-26]
	_ = x[SEMICOLON-27]
	_ = x[COMMA-28]
}

const _TokenType_name = "IDENT_KVAR_KFUNC_KIF_KELSE_KFOR_KBREAK_KCONTINUE_KRETURNNUMSTRINGEOFMINUSPLUSMULDIVLTLEQGTGEQLEFTPARENRIGHTPARENLEFTBRACERIGHTBRACEASSIGNEQUALSEMICOLONCOMMA"

var _TokenType_index = [...]uint8{0, 5, 10, 16, 20, 26, 31, 38, 48, 56, 59, 65, 68, 73, 77, 80, 83, 85, 88, 90, 93, 102, 112, 121, 131, 137, 142, 151, 156}

func (i TokenType) String() string {
	i -= 1
	if i < 0 || i >= TokenType(len(_TokenType_index)-1) {
		return "TokenType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _TokenType_name[_TokenType_index[i]:_TokenType_index[i+1]]
}