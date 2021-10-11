package model

import (
	"bufio"
	"cmdb/middleware"
	"cmdb/utils"
	"context"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

//var sudostr = " ansible_ssh_user=" + utils.WorkerUser + " ansible_ssh_pass=" + utils.WorkerPass + " ansible_sudo_pass=" + utils.WorkerSudoPass
var sudostr = ""

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
				middleware.SugarLogger.Errorf("ssh exec commands error !!!  %s ", err)
			}
			middleware.SugarLogger.Errorf(string(outs))
		}(v)
	}
}

func ExecLocalShell(id int, command string) string {
	timeout := 2
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout+1)*time.Hour)
	defer cancel()
	cmdarray := []string{"-c", command}
	cmd := exec.CommandContext(ctx, "bash", cmdarray...)
	out, err := cmd.CombinedOutput()
	var cmd_err string
	status := 1
	if err != nil {
		status = 0
		cmd_err = err.Error()
	}
	if ctx.Err() != nil {
		status = 0
		cmd_err = cmd_err + ctx.Err().Error()
	}
	//fmt.Printf("ctx.Err : [%v]\n", ctx.Err())
	//fmt.Printf("error   : [%v]\n", err)
	//fmt.Printf("out     : [%s]\n", string(out))
	success := string(out)
	UpdateRecords(id, status, success, cmd_err)
	return success + cmd_err
}

func GenerateAnsibleHosts() error {
	file, err := os.OpenFile(utils.AnsibleHosts, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		middleware.SugarLogger.Errorf("写入文件错误!!!%s", err)
	}
	defer file.Close()
	maps := AllHosts()
	for k, v := range maps {
		file.WriteString("[" + k + "]\n")
		for _, ip := range v {
			file.WriteString(ip + sudostr + "\n")
		}
		file.WriteString("[" + k + ":vars]\n")
		file.WriteString("ansible_ssh_user=" + utils.WorkerUser + "\n")
		file.WriteString("ansible_ssh_pass=" + utils.WorkerPass + "\n")
		file.WriteString("ansible_sudo_pass=" + utils.WorkerSudoPass + "\n")
	}
	return err
}

func AppendAnsibleHost(ips []string) error {
	file, err := os.OpenFile(utils.AnsibleHosts+"-tmplotus", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		middleware.SugarLogger.Errorf("追加文件错误!!!%s", err)
	}
	defer file.Close()
	file.WriteString("[" + "tmplotus" + "]\n")
	for _, v := range ips {
		file.WriteString(v + "\n")
	}
	file.WriteString("[" + "tmplotus" + ":vars]\n")
	file.WriteString("ansible_ssh_user=" + utils.WorkerUser + "\n")
	file.WriteString("ansible_ssh_pass=" + utils.WorkerPass + "\n")
	file.WriteString("ansible_sudo_pass=" + utils.WorkerSudoPass + "\n")
	return err
}

func GenerateClustersHosts() {
	maps := AllHosts()
	for k, v := range maps {
		//reg := regexp.MustCompile("-(lotus|chia|tmp)-.*$");
		//a := reg.ReplaceAllString(k, "")
		tmpfile := utils.AnsibleHosts + "-" + k
		//s, err := os.OpenFile(tmpfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		//defer  s.Close()
		t, err := os.OpenFile(tmpfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
		defer t.Close()
		if err != nil {
			middleware.SugarLogger.Errorf("写入文件错误!!!%s", err)
		}
		t.WriteString("[" + k + "]\n")
		for _, s := range v {
			t.WriteString(s + sudostr + "\n")
		}
		t.WriteString("[" + k + ":vars]\n")
		t.WriteString("ansible_ssh_user=" + utils.WorkerUser + "\n")
		t.WriteString("ansible_ssh_pass=" + utils.WorkerPass + "\n")
		t.WriteString("ansible_sudo_pass=" + utils.WorkerSudoPass + "\n")
	}
}

//func Hosts() map[string][]string {
//	var maps = make(map[string][]string, 0)
//	str := []string{}
//	ss := ""
//	clusters ,_:=GetClusters()
//	for _,v := range  clusters {
//		servers, _ := GetCluster(v.Cluster)
//		sort.Slice(servers, func(i, j int) bool { return servers[i].Cluster < servers[j].Cluster })
//		for _,s := range servers{
//			ss= v.Cluster
//			str = append(str,s.PrivateIpAddress+ " roles="+s.Label +" hostname="+s.Name)
//		}
//		maps[ss] = str
//		str = []string{}
//	}
//	return maps
//}

func AllHosts() map[string][]string {
	ansiblehost, _ := GetServers(0, 0)
	if err != nil {
		middleware.SugarLogger.Errorf("sql查询错误%s", err)
	}
	type server struct {
		PrivateIpAddress string `json:"private_ip_address"`
		Role             string `json:"role"` //cluster+Label
		Label            string `json:"label"`
		Name             string `json:"name"`
	}
	servers := []server{}
	for _, v := range ansiblehost {
		servers = append(servers, server{v.PrivateIpAddress, v.Cluster + "-" + v.Label, v.Label, v.Name})
	}
	sort.Slice(servers, func(i, j int) bool { return servers[i].Role < servers[j].Role })
	var maps = make(map[string][]string, 0)
	var worker, miner, storage, none []string
	for _, v := range servers {
		labels := v.Role
		str := "" + " roles=" + v.Label + " hostname=" + v.Name
		if _, ok := maps[labels]; !ok {
			miner = []string{}
			worker = []string{}
			storage = []string{}
			none = []string{}
		}
		switch labels {
		case "lotus-worker":
			worker = append(worker, v.PrivateIpAddress+str)
			sort.Strings(worker)
			maps[labels] = worker
		case "lotus-storage":
			storage = append(storage, v.PrivateIpAddress+str)
			sort.Strings(storage)
			maps[labels] = storage
		case "lotus-miner":
			miner = append(miner, v.PrivateIpAddress+str)
			sort.Strings(miner)
			maps[labels] = miner
		default:
			none = append(none, v.PrivateIpAddress+str)
			sort.Strings(none)
			maps[labels] = none
		}
	}
	return maps
}
