import base64

from Crypto.Cipher import AES
from Crypto.Hash import SHA512
from Registry import Registry

from .PrintTool import print_red, print_dict, print_yellow


class MobaXtermCipher:
    '''
    @author:        @clemthi on github 
    @repository:  https://github.com/clemthi/mobaxterm-decrypt-pass
    '''

    def __init__(self, master_password: bytes):
        self._key = SHA512.new(master_password).digest()[0:32]
        self._cipher = None

    def _generate_cipher(self) -> None:
        init_vector = AES.new(key=self._key, mode=AES.MODE_ECB).encrypt(b"\x00" * AES.block_size)
        self._cipher = AES.new(key=self._key, iv=init_vector, mode=AES.MODE_CFB, segment_size=8)

    def decrypt_password(self, cipher_text: str) -> bytes:
        self._generate_cipher()
        return self._cipher.decrypt(base64.b64decode(cipher_text)).decode('utf8')


def getEncFromIni(file: str) -> dict:
    """
    从配置文件中提取密码和凭据
    """
    print('[提示]---->从配置文件获取数据')
    with open(file, 'r', encoding='gb2312') as fr:
        data = fr.readlines()
        passwords = []
        credentials = []
        fi = 0
        ff = False
        ci = 0
        cf = False
        end = []
        ei = 0
        for v in data:
            if v == '\n':
                end.append(ei)
            if v == '[Passwords]\n':
                fi = ei
            if v == '[Credentials]\n':
                ci = ei
            ei += 1
        for i in end:
            if i > fi and not ff:
                tmp = data[fi + 1:i]
                for v in tmp:
                    passwords.append(v.strip())
                ff = True
            if i > ci and not cf and ci != 0:
                tmp = data[ci + 1:i]
                for v in tmp:
                    credentials.append(v.strip())
                cf = True
            if cf and ff:
                break
    return passwords, credentials


def getEncFromReg(file: str) -> dict:
    """
    从注册表中提取密码和凭据
    """
    passwords = []
    credentials = []
    print('[提示]---->从注册表文件获取数据')
    try:
        reg = Registry.Registry(file)
        pwdReg = reg.open("Software\\Mobatek\\MobaXterm\\P").values()
        for v in pwdReg:
            passwords.append(v.name() + '=' + v.value())
        creReg = reg.open("Software\\Mobatek\\MobaXterm\\C").values()
        for v in creReg:
            credentials.append(v.name() + '=' + v.value())
    except:
        print_red('[错误]---->注册表中未找到MobaXterm保存的密码和凭证信息，请从MobaXterm.ini配置文件中提取')
    return passwords, credentials


def analyzeMoba(file: str, masterkey: str, flag: int):
    passwords = []
    credentials = []
    cipher = MobaXtermCipher(masterkey.encode('utf8'))
    pwds = []
    cres = []
    if flag == 0:
        pwds, cres = getEncFromReg(file)
    elif flag == 1:
        pwds, cres = getEncFromIni(file)
    for pwd in pwds:
        password = {}
        tmp = pwd.split(':')
        if pwd.__contains__(':'):
            password.update({'协议与端口': tmp[0]})
        else:
            password.update({'协议与端口': ''})
        t = tmp[-1]
        index = t.index('=')
        connection = t[:index]
        encPwd = t[index + 1:]
        password.update({'连接账户与地址': connection})
        password.update({'保存的密码': cipher.decrypt_password(encPwd)})
        passwords.append(password)
    for cre in cres:
        credential = {}
        index = cre.index('=')
        name = cre[:index]
        credential.update({'凭据名': name})
        t = cre[index + 1:].split(':')
        credential.update({'用户名': t[0]})
        credential.update({'保存的密码': cipher.decrypt_password(t[1])})
        credentials.append(credential)
    dic = {}
    try:
        print_dict(passwords, passwords[0].keys(), )
        dic.update({'MobaXterm连接信息': passwords})
        print_dict(credentials, credentials[0].keys(), 'MovaXterm凭据信息')
        dic.update({'MovaXterm凭据信息': credentials})
    except:
        pass
    print_yellow('[提示]---->如果内容显示不全，请将终端全屏后再次执行')
    return dic
