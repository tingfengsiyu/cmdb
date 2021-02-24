package model

import (
	"bufio"
	"cmdb/middleware"
	"cmdb/utils"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"strings"
	"time"
)

type Connection struct {
	*ssh.Client
	password string
	sudopass string
}

func Connect(addr, user, password, sudopass string) (*Connection, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}

	return &Connection{conn, password, sudopass}, nil

}

func (conn *Connection) SendCommands(cmds ...string) ([]byte, error) {
	session, err := conn.NewSession()
	if err != nil {
		middleware.SugarLogger.Errorf("%s", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return []byte{}, err
	}

	in, err := session.StdinPipe()
	if err != nil {
		middleware.SugarLogger.Errorf("%s", err)
	}

	out, err := session.StdoutPipe()
	if err != nil {
		middleware.SugarLogger.Errorf("%s", err)
	}

	var output []byte

	go func(in io.WriteCloser, out io.Reader, output *[]byte) {
		var (
			line string
			r    = bufio.NewReader(out)
		)
		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}

			*output = append(*output, b)

			if b == byte('\n') {
				line = ""
				continue
			}

			line += string(b)

			if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
				_, err = in.Write([]byte(conn.password + "\n"))
				if err != nil {
					break
				}
			}
		}
	}(in, out, &output)

	cmd := strings.Join(cmds, "; ")
	_, err = session.Output(cmd)
	if err != nil {
		return []byte{}, err
	}

	return output, nil
}

func BatchSsh(msg chan []ScanMonitorPrometheus) {
	defer wg.Done()
	tmp, _ := <-msg
	for _, v := range tmp {
		var user, passwd, sudopasswd string
		if v.Label == "lotus-miner" || v.Label == "lotus-worker" {
			user = utils.WorkerUser
			passwd = utils.WorkerPass
			sudopasswd = utils.WorkerPass
		} else if v.Label == "lotus-storage" {
			user = utils.StorageUser
			passwd = utils.StoragePass
			sudopasswd = utils.StorageSudoPass
		} else {
			user = utils.WorkerUser
			passwd = utils.WorkerPass
			sudopasswd = utils.WorkerPass
		}
		conn, err := Connect("172.22.0.20:22", user, passwd, sudopasswd)
		if err != nil {
			middleware.SugarLogger.Errorf("%s", err)
		}
		sshdConfig := "sudo sed -i 's@PasswordAuthentication no@PasswordAuthentication yes@g' /etc/ssh/sshd_config"
		updatePass := fmt.Sprintf("sudo echo root:%s | chpasswd", utils.RootPass)
		updatePubKey := fmt.Sprintf("sudo grep ops /root/.ssh/authorized_keys || sudo sed -i '1i %s' /root/.ssh/authorized_keys ", utils.RootPub)
		_, err = conn.SendCommands(sshdConfig, updatePass, updatePubKey)
		if err != nil {
			middleware.SugarLogger.Errorf("ssh commands  %s ", err)
		}
		fmt.Println(time.Now())
	}
}
func BatchSsh2(v ScanMonitorPrometheus) {
	fmt.Println(time.Now())
	defer wg.Done()
	var user, passwd, sudopasswd string
	if v.Label == "lotus-miner" || v.Label == "lotus-worker" {
		user = utils.WorkerUser
		passwd = utils.WorkerPass
		sudopasswd = utils.WorkerPass
	} else if v.Label == "lotus-storage" {
		user = utils.StorageUser
		passwd = utils.StoragePass
		sudopasswd = utils.StorageSudoPass
	} else {
		user = utils.WorkerUser
		passwd = utils.WorkerPass
		sudopasswd = utils.WorkerPass
	}
	conn, err := Connect("172.22.0.20:22", user, passwd, sudopasswd)
	if err != nil {
		middleware.SugarLogger.Errorf("%s", err)
	}
	sshdConfig := "sudo sed -i 's@PasswordAuthentication no@PasswordAuthentication yes@g' /etc/ssh/sshd_config"
	updatePass := fmt.Sprintf("sudo echo root:%s | chpasswd", utils.RootPass)
	updatePubKey := fmt.Sprintf("sudo grep ops /root/.ssh/authorized_keys || sudo sed -i '1i %s' /root/.ssh/authorized_keys ", utils.RootPub)
	_, err = conn.SendCommands(sshdConfig, updatePass, updatePubKey)
	if err != nil {
		middleware.SugarLogger.Errorf("ssh commands  %s ", err)
	}
	fmt.Println(time.Now())

}

func OsInit() {
	monitorPrometheus := Prometheus_server()
	for _, v := range monitorPrometheus {
		go BatchSsh2(v)
		wg.Add(1)
		wg.Wait()
	}
}
