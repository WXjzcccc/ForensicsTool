import datetime
import json
import os
from base64 import b64decode

from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad

from .PrintTool import *


class UToolsCipher:
    _AES_KEY = ""
    _AES_IV = "UTOOLS0123456789"
    _SUPER_KEY = "e87e4c85ae6667a8c2da6ad72617ea6d029b4973e23e5c20fb745ccc74c9a9d6"
    clipboard_data = ""
    dec_path = ""

    def __init__(self, clipboard_data: str, key: str):
        self._AES_KEY = key
        self.clipboard_data = clipboard_data
        self.dec_path = self.clipboard_data + '_deccypt.txt'

    def decrypt(self, encrypted_data: str) -> str:
        cipher = AES.new(self._AES_KEY.encode("utf-8"), AES.MODE_CBC, self._AES_IV.encode("utf-8"))
        decrypted = cipher.decrypt(bytes.fromhex(encrypted_data))
        return unpad(decrypted, AES.block_size).decode("utf8")

    def write_decrypted_msg(self, msg: str, path: str) -> None:
        with open(path, 'w', encoding='utf8') as fw:
            for v in msg:
                fw.write(v + '\n')

    def analyze(self) -> list:
        dec_data = []
        for root, dirs, files in os.walk(self.clipboard_data):
            for file in files:
                if file == 'data':
                    _data = os.path.join(root, file)
                    with open(_data, 'r', encoding='utf8') as fr:
                        data = fr.readlines()
                        for v in data:
                            dec_data.append(self.decrypt(v))
                        self.write_decrypted_msg(dec_data, self.dec_path)
        return dec_data

    def my_timestamp(self, timestamp: str) -> str:
        delta = datetime.timedelta(milliseconds=int(timestamp))
        initial_time = datetime.datetime(1970, 1, 1)
        target_time = initial_time + delta
        formatted_date = target_time.strftime('%Y-%m-%d %H:%M:%S.%f')
        return formatted_date

    def parse(self):
        text_dict = []
        file_dict = []
        image_dict = []
        id = 0
        json_lst = self.analyze()
        for data in json_lst:
            json_data = json.loads(data)
            time = self.my_timestamp(json_data['timestamp'])
            hash = json_data['hash']
            id += 1
            if json_data["type"] == 'text':
                value = json_data['value']
                if len(value) > 100:
                    value = value[:100] + '...'
                text_dict.append({
                    "ID": id,
                    "内容": value,
                    "时间": time,
                    "哈希值": hash
                })
            elif json_data["type"] == 'files':
                values = json_data['value']
                name = ''
                ftype = ''
                path = ''
                for v in values:
                    name += v['name'] + '\n'
                    ftype += '文件' if v['isFile'] == True else '目录'
                    ftype += '\n'
                    path += v['path'] + '\n'
                file_dict.append({
                    "ID": id,
                    "类型": ftype[:-1],
                    "名称": name[:-1],
                    "路径": path[:-1],
                    "时间": time,
                    "哈希值": hash
                })
            elif json_data["type"] == 'image':
                size = json_data['size']
                image_dict.append({
                    "ID": id,
                    "文件名": hash,
                    "文件大小": size,
                    "时间": time
                })
        dic = {}
        if text_dict != []:
            print_dict(text_dict, text_dict[0].keys(), title='文本内容')
            dic.update({'文本内容': text_dict})
        if file_dict != []:
            print_dict(file_dict, file_dict[0].keys(), title='文件内容')
            dic.update({'文件内容': file_dict})
        if image_dict != []:
            print_dict(image_dict, image_dict[0].keys(), title='图片内容')
            dic.update({'图片内容': image_dict})
        print_yellow('[提示]---->时间之间有空白表示一条记录有多个文件！')
        print_yellow(
            '[提示]---->受数据长度影响，部分内容无法完整展示，解密后的数据已经保存到同目录下的_deccypt.txt结尾的文件中！')
        return dic

    def decrypt_super(self, encrypted_data: str, IV: str) -> str:
        cipher = AES.new(bytes.fromhex(self._SUPER_KEY), AES.MODE_CBC, bytes.fromhex(IV))
        decrypted = cipher.decrypt(bytes.fromhex(encrypted_data))
        return unpad(decrypted, AES.block_size).decode("utf8")

    def parseImages(self, path):
        _dir = self.clipboard_data + '/wxjzc_images'
        if not os.path.exists(_dir):
            if not os.path.isdir(_dir):
                os.mkdir(_dir)
        if not os.path.dirname(path).endswith('wxjzc_images'):
            with open(path, 'r', encoding='utf8') as fr:
                data = fr.read()
                img_data = data.split(',')[1]
                ext = data[11:].split(';')[0]
                fname = os.path.basename(path)
                with open(os.path.join(_dir, fname + '.' + ext), 'wb') as fw:
                    fw.write(b64decode(img_data))

    def parseSuper(self):
        text_dict = []
        file_dict = []
        image_dict = []
        id = 0
        for root, dirs, files in os.walk(self.clipboard_data):
            for file in files:
                if file.startswith('image_'):
                    self.parseImages(os.path.join(root, file))
                if file == 'data':
                    path = os.path.join(root, file)
                    with open(path, 'r', encoding='utf8') as fr:
                        with open(os.path.join(root, file + '_decrypt.txt'), 'w', encoding='utf8') as fw:
                            data = fr.readlines()
                            for v in data:
                                id += 1
                                line = v.split(' ')
                                time = self.my_timestamp(line[1])
                                hash = line[2].strip()
                                if line[0] == 'image':
                                    image_dict.append({
                                        "ID": id,
                                        "文件名": hash,
                                        "时间": time
                                    })
                                    fw.write(str({"类型": '图片', "时间": time, "哈希值": hash}) + '\n')
                                elif line[0] == 'file':
                                    values = json.loads(b64decode(line[3]).decode('utf-8'))
                                    name = ''
                                    ftype = ''
                                    path = ''
                                    for v in values:
                                        name += v['name'] + '\n'
                                        ftype += '文件' if v['isFile'] == True else '目录'
                                        ftype += '\n'
                                        path += v['path'] + '\n'
                                    file_dict.append({
                                        "ID": id,
                                        "类型": ftype[:-1],
                                        "名称": name[:-1],
                                        "路径": path[:-1],
                                        "时间": time,
                                        "哈希值": hash
                                    })
                                    fw.write(str({
                                        "类型": ftype[:-1],
                                        "名称": name[:-1],
                                        "路径": path[:-1],
                                        "时间": time,
                                        "哈希值": hash
                                    }) + '\n')
                                elif line[0] == 'text':
                                    value = line[3].strip()
                                    sp = value.split(':')
                                    iv = sp[0]
                                    enc_data = sp[1]
                                    msg = b64decode(self.decrypt_super(enc_data, iv)).decode('utf-8')
                                    text_dict.append({
                                        "ID": id,
                                        "内容": msg if len(msg) <= 100 else msg[:97] + '...',
                                        "时间": time,
                                        "哈希值": hash
                                    })
                                    fw.write(str({
                                        "内容": b64decode(self.decrypt_super(enc_data, iv)).decode('utf-8'),
                                        "时间": time,
                                        "哈希值": hash
                                    }) + '\n')
        dic = {}
        if text_dict != []:
            try:
                print_dict(text_dict, text_dict[0].keys(), title='文本内容')
                dic.update({'文本内容': text_dict})
            except:
                print_red('[错误]---->部分内容无法打印，请直接查看文件！')
        if file_dict != []:
            print_dict(file_dict, file_dict[0].keys(), title='文件内容')
            dic.update({'文件内容': file_dict})
        if image_dict != []:
            print_dict(image_dict, image_dict[0].keys(), title='图片内容')
            dic.update({'图片内容': image_dict})
        print_yellow('[提示]---->时间之间有空白表示一条记录有多个文件！')
        print_yellow(
            '[提示]---->受数据长度影响，部分内容无法完整展示，解密后的数据已经保存到同目录下的_deccypt.txt结尾的文件中！')
        return dic


def analyzeUTools(path: str, key: str):
    u = UToolsCipher(path, key)
    if key != 'super':
        print('[提示]---->正在提取uTools剪贴板数据')
        return u.parse()
    else:
        print('[提示]---->正在提取uTools超级剪贴板数据')
        return u.parseSuper()
