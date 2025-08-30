package engine

import (
	"bytes"
	"github.com/M09ic/go-ntlmssp"

	"github.com/chainreactors/gogo/v2/pkg"
	"github.com/chainreactors/utils/iutils"
)

var data = pkg.Decode("YmXgZhZgYGCoYNBgYGZgYNghsAPEZWAEY0aGBSAGAwPDAQjlBiJYYju6XsucFJz/goNBW8AjgYmBgYGLCaLAL8THNzg4AKyfvYljEQMaYGPcKMvAwMAPAAAA//8=")

func WMIScan(opt *pkg.RunnerOption, result *pkg.Result) {
	result.Port = "135"
	target := result.GetTarget()
	conn, err := pkg.NewSocket("tcp", target, opt.Delay)
	if err != nil {
		return
	}
	defer conn.Close()

	result.Open = true
	ret, err := conn.Request(data, 4096)
	if err != nil {
		return
	}

	if bytes.Index(ret, []byte("NTLMSSP")) != -1 {
		result.Protocol = "wmi"
		result.Status = "wmi"
		result.AddNTLMInfo(iutils.ToStringMap(ntlmssp.NTLMInfo(ret)), "wmi")
	}
}
