FROM registry.cn-hangzhou.aliyuncs.com/nqc/golang:latest as builder

# Create app directory
RUN mkdir -p /home/app
WORKDIR /home/app
COPY . /home/app/subscribe2clash

ENV GOPATH="/home/app:${GOPATH}"
# 安装 git
#RUN apk add git && git config --global http.sslVerify "false"
# && git config --global --add url."git@code.aliyun.com:".insteadOf "https://code.aliyun.com/"
#COPY .ssh /root/
#RUN chown 1000:1000 /root/.ssh/id_rsa
#go env -w GONOPROXY="code.aliyun.com,gopkg.in" &&  go env -w  GONOSUMDB="code.aliyun.com,gopkg.in" && go mod vendor &&
RUN cd subscribe2clash &&  go build -mod vendor -o ../bin/subscribe2clash *.go && \
  chmod +x ../bin/subscribe2clash

FROM registry.cn-hangzhou.aliyuncs.com/nqc/alpine:latest

WORKDIR /home/app
RUN mkdir -p /home/app/file

COPY --from=builder /home/app/bin ./bin

EXPOSE 10028
CMD [ "./bin/subscribe2clash" ]



