package engine

import (
	"github.com/chainreactors/fingers/common"
	. "github.com/chainreactors/gogo/v2/pkg"
	"github.com/chainreactors/logs"
)

var (
	TCP = "tcp"
	UDP = "udp"
)

func SocketFingerScan(opt *RunnerOption, result *Result) {
	// 如果是http协议,则判断cms,如果是tcp则匹配规则库.暂时不考虑udp
	var closureResp, finalResp []byte
	callback := func(f *common.Framework, v *common.Vuln) {
		if f != nil {
			result.Frameworks.Add(f)
			finalResp = closureResp
		}
		if v != nil {
			result.Vulns.Add(v)
		}
	}
	tcpsender := func(sendData []byte) ([]byte, bool) {
		target := result.GetTarget()
		logs.Log.Debugf("active detect: %s, data: %q", target, sendData)
		conn, err := NewSocket(TCP, target, opt.Delay)
		if err != nil {
			logs.Log.Debugf("active detect %s error, %s", target, err.Error())
			return nil, false
		}
		defer conn.Close()

		data, err := conn.QuickRequest(sendData, 1024)
		if err != nil {
			return nil, false
		}
		closureResp = data
		return data, true
	}

	//udpsender := func(sendData []byte) ([]byte, bool) {
	//	target := result.GetTarget()
	//	logs.Log.Debugf("active detect: , data: ", target, sendData)
	//	conn, err := NewSocket(UDP, target, RunOpt.Delay)
	//	if err != nil {
	//		logs.Log.Debugf("active detect %s error, %s", target, err.Error())
	//		return nil, false
	//	}
	//	defer conn.Close()
	//
	//	data, err := conn.QuickRequest(sendData, 1024)
	//	if err != nil {
	//		return nil, false
	//	}
	//
	//	return data, true
	//}

	if opt.VersionLevel > 0 {
		FingerEngine.SocketMatch(result.Content, result.Port, opt.VersionLevel, tcpsender, callback)
	} else {
		if group, ok := FingerEngine.SocketGroup[result.Port]; ok {
			frames, _ := group.Match(result.ToContent(), 1, tcpsender, callback, true)
			if len(frames) == 0 {
				FingerEngine.SocketMatch(result.Content, "", opt.VersionLevel, nil, callback)
			}
		}
	}

	if finalResp != nil {
		CollectSocketResponse(result, finalResp)
	}
	return
}
