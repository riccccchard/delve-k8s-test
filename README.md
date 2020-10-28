#  k8s 环境下使用delve 修改返回值测试

首先部署：

1. 在delveclient目录下打包docker镜像

   ```bash
   docker build --tag delveclient .
   ```

2. 在delveserver目录下打包docker镜像

   ```bash
   docker build --tag delveserver .
   ```

3. 在delveserver/deploy目录下，使用delveServer-deployment.yaml部署服务

   ```bash
   kubectl apply -f delveServer-deployment.yaml
   ```

   在delveServer的pod中，运行了两个容器，delveserver和delveclient

4. 通过kubectl exec进入pod，同时使用

   ```bash
   kubectl logs delveserver-cxqkm delveserver
   ```

   查看pod中delveServer container 的日志。

##  delve Server

delveserver启动一个delve server去attach目标进程

在delveServer/deploy/delveServer-deployment.yaml中，将其部署为daemonSet特权容器，添加SYS_PTRACE功能去ptrace其他进程；通过指明pid namespace为hostpid，使其能够获取node节点上的所有容器的pid。

delveServer启动一个http服务，监听node的3333端口，通过curl使其启动（注意要在容器内部，得exec进去）

```bash
curl -H "pid:8224" -H "address:127.0.0.1:7777" localhost:3333
```

pid为你想要attach的进程pid

address为delve server attach之后，监听客户端消息的端口，客户端delve client通过127.0.0.1:7777向delve server发送请求。



##  delve client

delveclient使用一个delve client与delveserver交互，修改db.Query返回值

delveclient启动一个http服务，监听Node的8888端口，通过curl使其启动（同样也要在容器内部）

```bash
curl -H "address:127.0.0.1:7777" localhost:8888
```

Address 为delve server监听的地址。



##  test

test启动一个http服务，执行sql查询。

通过http-deployment.yaml和http-service.yaml部署pod和service，这样就可以通过在本机访问pod的http服务

```bash
docker build --tag http_app
kubectl apply -f http-deployment.yaml
kubectl apply -f http-service.yaml
curl localhost:30307
```









