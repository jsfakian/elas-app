package docx

import (
	"bytes"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// PowerShell struct
type PowerShell struct {
	powerShell string
}

// New create new session
func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

func (p *PowerShell) execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	stdOut, stdErr = stdout.String(), stderr.String()
	return
}

func OpenDocx(filename string) {
	posh := New()

	log.Info("Filename: ", filename)
	cmd := `$Filename=` + strings.ReplaceAll(filename, "/", "\\") + `
	$Word=NEW-Object –comobject Word.Application
	$pids = Get-Process *Word* | Select-Object Handles, Id | sort-Object -Property Handles -Descending
	$Document=$Word.Documents.openNoRepairDialog($Filename)
	kill $pids[1].Id
	Exit;
	`

	stdOut, stdErr, err := posh.execute(cmd)
	log.Printf("ElevateProcessCmds:\nStdOut : '%s'\nStdErr: '%s'\nErr: %s", strings.TrimSpace(stdOut), stdErr, err)
}
