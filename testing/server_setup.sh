#!/bin/bash

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