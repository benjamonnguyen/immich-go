package names

import (
	"regexp"
	"strconv"
	"time"
)

// Samsung burst file name pattern
// #99  stack: Samsung #99
// 20231207_101605_001.jpg
// 20231207_101605_002.jpg
// 20231207_101605_xxx.jpg

var samsungRE = regexp.MustCompile(`^(\d{8}_\d{6})_(\d{3})(\..+)$`)

func (ic InfoCollector) Samsung(name string) (bool, NameInfo) {
	parts := samsungRE.FindStringSubmatch(name)
	if len(parts) == 0 {
		return false, NameInfo{}
	}
	info := NameInfo{
		Radical: parts[1],
		Base:    name,
		Ext:     parts[3],
		Type:    ic.SM.TypeFromExt(parts[3]),
		Kind:    KindBurst,
	}
	info.Index, _ = strconv.Atoi(parts[2])
	info.IsCover = info.Index == 1
	info.Taken, _ = time.ParseInLocation("20060102_150405", parts[1], ic.TZ)
	return true, info
}
