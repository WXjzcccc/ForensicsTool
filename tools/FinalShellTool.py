import jpype
from base64 import b64decode
import os
import json
from .PrintTool import print_red,print_dict,print_yellow

def get_relative_path(relative_path):
    """获取配置文件的绝对路径"""
    base_path = os.path.dirname(os.path.abspath(__file__))
    return os.path.join(base_path, relative_path)

def decrypt(dics):
    # 修改、编译自https://github.com/jas502n/FinalShellDecodePass
    clazz = get_relative_path('../lib/FinalShellDecodePass.class')
    print('[提示]---->正在解密连接密码')
    if os.path.exists(clazz) and os.path.isfile(clazz):
        pass
    else:
        print_red('[失败]---->原因[缺少依赖<FinalShellDecodePass.class>！]')
        return -1
    jpype.startJVM(jpype.getDefaultJVMPath())
    jpype.addClassPath(os.path.abspath(get_relative_path('../lib')))
    FinalShellDecodePass = jpype.JClass('FinalShellDecodePass')
    finalShellDecodePass = FinalShellDecodePass()
    for dic in dics:
        encPwd = dic['密码']
        b64c = b64decode(encPwd)
        randomKey = finalShellDecodePass.ranDomKey(b64c[:8])
        m = finalShellDecodePass.desDecode(b64c[8:], randomKey)
        dic.update({'密码':m})
    return dics

def analyzeFinalShell(folder :str):
    conns = []
    print('[提示]---->正在提取FinalShell连接信息')
    for root,dirs,files in os.walk(folder):
        for file in files:
            if file.endswith('_connect_config.json'):
                path = os.path.join(root,file)
                conn = {}
                try:
                    with open(path,'r') as fr:
                        con_json = json.loads(fr.read())
                        conn.update({'连接名':con_json['name']})
                        conn.update({'地址':con_json['host']})
                        conn.update({'端口':con_json['port']})
                        conn.update({'用户名':con_json['user_name']})
                        conn.update({'密码':con_json['password']})
                except:
                    with open(path,'r',encoding='utf8') as fr:
                        con_json = json.loads(fr.read())
                        conn.update({'连接名':con_json['name']})
                        conn.update({'地址':con_json['host']})
                        conn.update({'端口':con_json['port']})
                        conn.update({'用户名':con_json['user_name']})
                        conn.update({'密码':con_json['password']})
                conns.append(conn)
    conns = decrypt(conns)
    print_dict(conns,conns[0].keys(),'FinalShell连接信息')
    print_yellow('[提示]---->如果内容显示不全，请将终端全屏后再次执行')