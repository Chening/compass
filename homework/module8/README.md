### 第一部分

将httpserver部署到kubernetes集群中，做到如下要求：

公共部分

```bash
//制作1.8.1 镜像 
//代码主要修改点 
// 1. 添加配置文件传递
// 2. 添加优雅终止代码
# 本地打包镜像
set GOOS=linux
set  GOARCH=amd64
go build  .\httpserver.go  httpserver

docker build -t httpserver:v1.8.1 .

# 打包远程镜像
docker tag httpserver:v1.8.1 totorest/httpserver:v1.8.1

# 登录远程镜像仓库
docker login -u totorest

# 推送仓库
docker push totorest/httpserver:v1.8.1

# 推送完成镜像后退出
docker logout

# 获取认证 base 信息
[root@master-1 ~]# docker login
Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username: totorest
Password: 
WARNING! Your password will be stored unencrypted in /root/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store

// 将认证信息进行 base64处理 并且记录
// cat /root/.docker/config.json | base64 -w0
//
kubectl create secret generic cloudnative     --from-file=.dockerconfigjson=/root/.docker/config.json     --type=kubernetes.io/dockerconfigjson

// 查看创建的secret
kubectl get  secret

```

• 优雅启动

```yaml

          readinessProbe:
            httpGet:
              path: /healthz
              port: 19004
              scheme: HTTP
            initialDelaySeconds: 5
            periodSeconds: 3
   
```

• 优雅中止

```yaml
// 在程序添加对 SIGTERM 信号的监听处理

quit := make(chan os.Signal)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
logrus.Info("我被优雅终止了")
```

• 资源需求和QoS

```yaml

          resources:
            limits:
              memory: "200Mi"
              cpu: 200m
            requests:
              memory: "100Mi"
              cpu: 20m
```

• 探活

```yaml
同  优雅启动  处理
```

• 日常运维需求，日志等级

```yaml
// 添加日志等级的参数的传递  config.yml  中的参数配置如下
server:
  log: ERR

// 代码中添加如下部分 
// TODO more case
if logLevel == "INFO" {
	logrus.SetLevel(logrus.InfoLevel)
} else if logLevel == "ERR" {
	logrus.SetLevel(logrus.ErrorLevel)
}
```

• 配置和代码分离

```yaml
//1. 程序添加 config.yaml 的配置文件 将程序需要使用的参数放入文件中 并且在

// 创建configmap
touch config.yml

// 将如下配置写入  yaml 文件
apiVersion: v1
kind: ConfigMap
metadata:
 name: httpserver-config
data:
 config.yaml: |
     server:
      log: INFO

//创建 配置
kubectl apply -f config.yml
```

完整的deployment 如下

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: httpserver
  labels:
    app: httpserver
spec:
  replicas: 1
  selector:
    matchLabels:
      app: httpserver
  template:
    metadata:
      labels:
        app: httpserver
    spec:
      imagePullSecrets:
        - name: cloudnative
      containers:
        - name: httpserver
          image: totorest/httpserver:v1.8.1
          ports:
          - containerPort: 19004
          readinessProbe:
            httpGet:
              path: /healthz
              port: 19004
              scheme: HTTP
            initialDelaySeconds: 5
            periodSeconds: 3
          resources:
            limits:
              memory: "200Mi"
              cpu: 200m
            requests:
              memory: "100Mi"
              cpu: 20m
          volumeMounts:
          - name: config
            mountPath: /app
      volumes:
      - name: config
        configMap:
          name: httpserver-config
```

```yaml
kubectl apply -f httpserver-deployment.yaml

//查询刚刚创建的 pod

kubectl get pods

[root@master-1 008]# kubectl get pods
NAME                                      READY   STATUS    RESTARTS        AGE
httpserver-65c5889978-p5x9l               1/1     Running   0               30m

// 进入pod 内部 查看配置文件
kubectl exec  -it  httpserver-65c5889978-p5x9l  sh

/ # pwd
/
/ # ls
app         bin         dev         etc         home        httpserver  lib         media       mnt         opt         proc        root        run         sbin        srv         sys         tmp         usr         var
/ # cd app/
/app # ls
config.yaml
/app #

在 app 目录下 看到 config.yaml 配置文件
```