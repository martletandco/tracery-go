// Code generated by "stringer -type Type"; DO NOT EDIT.

package scan

import "strconv"

const _Type_name = "EOFErrorWordWhiteSpaceLeftBracketRightBracketLeftParenRightParenBackStrokeColonCommaOctoPeriod"

var _Type_index = [...]uint8{0, 3, 8, 12, 22, 33, 45, 54, 64, 74, 79, 84, 88, 94}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
