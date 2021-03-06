// Code generated by "stringer -type=OSType"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[OsUnkownType-0]
	_ = x[MacOs-1]
	_ = x[LinuxOs-2]
	_ = x[WindowsOs-3]
}

const _OSType_name = "OsUnkownTypeMacOsLinuxOsWindowsOs"

var _OSType_index = [...]uint8{0, 12, 17, 24, 33}

func (i OSType) String() string {
	if i < 0 || i >= OSType(len(_OSType_index)-1) {
		return "OSType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _OSType_name[_OSType_index[i]:_OSType_index[i+1]]
}
