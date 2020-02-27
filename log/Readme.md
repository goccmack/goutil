# Readme

Package log supports logging to a managed set of log files.

The logger is configured by a JSON file called `<component>.log.config`. `<component>` is the name of
the go binary executable (os.Executable()).
The logger looks for the
log config file first in the snooker
configuration directory (if it exists) and then in \$PWD
(see <https://github.com/goccmack/projects/GOUTIL/repos/snookeros/browse> for details).

The following is a JSON configuration structure containing the default logger paramerters:

    {
    	"RootDir": "/usr/local/var/log",
    	"NumFiles": 3,
    	"FileNumBytes": 1000000,
    	"Priority": "INFO",
    	"SuppressedFiles": ""
    }

If the working directory does not contain a log.config file the logger uses these parameters. All
fields of log.config are optional. The logger will used default values for missing parameters.

The logger reads log.config periodically. The logger uses changed parameters. The log.config can
be changed while the program is running and further logging reflects the changed log.config.

The logger initialises and closes automatically.

log.Panic(...) and log.Panicf(...) log a stack trace in addition to the log message, close the
open log file and call os.Exit(1).

The logger will automatically tag every log message with time, the source file and
line number of the call to log.

If the logger priority is DEBUG the logger can be configured to suppress the debug messages from one
or more files by providing a comma separated string to the suppressFilesDebug parameter
of log.Init(...). The components of the string correspond to file names without extension. E.g.:
"pkga,pkgb" will suppress debug messages from pkga.go and pkgb.go.

Package log supports four Priority levels in decreasing order of priority:
Panic, Warning, Info, Debug. The logger instance has two methods to log a message of each Priority:
<Priority> and <Priority>f, e.g.: Info and Infof. <Priority> takes a string parameter, while
<Priority>f takes a format string followed by a list of parameters. The format of the <Priorty>f
format string parameter is the same as for fmt.Printf.

The logger is initialised with the minimum priority level. The logger discards all messages logged with a
lower priority.

The directory log/examples demonstrates some featurs of log:

    examples/basic:
    	Shows the basic usage of the log package.

    examples/set_config:
    	Shows a change of log file size and debug level on a live program.

    examples/suppress_files:
    	Shows the suppression of debug message from a selection of files

examples/basic:

    log.config:
    {
    	"RootDir": "logs",
    }

Rootdir sets the logging directory to the directory "logs" under the working directory.
The other Config parameters will take their default values.

    package main

    import (
    	"github.com/goccmack/goutil/log/example/pkga"

    	"github.com/goccmack/goutil/log"
    )

    package main

    import (
    	"github.com/goccmack/goutil/log/examples/basic/pkga"

    	"github.com/goccmack/goutil/log"
    )

    func main() {
    	// Log something
    	log.Info("This message WILL appear in the log")

    	pkga.Go()
    	// Give pkga time to log a message
    	time.Sleep(time.Millisecond)

    	log.Debug("This message will NOT appear in the log")
    	log.Panic("This is a panic")
    }

    ======================

    package pkga

    import (
    	"github.com/goccmack/goutil/log"
    )

    func Go() {
    	// Log something
    	log.Debug("This is pkga")
    }

The resulting log is:

    File set configuration @ 2018-12-20T12:07:58.555278+01:00
    Maximum file size 1000000 bytes
    Maximum 3 files
    2018-12-20T12:07:58.555087+01:00 Log configuration:
    RootDir: logs
    NumFiles: 3
    NumBytes: 1000000
    Priority: INFO
    Suppress:
    2018-12-20T12:07:58.555464+01:00 [INFO] -main.go, line 14- This message WILL appear in the log
    2018-12-20T12:07:58.555482+01:00 [INFO] -pgka.go, line 9- This is pkga
    2018-12-20T12:07:58.555495+01:00 [PANIC] -main.go, line 19- This is a panic
    goroutine 1 [running]:
    github.com/goccmack/goutil/log.getPanicStackTrace(0x10bea8e, 0x3)
    	.../github.com/goccmack/goutil/log/interface.go:138 +0x74
    github.com/goccmack/goutil/log.Panic(0x10fac35, 0xf)
    	.../github.com/goccmack/goutil/log/interface.go:96 +0x22
    main.main()
    	.../github.com/goccmack/goutil/log/examples/basic/main.go:19 +0x8d
