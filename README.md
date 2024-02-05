 # 用法
1. `pip install -r requirements.txt`
2. `python3 ForensicsTool.py -h`
```
usage: ForensicsTool.py [-h] [-m MODE] [-f FILE] [-t TYPE] [-p PASSWORD] [--uin UIN] [--imei IMEI] [--wxid WXID] [--token TOKEN]
                        [--device DEVICE]

Forensics Tool

optional arguments:
  -h, --help            show this help message and exit
  -m MODE, --mode MODE
                        指定需要运行的模式:
                            [0]表示计算密钥，支持的type值为1-3
                            [1]表示解密数据库，支持的type值为1、2、4-7
                            [2]表示数据提取，支持的type值为8-10
  -f FILE, --file FILE  指定需要解密的数据库
  -t TYPE, --type TYPE
                        指定需要处理的内容:
                            [1]微信的EnMicroMsg.db
                            [2]微信的FTS5IndexMicroMsg_encrypt.db
                            [3]野火IM系应用的data
                            [4]高德的girf_sync.db
                            [5]钉钉的数据库
                            [6]SQLCipher4加密的数据库
                            [7]SQLCipher3加密的数据库
                            [8]Navicat连接信息提取，需指定-f为目标用户的注册表文件"NTUSER.DAT"
                            [9]MobaXterm连接信息解密，可以指定MobaXterm.ini配置文件或用户注册表文件"NTUSER.DAT"，解密需要给出主密码   
                            [10]Dbeaver连接信息解密，指定-f为目标文件data-sources.json和credentials-config.json的父目录
  -p PASSWORD, --password PASSWORD
                        解密的密码，处理钉钉和高德时不适用
  --uin UIN             微信用户的uin，可能是负值，在shared_prefs/auth_info_key_prefs.xml文件中_auth_uin的值
  --imei IMEI           微信获取到的IMEI或MEID，在shared_prefs/DENGTA_META.xml文件中IMEI_DENGTA的值，在高版本中通常是1234567890ABCDEF 
，可以为空
  --wxid WXID           数据库所属的wxid，一般情况下在解密EnMicroMsg.db的时候会一并提取，若无需要，请从shared_prefs/com.tencent.mm_pre
ferences.xml中提取login_weixin_username的值
  --token TOKEN         野火IM系应用的用户token，shared_prefs/config.xml的token的值
  --device DEVICE       钉钉解密需要的内容，通常在shared_prefs/com.alibaba.android.rimet_preferences.xml中带有数据库名的字段的值中出现
，如HUAWEI P40/armeabi-v7a/P40/qcom/HUAWEIP40
```

# 样例
![Alt text](static/image.png)

![Alt text](static/image-1.png)

![Alt text](static/image-2.png)

# 声明
本工具仅用于电子数据取证的学习与研究，请勿用于非法用途，否则后果自负。
