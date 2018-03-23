package utils

import (
	"log"
	"io/ioutil"
	"bytes"
	"os"
	"strings"
	"strconv"
	"io"
	"syscall"
	"path/filepath"
	"github.com/sirupsen/logrus"
)

type TelegrafReloader struct {
	name         string
	signal       syscall.Signal
	shouldReload bool
}

// Create new config and populate it from environment
func NewTelegrafReloader() *TelegrafReloader {
	return &TelegrafReloader{
		signal:       syscall.SIGHUP,
		shouldReload: false,
	}
}

func (tr *TelegrafReloader) Name() string {
	return tr.name
}

func (tr *TelegrafReloader) IsReloadRequested() bool {
	return tr.shouldReload
}

func (tr *TelegrafReloader) RequestReload() {
	tr.shouldReload = true
}

func (tr *TelegrafReloader) ReloadIfRequested() {
	if tr.shouldReload {
		tr.Reload()
	}
}
func (tr *TelegrafReloader) Reload() {
	tr.dispatchSignal()
	tr.shouldReload = false
}

func (tr *TelegrafReloader) dispatchSignal() {

	callback := func(path string, info os.FileInfo, err error) error {
		// We just return in case of errors, as they are likely due to insufficient
		// privileges. We shouldn't get any errors for accessing the information we
		// are interested in. Run as root (sudo) and log the error, in case you want
		// this information.
		if err != nil {
			return nil
		}

		// We are only interested in files with a path looking like /proc/<pid>/status.
		if strings.Count(path, "/") == 3 && strings.Contains(path, "/status") {

			// Let's extract the middle part of the path with the <pid> and
			// convert the <pid> into an integer. Log an error if it fails.
			pid, err := strconv.Atoi(path[6:strings.LastIndex(path, "/")])
			if err != nil {
				logrus.Debugf("failed to extract pid from path: %v", path)
				return nil
			}

			// The status file contains the name of the process in its first line.
			// The line looks like "Name: theProcess".
			// Log an error in case we cant read the file.
			f, err := ioutil.ReadFile(path)
			if err != nil {
				log.Println(err)
				return nil
			}

			// Extract the process name from within the first line in the buffer
			name := string(f[6:bytes.IndexByte(f, '\n')])

			if name == tr.name {
				logrus.Debugf("PID: %d, Name: %s will be signaled with %v", pid, name, tr.signal)
				proc, err := os.FindProcess(pid)
				if err != nil {
					logrus.Errorf("> Failed to signal, err: %v", err)
					log.Println(err)
				}

				proc.Signal(tr.signal)

				// Let's return a fake error to abort the walk through the rest of the /proc directory tree
				return io.EOF
			}
		}

		return nil
	}

	err := filepath.Walk("/proc", callback)
	if err != nil {
		if err == io.EOF {
			// Not an error, just a signal when we are done
			err = nil
		} else {
			log.Fatal(err)
		}
	}
}
