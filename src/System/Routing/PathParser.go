package Routing

import (
	"System/Routing/PathSegmentType"
	"strings"
)

type RoutSegment struct {
	SubSegments []RoutSubSegment
}
type RoutSubSegment struct {
	SegmentType int
	Name        string
}
type PathParser struct {
	segments     []RoutSegment
	routUrl      string
	haveCatchAll bool
}

func NewPathParser(segs []RoutSegment, url string, haveCatchAll bool) *PathParser {
	return &PathParser{segments: segs, routUrl: url, haveCatchAll: haveCatchAll}
}

func (this *PathParser) MatchSegment(routseg RoutSegment, pathSegIndex int, arrSegs []string, routData map[string]interface{}) bool {
	requestPathSegment := arrSegs[pathSegIndex]
	pathSegmentLen := len(requestPathSegment)
	paraEndIndex := pathSegmentLen - 1

	subCount := len(routseg.SubSegments)
	for subIndex := subCount - 1; subIndex >= 0; subIndex-- {
		subItem := routseg.SubSegments[subIndex]
		if paraEndIndex < 0 {
			return false
		}
		if subItem.SegmentType == PathSegmentType.CatchAll {
			var strValue string = ""
			for j, l := pathSegIndex, len(arrSegs); j < l; j++ {
				if j > pathSegIndex {
					strValue = strValue + "/"
				}
				strValue += arrSegs[j]
			}
			routData[subItem.Name] = strValue
			return true
		}

		var paramStartIndex = 0
		if subItem.SegmentType == PathSegmentType.Literal {
			namelen := len(subItem.Name)
			//长度的比较
			if paraEndIndex-namelen+1 < 0 {
				return false
			}
			paramStartIndex = paraEndIndex - namelen + 1
			//统一转换成小写，进行比较
			if strings.ToLower(requestPathSegment[paramStartIndex:paraEndIndex+1]) != strings.ToLower(subItem.Name) {
				return false
			}
			paraEndIndex = paramStartIndex - 1
			continue
		}
		//标准参数
		nextSubSegIndex := subIndex - 1
		if nextSubSegIndex < 0 {
			//当前参数是最后一个参数
			routData[subItem.Name] = requestPathSegment[0 : paraEndIndex+1]
			return true
		}
		if paraEndIndex == 0 {
			return false
		}
		nextSubItem := routseg.SubSegments[nextSubSegIndex]
		//下一个参数必是Literal，确定下一个参数的位置
		paramStartIndex = paraEndIndex - 1
		lastIndex := strings.LastIndex(requestPathSegment[0:paramStartIndex], nextSubItem.Name)
		//不匹配
		if lastIndex == -1 {
			return false
		}
		paramStartIndex = lastIndex + len(nextSubItem.Name)
		sectionValue := requestPathSegment[paramStartIndex:paraEndIndex]
		if sectionValue == "" {
			return false
		}
		routData[subItem.Name] = sectionValue
		paraEndIndex = paramStartIndex - 1
	}
	return true
}
func (this *PathParser) AddDefaults(routData map[string]interface{}, defaults map[string]interface{}) map[string]interface{} {
	if defaults != nil {
		for k, v := range defaults {
			_, ok := routData[k]
			if !ok {
				routData[k] = v
			}
		}
	}
	return routData
}

func (this *PathParser) Match(requestPath string, defaults map[string]interface{}) map[string]interface{} {
	requestPath = strings.Trim(requestPath, "/")
	ret := make(map[string]interface{})
	url := this.routUrl
	//如果路由没有参数，全是常量，直接比较
	if strings.ToLower(url) == strings.ToLower(requestPath) && strings.Index(url, "{") < 0 {
		return this.AddDefaults(ret, defaults)
	}
	routParse := RouteParser{}
	arrPathSegs := routParse.SplitUrlToSegmentString(requestPath)
	pathSegCount := len(arrPathSegs)
	routSegCount := len(this.segments)
	haveDefaults := (defaults != nil && len(defaults) > 0)
	//对URL的段数进行比较
	if !haveDefaults && ((this.haveCatchAll && pathSegCount < routSegCount) || (!this.haveCatchAll && pathSegCount != routSegCount)) {
		return nil
	}
	i := 0
	for _, segment := range this.segments {
		if i >= pathSegCount {
			break
		}
		if !this.MatchSegment(segment, i, arrPathSegs, ret) {
			return nil
		}
		i++
	}
	//如果请求路径，小于rout的段数，取默认值
	if i < routSegCount {
		if !haveDefaults {
			return nil
		}
		for ; i < routSegCount; i++ {
			segment := this.segments[i]
			//默认值只能有一个参数，否则认为不匹配
			if len(segment.SubSegments) != 1 {
				return nil
			}
			if segment.SubSegments[0].SegmentType != PathSegmentType.Standart {
				return nil
			}

			_, ok := defaults[segment.SubSegments[0].Name]
			//如果参数没有默认值，不匹配
			if !ok {
				return nil
			}
		}
	}
	return this.AddDefaults(ret, defaults)
}
