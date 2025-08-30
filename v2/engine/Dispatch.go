package engine

import (
	"sync/atomic"

	"github.com/chainreactors/gogo/v2/pkg"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/utils/iutils"
)

var (
	RunSum int32
)

func Dispatch(opt *pkg.RunnerOption, result *pkg.Result) {
	defer func() {
		if err := recover(); err != nil {
			logs.Log.Errorf("scan %s unexcept error, %v", result.GetTarget(), err)
			panic(err)
		}
	}()
	atomic.AddInt32(&RunSum, 1)
	if opt.ExcludeCIDRs != nil && opt.ExcludeCIDRs.ContainsString(result.Ip) {
		logs.Log.Debug("exclude ip: " + result.Ip)
		return
	}
	if result.Port == "137" || result.Port == "nbt" {
		NBTScan(opt, result)
		return
	} else if result.Port == "135" || result.Port == "wmi" {
		WMIScan(opt, result)
		return
	} else if result.Port == "oxid" {
		OXIDScan(opt, result)
		return
	} else if result.Port == "icmp" || result.Port == "ping" {
		ICMPScan(opt, result)
		return
	} else if result.Port == "snmp" || result.Port == "161" {
		SNMPScan(opt, result)
		return
	} else if result.Port == "445" || result.Port == "smb" {
		SMBScan(opt, result)
		if opt.Exploit == "ms17010" {
			MS17010Scan(opt, result)
		} else if opt.Exploit == "smbghost" || opt.Exploit == "cve-2020-0796" {
			SMBGhostScan(opt, result)
		} else if opt.Exploit == "auto" || opt.Exploit == "smb" {
			MS17010Scan(opt, result)
			SMBGhostScan(opt, result)
		}
		return
	} else if result.Port == "mssqlntlm" {
		MSSqlScan(opt, result)
		return
	} else if result.Port == "winrm" {
		WinrmScan(opt, result)
		return
	} else {
		InitScan(opt, result)
	}

	if !result.Open || result.SmartProbe {
		// 启发式探针或端口未OPEN,则直接退出, 不进行后续扫描
		return
	}

	// 指纹识别, 会根据versionlevel自动选择合适的指纹
	if result.IsHttp {
		HTTPFingerScan(opt, result)
	} else {
		SocketFingerScan(opt, result)
	}

	if result.Filter(opt.ScanFilters) {
		// 如果被过滤, 则停止后续扫描深度扫描
		return
	}
	//主动信息收集
	if opt.VersionLevel > 0 && result.IsHttp {
		// favicon指纹只有-v大于0并且为http服务才启用
		if result.HttpHosts != nil {
			hostScan(opt, result)
		}

		FaviconScan(opt, result)
		if result.Status != "404" {
			NotFoundScan(opt, result)
		}
	} else {
		// 如果versionlevel为0 ,或者非http服务, 则使用默认端口猜测指纹.
		if !result.IsHttp && result.NoFramework() {
			// 通过默认端口号猜测服务,不具备准确性
			result.GuessFramework()
		}
	}

	// 如果exploit参数不为none,则进行漏洞探测
	if opt.Exploit != "none" {
		NeutronScan(opt, result.GetHostBaseURL(), result)
	}

	result.Title = iutils.AsciiEncode(result.Title)
	return
}
