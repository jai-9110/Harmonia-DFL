* 所有實體機都要安裝
```
$ sudo apt-get update && sudo apt-get install -y apt-transport-https curl
# curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
$ cat <<EOF | sudo tee /etc/apt/sources.list.d/kubernetes.list
deb https://apt.kubernetes.io/kubernetes-xenial main
EOF
$ sudo apt-get update

$ sudo apt install -y kubelet=1.20.0-00 kubectl=1.20.0-00 kubeadm=1.20.0-00
$ sudo apt-mark hold kubelet kubeadm kubectl
```
> reset K8S  
> ```  
> $ sudo apt-get purge kubeadm kubectl kubelet kubernetes-cni kube*  
> $ sudo rm -rf ~/.kube  
> $ kubeadm reset

* Master、Worker都要
* 設定Registry，IP為Master的IP(Registry的位置在Master)  
叢集內所有Worker也都設定成Master的IP
```
sudo vi /etc/docker/daemon.json
```
```
{
    "live-restore": true,
    "group": "dockerroot",
    "insecure-registries": ["10.135.170.112:5000"]
} 
```

* 只有Master需要執行此行
```
$ systemctl restart docker
$ sudo docker run -d -p 5000:5000 -v ~/storage:/var/lib/registry --name registry registry:2 
```

* Master、Worker都要
* 額外設定 : IP設定為各節點的IP
```
$ sudo vi /etc/systemd/system/kubelet.service.d/10-kubeadm.conf 
```
```
Environment="KUBELET_EXTRA_ARGS=--node-ip=10.135.170.112”
```
```
$ systemctl restart kubelet
```
* 只有Master需要執行此行
```
$  sudo kubeadm config images pull 
```

# 安裝Calico
* 只有Master需安裝
```
# kubeadm init --apiserver-advertise-address=10.135.170.112 --pod-network-cidr=192.168.244.0/24

$ mkdir -p $HOME/.kube
$ sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
$ sudo chown $(id -u):$(id -g) $HOME/.kube/config
```
> 記下kubeadm join這行指令，到其他Worker中輸入此行指令，便可加入K8S叢集  
> token與sha256 24小時會過期，需重新申請  
> https://kubernetes.io/zh-cn/docs/setup/production-environment/tools/kubeadm/create-cluster-kubeadm/

```
$ kubectl apply -f https://docs.projectcalico.org/archive/v3.20/manifests/calico.yaml
check : $ kubectl get pods --all-namespaces
```
![image](https://github.com/jai-9110/Harmonia-DFL/blob/59b8812cdf66dc76bf114143dc0d80eb0756dd16/picture/get_pod1.png)
```
check : $ kubectl get node
```
![image](https://github.com/jai-9110/Harmonia-DFL/blob/3cc5922ed716eb8a66a75b1e8e7c1259eb21d04f/picture/get_node.png)

