// Code generated by "stringer -type=OnConflict"; DO NOT EDIT.

package dotfile

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Rename-0]
	_ = x[Replace-1]
	_ = x[Warn-2]
	_ = x[Fail-3]
}

const _OnConflict_name = "RenameReplaceWarnFail"

var _OnConflict_index = [...]uint8{0, 6, 13, 17, 21}

func (i OnConflict) String() string {
	if i < 0 || i >= OnConflict(len(_OnConflict_index)-1) {
		return "OnConflict(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OnConflict_name[_OnConflict_index[i]:_OnConflict_index[i+1]]
}