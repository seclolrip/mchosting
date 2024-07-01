#!/bin/bash

#Remove any existing Golang installations
sudo rm -rf go

# Install Golang 1.22.2
curl -O https://dl.google.com/go/go1.22.2.linux-amd64.tar.gz && \
mkdir go && \
tar -C $HOME/go -xzf go1.22.2.linux-amd64.tar.gz && \
rm -rf $HOME/go1.22.2.linux-amd64.tar.gz && \
echo 'export PATH="$PATH:$HOME/go/go/bin"' >> ~/.profile && \
echo 'export GOPATH="${HOME}/go/workspace"' >> ~/.profile && \
source ~/.profile

#Install Ubuntu Docker
sudo apt-get update && \
sudo apt-get install ca-certificates curl && \
sudo install -m 0755 -d /etc/apt/keyrings && \
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc && \
sudo chmod a+r /etc/apt/keyrings/docker.asc && \
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null && \
sudo apt-get update && \
sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin -y

#Install MongoDB
sudo apt install gnupg curl && \
curl -fsSL https://www.mongodb.org/static/pgp/server-7.0.asc | \
  sudo gpg -o /usr/share/keyrings/mongodb-server-7.0.gpg \
  --dearmor && \
echo "deb [ arch=amd64,arm64 signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg ] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list && \
sudo apt-get update && \
sudo apt-get install -y mongodb-org && \
#Disable default mongod config
sudo systemctl disable mongod
