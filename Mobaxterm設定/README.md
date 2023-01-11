1. 連線學校VPN
> 中原VPN教學 :   
> https://itouch.cycu.edu.tw/active_project/cycu2900h_04/cycu_01/new/SSL%20VPN%20%E5%AE%89%E8%A3%9D%E8%A8%AD%E5%AE%9A%E6%95%99%E5%AD%B8%E6%96%87%E4%BB%B6-201902.pdf
2. 點選Tunneling  
![image](https://github.com/jai-9110/Harmonia-DFL/blob/6d873a7046ea9ac211290a73b6d675d2151a2ac1/picture/Tunneling.png)
3. 點選右下角的New SSH tunnel
![image](https://github.com/jai-9110/Harmonia-DFL/blob/c3a36d240ee2fac43ae49ea29de93e536f2502c9/picture/NEW_SSH_Tunnel.png)
4. 左邊的Forwarded port為你的電腦對外轉發的接口(22)  
右下為實驗室的伺服器，IP為140.135.11.223，名稱為cdlab-junior，port為22  
右上為實驗室的路由器，第一欄為要連接的該台實體機的IP(eg. 10.135.170.112)，第二欄為所要連入的該台實體機對外的port num : 22
![image](https://github.com/jai-9110/Harmonia-DFL/blob/af3b57f177f230b9b7ee4c0d73cd2b8b7da18d63/picture/SSH_Tunnel-1.png)
![image](https://github.com/jai-9110/Harmonia-DFL/blob/e2536eacc805f4a3c47a77284e583229bcfcd6cc/picture/SSH_Tunnel-2.png)
5. 按Save保存，點選金鑰  
![image](https://github.com/jai-9110/Harmonia-DFL/blob/e2dcdc42937285e9ba5eb4e7872e0323796fcf59/picture/key.png)
> 金鑰請跟實驗室學長姐領取  
6. 點選Local network adapter(金鑰前的綠色選項)，選擇一個虛擬網卡的IP  
> 沒有虛擬網卡可透過Virtual box的主機網路管理員新增一張  
![image](https://github.com/jai-9110/Harmonia-DFL/blob/9acdcda48523064f5f647adccbceaee5abba1e6d/picture/Local_network_adapter.png)
7. 設定完便可按下Start，就可遠端連線至實驗室的實體機  
> 若在實驗室連線至實體機則不須連上學校VPN  
