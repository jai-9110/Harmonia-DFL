# 安裝Gitea
```
$ sudo nano gitea1_deployment.yml
$ sudo nano gitea2_deployment.yml
$ sudo nano gitea3_deployment.yml
```
[gitea1_deployment.yml](  
[gitea2_deployment.yml]( 
[gitea3_deployment.yml]
```
$ kubectl apply -f gitea1_deployment.yml
$ kubectl apply -f gitea2_deployment.yml
$ kubectl apply -f gitea3_deployment.yml
```
> 在三台實體機上各自建立Gitea server  
> 連線至Gitea前要先設定Mobaxterm的Tunneling
* port num根據上方的gitea.yml中所寫的進行設定，其餘設定與上方教學相同
![image](
