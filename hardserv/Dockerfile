FROM ubuntu:22.04

ARG JAVA_VERSION
ARG USERUSERNAME
ARG SSHKEY

SHELL ["/bin/bash", "-c"]

RUN apt purge openjdk-* || echo "No Java Packages Installed" && \
    apt update && apt upgrade -y && \
    apt install -y openjdk-${JAVA_VERSION}-jdk openssh-server sudo vim && \
    apt clean

RUN useradd -m -s /bin/bash $USERUSERNAME -G sudo && \
    mkdir /home/$USERUSERNAME/.ssh && \
    echo "$SSHKEY" >> /home/$USERUSERNAME/.ssh/authorized_keys && \
    chown -R $USERUSERNAME:$USERUSERNAME /home/$USERUSERNAME/.ssh && \
    chmod 700 /home/$USERUSERNAME/.ssh && \
    chmod 600 /home/$USERUSERNAME/.ssh/authorized_keys

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