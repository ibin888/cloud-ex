### 第一部分需求和答案：
##### 第一步: 测试一下镜像还能不能用
* 测试image是否可以使用
```shell
apiVersion: v1
kind: Pod
metadata:
  name: http01
spec:
  containers:
    - name: http01
      image: http01
      imagePullPolicy: IfNotPresent
      ports:
        - containerPort: 80

```
* 由于docker build 在k8s这台机器上，路由有问题，我更换到harbor的那台机器进行build
* 构建完成后上传到harbor 或者 docker hub
##### 第二步：需求
* 优雅启动:使用 startup 探针处理启动探测  
* 优雅终止 最好的方式是使用golang处理term信号，但我这里使用terminationGracePeriodSeconds: 60
* 资源需求和 QoS 保证`limit和request`
* 探活`livenessProbe`
* 日常运维需求，日志等级 ： 代码中用klog
* 配置和代码分离，这个程序应该就一个端口配置吧。用`configmap`，然后用volume或者env来引入

### 操作
```azure
make build .
docker tag http01 avaisa/http01
docker push avaisa/http01
kubectl create -f port.yaml
kubect create  -f dp.yaml
```


### 第二部分需求和答案：
除了将 httpServer 应用优雅的运行在 Kubernetes 之上，我们还应该考虑如何将服务发布给对内和对外的调用方。
来尝试用 Service, Ingress 将你的服务发布给集群外部的调用方吧。
在第一部分的基础上提供更加完备的部署 spec，包括（不限于）：

* 如何确保整个应用的高可用
  * 将第一部分的pod用deployment来处理，（3个副本）
* Service
  * 将deployment管理的pod 用service 分发出去
* Ingress
  * 7层类似nginx
* 如何通过证书保证 httpServer 的通讯安全
  * ingress 自带tls 证书 见ingress.yaml
  * 使用cert-manager 见ingress02.yaml

### 操作
```azure
kubect create  -f ca.yaml
kubectl create -f ingress02.yaml 或者kubectl create -f ingress.yaml
```

* 测试结果
```shell
k -n ingress-nginx get svc
curl -H "Host: cncamp.com" https://10.107.192.218 -v -k
```
