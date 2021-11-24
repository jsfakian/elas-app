package violation

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
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

func openDocx(filename string) {
	posh := New()

	cmd := `$Filename=` + filename + `
	$Word=NEW-Object –comobject Word.Application
	$Document=$Word.documents.open($Filename)
	`

	stdOut, stdErr, err := posh.execute(cmd)
	log.Printf("ElevateProcessCmds:\nStdOut : '%s'\nStdErr: '%s'\nErr: %s", strings.TrimSpace(stdOut), stdErr, err)
}
