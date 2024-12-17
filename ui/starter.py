import os

from PySide6.QtCore import Signal, QObject
from tools.DBeaverTool import analyzeDbeaer
from tools.DbPwdTool import DbPwdTool
from tools.DbTool import DbTool
from tools.FinalShellTool import analyzeFinalShell
from tools.HawkTool import analyzeHawk2
from tools.MetaMaskTool import analyzeMetaMask
from tools.MobaTool import analyzeMoba
from tools.NavicatTool import analyzeNavicat
from tools.UTools import analyzeUTools
from tools.WinTool import analyzeWin
from tools.XshellTool import analyzeXshell


def color_emit_str(signal, widget, text, color):
    signal.emit(widget, f'<span  style="color:{color};">{text}</span>')


def emit_data(signal, widget, data):
    signal.emit(widget, data)


def my_table_to_table(data: dict, head: list):
    my_dic = {}
    for key, value in data.items():
        lst = []
        lst.append(head)
        lst_data = []
        for key1, value1 in value.items():
            lst_data.append([key1, value1])
        lst.append(lst_data)
        my_dic[key] = lst
    return my_dic


def my_dict_to_table(data: dict):
    my_dic = {}
    for key, value in data.items():
        lst = []
        if value != []:
            names = list(value[0].keys())
            lst.append(names)
            lst_data = []
            for line in value:
                tmp = []
                print(line)
                for v in names:
                    tmp.append(line[v])
                lst_data.append(tmp)
            lst.append(lst_data)
            my_dic[key] = lst
    return my_dic


class Starter(QObject):
    passwd_calc_signal = Signal(object, str)
    decrypt_database_signal = Signal(object, str)
    analyze_file_signal = Signal(object, object)
    analyze_registry_signal = Signal(object, object)

    def passwd_calc(self, mission: int, uin: str = None, imei: str = None, wxid: str = None, token: str = None,
                    uid: str = None, widget: QObject = None):
        db_pwd = DbPwdTool()
        print(mission, uin, imei, wxid, token, uid)
        signal = self.passwd_calc_signal
        if mission == 1:
            if uin:
                if imei:
                    pwd = db_pwd.wechat(uin, imei)
                else:
                    pwd = db_pwd.wechat(uin)
                color_emit_str(signal, widget,
                               f'[√]成功获取到安卓微信EnMicroMsg.db数据库密钥【<span style="color:brown;">{pwd}</span>】',
                               "green")
            else:
                color_emit_str(signal, widget, "[×]请正确输入uin！", "red")
        elif mission == 2:
            if uin and wxid:
                if imei:
                    pwd = db_pwd.wechat_index(uin, wxid, imei)
                else:
                    pwd = db_pwd.wechat_index(uin, wxid)
                color_emit_str(signal, widget,
                               f'[√]成功获取到安卓微信FTS5IndexMicroMsg_encrypt.db数据库密钥【<span style="color:brown;">{pwd}</span>】',
                               "green")
            else:
                color_emit_str(signal, widget, "[×]请正确输入uin和wxid！", "red")
        elif mission == 3:
            if token:
                try:
                    pwd, flag = db_pwd.wildfire(token)
                    color_emit_str(signal, widget,
                                   f'[√]成功获取到野火IM系列应用data数据库密钥【<span style="color:brown;">{pwd}</span>】',
                                   "green")
                    if flag == 1:
                        color_emit_str(signal, widget,
                                       f'[+]请使用SQLCipher3进行解密',
                                       "blue")
                    else:
                        color_emit_str(signal, widget,
                                       f'[+]请使用SQLCipher4进行解密',
                                       "blue")
                except:
                    color_emit_str(signal, widget, "[×]解密token失败，请确认token是否输入正确！", "red")
            else:
                color_emit_str(signal, widget, "[×]请正确输入token！", "red")
        elif mission == 13:
            if uid:
                pwd = db_pwd.mostone(uid)
                color_emit_str(signal, widget,
                               f'[√]成功获取到默往msg.db数据库密钥【<span style="color:brown;">{pwd}</span>】',
                               "green")
            else:
                color_emit_str(signal, widget, "[×]请正确输入uid！", "red")
        elif mission == 16:
            if uid:
                pwd = db_pwd.tiktok(uid)
                color_emit_str(signal, widget,
                               f'[√]成功获取到抖音聊天数据库密钥【<span style="color:brown;">{pwd}</span>】',
                               "green")
        else:
            color_emit_str(signal, widget, "[×]不正确的任务！", "red")

    def decrypt_database(self, mission: int, filepath: str, password: str, widget: QObject):
        db_tool = DbTool()
        print(mission, filepath, password)
        signal = self.decrypt_database_signal
        ret, suc = None, None
        if not filepath and not password:
            color_emit_str(signal, widget, "[×]请输入路径和密钥！", "red")
        filepath = filepath.replace('"', '')
        if mission == 1:
            ret, wxid = db_tool.decrypt_EnMicroMsg(password, filepath)
            if wxid != "":
                color_emit_str(signal, widget, "[√]" + ret, "green")
                color_emit_str(signal, widget, f"[+]成功提取到wxid【{wxid}】", "green")
            else:
                color_emit_str(signal, widget, "[×]" + ret, "red")
        elif mission == 2:
            ret, suc = db_tool.decrypt_FTS5IndexMicroMsg(password, filepath)
        elif mission == 4:
            ret, suc = db_tool.decrypt_amap(filepath)
        elif mission == 5:
            ret, suc = db_tool.decrypt_dingtalk(password, filepath)
        elif mission == 6:
            ret, suc = db_tool.decrypt_SQLCipher4_default(password, filepath)
        elif mission == 7:
            ret, suc = db_tool.decrypt_SQLCipher3_default(password, filepath)
        elif mission == 15:
            ret, suc = db_tool.decrypt_wcdb_default(password, filepath)
        elif mission == 18:
            ret, suc = db_tool.decrypt_ntqq(password, filepath)
        else:
            color_emit_str(signal, widget, "不正确的任务", "red")
        if mission != 1:
            if suc == 1:
                color_emit_str(signal, widget, "[√]" + ret, "green")
            else:
                color_emit_str(signal, widget, "[×]" + ret, "red")

    def analyze_file(self, mission: int, filepath: str, password: str, widget: QObject):
        print(mission, filepath, password)
        signal = self.analyze_file_signal
        ret = None
        if not filepath and not password:
            emit_data(signal, widget, "[×]请输入路径和密钥！")
        filepath = filepath.replace('"', '')
        if mission == 8:
            data = analyzeNavicat(filepath)
            ret = my_dict_to_table(data)
        elif mission == 9:
            if password:
                try:
                    if filepath.endswith('.ini'):
                        data = analyzeMoba(filepath, password, 1)
                    else:
                        data = analyzeMoba(filepath, password, 0)
                    ret = my_dict_to_table(data)
                except Exception as e:
                    emit_data(signal, widget, str(e))
            else:
                emit_data(signal, widget, "请输入密钥！")
        elif mission == 10:
            try:
                data = analyzeDbeaer(filepath)
                ret = my_dict_to_table(data)
            except Exception as e:
                emit_data(signal, widget, str(e))
        elif mission == 11:
            try:
                data = analyzeFinalShell(filepath)
                ret = my_dict_to_table(data)
            except Exception as e:
                emit_data(signal, widget, str(e))
        elif mission == 12:
            if password:
                try:
                    data = analyzeXshell(filepath, password)
                    ret = my_dict_to_table(data)
                except Exception as e:
                    emit_data(signal, widget, str(e))
            else:
                emit_data(signal, widget, "请输入密钥！")
        elif mission == 14:
            try:
                data = analyzeUTools(filepath, password)
                ret = my_dict_to_table(data)
                emit_data(signal, widget, "解密后的内容已保存到同目录的data_decrypt.txt下，图片在wxjzc_images目录中")
            except Exception as e:
                emit_data(signal, widget, str(e))
        elif mission == 17:
            try:
                data = analyzeHawk2(filepath, password)
                ret = my_table_to_table(data, ['键', '值'])
            except Exception as e:
                emit_data(signal, widget, str(e))
        elif mission == 19:
            try:
                data = analyzeMetaMask(filepath)
                ret = my_dict_to_table(data)
            except Exception as e:
                emit_data(signal, widget, str(e))
        else:
            emit_data(signal, widget, "不正确的任务")
        emit_data(signal, widget, ret)

    def analyze_registry(self, filepath: str, widget: QObject):
        signal = self.analyze_registry_signal
        if not filepath:
            emit_data(signal, widget, "[×]请输入路径和密钥！")
        dic = {}
        if os.path.isdir(filepath):
            dirs = os.listdir(filepath)
            ntuser = []
            for v in dirs:
                if v.endswith('.DAT'):
                    ntuser.append(filepath + f'/{v}')
            dic.update({'NTUSER': ntuser})
            if 'SAM' in dirs:
                dic.update({'SAM': filepath + '/SAM'})
            if 'SYSTEM' in dirs:
                dic.update({'SYSTEM': filepath + '/SYSTEM'})
            if 'SOFTWARE' in dirs:
                dic.update({'SOFTWARE': filepath + '/SOFTWARE'})
        data1, data2 = analyzeWin(dic)
        ret1 = my_dict_to_table(data1)
        ret2 = my_table_to_table(data2, ['键', '值'])
        ret1.update(ret2)
        emit_data(signal, widget, ret1)
