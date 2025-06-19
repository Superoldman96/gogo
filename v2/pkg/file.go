package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/chainreactors/utils/encode"
	"github.com/chainreactors/utils/fileutils"
	"path"
	"strings"

	"github.com/chainreactors/logs"
	"github.com/chainreactors/utils/iutils"
)

func WriteSmartResult(file *fileutils.File, target string, ips []string) {
	var m map[string][]string = map[string][]string{}
	m[target] = ips
	marshal, err := json.Marshal(m)
	if err != nil {
		logs.Log.Error(err.Error())
	}

	if file != nil {
		file.SyncWrite(string(marshal) + "\n")
	}
}

func newFile(filename string, compress bool) (*fileutils.File, error) {
	file, err := fileutils.NewFile(filename, fileutils.ModeAppend, compress, false)
	if err != nil {
		return nil, err
	}

	var cursor int

	file.Encoder = func(i []byte) []byte {
		bs := encode.XorEncode(encode.MustDeflateCompress(i), fileutils.Key, cursor)
		cursor += len(bs)
		return bs
	}
	return file, nil
}

func getAutoFilename(config *Config, outtype string) string {
	var basename string
	target := strings.Replace(config.GetTargetName(), "/", ".", -1)
	target = strings.Replace(target, ":", "", -1)
	target = strings.Replace(target, "\\", "_", -1)
	if len(target) > 10 {
		if i := strings.IndexAny(target, "_"); i != -1 {
			target = target[:i]
		}
	}
	ports := strings.Replace(config.Ports, ",", "_", -1)
	basename = fmt.Sprintf("%s_%s_%s_%s", target, ports, config.Mod, outtype)
	return basename
}

var fileint = 1

func GetFilename(config *Config, name string) string {
	var basename string
	var basepath string

	if config.FilePath == "" {
		basepath = iutils.GetExcPath()
	} else {
		basepath = config.FilePath
	}

	if config.Filename != "" {
		return config.Filename
	}

	if config.Filenamef == "auto" {
		basename = path.Join(basepath, "."+getAutoFilename(config, name)+".dat")
	} else if config.Filenamef == "hidden" {
		if Win {
			basename = path.Join(basepath, "App_1634884664021088500_EC1B25B2-943.dat")
		} else {
			basename = path.Join(basepath, ".systemd-private-701215aa82634")
		}
	} else if config.Filenamef == "clear" {
		basename = path.Join(basepath, getAutoFilename(config, name)+".txt")
	} else {
		return config.Filename
	}

	if !fileutils.IsExist(basename) {
		return basename
	}

	for fileutils.IsExist(basename + iutils.ToString(fileint)) {
		fileint++
	}
	return basename + iutils.ToString(fileint)
}
