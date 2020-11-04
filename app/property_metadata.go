package app

import (
	"fmt"
	"strings"
)

type MetadataProperty struct{}

func init() {
	PropertyMap["metadata"] = &MetadataProperty{}
}

func (*MetadataProperty) Name() string {
	return "Metadata"
}

func (*MetadataProperty) Value(p *Process, _ bool) (string, error) {
	var metadataString []string
	for key, value := range p.Config.Metadata {
		metadataString = append(metadataString, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(metadataString, ","), nil
}

func (*MetadataProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*MetadataProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
