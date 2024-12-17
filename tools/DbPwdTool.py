import hashlib
from base64 import b64decode

from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad

from .PrintTool import *


class DbPwdTool:
    def __init__(self) -> None:
        pass

    def wildfire(self, token: str):
        """
        @token          用户登录的token，应用目录下shared_prefs/config.xml的token的值
        """
        # 目前新版的密钥是001122334455667778797A7B7C7D7E7F
        key = bytes.fromhex('001122334455667778797A7B7C7D7E7F')
        # 老版的密钥是7F7E7D7C7B7A79787766554433221100
        key2 = bytes.fromhex('7F7E7D7C7B7A79787766554433221100')
        iv = key
        iv2 = key2
        print('[提示]---->正在计算野火数据库密钥')
        flag = 0
        try:
            cipher = AES.new(key, AES.MODE_CBC, iv)
            msg = cipher.decrypt(b64decode(token))
            msg = unpad(msg, AES.block_size).decode("utf8")
        except:
            # 新版无法解密则使用老版解密
            cipher = AES.new(key2, AES.MODE_CBC, iv2)
            msg = cipher.decrypt(b64decode(token))
            msg = unpad(msg, AES.block_size).decode("gbk")
            flag = 1
        pwd = msg.split("|")[-1]
        print_green_key(f'[成功]---->获取到数据库密钥', pwd)
        if flag == 1:
            print_red('[注意]---->该数据库为老版本，请选择SQLCipher3进行解密!')
        else:
            print_red('[注意]---->该数据库为新版本，请选择SQLCipher4进行解密!')
        return pwd, flag

    def wechat(self, uin: str, imei: str = '1234567890ABCDEF'):
        """
        @uin            uin，可能是负值，在shared_prefs/auth_info_key_prefs.xml文件中_auth_uin的值
        @imei           微信获取到的IMEI或MEID，在shared_prefs/DENGTA_META.xml文件中IMEI_DENGTA的值，在高版本中通常是1234567890ABCDEF
        """
        print(f'[提示]---->正在计算微信数据库EnMicroMsg.db密钥')
        md5 = hashlib.md5()
        md5.update((imei + uin).encode())
        pwd = md5.hexdigest()[:7]
        print_green_key('[成功]---->获取到数据库密钥', pwd)
        return pwd

    def wechat_index(self, uin: str, wxid: str, imei: str = '1234567890ABCDEF'):
        """
        @uin            uin，可能是负值，在shared_prefs/auth_info_key_prefs.xml文件中_auth_uin的值
        @wxid           数据库所属的wxid，一般情况下在解密EnMicroMsg.db的时候会一并提取，若无需要，请从shared_prefs/com.tencent.mm_preferences.xml中提取login_weixin_username的值
        @imei           微信获取到的IMEI或MEID，在shared_prefs/DENGTA_META.xml文件中IMEI_DENGTA的值，在高版本中通常是1234567890ABCDEF
        """
        if uin.startswith('-'):
            uin = str(int(uin) + 4294967296)
        print(f'[提示]---->正在计算微信索引数据库密钥')
        md5 = hashlib.md5()
        md5.update((uin + imei + wxid).encode())
        pwd = md5.hexdigest()[:7]
        print_green_key(f'[成功]---->获取到索引数据库密钥', pwd)
        return pwd

    def mostone(self, uid: str):
        """
        @uid            uid，在shared_prefs/im.xml文件中的userId的值
        """
        print(f'[提示]---->正在计算默往数据库msg.db密钥')
        md5 = hashlib.md5()
        md5.update((uid).encode())
        pwd = md5.hexdigest()[:6].upper()
        print_green_key('[成功]---->获取到数据库密钥', pwd)
        print_red('[注意]---->请选择SQLCipher3进行解密!')
        return pwd

    def tiktok(self, uid: str):
        '''
        @uid            uid，一般在数据库文件名中就有
        '''
        print(f'[提示]---->正在计算抖音数据库密钥')
        pwd = f'byte{uid}imwcdb{uid}dance'
        print_green_key('[成功]---->获取到数据库密钥', pwd)
        print_red('[注意]---->请选择wcdb进行解密!')
        return pwd
