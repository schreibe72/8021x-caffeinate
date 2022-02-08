package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

type process struct {
	name string
	args []string
	cmd  *exec.Cmd
	wg   sync.WaitGroup
}

func newProcess(name string, args ...string) *process {
	return &process{name: name, args: args}
}

func (p *process) start() bool {
	output := false
	if p.cmd == nil {
		log.Printf("Start %s", p.name)
		output = true
		p.cmd = exec.Command(p.name, p.args...)
		stdout, err := p.cmd.StdoutPipe()
		if err != nil {
			log.Fatalf("Errror connection stdoutPipe %s: %s", p.name, err)
		}
		scanner := bufio.NewScanner(stdout)
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			for scanner.Scan() {
				log.Printf("Output [%s]: %s", p.name, scanner.Text())
			}
		}()
		if err := p.cmd.Start(); err != nil {
			log.Fatalf("Errror starting %s: %s", p.name, err)
		}
	}
	return output
}

func (p *process) stop() bool {
	if p.cmd == nil {
		return false
	}
	log.Printf("Stop %s", p.name)
	if err := p.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		log.Fatalf("Errror stopping %s: %s", p.name, err)
	}
	if err := p.cmd.Wait(); err != nil && err.Error() != "signal: terminated" {
		log.Printf("Wait error: %s", err)
	}
	p.wg.Wait()
	p.cmd = nil
	log.Printf("Stopped %s", p.name)
	return true
}

func (p *process) started() bool {
	if p.cmd != nil && p.cmd.Process != nil {
		_, err := os.FindProcess(p.cmd.Process.Pid)
		if err == nil {
			return true
		}
	}
	return false
}

func (p *process) findNameInProcesslist() bool {
	output := false
	cmd := exec.Command("ps", "-axc")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	// optionally, resize scanner's capacity for lines over 64K, see next example
	for scanner.Scan() {
		a := strings.Join(strings.Fields(scanner.Text())[3:], " ")
		if strings.Contains(a, p.name) {
			output = true
		}
	}
	cmd.Wait()
	return output
}
