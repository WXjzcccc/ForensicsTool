import binascii
import codecs
import hashlib

from Crypto.Cipher import Blowfish, AES
from Registry import Registry

from .PrintTool import print_yellow, print_dict

'''
用于11版本解密
'''


class Navicat11Cipher:
    DEFAULT_USER_KEY = "3DC5CA39"

    def __init__(self, user_key=None):
        if user_key is None:
            user_key = self.DEFAULT_USER_KEY
        self._key = self._init_key(user_key)
        self._encryptor = self._init_cipher(Blowfish.MODE_ECB)
        self._decryptor = self._init_cipher(Blowfish.MODE_ECB)
        self._iv = self._init_iv()

    def _init_key(self, user_key):
        sha1 = hashlib.sha1()
        sha1.update(user_key.encode("utf-8"))
        return sha1.digest()

    def _init_cipher(self, mode):
        return Blowfish.new(self._key, mode)

    def _init_iv(self):
        init_vec = codecs.decode("FFFFFFFFFFFFFFFF", "hex_codec")
        return self._encryptor.encrypt(init_vec)

    def _xor_bytes(self, a, b, l=None):
        if l != None:
            return bytes([x ^ y for x, y in zip(a, b)][:l])
        return bytes([x ^ y for x, y in zip(a, b)])

    def encrypt(self, in_data):
        out_data = bytearray(len(in_data))
        cv = bytearray(self._iv)
        blocks_len = len(in_data) // 8
        left_len = len(in_data) % 8
        for i in range(0, blocks_len):
            temp = in_data[i * 8:i * 8 + Blowfish.block_size]
            temp = self._xor_bytes(temp, cv)
            temp = self._encryptor.encrypt(temp)
            cv = self._xor_bytes(cv, temp)
            out_data[i * 8:i * 8 + Blowfish.block_size] = temp
        if left_len != 0:
            cv = self._encryptor.encrypt(cv)
            temp = in_data[blocks_len * 8:blocks_len * 8 + left_len]
            temp = self._xor_bytes(temp, cv)
            out_data[blocks_len * 8:] = temp

        return bytes(out_data)

    def decrypt(self, in_data):
        out_data = bytearray(len(in_data))
        cv = bytearray(self._iv)

        blocks_len = len(in_data) // 8
        left_len = len(in_data) % 8

        for i in range(0, blocks_len):
            temp = in_data[i * 8:i * 8 + Blowfish.block_size]
            temp = self._decryptor.decrypt(temp)
            temp = self._xor_bytes(temp, cv)
            out_data[i * 8:i * 8 + Blowfish.block_size] = temp
            cvb = bytearray(cv)
            for j in range(len(cv)):
                cvb[j] = cv[j] ^ in_data[i * 8 + j]
            cv = bytes(cvb)
        if left_len != 0:
            cv = self._encryptor.encrypt(cv)
            temp = in_data[blocks_len * 8:blocks_len * 8 + left_len]
            temp = self._xor_bytes(temp, cv)
            out_data[blocks_len * 8:] = temp
        return out_data

    def encrypt_string(self, input_string):
        in_data = input_string.encode("utf-8")
        out_data = self.encrypt(in_data)
        return codecs.encode(out_data, "hex_codec").decode("utf-8")

    def decrypt_string(self, hex_string):
        in_data = codecs.decode(hex_string, "hex_codec")
        out_data = self.decrypt(in_data)
        return out_data.decode("utf-8")


'''
用于12版本之后解密
'''


class Navicat12Cipher:
    _AesKey = b"libcckeylibcckey"
    _AesIV = b"libcciv libcciv "

    def encrypt_string(self, plaintext):
        try:
            cipher = AES.new(self._AesKey, AES.MODE_CBC, self._AesIV)
            ciphertext = cipher.encrypt(self._pad(plaintext.encode('utf-8')))
            return binascii.hexlify(ciphertext).decode()
        except Exception as e:
            return ""

    def decrypt_string(self, ciphertext):
        try:
            cipher = AES.new(self._AesKey, AES.MODE_CBC, self._AesIV)
            decrypted = cipher.decrypt(binascii.unhexlify(ciphertext))
            return self._unpad(decrypted).decode()
        except Exception as e:
            return ""

    def _pad(self, s):
        block_size = AES.block_size
        padding = block_size - len(s) % block_size
        return s + (padding * chr(padding)).encode()

    def _unpad(self, s):
        return s[:-s[-1]]


'''
解密函数封装
'''


def decryptPwd(pwd: str) -> str:
    password = ''
    try:
        dec12 = Navicat12Cipher()
        password = dec12.decrypt_string(pwd)
        if password == '':
            dec11 = Navicat11Cipher()
            password = dec11.decrypt_string(pwd)
    except:
        return ''
    return password


'''
读取mysql数据
'''


def get_mysql_info(server_keys: list) -> dict:
    info = {}
    print('[提示]---->正在读取MySQL数据')
    for v in server_keys:
        basic_info = {}
        try:
            basic_info.update({'连接名': v.name()})
            basic_info.update({'地址': v.value('Host').value()})
            basic_info.update({'端口': v.value('Port').value()})
            basic_info.update({'用户名': v.value('UserName').value()})
            basic_info.update({'密码': decryptPwd(v.value('Pwd').value())})
            basic_info.update({'缓存路径': v.value('QuerySavePath').value()})
        except:
            pass
        finally:
            info.update({v.name(): basic_info})
    return info


'''
读取mariadb数据
'''


def get_mariadb_info(server_keys: list) -> dict:
    info = {}
    print('[提示]---->正在读取MariaDB数据')
    for v in server_keys:
        basic_info = {}
        try:
            basic_info.update({'连接名': v.name()})
            basic_info.update({'地址': v.value('Host').value()})
            basic_info.update({'端口': v.value('Port').value()})
            basic_info.update({'用户名': v.value('UserName').value()})
            basic_info.update({'密码': decryptPwd(v.value('Pwd').value())})
            basic_info.update({'缓存路径': v.value('QuerySavePath').value()})
        except:
            pass
        finally:
            info.update({v.name(): basic_info})
    return info


'''
读取mongodb数据
'''


def get_mongodb_info(server_keys: list) -> dict:
    info = {}
    for v in server_keys:
        basic_info = {}
        try:
            basic_info.update({'连接名': v.name()})
            basic_info.update({'地址': v.value('Host').value()})
            basic_info.update({'端口': v.value('Port').value()})
            basic_info.update({'用户名': v.value('UserName').value()})
            basic_info.update({'密码': decryptPwd(v.value('Pwd').value())})
            basic_info.update({'缓存路径': v.value('QuerySavePath').value()})
        except:
            pass
        finally:
            info.update({v.name(): basic_info})
    return info


'''
读取sqlserver数据
'''


def get_sqlserver_info(server_keys: list) -> dict:
    info = {}
    print('[提示]---->正在读取SQLServer数据')
    for v in server_keys:
        basic_info = {}
        try:
            basic_info.update({'连接名': v.name()})
            basic_info.update({'地址': v.value('Host').value()})
            basic_info.update({'端口': v.value('Port').value()})
            basic_info.update({'用户名': v.value('UserName').value()})
            basic_info.update({'密码': decryptPwd(v.value('Pwd').value())})
            basic_info.update({'默认数据库': v.value('InitialDatabase').value()})
            basic_info.update({'缓存路径': v.value('QuerySavePath').value()})
            basic_info.update({'认证模式': v.value('MSSQLAuthenMode').value()[3:]})
        except:
            pass
        finally:
            info.update({v.name(): basic_info})
    return info


'''
读取oracle数据
'''


def get_oracle_info(server_keys: list) -> dict:
    info = {}
    print('[提示]---->正在读取Oracle数据')
    for v in server_keys:
        basic_info = {}
        try:
            basic_info.update({'连接名': v.name()})
            basic_info.update({'地址': v.value('Host').value()})
            basic_info.update({'端口': v.value('Port').value()})
            basic_info.update({'用户名': v.value('UserName').value()})
            basic_info.update({'密码': decryptPwd(v.value('Pwd').value())})
            basic_info.update({'缓存路径': v.value('QuerySavePath').value()})
        except:
            pass
        finally:
            info.update({v.name(): basic_info})
    return info


'''
读取pgsql数据
'''


def get_postgres_info(server_keys: list) -> dict:
    info = {}
    print('[提示]---->正在读取PostgreSQL数据')
    for v in server_keys:
        basic_info = {}
        try:
            basic_info.update({'连接名': v.name()})
            basic_info.update({'地址': v.value('Host').value()})
            basic_info.update({'端口': v.value('Port').value()})
            basic_info.update({'用户名': v.value('UserName').value()})
            basic_info.update({'密码': decryptPwd(v.value('Pwd').value())})
            basic_info.update({'缓存路径': v.value('QuerySavePath').value()})
        except:
            pass
        finally:
            info.update({v.name(): basic_info})
    return info


'''
读取sqlite数据
'''


def get_sqlite_info(server_keys: list) -> dict:
    info = {}
    print('[提示]---->正在读取SQLite数据')
    for v in server_keys:
        basic_info = {}
        try:
            basic_info.update({'连接名': v.name()})
            basic_info.update({'文件名': v.value('DatabaseFileName').value()})
            basic_info.update({'是否加密': '否' if v.value('SQLiteEncrypted').value() == 0 else '是'})
            basic_info.update({'是否保存密码': '否' if v.value('EncryptionSavePassword').value() == 'false' else '是'})
            basic_info.update({'密码': decryptPwd(v.value('EncryptionPassword').value())})
            basic_info.update({'缓存路径': v.value('QuerySavePath').value()})
        except:
            pass
        finally:
            info.update({v.name(): basic_info})
    return info


def get_navicat_connections(reg: str) -> dict:
    connections = {}
    database = []
    try:
        reg = Registry.Registry(reg)
        root_key = reg.open("Software\\PremiumSoft")
        for v in root_key.subkeys():
            if v.name() not in ['Data', 'NavicatPremium']:
                database.append(v.name())
        for db in database:
            try:
                db_key = reg.open(f"Software\\PremiumSoft\\{db}\\Servers")
                server_list = db_key.subkeys()
                connections.update({db: server_list})
            except:
                continue
    except:
        pass
    return connections


def analyzeNavicat(reg: str):
    print_yellow('<Navicat连接信息提取>')
    connections = get_navicat_connections(reg)
    dic = {}
    try:
        mysql = get_mysql_info(connections['Navicat'])
        count_mysql = 0
        records = []
        for v in mysql:
            record = {}
            for k in mysql[v]:
                record[k] = str(mysql[v][k])
            count_mysql += 1
            records.append(record)
        print_dict(records, records[0].keys(), title='MySQL连接信息')
        dic.update({"MySQL连接信息": records})
    except:
        pass
    try:
        mariadb = get_mariadb_info(connections['NavicatMARIADB'])
        records = []
        count_mariadb = 0
        for v in mariadb:
            record = {}
            for k in mariadb[v]:
                record[k] = str(mariadb[v][k])
            records.append(record)
            count_mariadb += 1
        print_dict(records, records[0].keys(), title='MariaDB连接信息')
        dic.update({"MariaDB连接信息": records})
    except:
        pass

    try:
        mongodb = get_mongodb_info(connections['NavicatMongoDB'])
        records = []
        count_mongodb = 0
        for v in mongodb:
            record = {}
            for k in mongodb[v]:
                record[k] = str(mongodb[v][k])
            records.append(record)
            count_mongodb += 1
        print_dict(records, records[0].keys(), title='MongoDB连接信息')
        dic.update({"MongoDB连接信息": records})
    except:
        pass

    try:
        sqlserver = get_sqlserver_info(connections['NavicatMSSQL'])
        records = []
        count_sqlserver = 0
        for v in sqlserver:
            record = {}
            for k in sqlserver[v]:
                record[k] = str(sqlserver[v][k])
            records.append(record)
            count_sqlserver += 1
        print_dict(records, records[0].keys(), title='SQLServer连接信息')
        dic.update({"SQLServer连接信息": records})
    except:
        pass

    try:
        oracle = get_oracle_info(connections['NavicatOra'])
        records = []
        count_oracle = 0
        for v in oracle:
            record = {}
            for k in oracle[v]:
                record[k] = str(oracle[v][k])
            records.append(record)
            count_oracle += 1
        print_dict(records, records[0].keys(), title='Oracle连接信息')
        dic.update({"Oracle连接信息": records})
    except:
        pass

    try:
        postgres = get_postgres_info(connections['NavicatPG'])
        records = []
        count_postgres = 0
        for v in postgres:
            record = {}
            for k in postgres[v]:
                record[k] = str(postgres[v][k])
            records.append(record)
            count_postgres += 1
        print_dict(records, records[0].keys(), title='PostgreSQL连接信息')
        dic.update({"PostgreSQL连接信息": records})
    except:
        pass

    try:
        sqlite = get_sqlite_info(connections['NavicatSQLITE'])
        records = []
        count_sqlite = 0
        for v in sqlite:
            record = {}
            for k in sqlite[v]:
                record[k] = str(sqlite[v][k])
            records.append(record)
            count_sqlite += 1
        print_dict(records, records[0].keys(), title='Sqlite连接信息')
        dic.update({"Sqlite连接信息": records})
    except:
        pass
    print_yellow('[提示]---->如果内容显示不全，请将终端全屏后再次执行')
    return dic
