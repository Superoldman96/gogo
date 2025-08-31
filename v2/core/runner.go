package core

import (
	"fmt"
	"github.com/chainreactors/fingers/common"
	"github.com/chainreactors/fingers/fingers"
	. "github.com/chainreactors/gogo/v2/engine"
	"github.com/chainreactors/logs"
	"github.com/chainreactors/parsers"
	"github.com/chainreactors/proxyclient"
	"github.com/chainreactors/utils"
	"github.com/chainreactors/utils/encode"
	"github.com/chainreactors/utils/fileutils"
	"golang.org/x/net/context"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	. "github.com/chainreactors/gogo/v2/pkg"
	neuhttp "github.com/chainreactors/neutron/protocols/http"
	"github.com/chainreactors/utils/iutils"
)

func NewRunner() *Runner {
	return &Runner{
		Config: Config{},
	}
}

var ver = ""
var logFile = ".sock.lock"

type Runner struct {
	MiscOption    `group:"Miscellaneous Options"`
	InputOption   `group:"Input Options"`
	OutputOption  `group:"Output Options"`
	SmartOption   `group:"Smart Options"`
	AdvanceOption `group:"Advance Options"`
	ConfigOption  `group:"Configuration Options"`

	start  time.Time
	Config Config
}

func (r *Runner) Prepare() bool {
	// 初始化日志工具
	if r.Quiet {
		logs.Log = logs.NewLogger(0)
		logs.Log.SetQuiet(true)
	} else {
		if r.Debug {
			logs.Log.SetLevel(logs.Debug)
		}
		logs.Log.SetFile(logFile)
		logs.Log.Init()
	}

	ExecuterOptions.Options.Timeout = r.Delay + r.HttpsDelay
	HttpTimeout = time.Duration(r.Delay+r.HttpsDelay) * time.Second
	Opt.PluginDebug = r.PluginDebug
	common.NoGuess = r.NoGuess
	fileutils.Key = []byte(r.Key)

	// 一些特殊的分支, 这些分支将会直接退出程序
	if r.Ver {
		fmt.Println(ver)
		return false
	}

	r.PrepareConfig()

	if r.Exclude != "" {
		r.Config.Excludes = utils.ParseCIDRs(strings.Split(r.Exclude, ","))
	} else if r.ExcludeList != "" {
		ips, err := fileutils.LoadFileToSlice(r.ExcludeList)
		if err != nil {
			logs.Log.Error(err.Error())
			return false
		}
		r.Config.Excludes = utils.ParseCIDRs(ips)
	}

	if r.Config.Excludes != nil {
		r.Config.RunnerOpt.ExcludeCIDRs = r.Config.Excludes
	}
	if r.FormatterFilename != "" {
		LoadNeutron("")
		var formatOut string
		if r.Outputf != Default {
			formatOut = r.Outputf
		} else if r.FileOutputf != Default {
			formatOut = r.FileOutputf
		} else {
			// 根据文件后缀自动识别格式
			autoFormat := getFormatByFileExtension(r.Config.Filename)
			switch autoFormat {
			case "csv", "json":
				formatOut = autoFormat
			default:
				formatOut = "color"
			}
		}
		FormatOutput(r.FormatterFilename, r.Config.Filename, formatOut, r.Config.Filenamef, r.Filters, r.FilterOr)
		return false
	}
	// 输出 Config
	if r.Printer != "" {
		printConfigs(r.Printer)
		return false
	}

	if len(r.Proxy) != 0 {
		var proxies []*url.URL
		for _, u := range r.Proxy {
			uri, err := url.Parse(u)
			if err != nil {
				logs.Log.Warnf("parse proxy error %s, skip proxy!", err.Error())
			} else {
				proxies = append(proxies, uri)
			}
		}
		dialer, err := proxyclient.NewClientChain(proxies)
		if err != nil {
			logs.Log.Warnf("parse proxy error %s, skip proxy!", err.Error())
		}
		neuhttp.DefaultTransport.DialContext = dialer.DialContext
		DefaultTransport.DialContext = dialer.DialContext
		ProxyDialTimeout = func(network, address string, duration time.Duration) (net.Conn, error) {
			ctx, _ := context.WithTimeout(context.Background(), duration)
			return dialer.DialContext(ctx, network, address)
		}
	}
	return true
}

func (r *Runner) Init() error {
	// 加载配置文件中的全局变量
	err := LoadPortConfig(r.PortConfig)
	if err != nil {
		return fmt.Errorf("load port config failed, %v", err)
	}

	err = LoadFinger(r.FingerFile)
	if err != nil {
		return fmt.Errorf("load finger config failed, %v", err)
	}

	err = LoadExtractor()
	if err != nil {
		return fmt.Errorf("extract config load FAIL!, %s", err.Error())
	}
	for _, e := range r.Extract {
		if reg, ok := ExtractRegexps[e]; ok {
			Extractors[e] = reg
		} else {
			Extractors[e] = []*parsers.Extractor{
				&parsers.Extractor{
					Name:            e,
					CompiledRegexps: []*regexp.Regexp{regexp.MustCompile(e)},
				},
			}
		}
	}

	if r.AttackType != "" {
		ExecuterOptions.Options.AttackType = r.AttackType
	}
	NeutronLoader(r.ExploitFile, r.Payloads)

	if r.Opsec {
		fingers.OPSEC = true
		r.Config.RunnerOpt.Opsec = true
		ExecuterOptions.Options.Opsec = true
	}
	return nil
}

func (r *Runner) PrepareConfig() {
	r.Config = Config{
		GOGOConfig: &parsers.GOGOConfig{
			IP:        r.IP,
			Ports:     r.Ports,
			ListFile:  r.ListFile,
			JsonFile:  r.JsonFile,
			Threads:   r.Threads,
			PortSpray: r.PortSpray,
			Mod:       r.Mod,
		},
		RunnerOpt: &RunnerOption{
			Delay:        r.Delay,
			HttpsDelay:   r.HttpsDelay,
			VersionLevel: setVersionLevel(r.Verbose),
			Exploit:      setExploit(r.ExploitName, r.Exploit),
			Debug:        r.Debug,
			Opsec:        r.Opsec,
			//SuffixStr:  r.SuffixStr,
		},
		IsListInput: r.IsListInput,
		IsJsonInput: r.IsJsonInput,
		PortProbe:   r.PortProbe,
		IpProbe:     r.IpProbe,
		NoSpray:     r.NoSpray,
		Filename:    r.Filename,
		FilePath:    r.FilePath,
		Compress:    !r.Compress,
		Tee:         r.Tee,
		Filters:     r.Filters,
		FilterOr:    r.FilterOr,
	}

	if r.FileOutputf == Default {
		// 根据文件后缀自动识别格式
		autoFormat := getFormatByFileExtension(r.Filename)
		if autoFormat == "dat" {
			r.Config.FileOutputf = "jl" // dat格式默认使用jsonlines
		} else {
			r.Config.Compress = false
			r.Config.FileOutputf = autoFormat
		}
	} else {
		r.Config.FileOutputf = r.FileOutputf
	}

	if r.Outputf == Default {
		r.Config.Outputf = "full"
	} else {
		r.Config.Outputf = r.Outputf
	}

	for _, filterStr := range r.OutputFilters {
		k, v, op := parseFilterString(filterStr)
		if op != "" {
			r.Config.OutputFilters = append(r.Config.OutputFilters, []string{k, v, op})
		}
	}

	for _, filterStr := range r.ScanFilters {
		k, v, op := parseFilterString(filterStr)
		if op != "" {
			r.Config.RunnerOpt.ScanFilters = append(r.Config.RunnerOpt.ScanFilters, []string{k, v, op})
		}
	}

	if r.AutoFile {
		r.Config.Filenamef = "auto"
	} else if r.HiddenFile {
		r.Config.Filenamef = "hidden"
	}

	if r.Ping {
		r.Config.AliveSprayMod = append(r.Config.AliveSprayMod, "icmp")
	}

}

func (r *Runner) Run() {
	r.start = time.Now()
	if r.WorkFlowName == "" && !r.IsWorkFlow {
		r.runWithCMD()
	} else {
		var workflowMap = WorkflowMap{}
		if r.IsWorkFlow {
			workflowMap["tmp"] = ParseWorkflowsFromInput(LoadFile(os.Stdin))
			r.WorkFlowName = "tmp"
		} else if fileutils.IsExist(r.WorkFlowName) {
			file, err := fileutils.Open(r.WorkFlowName)
			if err != nil {
				iutils.Fatal(err.Error())
			}
			workflowMap["tmp"] = ParseWorkflowsFromInput(LoadFile(file))
			r.WorkFlowName = "tmp"
		} else {
			if bs, ok := encode.DSLParser(r.WorkFlowName); ok {
				workflowMap["tmp"] = ParseWorkflowsFromInput(bs)
			} else {
				workflowMap = LoadWorkFlow()
			}
		}
		r.runWithWorkFlow(workflowMap)
	}
}

func (r *Runner) runWithCMD() {
	config := r.Config

	if config.Filename != "" {
		logs.Log.Warn("The result file has been specified, other files will not be created.")
	}

	if config.Filename == "" && config.IsBSmart() {
		config.SmartBFilename = GetFilename(&config, "bcidr")
	}
	if config.Filename == "" && config.IsSmart() {
		config.SmartCFilename = GetFilename(&config, "ccidr")
	}
	if config.Filename == "" && config.HasAlivedScan() {
		config.AlivedFilename = GetFilename(&config, "alived")
	}

	if config.Filenamef != "" {
		config.Filename = GetFilename(&config, config.FileOutputf)
	}

	preparedConfig, err := InitConfig(&config)
	if err != nil {
		iutils.Fatal(err.Error())
	}
	RunTask(*preparedConfig) // 运行
	r.Close(&config)
}

func (r *Runner) runWithWorkFlow(workflowMap WorkflowMap) {
	if workflows := workflowMap.Choice(r.WorkFlowName); len(workflows) > 0 {
		for _, workflow := range workflows {
			logs.Log.Important("workflow " + workflow.Name + " starting")
			config := workflow.PrepareConfig(r.Config)

			if config.Mod == SUPERSMARTB {
				config.FileOutputf = SUPERSMARTB
			}

			if config.Filename != "" {
				logs.Log.Warn("The result file has been specified, other files will not be created.")
			}

			if config.Filename == "" && config.IsBSmart() {
				config.SmartBFilename = GetFilename(config, "bcidr")
			}
			if config.Filename == "" && config.IsSmart() {
				config.SmartCFilename = GetFilename(config, "ccidr")
			}
			if config.Filename == "" && config.HasAlivedScan() {
				config.AlivedFilename = GetFilename(config, "alived")
			}
			if config.Filenamef != "" {
				config.Filename = GetFilename(config, config.FileOutputf)
			}

			if r.Verbose != nil {
				r.Config.RunnerOpt.VersionLevel = setVersionLevel(r.Verbose)
			} else {
				r.Config.RunnerOpt.VersionLevel = workflow.Verbose
			}

			if r.Config.RunnerOpt.Exploit != "none" {
				r.Config.RunnerOpt.Exploit = setExploit(r.ExploitName, r.Exploit)
			} else {
				r.Config.RunnerOpt.Exploit = workflow.Exploit
			}

			preparedConfig, err := InitConfig(config)
			if err != nil {
				iutils.Fatal(err.Error())
			}
			RunTask(*preparedConfig) // 运行
			r.Close(config)
		}
	} else {
		iutils.Fatal("not fount workflow " + r.WorkFlowName)
	}
}

func (r *Runner) Close(config *Config) {
	config.Close() // 关闭result与extract写入管道

	if r.HiddenFile {
		iutils.Chtime(config.Filename)
		if config.SmartBFile != nil && config.SmartBFile.IsInitialized() {
			iutils.Chtime(config.SmartBFilename)
		}
		if config.SmartCFile != nil && config.SmartCFile.IsInitialized() {
			iutils.Chtime(config.SmartBFilename)
		}
	}

	// 任务统计
	logs.Log.Importantf("Alived: %d, Total: %d", Opt.AliveSum, RunSum)
	logs.Log.Important("Time consuming: " + time.Since(r.start).String())

	// 输出文件名
	if config.File != nil && config.File.IsInitialized() {
		logs.Log.Importantf("Results: " + config.Filename)
	}
	if config.SmartBFile != nil && config.SmartBFile.IsInitialized() {
		logs.Log.Important("B CIDRs result: " + config.SmartBFilename)
	}
	if config.SmartCFile != nil && config.SmartCFile.IsInitialized() {
		logs.Log.Important("c CIDRs result: " + config.SmartCFilename)
	}
	if config.AliveFile != nil && config.AliveFile.IsInitialized() {
		logs.Log.Important("Alived result: " + config.AlivedFilename)
	}
	if fileutils.IsExist(config.Filename + "_extract") {
		logs.Log.Important("extractor result: " + config.Filename + "_extract")
	}
}

func printConfigs(t string) {
	if t == "port" {
		Printportconfig()
	} else if t == "nuclei" {
		PrintNeutronPoc()
	} else if t == "workflow" {
		PrintWorkflow()
	} else if t == "extract" {
		PrintExtract()
	} else {
		fmt.Println("choice port|nuclei|workflow|extract")
	}
}

func NeutronLoader(pocfile string, payloads []string) {
	ExecuterOptions.Options.VarsPayload = ParserCmdPayload(payloads)
	TemplateMap = LoadNeutron(pocfile)
}

func parseFilterString(s string) (k, v, op string) {
	if strings.Contains(s, "::") {
		kv := strings.Split(s, "::")
		return kv[0], kv[1], "::"
	} else if strings.Contains(s, "==") {
		kv := strings.Split(s, "==")
		return kv[0], kv[1], "=="
	} else if strings.Contains(s, "!=") {
		kv := strings.Split(s, "!=")
		return kv[0], kv[1], "!="
	} else if strings.Contains(s, "!:") {
		kv := strings.Split(s, "!:")
		return kv[0], kv[1], "!:"
	}
	return "", "", ""
}

func setVersionLevel(v []bool) int {
	if len(v) == 1 {
		return 1
	} else if len(v) == 2 {
		return 2
	}
	return 0
}

func setExploit(name string, enable bool) string {
	if name != "" {
		return name
	} else if enable {
		return "auto"
	}
	return "none"
}

// getFormatByFileExtension 根据文件扩展名自动识别输出格式
// 支持的格式: txt, csv, json, dat(默认)
func getFormatByFileExtension(filename string) string {
	if filename == "" {
		return "dat" // 默认格式
	}

	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".txt":
		return "full"
	case ".csv":
		return "csv"
	case ".json":
		return "jl"
	case ".dat":
		return "dat"
	default:
		return "dat" // 默认格式
	}
}
