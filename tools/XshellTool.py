import base64
import configparser
import os

from Crypto.Cipher import ARC4
from Crypto.Hash import SHA256

from .PrintTool import print_dict, print_yellow


# from https://github.com/dzxs/Xdecrypt
def decrypt_string(sid: str, encPwd: str):
    v1 = base64.b64decode(encPwd)
    v3 = ARC4.new(SHA256.new(sid.encode('ascii')).digest()).decrypt(v1[:len(v1) - 0x20])
    if SHA256.new(v3).digest() == v1[-32:]:
        return v3.decode('ascii')
    else:
        return None


def decrypt_string_new(sid: str, encPwd: str):
    # 7.1+
    v1 = base64.b64decode(encPwd)
    index = sid.index('S-1-5')
    username = sid[:index]
    ssid = sid[index:]
    key = (username[::-1] + ssid)[::-1]
    v3 = ARC4.new(SHA256.new(key.encode('ascii')).digest()).decrypt(v1[:len(v1) - 0x20])
    if SHA256.new(v3).digest() == v1[-32:]:
        return v3.decode('ascii')
    else:
        return None


def preDeal(file: str):
    print('[提示]---->正在处理配置文件')
    with open(file, 'rb') as fr:
        data = fr.read().replace(b'\x00', b'').replace(b'\xff', b'').replace(b'\xfe', b'')
        with open(file + '.deal', 'wb') as fw:
            fw.write(data)
    return file + '.deal'


def delDeal(file: str):
    print('[提示]---->正在删除临时文件')
    os.remove(file)


def analyzeXshell(foler: str, sid: str):
    conns = []
    print('[提示]---->正在提取XShell/XFtp连接信息')
    for root, dirs, files in os.walk(foler):
        for file in files:
            if file.endswith(r'.xsh'):
                conn = {}
                cf = configparser.ConfigParser()
                origin_file = os.path.join(root, file)
                deal_file = preDeal(origin_file)
                cf.read(deal_file, encoding='ansi')
                conn.update({'连接名': os.path.basename(file).replace(r'.xsh', '')})
                version = cf.get('SessionInfo', 'Version')
                conn.update({'地址': cf.get('CONNECTION', 'Host')})
                conn.update({'端口': cf.get('CONNECTION', 'Port')})
                if version != '7.0' and version.startswith('7.'):
                    conn.update({'密码': decrypt_string_new(sid, cf.get('CONNECTION:AUTHENTICATION', 'Password'))})
                else:
                    conn.update({'密码': decrypt_string(sid, cf.get('CONNECTION:AUTHENTICATION', 'Password'))})
                conn.update({'描述': cf.get('CONNECTION', 'Description')})
                delDeal(deal_file)
                conns.append(conn)
            if file.endswith(r'.xfp'):
                conn = {}
                cf = configparser.ConfigParser()
                origin_file = os.path.join(root, file)
                deal_file = preDeal(origin_file)
                cf.read(deal_file, encoding='ansi')
                conn.update({'连接名': os.path.basename(file).replace(r'.xfp', '')})
                version = cf.get('SessionInfo', 'Version')
                conn.update({'地址': cf.get('Connection', 'Host')})
                conn.update({'端口': cf.get('Connection', 'Port')})
                if version != '7.0' and version.startswith('7.'):
                    conn.update({'密码': decrypt_string_new(sid, cf.get('Connection', 'Password'))})
                else:
                    conn.update({'密码': decrypt_string(sid, cf.get('Connection', 'Password'))})
                conn.update({'描述': cf.get('Connection', 'Description')})
                delDeal(deal_file)
                conns.append(conn)
    dic = {}
    dic.update({'Xshell/Xftp连接信息': conns})
    print_dict(conns, conns[0].keys(), 'Xshell/Xftp连接信息')
    print_yellow('[提示]---->如果内容显示不全，请将终端全屏后再次执行')
    return dic
