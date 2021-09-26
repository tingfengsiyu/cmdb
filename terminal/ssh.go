package terminal

import (
	"bytes"
	"cmdb/model"
	"fmt"
	"golang.org/x/crypto/ssh"
	"time"
)

func NewSshClient(h model.ScanTerm) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 5,
		User:            h.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
		//HostKeyCallback: hostKeyCallBackFunc(h.Host),
	}
	pass, _ := model.ScryptPassw(h.Password)
	config.Auth = []ssh.AuthMethod{ssh.Password(pass)}
	addr := fmt.Sprintf("%s:%d", h.PrivateIpAddress, h.Port)
	c, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func runCommand(client *ssh.Client, command string) (stdout string, err error) {
	session, err := client.NewSession()
	if err != nil {
		//log.Print(err)
		return
	}
	defer session.Close()

	var buf bytes.Buffer
	session.Stdout = &buf
	err = session.Run(command)
	if err != nil {
		//log.Print(err)
		return
	}
	stdout = string(buf.Bytes())

	return
}
