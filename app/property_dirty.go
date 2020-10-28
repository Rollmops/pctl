package app

import (
	"strconv"
	"strings"
)

type DirtyProperty struct{}

func init() {
	PropertyMap["dirty"] = &DirtyProperty{}
}

func (*DirtyProperty) Name() string {
	return "Dirty"
}

func (*DirtyProperty) Value(p *Process, formatted bool) (string, error) {
	if p.RunningInfo != nil {
		if !formatted {
			return strconv.FormatBool(p.RunningInfo.DirtyInfo.IsDirty()), nil
		}
		if p.RunningInfo.DirtyInfo.IsDirty() {
			var dirtyParts []string
			if p.RunningInfo.DirtyInfo.DirtyCommand {
				dirtyParts = append(dirtyParts, "command changed")
			}
			if len(p.RunningInfo.DirtyInfo.DirtyMd5Hashes) > 0 {
				dirtyHashesString := strings.Join(p.RunningInfo.DirtyInfo.DirtyMd5Hashes, ",")
				dirtyParts = append(dirtyParts, "md5sum: "+dirtyHashesString)
			}
			return WarningColor(strings.Join(dirtyParts, " | ")), nil
		} else {
			return OkColor("-"), nil
		}
	}
	if formatted {
		return "", nil
	}
	return strconv.FormatBool(false), nil
}

func (*DirtyProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*DirtyProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
