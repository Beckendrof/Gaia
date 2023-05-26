/*
============================================================================================
This file contains utilities for the apostolis service. Please read them carefully.
Thank You --- @abhinowP
============================================================================================
*/

package utils

import (
	"io/ioutil"
	"math"
	"os/exec"
	"strings"
	"syscall"
)

func Split(r rune) bool {
	return r == '(' || r == ')' || r == ':'
}

// ApostolisHandleError is a function to handle errors in a better way.
// Returns true if any error and logs error.
func ApostolisHandleError(err error) bool {
	return err != nil
}

// ReturnExitCode is a function that returns exitcode with an exit message and status.
func ReturnExitCode(cmd *exec.Cmd) (out string, status string, exitcode int32) {
	/*
		Overview of this function:
		============================================================================================
		Summary:
		--------------------------------------------------------------------------------------------
		This function is used to return an exitcode with an exit message.


		args:
		--------------------------------------------------------------------------------------------
		cmd: the exec.Cmd variable passed with a custom cli command

		returns:
		--------------------------------------------------------------------------------------------
		out: error message
		message: successful or unsuccessful
		exitcode: int32 value of exitcode
		============================================================================================
	*/

	stdoutStderr, err := cmd.CombinedOutput() // Executes the command and returns output
	ApostolisHandleError(err)                 // logs any error during command execution

	exitcode = int32(cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()) // generates exitcode

	if exitcode == 0 {
		status = Successful.String() // status = "Successful"
	} else {
		status = Unsuccessful.String() // status = "Unsuccessful"
	}

	return string(stdoutStderr), status, exitcode // returns exitcode with output string and status
}

func Round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(Round(num*output)) / output
}

func GetNvidiaModel() string {
	data, err := ioutil.ReadFile("/proc/device-tree/compatible")
	if err != nil {
		GaiaLogger.Error("Failed to read Jetson model: %v\n", err.Error())
	}

	model := string(data)
	if strings.Contains(model, "xavier") {
		return "jetson-xavier"
	} else if strings.Contains(model, "orin") {
		return "jetson-orin"
	} else {
		return ""
	}
}
