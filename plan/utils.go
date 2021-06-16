package plan

import (
	"regexp"
	"strconv"
)

var (
	mapAddrRx  = regexp.MustCompile(`^([^\[]+)\["([^"]+)"]$`)
	listAddrRx = regexp.MustCompile(`^([^\[]+)\[(\d+)]$`)
)

func splitAddress(address string) (string, interface{}) {
	mapMatch := mapAddrRx.FindAllStringSubmatch(address, -1)

	if len(mapMatch) > 0 {
		return mapMatch[0][1], mapMatch[0][2]
	}

	listMatch := listAddrRx.FindAllStringSubmatch(address, -1)
	if len(listMatch) > 0 {
		idx, err := strconv.Atoi(listMatch[0][2])
		if err != nil {
			return "", nil
		}
		return listMatch[0][1], idx
	}

	return "", nil
}
