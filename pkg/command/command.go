package command

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

type Par struct {
	Commands []string
	Kill     chan bool
}

func (p *Par) Run() {
	var commands []*exec.Cmd

	for _, command := range p.Commands {
		commands = append(commands, RunCommand(command))
	}

	<-p.Kill

	for _, cmd := range commands {
		KillCommand(cmd)
	}
}

func RunCommand(command string) *exec.Cmd {
	cmd := exec.Command("sh", "-c", command)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	if err := cmd.Start(); err != nil {
		log.Fatal()
	}

	go func() {
		cmd.Wait()
	}()

	return cmd
}

func KillCommand(cmd *exec.Cmd) {
	if cmd.ProcessState == nil {
		if err := syscall.Kill(-cmd.Process.Pid, syscall.SIGINT); err != nil {
			log.Fatal(err)
		}
	}
}
