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
|---|:---:|:---:|:---:|:---:|:---:|  
| aggregator1 |	可讀 | 可寫	| 可讀 | 可讀 | 可讀 |  
| edge1	| 可讀 | 可讀 | 可寫 | - | - |  
| edge2	| -	| -	| -	| 可寫	| - |  
| edge3	| - | - | - | - | 可寫 |  
| *webhook* | http://mnist-aggregator1:9080 | http://mnist-edge1:9080 | http://mnist-aggregator1:9080 | http://mnist-aggregator1:9080 | http://mnist-aggregator1:9080 | 
| | http://mnist-edge1:9080 |			

| gitea2 | `train-plan2` | `global model2` | `local model1` | `local model2` | `local model3` |  
|---|:---:|:---:|:---:|:---:|:---:|  
| aggregator2 |	可讀 | 可寫	| 可讀 | 可讀 | 可讀 |  
| edge1	| - | - | 可寫 | - | - |  
| edge2	| 可讀	| 可讀	| -	| 可寫	| - |  
| edge3	| - | - | - | - | 可寫 |  
| *webhook* | http://mnist-aggregator2:9080 | http://mnist-edge2:9080 | http://mnist-aggregator2:9080 | http://mnist-aggregator2:9080 | http://mnist-aggregator2:9080 | 
| | http://mnist-edge2:9080 |		

| gitea3 | `train-plan3` | `global model3` | `local model1` | `local model2` | `local model3` |  
|---|:---:|:---:|:---:|:---:|:---:|  
| aggregator3 |	可讀 | 可寫	| 可讀 | 可讀 | 可讀 |  
| edge1	| - | - | 可寫 | - | - |  
| edge2	| -	| -	| -	| 可寫	| - |  
| edge3	| 可讀 | 可讀 | - | - | 可寫 |  
| *webhook* | http://mnist-aggregator3:9080 | http://mnist-edge3:9080 | http://mnist-aggregator3:9080 | http://mnist-aggregator3:9080 | http://mnist-aggregator3:9080 | 
| |http://mnist-edge3:9080 |		

4. 獲取GiteaUserToken，詳細方法可參考[Harmonia-FL Tutorial](https://github.com/jai-9110/Harmonia-FL/tree/main/%E5%AE%89%E8%A3%9DHarmonia)
> 所有用戶的token都要記下來  
```
$ sudo nano configs_DFL.yml // 上面記錄的token要寫到此檔案中
```
* token的順序與edgeMOdelRepos的順序相同，不可對調  
[configs_DFL.yml](https://github.com/jai-9110/Harmonia-DFL/blob/5a9ff4aeae9f4905f119aa5a7e697864ec1a7fb4/%E5%AE%89%E8%A3%9DHarmonia/configs_DFL.yml)

```
$ kubectl apply -f configs_DFL.yml
```

# 安裝Harmonia
```
$ git clone https://github.com/ailabstw/harmonia.git
$ sudo nano harmonia/src/protos/service.proto   // 版本更新
```
* 替換第三行為此內容
```
option go_package = “./;protos”;     
```

```
$ cd harmonia/src/steward/operator/edge/
$ ls
$ rm -rf payload.go
$ nano payload.go
```
[payload.go](https://github.com/jai-9110/Harmonia-DFL/blob/fe4bd87041fd9690925422a781e27e92f9eff56f/%E5%AE%89%E8%A3%9DHarmonia/DFL%E6%9E%B6%E6%A7%8B%E7%9A%84.go/operator/edge/payload.go)
```
$ rm -rf stateMachine.go
$ nano stateMachine.go
```
[stateMachine.go](https://github.com/jai-9110/Harmonia-DFL/blob/fe4bd87041fd9690925422a781e27e92f9eff56f/%E5%AE%89%E8%A3%9DHarmonia/DFL%E6%9E%B6%E6%A7%8B%E7%9A%84.go/operator/edge/stateMachine.go)
```
$ cd ..
$ rm -rf operator.go
$ nano operator.go
```
[operator.go](https://github.com/jai-9110/Harmonia-DFL/blob/fe4bd87041fd9690925422a781e27e92f9eff56f/%E5%AE%89%E8%A3%9DHarmonia/DFL%E6%9E%B6%E6%A7%8B%E7%9A%84.go/operator/operator.go)
```
$ cd ..
$ cd config
$ rm -rf config.go
$ nano config.go
```
[config.go](https://github.com/jai-9110/Harmonia-DFL/blob/fe4bd87041fd9690925422a781e27e92f9eff56f/%E5%AE%89%E8%A3%9DHarmonia/DFL%E6%9E%B6%E6%A7%8B%E7%9A%84.go/config/config.go)
```
$ cd iconfig
$ rm -rf iconfig.go
$ nano iconfig.go
```
[iconfig.go](https://github.com/jai-9110/Harmonia-DFL/blob/fe4bd87041fd9690925422a781e27e92f9eff56f/%E5%AE%89%E8%A3%9DHarmonia/DFL%E6%9E%B6%E6%A7%8B%E7%9A%84.go/config/iconfig.go)
```
$ cd ..
$ rm -rf steward.go
$ nano steward.go
```
[steward.go](https://github.com/jai-9110/Harmonia-DFL/blob/fe4bd87041fd9690925422a781e27e92f9eff56f/%E5%AE%89%E8%A3%9DHarmonia/DFL%E6%9E%B6%E6%A7%8B%E7%9A%84.go/steward.go)
```
$ cd ..
$ cd ..
$ sudo make all
```
* ##每次修改過程式碼都要再次執行sudo make all，而在make all之前要先執行make clean，才能make all
check : 檢查images 是否有operator fedavg
```
$ cp src/protos/python_protos/service_pb2_grpc.py examples/edge/
$ cp src/protos/python_protos/service_pb2.py examples/edge/
```
```
$ sudo docker build examples/edge --tag 10.135.170.112:5000/mnist-edge1
$ sudo docker build examples/edge --tag 10.135.170.112:5000/mnist-edge2
$ sudo docker build examples/edge --tag 10.135.170.112:5000/mnist-edge3

$ sudo docker tag harmonia/operator 10.135.170.112:5000/harmonia-operator-agg1
$ sudo docker tag harmonia/operator 10.135.170.112:5000/harmonia-operator-agg2
$ sudo docker tag harmonia/operator 10.135.170.112:5000/harmonia-operator-agg3
$ sudo docker tag harmonia/operator 10.135.170.112:5000/harmonia-operator-edge1
$ sudo docker tag harmonia/operator 10.135.170.112:5000/harmonia-operator-edge2
$ sudo docker tag harmonia/operator 10.135.170.112:5000/harmonia-operator-edge3

$ sudo docker tag harmonia/fedavg 10.135.170.112:5000/harmonia-fedavg1
$ sudo docker tag harmonia/fedavg 10.135.170.112:5000/harmonia-fedavg2
$ sudo docker tag harmonia/fedavg 10.135.170.112:5000/harmonia-fedavg3

$ sudo docker push 10.135.170.112:5000/mnist-edge1
$ sudo docker push 10.135.170.112:5000/mnist-edge2
$ sudo docker push 10.135.170.112:5000/mnist-edge3

$ sudo docker push 10.135.170.112:5000/harmonia-operator-agg1
$ sudo docker push 10.135.170.112:5000/harmonia-operator-agg2
$ sudo docker push 10.135.170.112:5000/harmonia-operator-agg3

$ sudo docker push 10.135.170.112:5000/harmonia-operator-edge1
$ sudo docker push 10.135.170.112:5000/harmonia-operator-edge2
$ sudo docker push 10.135.170.112:5000/harmonia-operator-edge3

$ sudo docker push 10.135.170.112:5000/harmonia-fedavg1
$ sudo docker push 10.135.170.112:5000/harmonia-fedavg2
$ sudo docker push 10.135.170.112:5000/harmonia-fedavg3
```
```
$ sudo nano mnist_DFL.yml
```
[mnist_DFL.yml](https://github.com/jai-9110/Harmonia-DFL/blob/5a4ac31e45ad051eb9876fdd657e12785cf593c4/%E5%AE%89%E8%A3%9DHarmonia/mnist_DFL.yml)
```
$ kubectl apply -f mnist_DFL.yml
check : $ kubectl get po
```
![image](https://github.com/jai-9110/Harmonia-DFL/blob/ca4cc4056ed36ef5e60ebfd0230cf78df0e74a67/picture/DFL_pod.png)

# 開始訓練
```
$ git clone http://127.0.0.1:3001/gitea1/global-model1.git
$ pushd global-model1

$ git commit -m "pretrained model1" --allow-empty
$ git push origin master

$ popd
$ rm -rf global-model1
```
---------------------------------------------------------
```
$ git clone http://127.0.0.1:3002/gitea2/global-model2.git
$ pushd global-model2

$ git commit -m "pretrained model2" --allow-empty
$ git push origin master

$ popd
$ rm -rf global-model2
```
------------------------------------------------------------
```
$ git clone http://127.0.0.1:3003/gitea3/global-model3.git
$ pushd global-model3

$ git commit -m "pretrained model3" --allow-empty
$ git push origin master

$ popd
$ rm -rf global-model3
```
-------------------------------------------------------
```
$ git clone http://127.0.0.1:3001/gitea1/train-plan1.git
$ pushd train-plan1

$ cat > plan.json << EOF
{
    "name": "MNIST",
    "round": 10,
    "edge": 3,
    "EpR": 1,
    "timeout": 90,
    "pretrainedModel": "master"
}
EOF

$ git add plan.json
$ git commit -m "train plan1 commit”
$ git push origin master

$ popd
$ rm -rf train-plan1
```
-------------------------------------------------------
```
$ git clone http://127.0.0.1:3002/gitea2/train-plan2.git
$ pushd train-plan2

$ cat > plan.json << EOF
{
    "name": "MNIST",
    "round": 12,
    "edge": 3,
    "EpR": 1,
    "timeout": 90,
    "pretrainedModel": "master"
}
EOF

$ git add plan.json
$ git commit -m "train plan2 commit”
$ git push origin master

$ popd
$ rm -rf train-plan2
```
----------------------------------------------------------------
```
$ git clone http://127.0.0.1:3003/gitea3/train-plan3.git
$ pushd train-plan3

$ cat > plan.json << EOF
{
    "name": "MNIST",
    "round": 15,
    "edge": 3,
    "EpR": 1,
    "timeout": 90,
    "pretrainedModel": "master"
}
EOF

$ git add plan.json
$ git commit -m "train plan3 commit”
$ git push origin master

$ popd
$ rm -rf train-plan3
```
* 查看pod運行情況(log)
```
$ kubectl logs <pod name> <container name> -f
```
