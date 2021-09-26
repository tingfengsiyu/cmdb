package model

import (
	"strings"
)

var mac = ""

func ScanHardWareInfo() {
	servers, _ := GetServers(0, 0)
	for _, v := range servers {
		go func(ip string) {
			cmd := "/data/ops/script/execshell.sh  " + ip + ` 'free -h |head -2 |tail -1 |awk "{print \$2}" ; lscpu |grep "Model name" |tail -1 |awk -F ":" "{print \$2}"  ;
fdisk -l  |grep Disk  |grep dev |egrep -v "md|mapper"  |awk "{print \$2,\$3,\$4}"|sed "s@/dev/@@g" | tr "\n" " ";echo ;
if [ -f /usr/bin/nvidia-smi  ]; then nvidia-smi -L |awk "{print \$5}" ; else echo nogpu; fi '`
			outs := ExecLocalShell(0, cmd)
			s := strings.Split(outs, "\n")
			db.Model(&Server{}).Where("private_ip_address =? ", ip).Updates(Server{Memory: s[3], Cpu: strings.TrimSpace(s[4]), Disk: s[5], Mac: mac, Gpu: s[6] + " " + s[len(s)-1]})
		}(v.PrivateIpAddress)
	}
}

/*
if  hostnamectl |grep -i centos &> /dev/null  ; then dev=$(grep  $IP /etc/sysconfig/network-scripts/* |awk -F'[:-]' '{print $3}'); else dev=$(cat /etc/netplan/00-installer-config.yaml  |grep 'en.*:$'|sed 's@:@@g') ; fi
ip  addr show $dev |grep link/eth|awk '{print $2}'
*/
