import base64
import xml.etree.ElementTree as ET
from Crypto.Cipher import AES
from .PrintTool import print_table,print_yellow,print_red

class HawkTool:
    key_bytes = None

    def __init__(self,key_text):
        self.key_bytes = base64.b64decode(key_text.encode('utf8'))

    def decrypt(self,entity:str ,cipher_text:str) -> str:
        '''
        @entity: 键名
        @cipher_text: Hawk2.xml中的密文
        '''
        cipher_bytes = base64.b64decode(cipher_text.encode('utf8'))
        entity_bytes = entity.encode('utf8')
        iv = cipher_bytes[2:14]
        tag = cipher_bytes[-16:]
        data = cipher_bytes[:2]+entity_bytes
        enc_data = cipher_bytes[14:-16]
        cipher = AES.new(self.key_bytes,AES.MODE_GCM,iv)
        cipher.update(data)
        decrypt_data = cipher.decrypt_and_verify(enc_data,tag)
        return decrypt_data.decode('utf8')
    
    def parseHawk2(self,hawk2_path:str) -> dict:
        tree = ET.parse(hawk2_path)
        root = tree.getroot()
        result = {}
        for child in root:
            if child.tag == 'string' and child.text.index('##0V@') > 0:
                cipher_text = child.text.split('##0V@')[1]
                result[child.attrib['name']] = self.decrypt(child.attrib['name'],cipher_text)
        return result

def analyzeHawk2(hawk2_path:str,key_text:str):
    try:
        hawk2 = HawkTool(key_text)
        print_yellow('[提示]---->正在解密Hawk2.xml')
        result = hawk2.parseHawk2(hawk2_path)
        print_table(result,title='Hawk2解析结果')
    except:
        print_red('[错误]---->Hawk2解析失败，请检查密钥和文件路径是否正确')