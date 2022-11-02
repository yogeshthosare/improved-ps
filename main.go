package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// Process - represents a running process detected on the host machine
type Process struct {
	Uid        int
	Pid        int
	TheCmdline string
	TheFields  []string
}

func consolidateExecErrInfo(err error, stdout *bytes.Buffer, stderr *bytes.Buffer) error {
	return errors.New("Error executing exec: " + err.Error() + " : " + stdout.String() + " : " + stderr.String())
}

func RunCmd(cmd *exec.Cmd) (stdout *bytes.Buffer, stderr *bytes.Buffer, err error) {

	cmdOutput := &bytes.Buffer{}
	cmdErrOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	cmd.Stderr = cmdErrOutput

	err = cmd.Run()
	if err != nil { // get all the information into the error string
		err = consolidateExecErrInfo(err, cmdOutput, cmdErrOutput)
	}

	return cmdOutput, cmdErrOutput, err
}

func main() {
	var processes []*Process
	processes = make([]*Process, 0)
	// With the ps command, we'll want just the pid and the command line
	// Note that "cmd" does not work on Mac, but "command" works on both Mac & Linux
	// -A select all processes
	// -o format specfier in results, in our case  pid,command
	exe_string := "/sbin/metadataagent"
	cmd := exec.Command("/bin/ps", "-Ao", "uid,pid,command")

	cmdOutput, cmdErrOutput, err := RunCmd(cmd)

	if err != nil {
		err = errors.New("Error executing ps : " + err.Error() + " : " + cmdErrOutput.String() + " : " + cmdErrOutput.String())
	}

	output := cmdOutput.Bytes()

	if len(output) > 0 {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, exe_string) { // is proc a metadataagent?
				filteredLine := strings.Join(strings.Fields(line), " ")
				lineTrimmed := strings.SplitN(filteredLine, " ", 3)
				uid, err := strconv.Atoi(lineTrimmed[0])
				if err != nil {
					fmt.Println("Error in ps", err)
					os.Exit(0)
				}
				pid, err := strconv.Atoi(lineTrimmed[1])
				if err != nil {
					fmt.Println("Error in ps", err)
					os.Exit(0)
				}
				cmdline := lineTrimmed[2]
				fields := make([]string, 0)
				for i, field := range strings.Split(cmdline, " ") {
					if i > 1 {
						fields = append(fields, field)
					}
				}

				processes = append(processes, &Process{
					Uid:        uid,
					Pid:        pid,
					TheCmdline: cmdline,
					TheFields:  fields,
				})
			}
		}
	}

	for i, proc := range processes {
		fmt.Printf("Uid %d, %+v\n", i, proc.Uid)
		fmt.Printf("Pid %d, %+v\n", i, proc.Pid)
		fmt.Printf("TheCmdLine %d, %+v\n", i, proc.TheCmdline)
		fmt.Printf("TheFields %d, %+v\n", i, proc.TheFields)
		if i == 2 {
			break
		}
	}
}
