FROM ubuntu:22.04

ARG JAVA_VERSION
ARG MCFILES
ARG SSHKEY

SHELL ["/bin/bash", "-c"]

RUN apt purge openjdk-* || echo "No Java Packages Installed" && \
    apt update && apt upgrade -y && \
    apt install -y openjdk-${JAVA_VERSION}-jdk openssh-server sudo vim && \
    apt clean

RUN useradd -m -s /bin/bash mchosted -G sudo && \
    mkdir /home/mchosted/.ssh && \
    mkdir /home/mchosted/mcserver && \
    echo "$SSHKEY" >> /home/mchosted/.ssh/authorized_keys && \
    chown -R mchosted:mchosted /home/mchosted/.ssh && \
    chmod 700 /home/mchosted/.ssh && \
    chmod 600 /home/mchosted/.ssh/authorized_keys

COPY ${MCFILES} /home/mchosted/mcserver

ENV JAVA_HOME=/usr/lib/jvm/java-${JAVA_VERSION}-openjdk-amd64

RUN java -version

RUN sed -i 's/^#PermitRootLogin.*/PermitRootLogin no/' /etc/ssh/sshd_config && \
    sed -i 's/^#PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config && \
    sed -i 's/^UsePAM.*/UsePAM yes/' /etc/ssh/sshd_config && \
    sed -i 's/^PubkeyAuthentication.*/PubkeyAuthentication yes/' /etc/ssh/sshd_config && \
    mkdir /var/run/sshd

RUN history -c && \
    rm -rf /root/.bash_history

EXPOSE 22

CMD ["/usr/sbin/sshd", "-D"]