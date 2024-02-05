import argparse,textwrap
from tools.DbPwdTool import DbPwdTool
from tools.DbTool import DbTool
from tools.NavicatTool import analyzeNavicat
from tools.PrintTool import print_red
from tools.MobaTool import analyzeMoba
from tools.DBeaverTool import analyzeDbeaer
import sys
import rich
import os

def check_arg(arg):
    if arg is None or arg == '':
        return False
    return True

# 创建 ArgumentParser 对象
parser = argparse.ArgumentParser(formatter_class=argparse.RawTextHelpFormatter,description='''
Forensics Tool''')

# 添加命令行参数
parser.add_argument('-m','--mode', type=int,help='''
指定需要运行的模式:
    [0]表示计算密钥，支持的type值为1-3
    [1]表示解密数据库，支持的type值为1、2、4-7
    [2]表示数据提取，支持的type值为8-10''')
parser.add_argument('-f', '--file', type=str,help='指定需要解密的数据库')
parser.add_argument('-t', '--type', type=int,help='''
指定需要处理的内容:
    [1]微信的EnMicroMsg.db
    [2]微信的FTS5IndexMicroMsg_encrypt.db
    [3]野火IM系应用的data
    [4]高德的girf_sync.db
    [5]钉钉的数据库
    [6]SQLCipher4加密的数据库
    [7]SQLCipher3加密的数据库
    [8]Navicat连接信息提取，需指定-f为目标用户的注册表文件"NTUSER.DAT"
    [9]MobaXterm连接信息解密，可以指定MobaXterm.ini配置文件或用户注册表文件"NTUSER.DAT"，解密需要给出主密码
    [10]Dbeaver连接信息解密，指定-f为目标文件data-sources.json和credentials-config.json的父目录''')
parser.add_argument('-p', '--password', type=str, help='解密的密码，处理钉钉和高德时不适用')
parser.add_argument('--uin', type=str, help='微信用户的uin，可能是负值，在shared_prefs/auth_info_key_prefs.xml文件中_auth_uin的值')
parser.add_argument('--imei', type=str, help='微信获取到的IMEI或MEID，在shared_prefs/DENGTA_META.xml文件中IMEI_DENGTA的值，在高版本中通常是1234567890ABCDEF，可以为空')
parser.add_argument('--wxid', type=str, help='数据库所属的wxid，一般情况下在解密EnMicroMsg.db的时候会一并提取，若无需要，请从shared_prefs/com.tencent.mm_preferences.xml中提取login_weixin_username的值')
parser.add_argument('--token', type=str, help='野火IM系应用的用户token，shared_prefs/config.xml的token的值')
parser.add_argument('--device', type=str, help='钉钉解密需要的内容，通常在shared_prefs/com.alibaba.android.rimet_preferences.xml中带有数据库名的字段的值中出现，如HUAWEI P40/armeabi-v7a/P40/qcom/HUAWEIP40')

# 解析命令行参数
args = parser.parse_args()
rich.print("""
   ____                      _           ______          __
  / __/__  _______ ___  ___ (_)______   /_  __/__  ___  / /
 / _// _ \/ __/ -_) _ \(_-</ / __(_-<    / / / _ \/ _ \/ / 
/_/  \___/_/  \__/_//_/___/_/\__/___/   /_/  \___/\___/_/  
                                                            Author: WXjzc""")
if args.mode == 0:
    if check_arg(args.file):
        print_red('[错误]---->该模式不支持file！')
        sys.exit(-1)
    dbPwdTool = DbPwdTool()
    uin = args.uin
    imei = args.imei
    wxid = args.wxid
    token = args.token
    if args.type == 1:
        if check_arg(uin):
            if check_arg(imei):
                dbPwdTool.wechat(uin,imei)
            else:
                dbPwdTool.wechat(uin)
        else:
            print_red('[错误]---->请给出uin！')
    elif args.type == 2:
        if check_arg(uin) and check_arg(wxid):
            if check_arg(imei):
                dbPwdTool.wechat_index(uin,wxid,imei)
            else:
                dbPwdTool.wechat_index(uin,wxid)
        else:
            print_red('[错误]---->请给出uin和wxid！')
    elif args.type == 3:
        if check_arg(token):
            dbPwdTool.wildfire(token)
        else:
            print_red('[错误]---->请给出token！')
    else:
        print_red('[错误]---->不支持的type值！')
        sys.exit(-1)
elif args.mode == 1:
    dbTool = DbTool()
    device = args.device
    key = args.password
    file = args.file
    if not check_arg(file):
        print_red('[错误]---->请给出file参数！')
        sys.exit(-1)
    if args.type == 1:
        if check_arg(key):
            dbTool.decrypt_EnMicroMsg(key,file)
        else:
            print_red('[错误]---->请给出解密密码！')
    elif args.type == 2:
        if check_arg(key):
            dbTool.decrypt_FTS5IndexMicroMsg(key,file)
        else:
            print_red('[错误]---->请给出解密密码！')
    elif args.type == 4:
        dbTool.decrypt_amap(file)
    elif args.type == 5:
        if check_arg(device) and len(device.split('/')) == 5:
            dbTool.decrypt_dingtalk(device,file)
        else:
            print_red('[错误]---->请给出正确格式的device！')
    elif args.type == 6:
        if check_arg(key):
            dbTool.decrypt_SQLCipher4_default(key,file)
        else:
            print_red('[错误]---->请给出解密密码！')
    elif args.type == 7:
        if check_arg(key):
            dbTool.decrypt_SQLCipher3_default(key,file)
        else:
            print_red('[错误]---->请给出解密密码！')
    else:
        print_red('[错误]---->不支持的type值！')
elif args.mode == 2:
    file = args.file
    password = args.password
    if args.type == 8:
        f = open(file,'rb')
        head = f.read(4)
        f.close()
        if head == b'regf' and os.path.basename(file) == 'NTUSER.DAT':
            analyzeNavicat(file)
        else:
            print_red('[错误]---->不是NTUSER注册表文件！')
    elif args.type == 9:
        if check_arg(password):
            f = open(file,'rb')
            head = f.read(4)
            f.close()
            try:
                if head == b'regf' and os.path.basename(file) == 'NTUSER.DAT':
                    analyzeMoba(file,password,0)
                elif os.path.basename(file).startswith('MobaXterm') and os.path.basename(file).endswith('.ini'):
                    analyzeMoba(file,password,1)
                else:
                    print_red('[错误]---->不是待分析的文件，请给出NTUSER注册表文件或MobaXterm.ini配置文件！')
            except:
                print_red('[错误]---->masterkey不正确！')
        else:
            print_red('[错误]---->必须给出-p参数，MobaXterm的masterkey！')
    elif args.type == 10:
        if os.path.isdir(file):
            dirs = os.listdir(file)
            if 'credentials-config.json' in dirs and 'data-sources.json' in dirs:
                analyzeDbeaer(file)
            else:
                print_red('[错误]---->目录中找不到credentials-config.json和data-sources.json文件')
        else:
            print_red('[错误]---->请指定目录，而不是文件')
    else:
        print_red('[错误]---->不支持的type值！')
