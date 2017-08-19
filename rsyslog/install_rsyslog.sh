#!/usr/bin/env bash

echo 'installing dependencies'
yum update

cat <<EOT >> /etc/yum.repos.d/rsyslog.repo
[v8-stable]
name=Adiscon CentOS-6 - local packages for \$basearch
baseurl=http://rpms.adiscon.com/v8-stable/epel-6/\$basearch
enabled=0
gpgcheck=0
gpgkey=http://rpms.adiscon.com/RPM-GPG-KEY-Adiscon
protect=1
EOT

sudo yum install -y json-c
sudo yum install -y rsyslog --enablerepo=v8-stable --setopt=v8-stable.priority=1
sudo yum install -y rsyslog-mmjsonparse --enablerepo=v8-stable --setopt=v8-stable.priority=1

