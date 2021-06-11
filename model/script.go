package model

import (
	"bufio"
	"cmdb/middleware"
	"cmdb/utils"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

var sudostr = " ansible_ssh_user=" + utils.WorkerUser + " ansible_ssh_pass=" + utils.WorkerPass + " ansible_sudo_pass=" + utils.WorkerSudoPass

type BatchIpStruct struct {
	SourceStartIp     string `json:"source_start_ip" validate:"required,min=10,max=12" `
	SourceGateway     string `json:"source_gateway" validate:"required,min=10,max=10" `
	SourceEndNumber   string `json:"source_end_number" validate:"required,gte=2"  `
	TargetStartIP     string `json:"target_start_ip" validate:"required,min=10,max=12" `
	TargetGateway     string `json:"target_gateway" validate:"required,min=10,max=10" `
	TargetClusterName string `json:"target_cluster_name" validate:"required,min=4,max=50"`
}

type UpdateClusterStruct struct {
	SourceStartIp     string `json:"source_start_ip" validate:"required,min=10,max=12" `
	SourceEndNumber   string `json:"source_end_number" validate:"required,gte=2"  `
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
	Operating         string `json:"operating" validate:"required,min=1,max=3"`
}

type ansibleStruct struct {
	PrivateIpAddress string `json:"private_ip_address"`
	Label            string `json:"label"`
	Cluster          string `json:"cluster"`
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
			user = utils.WorkerUser
			passwd = utils.WorkerPass
			sudopasswd = utils.WorkerSudoPass
			outs, err := SshCommands(user, passwd, v.PrivateIpAddress+":"+"22", sudopasswd, "hostnamectl set-hostname "+v.Name)
			if err != nil {
				fmt.Println("ssh exec commands error !!!  %s ", err)
			}
			fmt.Println(string(outs))
		}(v)
	}
}

func ExecLocalShell(command string) {
	cmd := exec.Command("/bin/bash", "-c", utils.ScriptDir+command)
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Exec shell error !!!", time.Now(), err.Error())
	}
	fmt.Println(string(output))
	fmt.Println("Exec shell success !!!", time.Now())
	/*
		cmd := exec.Command("/bin/bash", "-c", command)

		stdin, _ := cmd.StdinPipe()
		stdout, _ := cmd.StdoutPipe()

		if err := cmd.Start(); err != nil {
			fmt.Println("Execute failed when Star`21wt:" + err.Error())
			return
		}

		stdin.Close()

		out_bytes, _ := ioutil.ReadAll(stdout)
		stdout.Close()

		fmt.Println(string(out_bytes))
	*/

}

func GenerateAnsibleHosts() error {

	var ansiblehost = []ansibleStruct{}

	err = db.Model(&Server{}).Select("private_ip_address,label,cluster").Scan(&ansiblehost).Error
	if err != nil {
		middleware.SugarLogger.Errorf("sql查询错误%s", err)
	}
	type server struct {
		PrivateIpAddress string `json:"private_ip_address"`
		Role             string `json:"Role"` //cluster+Label
	}
	servers := []server{}
	for _, v := range ansiblehost {
		servers = append(servers, server{v.PrivateIpAddress, v.Cluster + "-" + v.Label})
	}
	sort.Slice(servers, func(i, j int) bool { return servers[i].Role < servers[j].Role })
	var worker, miner, storage, none []string
	var maps = make(map[string][]string, 0)
	for _, v := range servers {
		if _, ok := maps[v.Role]; !ok {
			miner = []string{}
			worker = []string{}
			storage = []string{}
			none = []string{}
		}

		switch v.Role {
		case "lotus-worker":
			worker = append(worker, v.PrivateIpAddress)
			sort.Strings(worker)
			maps[v.Role] = worker
		case "lotus-storage":
			storage = append(storage, v.PrivateIpAddress)
			sort.Strings(storage)
			maps[v.Role] = storage
		case "lotus-miner":
			miner = append(miner, v.PrivateIpAddress)
			sort.Strings(miner)
			maps[v.Role] = miner
		default:
			none = append(none, v.PrivateIpAddress)
			sort.Strings(none)
			maps[v.Role] = none
		}

	}
	file, err := os.OpenFile(utils.AnsibleHosts, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		middleware.SugarLogger.Errorf("写入文件错误!!!%s", err)
		return nil
	}
	defer file.Close()
	for k, v := range maps {
		file.WriteString("[" + k + "]\n")
		for _, ip := range v {
			file.WriteString(ip + sudostr + "\n")
		}
	}
	return err
}

func AppendAnsibleHosts(ips []string, cluster string) error {
	file, err := os.OpenFile(utils.AnsibleHosts, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		middleware.SugarLogger.Errorf("写入文件错误!!!%s", err)
		return nil
	}
	defer file.Close()
	file.WriteString("[" + cluster + "-tmpworker]\n")
	for _, ip := range ips {
		file.WriteString(ip + sudostr + "\n")
	}
	//sync  target cluster ansible hosts
	SyncTargetHosts(ips, cluster)
	return err
}

func SyncTargetHosts(ips []string, cluster string) error {
	tmpfile := utils.AnsibleHosts + ".tmp"
	file, err := os.OpenFile(tmpfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		middleware.SugarLogger.Errorf("写入文件错误!!!%s", err)
		return nil
	}
	defer file.Close()

	//sync  target cluster ansible hosts
	servers, _ := GetCluster(cluster)
	sort.Slice(servers, func(i, j int) bool { return servers[i].Label < servers[j].Label })
	var worker, miner, storage, none []string
	var maps = make(map[string][]string, 0)
	for _, v := range servers {
		if _, ok := maps[v.Label]; !ok {
			miner = []string{}
			worker = []string{}
			storage = []string{}
			none = []string{}
		}
		switch v.Label {
		case "lotus-worker":
			worker = append(worker, v.PrivateIpAddress)
			sort.Strings(worker)
			maps[v.Label] = worker
		case "lotus-storage":
			storage = append(storage, v.PrivateIpAddress)
			sort.Strings(storage)
			maps[v.Label] = storage
		case "lotus-miner":
			miner = append(miner, v.PrivateIpAddress)
			sort.Strings(miner)
			maps[v.Label] = miner
		default:
			none = append(none, v.PrivateIpAddress)
			sort.Strings(none)
			maps[v.Label] = none
		}
	}
	for k, v := range maps {
		file.WriteString("[" + k + "]\n")
		for _, ip := range v {
			file.WriteString(ip + sudostr + "\n")
		}
	}
	//
	file.WriteString("[" + "addworker]\n")
	for _, ip := range ips {
		file.WriteString(ip + sudostr + "\n")
	}
	//
	cmd := "scp.sh " + cluster + "-*miner " + tmpfile
	ExecLocalShell(cmd)
	cmd = "execshell.sh " + cluster + "-*miner " + " mv  " + tmpfile + "  " + utils.AnsibleHosts
	ExecLocalShell(cmd)
	return err
}
