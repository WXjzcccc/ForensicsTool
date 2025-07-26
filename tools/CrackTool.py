import os
import subprocess
from PySide6.QtCore import QProcess, Slot, Signal, QObject
from .PrintTool import print_red, print_yellow

def color_emit_str(signal, widget, text, color):
    signal.emit(widget, f'<span  style="color:{color};">{text}</span>')

class CrackTool:
    def __init__(self, signal: Signal, widget: QObject):
        self.qt_process = QProcess()
        self.signal = signal
        self.widget = widget
        self.qt_process.readyReadStandardOutput.connect(self.handle_output)
        self.qt_process.readyReadStandardError.connect(self.handle_error)
        self.qt_process.finished.connect(self.on_finished)
        self.qt_process.started.connect(self.on_started)


    def get_relative_path(self, relative_path):
        """获取配置文件的绝对路径"""
        base_path = os.path.dirname(os.path.abspath(__file__))
        return os.path.join(base_path, relative_path)


    def crack_airdrop(self, mission: int, target: str, mac: str, region: str, length: str):
        forensics_cracker = self.get_relative_path("../lib/ForensicsCracker.exe")
        if os.path.exists(forensics_cracker) and os.path.isfile(forensics_cracker):
            forensics_cracker = os.path.abspath(forensics_cracker)
        else:
            print_red('[失败]---->原因[缺少程序<ForensicsCracker.exe>！]')
            ret = '缺少程序<ForensicsCracker.exe>！'
            return ret, ''
        print_yellow('<AirDrop破解>')
        cmd = ["-m",str(mission),"-target",target]
        if mac.find(",") != -1:
            cmd.append('-mac')
            cmd.append(mac)
        else:
            try:
                int(mac)
                cmd.append('-mac')
                cmd.append(mac)
            except ValueError:
                pass
        try:
            int(region)
            cmd.append('-region')
            cmd.append(region)
        except ValueError:
            pass
        try:
            int(length)
            cmd.append('-length')
            cmd.append(length)
        except ValueError:
            pass
        self.qt_process.start(forensics_cracker,cmd)


    def crack_wx_uin(self,mission: int, target: str):
        forensics_cracker = self.check_exe()
        if forensics_cracker == '':
            color_emit_str(self.signal,self.widget,"[×]缺少程序<ForensicsCracker.exe>！","red")
        color_emit_str(self.signal,self.widget,"<开始爆破微信UIN>","yellow")
        self.qt_process.start(forensics_cracker,
            ["-m", str(mission), "-target", target])

    def check_exe(self):
        forensics_cracker = self.get_relative_path("../lib/ForensicsCracker.exe")
        if os.path.exists(forensics_cracker) and os.path.isfile(forensics_cracker):
            forensics_cracker = os.path.abspath(forensics_cracker)
            return forensics_cracker
        else:
            print_red('[失败]---->原因[缺少程序<ForensicsCracker.exe>！]')
            return ''


    @Slot()
    def handle_output(self):
        data = self.qt_process.readAllStandardOutput().data().decode()
        data = data.replace("\n","<br>")
        color_emit_str(self.signal,self.widget,data.rstrip("<br>"),"green")

    @Slot()
    def handle_error(self):
        error = self.qt_process.readAllStandardError().data().decode()
        error = error.replace("\n","<br>")
        color_emit_str(self.signal,self.widget,error.rstrip("<br>"),"red")

    @Slot(int)
    def on_finished(self, exit_code):
        color_emit_str(self.signal,self.widget,"[√]爆破结束","green")

    @Slot()
    def on_started(self):
        crack_type = self.qt_process.arguments()[1]
        if crack_type == '1':
            color_emit_str(self.signal, self.widget, "[+]开始爆破AirDrop", "brown")
        elif crack_type == '2':
            color_emit_str(self.signal, self.widget, "[+]开始爆破微信UIN", "brown")