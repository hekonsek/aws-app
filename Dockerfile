FROM fedora:29

RUN dnf install -y git

ADD awsom /usr/bin/awsom