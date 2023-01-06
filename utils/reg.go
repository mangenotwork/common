package utils

import "regexp"

func RegFindAll(regStr, rest string) [][]string {
	reg := regexp.MustCompile(regStr)
	List := reg.FindAllStringSubmatch(rest, -1)
	reg.FindStringSubmatch(rest)
	return List
}

func RegFindAllTxt(regStr, rest string) (dataList []string) {
	reg := regexp.MustCompile(regStr)
	resList := reg.FindAllStringSubmatch(rest, -1)
	for _, v := range resList {
		if len(v) < 1 {
			continue
		}
		dataList = append(dataList, v[1])
	}
	return
}
