import subprocess
import os
from Crypto.Cipher import AES
import sqlite3
import hashlib
from .PrintTool import *

class DbTool:
    def __init__(self) -> None:
        pass

    def get_relative_path(self,relative_path):
        """获取配置文件的绝对路径"""
        base_path = os.path.dirname(os.path.abspath(__file__))
        return os.path.join(base_path, relative_path)

    def decrypt_EnMicroMsg(self,key :str,path: str) -> int:
        """
        @key:       数据库密钥
        @path:      需要解密的数据库路径
        """
        sqlite3 = self.get_relative_path('../lib/sqlite3.exe')
        if os.path.exists(sqlite3) and os.path.isfile(sqlite3):
            sqlite3 = os.path.abspath(sqlite3)
        else:
            print_red('[失败]---->原因[缺少依赖<sqlite3.exe>！]')
            return -1
        print_yellow(f'<微信数据库解密>')
        print_yellow_key(f'[提示]---->正在解密微信数据库EnMicroMsg.db，解密密钥为',key)
        work_dir = os.getcwd()
        floder = os.path.realpath(path)
        floder = '\\'.join(floder.split('\\')[:-1])
        os.chdir(floder)
        # PRAGMA cipher_compatibility = 1是专为微信数据库设计的参数，需要SQLCipher版本至少为4.1
        result = subprocess.run([sqlite3,path,f"PRAGMA key = '{key}';PRAGMA cipher_compatibility = 1;ATTACH DATABASE 'EnMicroMsg_dec.db' AS EnMicroMsg_dec KEY '';SELECT sqlcipher_export('EnMicroMsg_dec');DETACH DATABASE EnMicroMsg_dec;"],capture_output=True, text=True)
        out = result.stdout
        err = result.stderr
        os.chdir(work_dir)
        if out.strip() == 'ok' and err == '':
            print_green('[成功]---->EnMicroMsg.db解密成功，解密后数据库在源文件同目录下EnMicroMsg_dec.db')
            wxid = self.getWXid(path)
            print_green_key(f'[成功]---->提取到wxid为',wxid)
            return 1
        else:
            print_red('[失败]---->EnMicroMsg.db解密失败')
            if err.strip().__contains__('already exists'):
                print_red(f'[失败]---->原因[数据库已解密过，解密后的文件已存在！]')
                try:
                    wxid = self.getWXid(path)
                    print_green_key(f'[成功]---->提取到wxid为',wxid)
                except:
                    pass
            elif err.strip().__contains__('file is not a database'):
                print_red(f'[失败]---->原因[解密密钥或解密参数不正确！]')
            else:
                print_red(f'[失败]---->原因[{err.strip()}]')
        return 0
    
    def decrypt_FTS5IndexMicroMsg(self,key :str,path: str) -> int:
        """
        @key:       数据库密钥
        @path:      需要解密的数据库路径
        """
        sqlite3 = self.get_relative_path('../lib/sqlite3.exe')
        if os.path.exists(sqlite3) and os.path.isfile(sqlite3):
            sqlite3 = os.path.abspath(sqlite3)
        else:
            print_red('[失败]---->原因[缺少依赖<sqlite3.exe>！]')
            return -1
        print_yellow(f'<微信数据库解密>')
        print_yellow_key(f'[提示]---->正在解密微信数据库FTS5IndexMicroMsg_encrypt.db，解密密钥为',key)
        work_dir = os.getcwd()
        floder = os.path.realpath(path)
        floder = '\\'.join(floder.split('\\')[:-1])
        os.chdir(floder)
        # 还需要设置page_size为4096，不过这是SQLCipher4的默认值，就不加了PRAGMA cipher_page_size = 4096;
        result = subprocess.run([sqlite3,path,f"PRAGMA key = '{key}';PRAGMA kdf_iter = 64000;PRAGMA cipher_kdf_algorithm = PBKDF2_HMAC_SHA1;PRAGMA cipher_hmac_algorithm = HMAC_SHA1;ATTACH DATABASE 'FTS5IndexMicroMsg.db' AS FTS5IndexMicroMsg KEY '';SELECT sqlcipher_export('FTS5IndexMicroMsg');DETACH DATABASE FTS5IndexMicroMsg;"],capture_output=True, text=True)
        out = result.stdout
        err = result.stderr
        os.chdir(work_dir)
        if out.strip() == 'ok' and err == '':
            print_green('[成功]---->FTS5IndexMicroMsg_encrypt.db解密成功，解密后数据库在源文件同目录下FTS5IndexMicroMsg.db')
            return 1
        else:
            print_red('[失败]---->FTS5IndexMicroMsg_encrypt.db解密失败')
            if err.strip().__contains__('already exists'):
                print_red(f'[失败]---->原因[数据库已解密过，解密后的文件已存在！]')
            elif err.strip().__contains__('file is not a database'):
                print_red(f'[失败]---->原因[解密密钥或解密参数不正确！]')
            else:
                print_red(f'[失败]---->原因[{err.strip()}]')
        return 0
    
    class DBHelper:
        def __init__(self,db_path: str) -> None:
            self.conn = sqlite3.connect(db_path)
            self.cursor = self.conn.cursor()

        def select(self,sql: str) -> list:
            """
            @sql:       要执行的SQL语句
            """
            self.cursor.execute(sql)
            rows = self.cursor.fetchall()
            return rows

        def close(self) -> None:
            self.cursor.close()
            self.conn.close()

    def getWXid(self,path :str) -> str:
        dbpath = path.replace('EnMicroMsg.db','EnMicroMsg_dec.db')
        db = self.DBHelper(dbpath)
        wxid = db.select('select value from userinfo where id = 2;')[0][0]
        db.close()
        return wxid

    def decrypt_SQLCipher4_default(self,key :str,path :str) -> int:
        """
        @key:       数据库密钥
        @path:      需要解密的数据库路径
        """
        sqlite3 = self.get_relative_path('../lib/sqlite3.exe')
        if os.path.exists(sqlite3) and os.path.isfile(sqlite3):
            sqlite3 = os.path.abspath(sqlite3)
        else:
            print_red('[失败]---->原因[缺少依赖<sqlite3.exe>！]')
            return -1
        bname = os.path.basename(path)
        name = bname.split('.')[0]
        work_dir = os.getcwd()
        floder = os.path.realpath(path)
        floder = '\\'.join(floder.split('\\')[:-1])
        os.chdir(floder)
        print_yellow(f'<SQLCipher4版本数据库解密>')
        print_yellow_key(f'[提示]---->正在解密数据库{bname}，解密密钥为',key)
        result = subprocess.run([sqlite3,path,f"PRAGMA key = '{key}';ATTACH DATABASE '{name}_dec.db' AS {name}_dec KEY '';SELECT sqlcipher_export('{name}_dec');DETACH DATABASE {name}_dec;"],capture_output=True, text=True)
        out = result.stdout
        err = result.stderr
        os.chdir(work_dir)
        if out.strip() == 'ok' and err == '':
            print_green(f'[成功]---->{bname}解密成功，解密后数据库在源文件同目录下{name}_dec.db')
            return 1
        else:
            print_red(f'[失败]---->{bname}解密失败')
            if err.strip().__contains__('already exists'):
                print_red(f'[失败]---->原因[数据库已解密过，解密后的文件已存在！]')
            elif err.strip().__contains__('file is not a database'):
                print_red(f'[失败]---->原因[解密密钥或解密参数不正确！]')
            else:
                print_red(f'[失败]---->原因[{err.strip()}]')
        return 0
    
    def decrypt_wcdb_default(self,key :str,path :str) -> int:
        """
        @key:       数据库密钥
        @path:      需要解密的数据库路径
        """
        sqlite3 = self.get_relative_path('../lib/sqlite3.exe')
        if os.path.exists(sqlite3) and os.path.isfile(sqlite3):
            sqlite3 = os.path.abspath(sqlite3)
        else:
            print_red('[失败]---->原因[缺少依赖<sqlite3.exe>！]')
            return -1
        bname = os.path.basename(path)
        name = bname.split('.')[0]
        work_dir = os.getcwd()
        floder = os.path.realpath(path)
        floder = '\\'.join(floder.split('\\')[:-1])
        os.chdir(floder)
        print_yellow(f'<wcdb数据库解密>')
        print_yellow_key(f'[提示]---->正在解密数据库{bname}，解密密钥为',key)
        result = subprocess.run([sqlite3,path,f"PRAGMA key = '{key}';PRAGMA cipher_page_size = 4096;PRAGMA kdf_iter = 64000;PRAGMA cipher_kdf_algorithm = PBKDF2_HMAC_SHA1;PRAGMA cipher_hmac_algorithm = HMAC_SHA1;ATTACH DATABASE '{name}_dec.db' AS {name}_dec KEY '';SELECT sqlcipher_export('{name}_dec');DETACH DATABASE {name}_dec;"],capture_output=True, text=True)
        out = result.stdout
        err = result.stderr
        os.chdir(work_dir)
        if out.strip() == 'ok' and err == '':
            print_green(f'[成功]---->{bname}解密成功，解密后数据库在源文件同目录下{name}_dec.db')
            return 1
        else:
            print_red(f'[失败]---->{bname}解密失败')
            if err.strip().__contains__('already exists'):
                print_red(f'[失败]---->原因[数据库已解密过，解密后的文件已存在！]')
            elif err.strip().__contains__('file is not a database'):
                print_red(f'[失败]---->原因[解密密钥或解密参数不正确！]')
            else:
                print_red(f'[失败]---->原因[{err.strip()}]')
        return 0

    def decrypt_SQLCipher3_default(self,key :str,path :str) -> int:
        """
        @key:       数据库密钥
        @path:      需要解密的数据库路径
        """
        sqlite3 = self.get_relative_path('../lib/sqlite3.exe')
        if os.path.exists(sqlite3) and os.path.isfile(sqlite3):
            sqlite3 = os.path.abspath(sqlite3)
        else:
            print_red('[失败]---->原因[缺少依赖<sqlite3.exe>！]')
            return -1
        bname = os.path.basename(path)
        name = bname.split('.')[0]
        work_dir = os.getcwd()
        floder = os.path.realpath(path)
        floder = '\\'.join(floder.split('\\')[:-1])
        os.chdir(floder)
        print_yellow(f'<SQLCipher3版本数据库解密>')
        print_yellow_key(f'[提示]---->正在解密数据库{bname}，解密密钥为',key)
        result = subprocess.run([sqlite3,path,f"PRAGMA key = '{key}';PRAGMA cipher_page_size = 1024;PRAGMA kdf_iter = 64000;PRAGMA cipher_kdf_algorithm = PBKDF2_HMAC_SHA1;PRAGMA cipher_hmac_algorithm = HMAC_SHA1;ATTACH DATABASE '{name}_dec.db' AS {name}_dec KEY '';SELECT sqlcipher_export('{name}_dec');DETACH DATABASE {name}_dec;"],capture_output=True, text=True)
        out = result.stdout
        err = result.stderr
        os.chdir(work_dir)
        if out.strip() == 'ok' and err == '':
            print_green(f'[成功]---->{bname}解密成功，解密后数据库在源文件同目录下{name}_dec.db')
            return 1
        else:
            print_red(f'[失败]---->{bname}解密失败')
            if err.strip().__contains__('already exists'):
                print_red(f'[失败]---->原因[数据库已解密过，解密后的文件已存在！]')
            elif err.strip().__contains__('file is not a database'):
                print_red(f'[失败]---->原因[解密密钥或解密参数不正确！]')
            else:
                print_red(f'[失败]---->原因[{err.strip()}]')
        return 0

    def decrypt_amap(self,file :str):
        if not os.path.exists(file):
            return '文件不存在！'
        print_yellow('<高德地图数据库解密>')
        print('[提示]---->正在获取数据库文件大小')
        size = hex(os.path.getsize(file)//1024)[2:]
        t = ''
        for i in range(8-len(size)):
            t += '0'
        # 第29-32字节，表示数据库的大小，取文件大小即可，补0
        size = t + size
        # 固定解密密钥
        key = 'a4a11bb9ef4b2f4c'
        # SQLite3固定文件头，16字节
        head = '53514C69746520666F726D6174203300'
        print_yellow_key(f'[提示]---->正在解密高德数据库girf_sync.db，解密密钥为',key)
        path = os.path.abspath(file).replace('girf_sync.db','girf_sync_dec.db')
        sqlite3 = self.get_relative_path('../lib/sqlite3.exe')
        if os.path.exists(sqlite3) and os.path.isfile(sqlite3):
            sqlite3 = os.path.abspath(sqlite3)
        else:
            print_red('[失败]---->原因[缺少依赖<sqlite3.exe>！]')
            return -1
        with open(file,'rb') as fr:
            data = fr.read()
            # 获取页大小，读写、缓存等信息，原数据明文，非加密，在第3个8字节数据中给出
            prop = data[16:24].hex()
            cipher = AES.new(key.encode('utf8'), AES.MODE_ECB)
            dec_data = cipher.decrypt(data)
            # 在第93-96字节存在第25-28字节的备份数据，拿来即可
            print('[提示]---->正在修补数据库文件头')
            magic = dec_data[92:96].hex()+size
            # 拼出新的32字节文件头
            true_head = head+prop+magic
            dec_data = bytes.fromhex(true_head)+dec_data[32:]
            with open(path,'wb') as fw:
                fw.write(dec_data)
            work_dir = os.getcwd()
            floder = os.path.realpath(path)
            floder = '\\'.join(floder.split('\\')[:-1])
            os.chdir(floder)
            result = subprocess.run([sqlite3,path,".tables"],capture_output=True, text=True)
            err = result.stderr
            if err != '':
                print_red('[失败]---->girf_sync.db解密失败，未知原因，可能是更换了密钥！')
            else:
                print_green('[成功]---->girf_sync.db解密成功，解密后数据库在源文件同目录下girf_sync_dec.db')
            os.chdir(work_dir)

    def decrypt_dingtalk(self,device :str, path :str):
        """
        @device             设备信息，通常在shared_prefs/com.alibaba.android.rimet_preferences.xml中带有数据库的字段的值中出现，如HUAWEI P40/armeabi-v7a/P40/qcom/HUAWEIp40
        """
        if not os.path.exists(path):
            return '文件不存在！'
        cpu = ['armeabi','armeabi-v7a','arm64-v8a','x86','x86_64','mips','mips64']
        info = device.split('/')
        cpu.remove(info[1])
        cpu.insert(0,info[1])
        print_yellow('<钉钉数据库解密>')
        bname = os.path.basename(path)
        wpath = os.path.abspath(path).replace(bname,bname.replace('.db','_dec.db'))
        for v in cpu:
            t = info
            t[1] = v
            tmp_key = '/'.join(t)
            md5 = hashlib.md5()
            md5.update(tmp_key.encode('utf8'))
            key = md5.hexdigest()[:16]
            print_yellow_key(f'[提示]---->正在解密钉钉数据库{bname}，解密密钥为',key)
            with open(path,'rb') as fr:
                data = fr.read()
                cipher = AES.new(key.encode('utf8'), AES.MODE_ECB)
                dec_data = cipher.decrypt(data)
                head = dec_data[:16]
                with open(wpath,'wb') as fw:
                    fw.write(dec_data)
                if head == b'SQLite format 3\0':
                    print_green(f'[成功]---->{bname}解密成功，解密后数据库在源文件同目录下{bname.replace(".db","_dec.db")}')
                    return
                else:
                    print_red(f'[失败]---->{bname}解密失败，解密密钥不正确！')
        print_red(f'[失败]---->{bname}解密失败，未能找到密钥，请确认输入的内容！')