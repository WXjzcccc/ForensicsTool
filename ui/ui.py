from typing import List

from PySide6.QtCore import Slot, QObject, QProcess
from PySide6.QtWidgets import QHBoxLayout, QComboBox, QLabel, QGridLayout, QLineEdit, QTextEdit, QWidget, QPushButton, \
    QVBoxLayout, QTableWidget, QTableWidgetItem, QTabWidget, QMessageBox

from .starter import Starter


def show_message_box(title: str, message: str):
    message_box = QMessageBox()
    message_box.setWindowTitle(title)
    message_box.setText(message)
    message_box.setIcon(QMessageBox.Icon.Warning)
    message_box.exec()


class MyLineEdit(QLineEdit):
    def __init__(self, parent=None, acceptDrops: bool = False):
        super().__init__(parent)
        self.setAcceptDrops(acceptDrops)

    def dragEnterEvent(self, event):
        if event.mimeData().hasUrls():
            event.acceptProposedAction()

    def dropEvent(self, event):
        urls = event.mimeData().urls()
        if urls:
            file_path = urls[0].toLocalFile()
            self.setText(file_path)


class UI(QObject):
    def __init__(self):
        super().__init__()
        self.ui_widgets = {
            "密钥计算": None,
            "数据库解密": None,
            "数据提取": None,
            "注册表分析": None,
            "暴力破解": None,
            "关于": None,
        }
        self.missions = {
            1: '微信的EnMicroMsg.db',
            2: '微信的FTS5IndexMicroMsg_encrypt.db',
            3: '野火IM系应用的data',
            4: '高德的girf_sync.db',
            5: '钉钉的数据库，需要提供shared_prefs/com.alibaba.android.rimet_preferences.xml中带有数据库名的字段的值中出现，如HUAWEI P40/armeabi-v7a/P40/qcom/HUAWEIP40',
            6: 'SQLCipher4加密的数据库',
            7: 'SQLCipher3加密的数据库',
            8: 'Navicat连接信息提取，需指定-f为目标用户的注册表文件"NTUSER.DAT"',
            9: 'MobaXterm连接信息解密，可以指定MobaXterm.ini配置文件或用户注册表文件"NTUSER.DAT"，解密需要给出主密码',
            10: 'Dbeaver连接信息解密，指定-f为目标文件data-sources.json和credentials-config.json的父目录',
            11: 'FinalShell连接信息解密，指定-f为目标文件夹conn，需要确保已经配置了JAVA_HOME环境变量',
            12: 'XShell、XFtp连接信息解密，指定-f为目标文件夹session，并提供-p参数，值为计算机的用户名+sid',
            13: '默往APP的msg.db，计算密钥时提供--uid参数',
            14: '提取uTools的剪贴板数据，指定-f参数为剪贴板数据目录，-p为解密密钥，解密超级剪贴板的数据请指定-p参数值为super',
            15: 'wcdb加密的数据库',
            16: '抖音的聊天数据库，计算密钥时提供--uid参数',
            17: 'Hawk2.xml数据解密，指定-f参数为文件路径，-p参数为同目录下的crypto.KEY_256.xml或crypto.KEY_128.xml中的base64值',
            18: 'ntqq数据库解密，指定-f参数位文件路径，--uid参数为对应qq号的uid',
            19: 'MetaMask解析，需要指定-f参数为persist-root文件路径',
        }
        self.params = {
            "password": "解密的密码，处理高德时不适用",
            "uin": "微信用户的uin，可能是负值，在shared_prefs/auth_info_key_prefs.xml文件中_auth_uin的值",
            "imei": "微信获取到的IMEI或MEID，在shared_prefs/DENGTA_META.xml文件中IMEI_DENGTA的值，在高版本中通常是1234567890ABCDEF，可以为空",
            "wxid": "数据库所属的wxid，一般情况下在解密EnMicroMsg.db的时候会一并提取，若无需要，请从shared_prefs/com.tencent.mm_preferences.xml中提取login_weixin_username的值",
            "token": "野火IM系应用的用户token，shared_prefs/config.xml的token的值",
            "uid": "默往（通常在shared_prefs/im.xml中的userId的值）、抖音（数据库文件名中的id）计算密钥需要的内容、QQ（msf_mmkv_file中QQ号对应的uid）",
            "file": "指定需要处理的文件",
        }
        self.crack_missions = {
            1: "AirDrop手机号爆破",
            2: "微信UIN爆破"
        }
        self.crack_params = {
            "target": "目标列表，以,进行分隔",
            "mac": "AirDrop要爆破的手机号段，以,进行分隔，如138,139",
            "region": "AirDrop要爆破的区号，如86、85、1",
            "length": "AirDrop要爆破的手机号长度（去除区号和号段）",
        }
        self.starter = Starter()

    def getComboBox(self, ids: List[int], combo_type: int = 1) -> QComboBox:
        combo = QComboBox()
        if combo_type == 1:
            for _id in ids:
                combo.addItem(self.missions[_id], _id)
        elif combo_type == 2:
            for _id in ids:
                combo.addItem(self.crack_missions[_id], _id)
        combo.setMaximumWidth(600)
        return combo

    def getInputGroup(self, param: str, accept_drops: bool = False, param_type: int = 1) -> (QLabel, QLineEdit):
        label = QLabel()
        label.setText(param)
        line_edit = MyLineEdit(acceptDrops=accept_drops)
        if param_type == 1:
            line_edit.setPlaceholderText(self.params.get(param))
            line_edit.setToolTip(self.params.get(param))
        elif param_type == 2:
            line_edit.setPlaceholderText(self.crack_params.get(param))
            line_edit.setToolTip(self.crack_params.get(param))
        return label, line_edit

    def init_passwd_calc(self):
        widget = QWidget()
        layout = QGridLayout()
        mission_label = QLabel("请选择任务")
        combo = self.getComboBox([1, 2, 3, 13, 16])
        layout.addWidget(mission_label, 0, 0, 1, 1)
        layout.addWidget(combo, 0, 1, 1, 3)
        uin_label, uin_line = self.getInputGroup("uin")
        layout.addWidget(uin_label, 1, 0, 1, 1)
        layout.addWidget(uin_line, 1, 1, 1, 1)
        imei_label, imei_line = self.getInputGroup("imei")
        layout.addWidget(imei_label, 1, 2, 1, 1)
        layout.addWidget(imei_line, 1, 3, 1, 1)
        wxid_label, wxid_line = self.getInputGroup("wxid")
        layout.addWidget(wxid_label, 2, 0, 1, 1)
        layout.addWidget(wxid_line, 2, 1, 1, 1)
        token_label, token_line = self.getInputGroup("token")
        layout.addWidget(token_label, 2, 2, 1, 1)
        layout.addWidget(token_line, 2, 3, 1, 1)
        uid_label, uid_line = self.getInputGroup("uid")
        layout.addWidget(uid_label, 3, 0, 1, 1)
        layout.addWidget(uid_line, 3, 1, 1, 1)
        button = QPushButton("开始")
        button.clicked.connect(
            lambda: self.handle_passwd_calc(int(combo.currentData()), uin_line.text(), imei_line.text(),
                                            wxid_line.text(), token_line.text(), uid_line.text(), output_field))
        self.starter.passwd_calc_signal.connect(self.update_text_out)
        hbox = QHBoxLayout()
        clear_button = QPushButton("清空输出")
        clear_button.clicked.connect(lambda: output_field.clear())
        hbox.addWidget(button)
        hbox.addWidget(clear_button)
        layout.addLayout(hbox, 3, 2, 1, 2)
        output_field = QTextEdit()
        output_field.setReadOnly(True)
        layout.addWidget(output_field, 4, 0, 1, 4)
        widget.setLayout(layout)
        self.ui_widgets["密钥计算"] = widget

    def handle_passwd_calc(self, mission: int, uin: str, imei: str, wxid: str, token: str, uid: str,
                           output_field: QTextEdit):
        self.starter.passwd_calc(mission, uin, imei, wxid, token, uid, output_field)

    @Slot(QTextEdit, str)
    def update_text_out(self, output_field: QTextEdit, out: str):
        output_field.append(out)

    def init_decrypt_database(self):
        widget = QWidget()
        layout = QGridLayout()
        mission_label = QLabel("请选择任务")
        combo = self.getComboBox([1, 2, 4, 5, 6, 7, 15, 18])
        layout.addWidget(mission_label, 0, 0, 1, 1)
        layout.addWidget(combo, 0, 1, 1, 3)
        file_label, file_line = self.getInputGroup("file", True)
        file_line.setAcceptDrops(True)
        layout.addWidget(file_label, 1, 0, 1, 1)
        layout.addWidget(file_line, 1, 1, 1, 1)
        passwd_label, passwd_line = self.getInputGroup("password")
        layout.addWidget(passwd_label, 1, 2, 1, 1)
        layout.addWidget(passwd_line, 1, 3, 1, 1)
        button = QPushButton("开始")
        button.clicked.connect(
            lambda: self.handle_decrypt_database(int(combo.currentData()), file_line.text(), passwd_line.text(),
                                                 output_field))
        hbox = QHBoxLayout()
        clear_button = QPushButton("清空输出")
        clear_button.clicked.connect(lambda: output_field.clear())
        hbox.addWidget(button)
        hbox.addWidget(clear_button)
        layout.addLayout(hbox, 2, 0, 1, 4)
        output_field = QTextEdit()
        output_field.setReadOnly(True)
        layout.addWidget(output_field, 3, 0, 1, 4)
        widget.setLayout(layout)
        self.starter.decrypt_database_signal.connect(self.update_text_out)
        self.ui_widgets["数据库解密"] = widget

    def handle_decrypt_database(self, mission: int, path: str, password: str, output_field: QTextEdit):
        self.starter.decrypt_database(mission, path, password, output_field)

    def init_analyze_file(self):
        widget = QWidget()
        layout = QGridLayout()
        mission_label = QLabel("请选择任务")
        combo = self.getComboBox([8, 9, 10, 11, 12, 14, 17, 19])
        layout.addWidget(mission_label, 0, 0, 1, 1)
        layout.addWidget(combo, 0, 1, 1, 3)
        file_label, file_line = self.getInputGroup("file", True)
        file_line.setAcceptDrops(True)
        layout.addWidget(file_label, 1, 0, 1, 1)
        layout.addWidget(file_line, 1, 1, 1, 1)
        passwd_label, passwd_line = self.getInputGroup("password")
        layout.addWidget(passwd_label, 1, 2, 1, 1)
        layout.addWidget(passwd_line, 1, 3, 1, 1)
        button = QPushButton("开始")
        self.starter.analyze_file_signal.connect(self.update_table_out)
        button.clicked.connect(
            lambda: self.handle_analyze_file(int(combo.currentData()), file_line.text(), passwd_line.text(), tab))
        layout.addWidget(button, 2, 0, 1, 4)
        tab = QTabWidget()
        layout.addWidget(tab, 3, 0, 1, 4)
        widget.setLayout(layout)
        self.ui_widgets["数据提取"] = widget

    def handle_analyze_file(self, mission: int, path: str, password: str, tab: QTabWidget):
        tab.clear()
        self.starter.analyze_file(mission, path, password, tab)

    @Slot(QTextEdit, object)
    def update_table_out(self, tab: QTabWidget, data: object):
        if type(data) == str:
            show_message_box("警告", data)
        if type(data) == dict:
            for key, value in data.items():
                self.add_table(key, value, tab)

    def init_analyze_registry(self):
        widget = QWidget()
        layout = QGridLayout()
        file_label, file_line = self.getInputGroup("file", True)
        file_label = QLabel()
        file_label.setText("file")
        file_line = MyLineEdit(acceptDrops=True)
        file_line.setPlaceholderText("请提供包含SAM、SOFTWARE、SYSTEM注册表文件的目录")
        file_line.setToolTip("请提供包含SAM、SOFTWARE、SYSTEM注册表文件的目录")
        file_line.setAcceptDrops(True)
        layout.addWidget(file_label, 0, 0, 1, 1)
        layout.addWidget(file_line, 0, 1, 1, 1)
        button = QPushButton("开始")
        self.starter.analyze_registry_signal.connect(self.update_table_out)
        layout.addWidget(button, 0, 2, 1, 2)
        tab = QTabWidget()
        layout.addWidget(tab, 1, 0, 1, 4)
        button.clicked.connect(lambda: self.handle_analyze_registry(file_line.text(), tab))
        widget.setLayout(layout)
        self.ui_widgets["注册表分析"] = widget

    def handle_analyze_registry(self, path: str, tab: QTabWidget):
        tab.clear()
        self.starter.analyze_registry(path, tab)

    def init_forensics_crack(self):
        widget = QWidget()
        layout = QGridLayout()
        mission_label = QLabel("请选择任务")
        combo = self.getComboBox([1,2],2)
        layout.addWidget(mission_label, 0, 0, 1, 1)
        layout.addWidget(combo, 0, 1, 1, 3)
        target_label, target_line = self.getInputGroup("target", param_type=2)
        layout.addWidget(target_label, 1, 0, 1, 1)
        layout.addWidget(target_line, 1, 1, 1, 1)
        mac_label, mac_line = self.getInputGroup("mac", param_type=2)
        layout.addWidget(mac_label, 1, 2, 1, 1)
        layout.addWidget(mac_line, 1, 3, 1, 1)
        region_label, region_line = self.getInputGroup("region", param_type=2)
        layout.addWidget(region_label, 2, 0, 1, 1)
        layout.addWidget(region_line, 2, 1, 1, 1)
        length_label, length_line = self.getInputGroup("length", param_type=2)
        layout.addWidget(length_label, 2, 2, 1, 1)
        layout.addWidget(length_line, 2, 3, 1, 1)
        button = QPushButton("开始")
        button.clicked.connect(
            lambda: self.handle_forensics_crack(int(combo.currentData()), target_line.text(), mac_line.text(),
                                                 region_line.text(), length_line.text(), output_field))
        hbox = QHBoxLayout()
        stop_button = QPushButton("停止")
        stop_button.clicked.connect(lambda: self.starter.stop_crack())
        clear_button = QPushButton("清空输出")
        clear_button.clicked.connect(lambda: output_field.clear())
        hbox.addWidget(button)
        hbox.addWidget(stop_button)
        hbox.addWidget(clear_button)
        layout.addLayout(hbox, 3, 0, 1, 4)
        output_field = QTextEdit()
        output_field.setReadOnly(True)
        layout.addWidget(output_field, 4, 0, 1, 4)
        widget.setLayout(layout)
        self.starter.forensics_crack_signal.connect(self.update_text_out)
        self.ui_widgets["暴力破解"] = widget

    def handle_forensics_crack(self, mission: int, target: str, mac: str, region: str, length: str, output_field: QTextEdit):
        self.starter.forensics_crack(mission, target, mac, region, length, output_field)

    def init_about(self):
        widget = QWidget()
        layout = QVBoxLayout()

        title_label = QLabel("<h3>关于软件</h3>")
        intro_label = QLabel('''个人用于学习电子取证、提升效率的小软件。<br>
        开源地址：<a href=\"https://github.com/WXjzcccc/ForensicsTool\">https://github.com/WXjzcccc/ForensicsTool</a>''')
        intro_label.setOpenExternalLinks(True)

        author_label = QLabel("<h3>关于作者</h3>")
        author_text = QLabel('''邮箱：wxjzcroot@gmail.com<br>
        GitHub：<a href=\"https://github.com/WXjzcccc\">https://github.com/WXjzcccc</a><br>
        博客：<a href=\"https://www.cnblogs.com/WXjzc\">https://www.cnblogs.com/WXjzc</a>''')
        author_text.setOpenExternalLinks(True)

        contributors_label = QLabel("<h3>贡献者</h3>")
        contributors_text = QLabel('''[1]WXjzc,<a href=\"https://github.com/WXjzcccc\">https://github.com/WXjzcccc</a><br>
        [2]b3nguang,<a href=\"https://github.com/b3nguang\">https://github.com/b3nguang</a>（MetaMask解析）''')
        contributors_text.setOpenExternalLinks(True)

        license_label = QLabel("<h3>软件许可</h3>")
        license_text = QLabel("GPL3.0")

        opensource_label = QLabel("<h3>使用的开源软件</h3>")
        opensource_text = QLabel('''<a href=\"https://github.com/jas502n/FinalShellDecodePass\">https://github.com/jas502n/FinalShellDecodePass</a><br>
        <a href=\"https://github.com/dzxs/Xdecrypt\">https://github.com/dzxs/Xdecrypt</a><br>
        <a href=\"https://github.com/clemthi/mobaxterm-decrypt-pass\">https://github.com/clemthi/mobaxterm-decrypt-pass</a>''')
        opensource_text.setOpenExternalLinks(True)
        layout.addWidget(title_label)
        layout.addWidget(intro_label)
        layout.addWidget(author_label)
        layout.addWidget(author_text)
        layout.addWidget(contributors_label)
        layout.addWidget(contributors_text)
        layout.addWidget(license_label)
        layout.addWidget(license_text)
        layout.addWidget(opensource_label)
        layout.addWidget(opensource_text)
        widget.setLayout(layout)
        self.ui_widgets["关于"] = widget

    def init_ui(self) -> dict:
        self.init_passwd_calc()
        self.init_decrypt_database()
        self.init_analyze_file()
        self.init_analyze_registry()
        self.init_forensics_crack()
        self.init_about()
        return self.ui_widgets

    def add_table(self, name, data, tab_widget):
        # 创建新 Tab 页
        tab = QWidget()
        tab_layout = QVBoxLayout(tab)
        # 创建表格
        table = QTableWidget()
        table_head = data[0]
        table_data = data[1]
        x = len(table_head)
        y = len(table_data)
        table.setColumnCount(x)
        table.setRowCount(y)
        table.setHorizontalHeaderLabels(table_head)
        table.setVerticalHeaderLabels([f"{i}" for i in range(1, y + 1)])
        # 填充表格内容
        for row in range(y):
            for col in range(x):
                msg = str(table_data[row][col])
                item = QTableWidgetItem(msg)
                item.setToolTip(msg)
                table.setItem(row, col, item)
        # 设置表格的宽度自适应和最小宽度
        table.horizontalHeader().setStretchLastSection(True)
        table.horizontalHeader().setDefaultSectionSize(120)
        tab_layout.addWidget(table)
        tab_widget.addTab(tab, name)
