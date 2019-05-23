FROM fedora:30

RUN dnf install -y --nogpgcheck git

ADD awsom /usr/bin/awsom

ENTRYPOINT ["/usr/bin/awsom"]