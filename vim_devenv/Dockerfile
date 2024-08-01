FROM debian:12

RUN apt-get -y update && \
    apt-get -y install build-essential git autoconf make texinfo gnutls-bin gcc wget curl

RUN cd /tmp && \
    wget https://go.dev/dl/go1.22.5.linux-amd64.tar.gz && \
    rm -rf /opt/go && \
    mkdir -p /opt && \
    tar -C /opt -xzf go1.22.5.linux-amd64.tar.gz && \
    rm -rf /tmp/*

ENV PATH="$PATH:/opt/go/bin"

RUN go install golang.org/x/tools/gopls@latest && \
    go install mvdan.cc/gofumpt@latest && \
    go install github.com/spf13/cobra-cli@latest && \
    go install github.com/go-delve/delve/cmd/dlv@latest && \
    go install golang.org/x/tools/cmd/goimports@latest
ENV PATH="$PATH:/root/go/bin"

RUN mkdir -p /tmp/neovim && cd /tmp/neovim && \
    wget 'https://github.com/neovim/neovim/releases/download/v0.10.0/nvim-linux64.tar.gz' && \
    tar xzvf nvim-linux64.tar.gz && \
    mkdir -p /opt && mv -v nvim-linux64 /opt && \
    cd / && rm -rf /tmp/neovim
ENV PATH="$PATH:/opt/nvim-linux64/bin"

RUN mkdir -p /root/.config/nvim
COPY init.lua /root/.config/nvim/init.lua
