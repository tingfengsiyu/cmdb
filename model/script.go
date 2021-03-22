package model

import (
	"bufio"
	"cmdb/middleware"
	"cmdb/utils"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os/exec"
	"strings"
)

type BatchIpStruct struct {
	SourceStartIp     string `json:"source_start_ip" validate:"required,min=10,max=12" `
	SourceGateway     string `json:"source_gateway" validate:"required,min=10,max=10" `
	SourceEndNumber   string `json:"source_end_number" validate:"required,gte=2"  `
	TargetStartIP     string `json:"target_start_ip" validate:"required,min=10,max=12" `
	TargetGateway     string `json:"target_gateway" validate:"required,min=10,max=10" `
	TargetClusterName string `json:"target_cluster_name" validate:"required,min=4,max=50"`
}

type OsInitStruct struct {
	InitUser     string `json:"init_user" validate:"required,min=10,max=10" `
	InitPass     string `json:"init_pass" validate:"required,min=4,max=50"`
	Role         string `json:"role" validate:"required,min=4,max=10"`
	StorageMount StorageMountStruct
}

type StorageMountStruct struct {
	InitStartIP       string `json:"init_start_ip" validate:"required,min=10,max=12" `
	InitEndNumber     string `json:"init_end_number" validate:"required,gte=2" `
	StorageStartIP    string `json:"storage_start_ip" validate:"required,min=4,max=50"`
	StorageStopnumber string `json:"storage_stop_number" validate:"required,min=1,max=3"`
}

type Connection struct {
	*ssh.Client
	password string
	sudopass string
}

func SshCommands(user, password, addr, sudopass string, cmds ...string) ([]byte, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		middleware.SugarLogger.Errorf("sshconnection %s", err)
		return nil, err
	}
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
				_, err = in.Write([]byte(sudopass + "\n"))
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

func UpdateHostName() {
	servers, _ := GetServers(0, 0)
	for _, v := range servers {
		go func(v Server) {
			var user, passwd, sudopasswd string
			user = "root"
			passwd = utils.RootPass
			sudopasswd = utils.RootPass
			hostname := fmt.Sprintf("hostnamectl set-hostname %s ", v.Name)
			outs, err := SshCommands(user, passwd, v.PrivateIpAddress+":"+"22", sudopasswd, hostname)
			if err != nil {
				middleware.SugarLogger.Errorf("ssh commands  %s ", err)
			}
			middleware.SugarLogger.Infof("%s ", string(outs))
		}(v)
	}
}

func ExecLocalShell(command string) {
	cmd := exec.Command("/bin/bash", "-c", command)

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()

	if err := cmd.Start(); err != nil {
		fmt.Println("Execute failed when Start:" + err.Error())
		return
	}

	stdin.Close()

	out_bytes, _ := ioutil.ReadAll(stdout)
	stdout.Close()

	fmt.Println(string(out_bytes))
}

func GenerateAnsibleHosts() {

}
