package Routing

import (
	"System/Routing/PathSegmentType"
	"errors"
	"strings"
)

type RouteParser struct {
}

func (this *RouteParser) SplitUrlToSegmentString(url string) []string {
	var result []string
	if url != "" {
		var index int
		for i, j := 0, len(url); i < j; i = index + 1 {
			index = strings.Index(url[i:], "/")
			if index == -1 {
				str := url[i:]
				if len(str) > 0 {
					result = append(result, str)
				}
				return result
			}
			index = index + i
			item := url[i:index]
			if len(item) > 0 {
				result = append(result, item)
			}
		}
	}
	return result
}
func (this *RouteParser) ParseUrl(url string) (*PathParser, error) {
	var segments []RoutSegment
	if len(url) > 0 {
		if strings.Index(url, "//") != -1 {
			return nil, errors.New("路径中不能出现连续的分隔符//")
		}
		if url[0] == '~' || url[0] == '/' {
			return nil, errors.New("不能以'~','/','~/'开头")
		}
	} else {
		return NewPathParser(segments, url, false), nil
	}

	segmentStrings := this.SplitUrlToSegmentString(url)
	segCount := len(segmentStrings)
	haveSegmentWithCatchAll := false
	for i := 0; i < segCount; i++ {
		if haveSegmentWithCatchAll {
			return nil, errors.New("catch all必须在路由的最后一部分")
		}
		seg := segmentStrings[i]
		seglen := len(seg)

		var subSegments []RoutSubSegment
		if strings.Index(seg, "{}") != -1 {
			return nil, errors.New("参数名称不能为空")
		}
		if strings.IndexAny(seg, "{}") == -1 {
			subItem := RoutSubSegment{SegmentType: PathSegmentType.Literal, Name: seg}
			subSegments = append(subSegments, subItem)
			item := RoutSegment{SubSegments: subSegments}
			segments = append(segments, item)
			continue
		}
		var tem string
		var from, start int = 0, 0
		for from < seglen {

			start = strings.Index(seg[from:], "{")
			if start >= seglen-2 {
				return nil, errors.New("未结束的URL参数，缺少}")
			}
			if start < 0 {
				/*判断是否有},如果有说明参数错误*/
				if strings.Index(seg[from:], "}") >= from {
					return nil, errors.New("URL参数错误，缺少{")
				}
				tmp := seg[from:]
				subItem := RoutSubSegment{SegmentType: PathSegmentType.Literal, Name: tmp}
				subSegments = append(subSegments, subItem)
				from += len(tem)
				break
			}
			start = start + from

			if from == 0 && start > 0 {
				strSeg := seg[from:start]
				subItem := RoutSubSegment{SegmentType: PathSegmentType.Literal, Name: strSeg}
				subSegments = append(subSegments, subItem)
			}

			end := strings.Index(seg[start:], "}")

			if end < 0 {
				return nil, errors.New("没有关才的URL参数，缺少}")
			}

			end = end + start

			next := strings.Index(seg[end:], "{")
			if next+end == end+1 {
				return nil, errors.New("不允许两个连续的参数在一起，中间必须有常量字符隔开")
			}
			if next == -1 {
				next = seglen
			} else {
				next = end + next
			}
			paramName := seg[start+1 : end]
			var segtype int
			if paramName[0] == '*' {
				haveSegmentWithCatchAll = true
				segtype = PathSegmentType.CatchAll
				paramName = paramName[1:]
			} else {
				segtype = PathSegmentType.Standart
			}
			temSubSeg := RoutSubSegment{SegmentType: segtype, Name: paramName}
			subSegments = append(subSegments, temSubSeg)

			if end < seglen-1 {
				token := seg[end+1 : next]
				subItem := RoutSubSegment{SegmentType: PathSegmentType.Literal, Name: token}
				subSegments = append(subSegments, subItem)
				end = end + len(token)
			}
			from = end + 1
		}
		pathseg := RoutSegment{SubSegments: subSegments}
		segments = append(segments, pathseg)
	}
	return NewPathParser(segments, url, haveSegmentWithCatchAll), nil
}
