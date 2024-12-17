import json

from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad

from .PrintTool import print_dict


def decrypt(bstr):
    key = bytes.fromhex('babb4a9f774ab853c96c2d653dfe544a')
    iv = bytes.fromhex('0' * 32)
    cipher = AES.new(key, mode=AES.MODE_CBC, iv=iv)
    return unpad(cipher.decrypt(bstr), AES.block_size)


def analyzeDbeaer(folder: str):
    dics = []
    print('[提示]---->正在解密文件')
    with open(folder + '/credentials-config.json', 'rb') as fr:
        encData = fr.read()
        pwdData = decrypt(encData)[16:].decode('utf8')
        pwdJson = json.loads(pwdData)
    with open(folder + '/data-sources.json', 'r', encoding='utf-8') as fr:
        infoJson = json.loads(fr.read())["connections"]
    print('[提示]---->正在解析文件')
    for key in infoJson.keys():
        dic = {}
        info = infoJson.get(key)
        pwd = pwdJson.get(key)
        dic.update({'数据库类型': info.get("provider")})
        dic.update({'连接名': info.get("name")})
        dic.update({'地址': info.get("configuration").get("host")})
        dic.update({'端口': info.get("configuration").get("port")})
        dic.update({'数据库名': info.get("configuration").get("database")})
        dic.update({'账户': pwd.get("#connection").get("user")})
        dic.update({'密码': pwd.get("#connection").get("password")})
        dics.append(dic)
    dic = {}
    dic.update({"DBeaver连接信息": dics})
    print_dict(dics, dics[0].keys(), 'DBeaver连接信息')
    return dic
