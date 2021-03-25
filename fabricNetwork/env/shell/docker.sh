#!/bin/sh

installType=$1

dockerTar=$2
distPath="/usr/bin"

if [ "${installType}" == "" ]; then
    echo 'default online'
    installType="online"
fi

function installOnLine() {
    sudo yum install -y yum-utils device-mapper-persistent-data lvm2
    sudo yum-config-manager --add-repo http://mirrors.aliyun.com/docker-ce/linux/centos/docker-ce.repo
    sudo yum makecache fast
    sudo yum install -y docker-ce-19.03.9-3.el7
    mkdir -p /etc/docker
    cat >/etc/docker/daemon.json <<EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"],
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "100m"
  }

}
EOF
    # mkdir -p /etc/systemd/system/docker.service.d
    systemctl daemon-reload && systemctl restart docker && systemctl enable docker
    #比如向a.txt文件匹配1234字符串的行的后面添加hahaha
    # sed -i '/1234/a\hahaha' a.txt
    sed -i '/ExeccStart=/a\ -H tcp://0.0.0.0:2375' /usr/lib/systemd/system/docker.service
    systemctl daemon-reload && systemctl restart docker
}

function installOffline() {
    ##TODO check $dockerTar
    tar -zxvf $dockerTar
    cp docker/* $distPath

    cat >/etc/systemd/system/docker.service <<EOF
[Unit]

Description=Docker Application Container Engine

Documentation=https://docs.docker.com

After=network-online.target firewalld.service

Wants=network-online.target

 

[Service]

Type=notify

# the default is not to use systemd for cgroups because the delegate issues still

# exists and systemd currently does not support the cgroup feature set required

# for containers run by docker

ExecStart=/usr/bin/dockerd -H tcp://0.0.0.0:2375 -H unix://var/run/docker.sock

ExecReload=/bin/kill -s HUP $MAINPID

# Having non-zero Limit*s causes performance problems due to accounting overhead

# in the kernel. We recommend using cgroups to do container-local accounting.

LimitNOFILE=infinity

LimitNPROC=infinity

LimitCORE=infinity

# Uncomment TasksMax if your systemd version supports it.

# Only systemd 226 and above support this version.

#TasksMax=infinity

TimeoutStartSec=0

# set delegate yes so that systemd does not reset the cgroups of docker containers

Delegate=yes

# kill only the docker process, not all processes in the cgroup

KillMode=process

# restart the docker process if it exits prematurely

Restart=on-failure

StartLimitBurst=3

StartLimitInterval=60s

 

[Install]

WantedBy=multi-user.target
EOF
    chmod +x /etc/systemd/system/docker.service
    systemctl daemon-reload && systemctl restart docker && systemctl enable docker
}
if [ "${installType}" == "online" ]; then
    echo 'exec installOnLine '
    installOnLine
elif [ "${installType}" == "offline" ]; then
    echo 'exec installOffline'
    if [ "$dockerTar" == "" ]; then
        exit 1
    fi
    installOffline
fi
