from Crypto.Cipher import AES
from Crypto.Util.Padding import unpad
import os
import json
import datetime
from .PrintTool import *

class UToolsCipher:
    _AES_KEY = ""
    _AES_IV = "UTOOLS0123456789"
    clipboard_data = ""
    dec_path = ""
    def __init__(self,clipboard_data :str, key :str):
        self._AES_KEY = key
        self.clipboard_data = clipboard_data
        self.dec_path = self.clipboard_data+'_deccypt.txt'
        
    def decrypt(self,encrypted_data :str) -> str:
        cipher = AES.new(self._AES_KEY.encode("utf-8"), AES.MODE_CBC, self._AES_IV.encode("utf-8"))
        decrypted = cipher.decrypt(bytes.fromhex(encrypted_data))
        return unpad(decrypted,AES.block_size).decode("utf8")
    
    def write_decrypted_msg(self,msg :str,path :str) -> None:
        with open(path,'w',encoding='utf8') as fw:
            for v in msg:
                fw.write(v+'\n')
    
    def analyze(self) -> list:
        dec_data = []
        for root,dirs,files in os.walk(self.clipboard_data):          
            for file in files:
                if file == 'data':
                    _data = os.path.join(root,file)
                    with open(_data,'r',encoding='utf8') as fr:
                        data = fr.readlines()
                        for v in data:
                            dec_data.append(self.decrypt(v))
                        self.write_decrypted_msg(dec_data,self.dec_path)
        return dec_data
    
    def my_timestamp(self,timestamp :str) -> str:
        delta = datetime.timedelta(milliseconds=int(timestamp))
        initial_time = datetime.datetime(1970, 1, 1)
        target_time = initial_time + delta
        formatted_date = target_time.strftime('%Y-%m-%d %H:%M:%S.%f')
        return formatted_date
        
    def parse(self) -> list:
        text_dict = []
        file_dict = []
        image_dict = []
        json_lst = self.analyze()
        for data in json_lst:
            json_data = json.loads(data)
            time = self.my_timestamp(json_data['timestamp'])
            hash = json_data['hash']
            if json_data["type"] == 'text':
                value = json_data['value']
                if len(value) > 100:
                    value = value[:100]+'...'
                text_dict.append({
                    "内容":value,
                    "时间":time,
                    "哈希值":hash
                })
            elif json_data["type"] == 'files':
                values = json_data['value']
                name = ''
                ftype = ''
                path = ''
                for v in values:
                    name += v['name']+'\n'
                    ftype += '文件' if v['isFile'] == True else '目录'
                    ftype += '\n'
                    path += v['path']+'\n'
                file_dict.append({
                    "类型":ftype[:-1],
                    "名称":name[:-1],
                    "路径":path[:-1],
                    "时间":time,
                    "哈希值":hash
                })
            elif json_data["type"] == 'image':
                size = json_data['size']
                image_dict.append({
                    "文件名":hash,
                    "文件大小":size,
                    "时间":time
                })
        print_dict(text_dict,text_dict[0].keys(),title='文本内容')
        print_dict(file_dict,file_dict[0].keys(),title='文件内容')
        print_dict(image_dict,image_dict[0].keys(),title='图片内容')
        print_yellow('时间之间有空白表示一条记录有多个文件！')
        print_yellow('受数据长度影响，部分内容无法完整展示，解密后的数据已经保存到同目录下的_deccypt.txt结尾的文件中！')
        
def analyzeUTools(path :str,key :str):
    print('[提示]---->正在提取uTools剪贴板数据')
    u = UToolsCipher(path,key)
    u.parse()