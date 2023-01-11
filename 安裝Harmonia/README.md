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

## 設定Gitea(以gitea1為例)
1. 註冊gitea1、aggregator1、edge1、edge2、edge3    // 電子信箱可隨便打，不要重複即可
![image](https://github.com/jai-9110/Harmonia-DFL/blob/7b9eaae0d6cbad22babca045e2468f357fabd504/picture/%E8%A8%BB%E5%86%8A%E5%B8%B3%E8%99%9F.png)
2. 登入gitea1，點選右上角的新增儲存庫
![image](https://github.com/jai-9110/Harmonia-DFL/blob/6356d46845b3a1aadb39aa727655878e84625e6f/picture/%E6%96%B0%E5%A2%9E%E5%84%B2%E5%AD%98%E5%BA%AB.png)
3. 設定儲存庫，詳細設定可參考Harmonia-FL Tutorial
| gitea1 | `train-plan1` | `global model1` | `local model1` | `local model2` | `local model3` |
| --- | :---: | :---: | :---: | :---: | :---: |
| aggregator1 |	可讀 | 可寫	| 可讀 | 可讀 | 可讀 |
| edge1	| 可讀 | 可讀 | 可寫 | - | - |
| edge2	| -	| -	| -	| 可寫	| - | 
| edge3	| - | - | - | - | 可寫 | 
| *webhook* | http://mnist-aggregator1:9080 | http://mnist-edge1:9080 | http://mnist-aggregator1:9080 | http://mnist-aggregator1:9080 | http://mnist-aggregator1:9080 |
| | http://mnist-edge1:9080 |				
