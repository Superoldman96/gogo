package pkg

import (
	"bytes"
	"net/http"
	"strings"

	"github.com/chainreactors/parsers"
	"github.com/chainreactors/utils/iutils"
)

func CollectSocketResponse(result *Result, socketContent []byte) {
	if ishttp, _ := GetStatusCode(socketContent); ishttp {
		result.Protocol = "http"
		CollectParsedResponse(result, parsers.NewResponseWithRaw(socketContent))
	} else {
		result.Content = socketContent
		if title := parsers.MatchTitle(socketContent); title != "" {
			result.HasTitle = true
			result.Title = title
		} else {
			result.Title = parsers.MatchCharacter(socketContent)
		}
		result.AddExtracts(Extractors.Extract(string(socketContent), true))
	}
}

func CollectHttpResponse(result *Result, resp *http.Response) {
	if resp == nil {
		return
	}

	CollectParsedResponse(result, parsers.NewResponse(resp, int64(DefaultMaxSize*2)))
}

func CollectParsedResponse(result *Result, resp *parsers.Response) {
	if resp == nil {
		return
	}
	result.IsHttp = true
	result.Httpresp = resp

	// tolower 忽略大小写
	for i, c := range result.Httpresp.History {
		result.Httpresp.History[i].Raw = bytes.ToLower(c.Raw)
	}
	result.Content = bytes.ToLower(result.Httpresp.Raw)
	result.Status = iutils.ToString(resp.Resp.StatusCode)
	result.Midware = result.Httpresp.Server
	result.Title = result.Httpresp.Title
	result.HasTitle = result.Httpresp.HasTitle
	result.AddExtracts(Extractors.Extract(string(result.Httpresp.Raw), true))
}

// GetStatusCode 从socket中获取http状态码
func GetStatusCode(content []byte) (bool, string) {
	if len(content) > 12 && bytes.HasPrefix(content, []byte("HTTP")) {
		return true, string(content[9:12])
	}
	return false, "tcp"
}

func FormatCertDomains(domains []string) []string {
	var hosts []string
	for _, domain := range domains {
		if strings.HasPrefix(domain, "*.") {
			domain = strings.Trim(domain, "*.")
		}
		hosts = append(hosts, domain)
	}
	return iutils.StringsUnique(hosts)
}

func CollectTLS(result *Result, resp *http.Response) {
	result.Host = strings.Join(resp.TLS.PeerCertificates[0].DNSNames, ",")
	if len(resp.TLS.PeerCertificates[0].DNSNames) > 0 {
		// 经验公式: 通常只有cdn会绑定超过2个host, 正常情况只有一个host或者带上www的两个host
		result.HttpHosts = append(result.HttpHosts, FormatCertDomains(resp.TLS.PeerCertificates[0].DNSNames)...)
	}
}
