# 安裝Gitea
```
$ sudo nano gitea1_deployment.yml
$ sudo nano gitea2_deployment.yml
$ sudo nano gitea3_deployment.yml
```
[gitea1_deployment.yml](https://github.com/jai-9110/Harmonia-DFL/blob/889284300bcd6de6b3d9fe9a0e8467a9335db4a3/%E5%AE%89%E8%A3%9DHarmonia/gitea1_deployment.yml)    
[gitea2_deployment.yml](https://github.com/jai-9110/Harmonia-DFL/blob/889284300bcd6de6b3d9fe9a0e8467a9335db4a3/%E5%AE%89%E8%A3%9DHarmonia/gitea2_deployment.yml)  
[gitea3_deployment.yml](https://github.com/jai-9110/Harmonia-DFL/blob/889284300bcd6de6b3d9fe9a0e8467a9335db4a3/%E5%AE%89%E8%A3%9DHarmonia/gitea3_deployment.yml)  
```
$ kubectl apply -f gitea1_deployment.yml
$ kubectl apply -f gitea2_deployment.yml
$ kubectl apply -f gitea3_deployment.yml
```
> 在三台實體機上各自建立Gitea server  
> 連線至Gitea前要先設定Mobaxterm的Tunneling
* port num根據上方的gitea.yml中所寫的進行設定，其餘設定與上方教學相同
![image](https://github.com/jai-9110/Harmonia-DFL/blob/88e04ef46d8b844352e136923f705c72473b7d53/picture/gitea_tunnel.png)

```
$ kubectl port-forward --address 0.0.0.0 service/harmonia-gitea1 3001
$ kubectl port-forward --address 0.0.0.0 service/harmonia-gitea2 3002
$ kubectl port-forward --address 0.0.0.0 service/harmonia-gitea3 3003
```
> 開啟三個shell，執行上方三行指令  
> 開啟三個瀏覽器，分別輸入192.168.71.10:3001、192.168.71.10:3002、192.168.71.10:3003，往下拉後直接安裝  
