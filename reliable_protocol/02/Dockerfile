FROM debian:latest

# 设置环境变量
ENV GOLANG_VERSION 1.21.1
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

# 安装基础工具和 Go
RUN apt-get update && apt-get install -y \
    wget \
    git \
    gcc \
    iproute2 \
    net-tools \
    iputils-ping \
    && rm -rf /var/lib/apt/lists/* \
    && wget -P /tmp "https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz" \
    && tar -C /usr/local -xzf "/tmp/go${GOLANG_VERSION}.linux-amd64.tar.gz" \
    && rm "/tmp/go${GOLANG_VERSION}.linux-amd64.tar.gz"

# 创建工作目录
WORKDIR /app

# 复制 Go 项目文件
COPY . .

# 编译 Go 程序
RUN go mod download \
    && go build -o /cmd/client /internal/client && go build -o /cmd/server /internal/server

# 暴露端口（根据你的应用需要修改）
EXPOSE 8080

# 设置启动命令