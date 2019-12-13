// +build !windows

// TODO: Make all these methods POST

package v1

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"zxq.co/ripple/rippleapi/common"
)

// MetaRestartGET restarts the API with Zero Downtime™.
func MetaRestartGET(md common.MethodData) common.CodeMessager {
	proc, err := os.FindProcess(syscall.Getpid())
	if err != nil {
		return common.SimpleResponse(500, "couldn't find process. what the fuck?")
	}
	go func() {
		time.Sleep(time.Second)
		proc.Signal(syscall.SIGUSR2)
	}()
	return common.SimpleResponse(200, "brb")
}

var upSince = time.Now()

type metaUpSinceResponse struct {
	common.ResponseBase
	Code  int   `json:"code"`
	Since int64 `json:"since"`
}

// MetaUpSinceGET retrieves the moment the API application was started.
// Mainly used to get if the API was restarted.
func MetaUpSinceGET(md common.MethodData) common.CodeMessager {
	return metaUpSinceResponse{
		Code:  200,
		Since: int64(upSince.UnixNano()),
	}
}

// MetaUpdateGET updates the API to the latest version, and restarts it.
func MetaUpdateGET(md common.MethodData) common.CodeMessager {
	if f, err := os.Stat(".git"); err == os.ErrNotExist || !f.IsDir() {
		return common.SimpleResponse(500, "instance is not using git")
	}
	go func() {
		if !execCommand("git", "pull", "origin", "master") {
			return
		}
		// go get
		//        -u: update all dependencies
		//        -d: stop after downloading deps
		if !execCommand("go", "get", "-v", "-u", "-d") {
			return
		}
		if !execCommand("bash", "-c", "go build -v -ldflags \"-X main.Version=`git rev-parse HEAD`\"") {
			return
		}

		proc, err := os.FindProcess(syscall.Getpid())
		if err != nil {
			log.Println(err)
			return
		}
		proc.Signal(syscall.SIGUSR2)
	}()
	return common.SimpleResponse(200, "Started updating! "+surpriseMe())
}

func execCommand(command string, args ...string) bool {
	cmd := exec.Command(command, args...)
	cmd.Env = os.Environ()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
		return false
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return false
	}
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return false
	}
	data, err := ioutil.ReadAll(stderr)
	if err != nil {
		log.Println(err)
		return false
	}
	// Bob. We got a problem.
	if len(data) != 0 {
		log.Println(string(data))
	}
	io.Copy(os.Stdout, stdout)
	cmd.Wait()
	stdout.Close()
	return true
}
