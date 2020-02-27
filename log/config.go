package log

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	logConfigFileSuffix = "logging.config"
)

type jsonConfig struct {
	RootDir         string `json:",omitempty"`
	NumFiles        *int   `json:",omitempty"`
	FileNumBytes    *int   `json:",omitempty"`
	Priority        string `json:",omitempty"`
	SuppressedFiles string `json:",omitempty"`
}

// Config contains the logger configuration. It is read from the JSON file log.config in the
// working directory of from the default values if log.config does not exist.
type Config struct {
	RootDir      string
	FileName     string
	NumFiles     int
	FileNumBytes int
	Priority     Priority
	// comma separated list of files whose DEBUG messages are suppressed
	SuppressedFiles string
}

// Clone returns a deep copy of c
func (c *Config) Clone() *Config {
	return &Config{
		RootDir:         c.RootDir,
		FileName:        c.FileName,
		NumFiles:        c.NumFiles,
		FileNumBytes:    c.FileNumBytes,
		Priority:        c.Priority,
		SuppressedFiles: c.SuppressedFiles,
	}
}

// Equal returns true iff all fields of c are equal to c1
func (c *Config) Equal(c1 *Config) bool {
	if c.RootDir != c1.RootDir ||
		c.FileName != c1.FileName ||
		c.NumFiles != c1.NumFiles ||
		c.FileNumBytes != c1.FileNumBytes ||
		c.Priority != c1.Priority {

		return false
	}

	return true
}

// String returns a formatted string of c.
func (c *Config) String() string {
	return fmt.Sprintf("Config{%s,%s,%d,%d,%s}",
		c.RootDir, c.FileName, c.NumFiles, c.FileNumBytes, c.Priority)
}

// ToJSON returns the JSON format of log.config of c.
// The following can be used to generate the default JSON for log.config:
//		fmt.Println(log.DefaultConfig().ToJSON())
//
// 		{
// 		    "RootDir": "/usr/local/var/log",
// 		    "NumFiles": 3,
// 		    "FileNumBytes": 1000000,
// 		    "Priority": "INFO",
// 		    "SuppressedFiles": ""
// 		}
func (c *Config) ToJSON() string {
	jc := &jsonConfig{
		RootDir:      c.RootDir,
		NumFiles:     &c.NumFiles,
		FileNumBytes: &c.FileNumBytes,
		Priority:     c.Priority.String(),
	}
	b, err := json.Marshal(jc)
	if err != nil {
		panic(err)
	}
	var out bytes.Buffer
	json.Indent(&out, b, "", "    ")
	return out.String()
}

var (
	// Name of the executeable file; will be used to create logfile names
	fileName = getFileName()
)

const (
	// DefaultLogRootDir determines the directory for log files if not specified in log.config
	DefaultLogRootDir = "/usr/local/var/log"
	// DefaultNumFiles determines the maximum number of log files if not specified in log.config
	DefaultNumFiles = 3
	// DefaultLogFileNumBytes determines the maximum log file size if not specified in log.config
	DefaultLogFileNumBytes = 1000000
	// DefaultPriority determines the logger priority if not specified in log.config
	DefaultPriority = INFO
	// DefaultSuppressedFiles determines the suppressed files if not specified in log.config
	DefaultSuppressedFiles = ""
)

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		RootDir:         DefaultLogRootDir,
		FileName:        fileName,
		NumFiles:        DefaultNumFiles,
		FileNumBytes:    DefaultLogFileNumBytes,
		Priority:        DefaultPriority,
		SuppressedFiles: DefaultSuppressedFiles,
	}
}

func getConfigFile() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	files, err := ioutil.ReadDir(cwd)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), logConfigFileSuffix) {
			return f.Name()
		}
	}
	return ""
}

func getFileName() string {
	pth, err := os.Executable()
	if err != nil {
		panic(err)
	}
	_, fname := path.Split(pth)
	return fname
}

func jsonToConfig(jc *jsonConfig) *Config {
	c := new(Config)
	if jc.RootDir == "" {
		c.RootDir = DefaultLogRootDir
	} else {
		c.RootDir = jc.RootDir
	}
	c.FileName = fileName
	if jc.NumFiles == nil {
		c.NumFiles = DefaultNumFiles
	} else {
		c.NumFiles = *jc.NumFiles
	}
	if jc.FileNumBytes == nil {
		c.FileNumBytes = DefaultLogFileNumBytes
	} else {
		c.FileNumBytes = *jc.FileNumBytes
	}
	if jc.Priority == "" {
		c.Priority = DefaultPriority
	} else {
		if p, err := ToPriority(jc.Priority); err != nil {
			fmt.Fprintf(os.Stderr, "Invalid priority string: %s\n", jc.Priority)
			c.Priority = DefaultPriority
		} else {
			c.Priority = p
		}
	}
	c.SuppressedFiles = jc.SuppressedFiles
	return c
}

func readConfigFile(warnIfNoCfg bool) *Config {
	cfgFile := getConfigFile()
	if cfgFile == "" {
		fmt.Fprintln(os.Stderr, "No logging config file found. Using defaults")
		return DefaultConfig()
	}

	data, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		if warnIfNoCfg {
			fmt.Fprintf(os.Stderr, "Warning reading %s: %s\n", cfgFile, err)
		}
		return DefaultConfig()
	}
	jc := new(jsonConfig)
	if err := json.Unmarshal(data, &jc); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %s\n", cfgFile, err)
		return DefaultConfig()
	}
	c := jsonToConfig(jc)
	return c

}
