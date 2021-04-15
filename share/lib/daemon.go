package lib

import (
    "flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

var (
	APP_NAME      string
	APP_VERSION   string
	ENV_VAR_NAME  string
	ENV_VAR_VALUE string
	PID_NAME      string
)

func RunAsDaemon(name, version, envVarName, envVarValue, pidName string) {

	var nc, stop, restart, debug bool

	flag.Usage = func() {
		fmt.Printf(" Usage: %s [OPTIONS]\n\n", os.Args[0])
		fmt.Printf(" These are the optional arguments you can pass to " + APP_NAME + ":\n")
		fmt.Printf(" -help     -- this message\n")
		fmt.Printf(" -version  -- print the version and exit\n")
		fmt.Printf(" -debug    -- debug mode to run the program\n")
		fmt.Printf(" -start       -- do not output to a console and background\n")
		fmt.Printf(" -restart  -- restart program\n\n")
		fmt.Printf(" -stop     -- stop program\n")

	}

	APP_NAME = name
	APP_VERSION = version
	ENV_VAR_NAME = envVarName
	ENV_VAR_VALUE = envVarValue
	PID_NAME = pidName

	if len(os.Args) == 1 {
		flag.Usage()
		os.Exit(0)
	}

	// Check args
	for _, arg := range os.Args[1:] {
		if arg == "-debug" {
			debug = true
		}
		if arg == "-help" {
			flag.Usage()
			os.Exit(0)
		}
		if arg == "-version" {
			fmt.Println(APP_VERSION)
			os.Exit(0)
		}
		if arg == "-stop" {
			stop = true
		}
		if arg == "-nc" {
			nc = true
		}
		if arg == "-restart" {
			restart = true
		}
		if arg == "-g" {
			runtime.GOMAXPROCS(runtime.NumCPU())
		}
		if arg == "-g2" {
			runtime.GOMAXPROCS(runtime.NumCPU() * 2)
		}
	}

	if !debug && !stop && !nc && !restart {
		flag.Usage()
		os.Exit(0)
	}

	if debug && stop {
		flag.Usage()
		os.Exit(0)
	}

	if debug && !nc {
		return
	}

	if !CheckDir(GetPath("run/")) {

		fmt.Println("run directory is create Faild")
		os.Exit(0)
	}

	if !CheckDir(GetPath("logs/")) {
		fmt.Println("logs directory is create Faild")
		os.Exit(0)
	}

	if nc {
		if r, err := ioutil.ReadFile(GetPath("run/" + PID_NAME + ".pid")); err == nil {
			if pid, _ := strconv.Atoi(string(r)); pid > 0 {
				if _, err := os.FindProcess(pid); err == nil {
					fmt.Println("Already exists\n")
					os.Exit(0)
				} else {
					fmt.Println("%v", err)
				}
			}
		}
		daemon()
	}

	if stop {
		killProcess(PID_NAME)
		os.Exit(0)
	}

	if restart {
		killProcess(PID_NAME)
		reborn()
	}

	savePIDFile(PID_NAME, syscall.Getpid())
}

func daemon() {
	ppid := os.Getppid()
	if ppid == 1 {
		return
	}

	if !wasReborn() {
		reborn()
	}
}

func reborn() {
	var path string
	var err error
	if path, err = filepath.Abs(os.Args[0]); err != nil {
		return
	}

	cmd := exec.Command(path, "-nc")

	// Prepare environment variables
	envVar := fmt.Sprintf("%s=%s", ENV_VAR_NAME, ENV_VAR_VALUE)
	cmd.Env = append(os.Environ(), envVar)

	// Set output
	setOutput(cmd)

	if err := cmd.Start(); err != nil {
		os.Exit(-1)
	} else {
		fmt.Printf("Starting %s:    [ ok ]\n", APP_NAME)
	}
	os.Exit(0)
}

// func WasReborn, return true if the process has environment
// variable _GO_DAEMON=1 (child process).
func wasReborn() bool {
	return os.Getenv(ENV_VAR_NAME) == ENV_VAR_VALUE
}

func setOutput(cmd *exec.Cmd) {
	fileName := fmt.Sprintf("logs/%s.log", time.Now().Format("2006-01-02"))
	outFile, err := os.OpenFile(GetPath(fileName), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		cmd.Stdout = outFile
	}
}

func killProcess(name string) {
	if r, err := ioutil.ReadFile(GetPath("run/" + name + ".pid")); err == nil {
		pid, _ := strconv.Atoi(string(r))
		if pid > 0 {
			if p, err := os.FindProcess(pid); err == nil {
				if err = p.Signal(syscall.SIGKILL); err != nil {
					fmt.Printf("Failed to kill %s : %v, Error: %v\n", name, pid, err)
				} else {
					fmt.Printf("Stopping %s:    [ ok ]\n", APP_NAME)
				}
				removePIDFile(name)
			} else {
				fmt.Printf("%v\n", err)
			}
		}
	} else {
		fmt.Printf("Stopping %s:    [ Not Exist ]\n", APP_NAME)
	}
}

func savePIDFile(name string, pid int) {
	if err := ioutil.WriteFile(GetPath("run/"+name+".pid"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		fmt.Printf("Write %s PID file Error: %v\n", name, err)
	}
}

func removePIDFile(name string) {
	if err := os.Remove(GetPath("run/" + name + ".pid")); err != nil {
		fmt.Printf("Remove %s PID file Error: %v\n", name, err)
	}
}
