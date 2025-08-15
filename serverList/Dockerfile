#使用buildx构建 避免跨平台问题
#docker buildx create --use
#docker buildx build --platform=linux/amd64 -t rosehtx/serverlist:v1 .

#如果构建的时候连接https://auth.docker.io/token异常可通过构建本地镜像源来处理
#先构建本地镜像源服务 docker pull registry:2.8.1
#docker run -d -p 5000:5000 --restart=always --name registry registry:2.8.1
#给自己需要的镜像打tag并推送到本地镜像源
#docker tag golang:1.18.2 127.0.0.1:5000/golang:1.18.2
#docker push 127.0.0.1:5000/golang:1.18.2
#(可通过 curl -X GET http://127.0.0.1:5000/v2/_catalog查看有哪些镜像推送到本地源)

#继续用命令构建如果还是提示链接异常看如下（感觉上面的链接超时也是这问题，可以试试不构建本地镜像直接如下尝试）
#buildx(docker buildx ls命令查看) 默认使用的是docker desktop内部的虚拟环境/WSL2执行构建,会链接不上127.0.0.1
#需要让buildx使用本地模式 docker context use default => docker buildx use default

# 通过本地镜像源来进行构建  127.0.0.1:5000 是构建的本地私有镜像源
# 推不了远程docker hub可以重新登录下 docker login
# 构建阶段
FROM 127.0.0.1:5000/golang:1.18.2 AS builder
WORKDIR /app
COPY . .
# 编译项目 得加CGO_ENABLED=0 不然k8s运行pod的时候会失败,找不到serverList可执行文件,alpine的问题
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o serverList .

# 运行阶段
FROM 127.0.0.1:5000/alpine:latest
WORKDIR /app
COPY --from=builder /app/serverList .
RUN chmod +x serverList
CMD ["./serverList"]