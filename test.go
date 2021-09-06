package main

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

type SSHConnection struct {
	session *ssh.Session
	stdin   io.Writer
	stdout  io.Reader
	options *SSHOptions
}

type Windows struct {
	Width  int
	Height int
}

type ServerConnection interface {
	io.ReadWriteCloser
	SetWinSize(width, height int) error
	KeepAlive() error
}
type SSHOption func(*SSHOptions)

type SSHOptions struct {
	charset string
	win     Windows
	term    string
}
type SSHClientOptions struct {
	Host         string
	Port         string
	Username     string
	Password     string
	PrivateKey   string
	Passphrase   string
	Timeout      int
	keyboardAuth ssh.KeyboardInteractiveChallenge
	PrivateAuth  ssh.Signer

	proxySSHClientOptions []SSHClientOptions
}

func test() {
	client, err := ssh.Dial("tcp", "127.0.0.1:22", &ssh.ClientConfig{
		User:            "root",
		Auth:            []ssh.AuthMethod{ssh.Password("bs1QB7@?JC/qkBF")},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		log.Fatalf("SSH dial error: %s", err.Error())
	}
	// 建立新会话
	session, err := client.NewSession()
	defer session.Close()
	if err != nil {
		log.Fatalf("new session error: %s", err.Error())
	}
	session.Stdout = os.Stdout // 会话输出关联到系统标准输出设备
	session.Stderr = os.Stderr // 会话错误输出关联到系统标准错误输出设备
	session.Stdin = os.Stdin   // 会话输入关联到系统标准输入设备
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // 禁用回显（0禁用，1启动）
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, //output speed = 14.4kbaud
	}
	if err = session.RequestPty("linux", 32, 160, modes); err != nil {
		log.Fatalf("request pty error: %s", err.Error())
	}
	stdins, _ := session.StdinPipe()
	stdouts, _ := session.StdoutPipe()
	fmt.Println(stdins, stdouts)

	if err = session.Shell(); err != nil {
		log.Fatalf("start shell error: %s", err.Error())
	}
	if err = session.Wait(); err != nil {
		log.Fatalf("return error: %s", err.Error())
	}
}

type SSHClient struct {
	*ssh.Client
	Cfg         *SSHClientOptions
	ProxyClient *SSHClient

	sync.Mutex

	traceSessionMap map[*ssh.Session]time.Time

	refCount int32
}

func NewSSHConnection(sess *ssh.Session, opts ...SSHOption) (*SSHConnection, error) {
	//client, _ := ssh.Dial("tcp", "127.0.0.1:22", &ssh.ClientConfig{
	//	User:            "root",
	//	Auth:            []ssh.AuthMethod{ssh.Password("bs1QB7@?JC/qkBF")},
	//	HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	//})

	options := &SSHOptions{
		charset: "utf8",
		win: Windows{
			Width:  80,
			Height: 120,
		},
		term: "xterm",
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4 kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4 kbaud
	}
	err := sess.RequestPty(options.term, options.win.Height, options.win.Width, modes)
	if err != nil {
		return nil, err
	}
	stdin, err := sess.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := sess.StdoutPipe()
	if err != nil {
		return nil, err
	}

	err = sess.Shell()
	if err != nil {
		return nil, err
	}
	return &SSHConnection{
		session: sess,
		stdin:   stdin,
		stdout:  stdout,
		options: options,
	}, nil
}

func main() {
	test()

}
