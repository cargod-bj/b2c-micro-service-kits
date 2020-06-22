package log

import (
	"fmt"
	"github.com/cargod-bj/b2c-common/utils"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	oplogging "github.com/op/go-logging"
	"io"
	"os"
	"strings"
	"time"
)

var (
	defaultFormatter = `%{time:2006/01/02 - 15:04:05.000} %{longfile} %{color:bold}▶ [%{level:.6s}] %{message}%{color:reset}`
	confModule       string
	confPrefix       string
	confStdoutLevel  string
	confIsLog2Stdout bool
	confFileLevel    string
	confLogDir       string
	confIsLog2File   bool
	Logger           *oplogging.Logger
	hasInit          bool = false
)

// 初始化log模块。
//
// ## 参数介绍
// 	1. module log的模块名称
// 	2. prefix log的前缀
// 	3. isLog2Stdout 是否同步输出到 命令行
// 	4. stdOutLevel 输出到 命令行 的级别，支持："DEBUG"，"INFO"，"NOTICE"，"WARNING"，"ERROR"，"CRITICAL"，如果为空则默认为 DEBUG
// 	5. isLog2File 是否同步输出到 文件
// 	6. fileLevel 输出到 文件 的级别，支持："DEBUG"，"INFO"，"NOTICE"，"WARNING"，"ERROR"，"CRITICAL"，如果为空则默认为 DEBUG
// 	6. logDir 输出到 文件 的文件目录，如果为空则不能输出到文件
//
func InitLog(module, prefix, stdOutLevel, fileLevel, logDir string, isLog2Stdout, isLog2File bool) {
	confModule = module
	confPrefix = prefix
	confStdoutLevel = stdOutLevel
	confFileLevel = fileLevel
	confIsLog2Stdout = isLog2Stdout
	confLogDir = logDir
	confIsLog2File = isLog2File
	if confFileLevel == "" {
		confFileLevel = "DEBUG"
	}
	if confStdoutLevel == "" {
		confStdoutLevel = "DEBUG"
	}
	if confLogDir == "" && confIsLog2File {
		panic("the logDir is nil and the isLog2File is true, they are conflict")
	}

	if confPrefix == "" {
		panic("the prefix must be not empty")
	}

	Logger = oplogging.MustGetLogger(confModule)
	Logger.ExtraCalldepth = 1
	var backends []oplogging.Backend
	backends = registerStdout(backends)
	backends = registerFile(backends)

	backendLeveled := oplogging.AddModuleLevel(oplogging.MultiLogger(backends...))
	Logger.SetBackend(backendLeveled)

	hasInit = true
}

func SetLog(logger *oplogging.Logger) {
	hasInit = true
	Logger = logger
}

func registerStdout(backends []oplogging.Backend) []oplogging.Backend {
	if !confIsLog2Stdout {
		return backends
	}
	level, err := oplogging.LogLevel(confStdoutLevel)
	if err != nil {
		fmt.Println(err)
	}
	backends = append(backends, createBackend(os.Stdout, level))

	return backends
}

func registerFile(backends []oplogging.Backend) []oplogging.Backend {
	if confLogDir == "" {
		return backends
	}
	fmt.Println("log directory:" + confLogDir)
	tempDir := confLogDir
	if ok, _ := utils.PathExists(tempDir); !ok {
		// directory not exist
		currentDir, _ := os.Getwd()
		tempDir = currentDir + tempDir
		if strings.HasPrefix(tempDir, string(os.PathSeparator)) && strings.HasSuffix(currentDir, string(os.PathSeparator)) {
			tempDir = currentDir + tempDir[1:]
		}
		fmt.Println("create log directory:" + tempDir)
		err := os.MkdirAll(tempDir, os.ModePerm)
		if err != nil {
			fmt.Println("create log directory failed")
		}
	}
	logFileDir := tempDir + string(os.PathSeparator)
	fileWriter, err := rotatelogs.New(
		logFileDir+"%Y-%m-%d-%H-%M.log",
		// generate soft link, point to latest log file
		rotatelogs.WithLinkName(logFileDir+"latest_log.log"),
		// maximum time to save log files
		rotatelogs.WithMaxAge(7*24*time.Hour),
		// time period of log file switching
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	level, err := oplogging.LogLevel(confFileLevel)
	if err != nil {
		fmt.Println(err)
	}
	return append(backends, createBackend(fileWriter, level))
}

func createBackend(w io.Writer, level oplogging.Level) oplogging.Backend {
	backend := oplogging.NewLogBackend(w, confPrefix, 0)
	//backend := oplogging.NewLogBackend(w, confPrefix, log.LstdFlags | log.Lmicroseconds | log.Llongfile)
	stdoutWriter := false
	if w == os.Stdout {
		stdoutWriter = true
	}
	format := getLogFormatter(stdoutWriter)
	backendLeveled := oplogging.AddModuleLevel(oplogging.NewBackendFormatter(backend, format))
	backendLeveled.SetLevel(level, confModule)
	return backendLeveled
}

func getLogFormatter(stdoutWriter bool) oplogging.Formatter {
	pattern := defaultFormatter
	if !stdoutWriter {
		// Color is only required for console output
		// Other writers don't need %{color} tag
		pattern = strings.Replace(pattern, "%{color:bold}", "", -1)
		pattern = strings.Replace(pattern, "%{color:reset}", "", -1)
	}
	if !confIsLog2File {
		// Remove %{logfile} tag
		pattern = strings.Replace(pattern, "%{longfile}", "", -1)
	}
	return oplogging.MustStringFormatter(pattern)
}
