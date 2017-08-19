#!/usr/bin/env bash

echo 'installing dependencies'
yum update

wget -c --header "Cookie: oraclelicense=accept-securebackup-cookie" http://download.oracle.com/otn-pub/java/jdk/8u131-b11/d54c1d3a095b4ff2b6607d096fa80163/jdk-8u131-linux-x64.rpm

yum localinstall -y jdk-8u131-linux-x64.rpm
rm /home/ec2-user/jdk-8u131-linux-x64.rpm

rpm --import https://artifacts.elastic.co/GPG-KEY-elasticsearch
yum install -y logstash-1.4.2-1_2c0f5a1.noarch.rpm

cat <<EOT >> /etc/yum.repos.d/logstash.repo
[logstash-5.x]
name=Elastic repository for 5.x packages
baseurl=https://artifacts.elastic.co/packages/5.x/yum
gpgcheck=1
gpgkey=https://artifacts.elastic.co/GPG-KEY-elasticsearch
enabled=1
autorefresh=1
type=rpm-md
EOT

echo 'instaling logstash'
yum install -y logstash