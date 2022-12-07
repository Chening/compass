### **第二部分**

• service
• ingress
• 如何保证应用高可用
• 如何通过证书保证httpserver通讯安全

在 httpserver.go 的同级目录创建 Dockerfile 文件

---

```yaml
# 安装 ingress controller
# 对k8s源的镜像转换镜像源
https://github.com/anjia0532/gcr.io_mirror

wget https://raw.githubusercontent.com/anjia0532/gcr.io_mirror/master/pull-k8s-image.sh
chmod +x pull-k8s-image.sh

./pull-k8s-image.sh  k8s.gcr.io/ingress-nginx/controller:v1.0.0
./pull-k8s-image.sh  k8s.gcr.io/ingress-nginx/kube-webhook-certgen:v1.0

# 创建 ingress controller
# 修改 k8s.gcr.io 为  anjia0532/google-containers
kubectl apply -f nginx-ingress-deployment.yaml

# 创建完成之后 会生成3个pod 1个service

kubectl get pods -n ingress-nginx

kubectl get svc -n ingress-nginx
NAME                                 TYPE        CLUSTER-IP    EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             NodePort    10.96.4.226   <none>        80:32733/TCP,443:30359/TCP   40m
ingress-nginx-controller-admission   ClusterIP   10.96.44.56   <none>        443/TCP                      40m

# 创建证书
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=cncamp.com/O=cncamp"
kubectl create secret tls cncamp-tls --cert=./tls.crt --key=./tls.key

kubectl get secret

# 创建 service
kubectl apply -f   httpserver-service.yaml

# 测试 service 
[root@master-1 008]# kubectl get svc
NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
httpsvc      ClusterIP   10.96.100.38   <none>        80/TCP    11m
kubernetes   ClusterIP   10.96.0.1      <none>        443/TCP   80d

[root@master-1 008]# curl  10.96.100.38/healthz

This is my status eq 200 page!

# 创建 ingress
httpserver-ingress.yaml

# 测试 ingress
curl -H "Host: cncamp.com" https://10.96.4.226/healthz -k
```