package fingers

import (
	"getitle/v1/pkg/dsl"
	"regexp"
	"strings"
)

func compiledMatch(reg *regexp.Regexp, s string) (string, bool) {
	matched := reg.FindStringSubmatch(s)
	if matched == nil {
		return "", false
	}
	if len(matched) == 1 {
		return "", true
	} else {
		return strings.TrimSpace(matched[1]), true
	}
}

func FingerMatcher(finger *Finger, content string) (*Framework, *Vuln, bool) {
	// 只进行被动的指纹判断, 将无视rules中的senddata字段
	for i, rule := range finger.Rules {
		var ishttp bool
		if finger.Protocol == "http" {
			ishttp = true
		}
		hasFrame, hasVuln, res := RuleMatcher(rule, content, ishttp)
		if hasFrame {
			frame, vuln := finger.ToResult(hasFrame, hasVuln, res, i)
			if frame.Version == "" && rule.Regexps.CompiledVersionRegexp != nil {
				for _, reg := range rule.Regexps.CompiledVersionRegexp {
					res, _ := compiledMatch(reg, content)
					if res != "" {
						frame.Version = res
						break
					}
				}
			}
			return frame, vuln, true
		}
	}
	return nil, nil, false
}

func RuleMatcher(rule *Rule, content string, ishttp bool) (bool, bool, string) {
	// 漏洞匹配优先
	if rule.Regexps == nil {
		return false, false, ""
	}
	for _, reg := range rule.Regexps.CompiledVulnRegexp {
		res, ok := compiledMatch(reg, content)
		if ok {
			return true, true, res
		}
	}

	var body, header string
	if ishttp {
		cs := strings.Index(content, "\r\n\r\n")
		if cs != -1 {
			body = content[cs+4:]
			header = content[:cs]
		}
	} else {
		body = content
	}

	// body匹配
	for _, bodyReg := range rule.Regexps.Body {
		if strings.Contains(body, bodyReg) {
			return true, false, ""
		}
	}

	// 正则匹配
	for _, reg := range rule.Regexps.CompliedRegexp {
		res, ok := compiledMatch(reg, content)
		if ok {
			return true, false, res
		}
	}

	// MD5 匹配
	for _, md5s := range rule.Regexps.MD5 {
		if md5s == dsl.Md5Hash([]byte(content)) {
			return true, false, ""
		}
	}

	// mmh3 匹配
	for _, mmh3s := range rule.Regexps.MMH3 {
		if mmh3s == dsl.Mmh3Hash32([]byte(content)) {
			return true, false, ""
		}
	}

	// http头匹配, http协议特有的匹配
	if !ishttp {
		return false, false, ""
	}

	for _, headerReg := range rule.Regexps.Header {
		if strings.Contains(header, headerReg) {
			return true, false, ""
		}
	}
	return false, false, ""
}