#!/bin/bash

# Setup server
sudo hostnamectl set-hostname reaction.pics

# Install go
sudo apt-get update
sudo apt-get install golang-go

# Clone repository
mkdir -p ~/gocode/src/github.com/albertyw/
mkdir -p ~/gocode/bin/
mkdir -p ~/gocode/pkg/
git clone git@github.com:albertyw/reaction-pics ~/gocode/src/github.com/albertyw/reaction-pics
ln -s gocode/src/github.com/albertyw/reaction-pics/ reaction-pics

# Install nginx
sudo add-apt-repository ppa:nginx/stable
sudo apt-get update
sudo apt-get install -y nginx

# Configure nginx
sudo rm -r /etc/nginx/sites-available
sudo rm -r /etc/nginx/sites-enabled
sudo ln -s ~/gocode/src/github.com/albertyw/reaction-pics/config/sites-enabled /etc/nginx/sites-enabled
sudo rm -r /var/www/html

# Secure nginx
sudo mkdir /etc/nginx/ssl
sudo openssl dhparam -out /etc/nginx/ssl/dhparams.pem 2048
# Copy server.key and server.pem to /etc/nginx/ssl
sudo service nginx restart

# Compile
wget https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz
tar xvf glide-v0.12.3-linux-amd64.tar.gz
mv linux-amd64/glide ~/gocode/bin/
rm -r linux-amd64
rm ~/glide-v0.12.3-linux-amd64.tar.gz
cd ~/gocode/src/github.com/albertyw/reaction-pics
glide install
go build