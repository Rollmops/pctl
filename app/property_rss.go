package app

import "fmt"

type RssProperty struct{}

func init() {
	PropertyMap["rss"] = &RssProperty{}
}

func (*RssProperty) Name() string {
	return "Rss"
}

func (*RssProperty) Value(p *Process, _ bool) (string, error) {
	var rss string
	if p.Info != nil && p.IsRunning() {
		memoryInfo, err := p.Info.GoPsutilProcess.MemoryInfo()
		if err != nil {
			rss = "error"
		} else {
			rss = ByteCountIEC(memoryInfo.RSS)
		}
	}
	return rss, nil
}

func (*RssProperty) FormattedSumValue(processList ProcessList) (string, error) {
	var rssSum uint64
	for _, p := range processList {
		if p.Info != nil && p.IsRunning() {
			memoryInfo, err := p.Info.GoPsutilProcess.MemoryInfo()
			if err == nil {
				rssSum += memoryInfo.RSS
			}
		}
	}
	return fmt.Sprintf("Σ %s", ByteCountIEC(rssSum)), nil
}
func (*RssProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
