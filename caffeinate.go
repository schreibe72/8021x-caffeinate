package main

import (
	"log"
	"os/exec"
	"syscall"
)

type caffeinate struct {
	cmd *exec.Cmd
}

func newCaffeinate() caffeinate {
	return caffeinate{}
}

func (c *caffeinate) start() bool {
	output := false
	if c.cmd == nil {
		log.Println("Start caffeinate")
		output = true
		c.cmd = exec.Command("caffeinate", "-s")
		go func() {
			switch err := c.cmd.Run(); {
			case err.Error() == "signal: terminated":
				log.Println("caffeinate terminated")
			case err == nil:
				log.Println("caffeinate ended")
			default:
				log.Fatalf("Errror starting caffeinate: %s", err)
			}
		}()
	}
	return output
}

func (c *caffeinate) stop() bool {
	if c.cmd == nil {
		return false
	}
	log.Println("Stop caffeinate")
	if err := c.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		log.Fatalf("Errror stopping caffeinate: %s", err)
	}
	if _, err := c.cmd.Process.Wait(); err != nil {
		log.Printf("Wait error: %s", err)
	}
	c.cmd = nil
	log.Println("Stopped caffeinate")
	return true
}
