FROM ubuntu:22.04

ARG JAVA_VERSION
ARG USERUSERNAME
ARG USERUSERPWSSHKEY

SHELL ["bash"]

RUN useradd -m -s /bin/bash $USERUSERNAME -G sudo && \
    mkdir -p /home/$USERUSERNAME/.ssh && \
    echo "$USERUSERPWSSHKEY" > /home/$USERUSERNAME/.ssh/authorized_keys && \
    chown -R $USERUSERNAME:$USERUSERNAME /home/$USERUSERNAME/.ssh && \
    chmod 700 /home/$USERUSERNAME/.ssh && \
    chmod 600 /home/$USERUSERNAME/.ssh/authorized_keys

RUN sed -i 's/^#PermitRootLogin.*/PermitRootLogin no/' /etc/ssh/sshd_config

RUN history -c && \
    rm -rf /root/.bash_history

USER $USERNAME

RUN apt purge java

RUN apt update && apt upgrade -y && \
    apt install -y openjdk-${JAVA_VERSION}-jdk

ENV JAVA_HOME /usr/lib/jvm/java-${JAVA_VERSION}-openjdk-amd64

RUN java -version

RUN history -c && \
    rm -rf /home/$USERUSERNAME/.bash_history

EXPOSE 22

CMD ["bash"]
