import os

from Registry import Registry
from .PrintTool import *
import sys
import datetime
from construct import Struct, Bytes, CString
import struct
import pytz


class WinTool:
    _system = ''
    _ntuser = []
    _sam = ''
    _software = ''

    def __init__(self, hives: dict):
        self._system = hives['SYSTEM'] if 'SYSTEM' in hives.keys() else ''
        self._ntuser = hives['NTUSER'] if 'NTUSER' in hives.keys() else []
        self._sam = hives['SAM'] if 'SAM' in hives.keys() else ''
        self._software = hives['SOFTWARE'] if 'SOFTWARE' in hives.keys() else ''

    def timestamp(self, time_str: str) -> str:
        dt = datetime.datetime.fromtimestamp(time_str)
        formatted_date = dt.strftime('%Y-%m-%d %H:%M:%S')
        return formatted_date

    def timestamp_hex(self, time_str: str) -> str:
        dt = datetime.datetime.fromtimestamp(int.from_bytes(time_str, byteorder='little') / 10 ** 7 - 11644473600)
        formatted_date = dt.strftime('%Y-%m-%d %H:%M:%S')
        return formatted_date

    def get_timezone(self, ) -> str:
        if self._system == '':
            return '没有给出SYSTEM注册表文件!'
        reg = Registry.Registry(self._system)
        root_key = reg.open(r'ControlSet001\Control\TimeZoneInformation')
        standard_name = root_key.value('StandardName').value()
        reg = Registry.Registry(self._software)
        root_key = reg.open(r'Microsoft\Windows NT\CurrentVersion\Time Zones')
        for sub_key in root_key.subkeys():
            if sub_key.value('MUI_Std').value() == standard_name:
                return sub_key.value('Display').value()
        return ''

    def get_last_shutdown_time(self) -> str:
        if self._system == '':
            return '没有给出SYSTEM注册表文件!'
        reg = Registry.Registry(self._system)
        root_key = reg.open(r'ControlSet001\Control\Windows')
        shutdown_time = root_key.value('ShutdownTime').value()
        return self.timestamp_hex(shutdown_time)

    def timestamp_to_datetime(self, converted_time, origin_timezone='UTC', target_timezone='Asia/Shanghai'):
        _timezone = pytz.timezone(origin_timezone)
        converted_time = _timezone.localize(converted_time)
        converted_time = converted_time.astimezone(pytz.timezone(target_timezone))
        date_str = converted_time.strftime('%Y-%m-%d %H:%M:%S')
        return date_str

    def get_relative_path(self,relative_path):
        """获取配置文件的绝对路径"""
        base_path = os.path.dirname(os.path.abspath(__file__))
        return os.path.join(base_path, relative_path)

    def windows_file_time_to_datetime(self, timestamp, origin_timezone='UTC', target_timezone='Asia/Shanghai'):
        try:
            base_time = datetime.datetime(1601, 1, 1)
            delta = datetime.timedelta(microseconds=timestamp / 10)
            converted_time = base_time + delta
            return self.timestamp_to_datetime(converted_time, origin_timezone, target_timezone)
        except Exception as e:
            print_red(str(e))
            return '解析失败'

    def get_nt_hash(self) -> dict:
        if self._sam == '':
            return '没有给出SAM注册表文件!'
        if self._system == '':
            return '没有给出SYSTEM注册表文件!'
        import subprocess
        nthash = self.get_relative_path('../lib/ntHashFromRegFile.exe')
        if os.path.exists(nthash) and os.path.isfile(nthash):
            nthash = os.path.abspath(nthash)
        else:
            print_red('[失败]---->原因[缺少依赖<ntHashFromRegFile.exe>！]')
            return {}
        command = [nthash, '--sam', self._sam, '--system', self._system]
        result = subprocess.run(command, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        stdout = result.stdout.decode('utf-8')
        lines = stdout.splitlines()
        ret = []
        for line in lines:
            tmp = line.split(':')
            username = tmp[0]
            uid = tmp[1]
            lm_hash = tmp[2]
            nt_hash = tmp[3]
            ret.append({
                "username":username,
                    "uid":uid,
                    "lm_hash":lm_hash,
                    "nt_hash":nt_hash
            })
        return ret


    def parse_f_key(self, binary_data):
        # 从二进制数据中提取特定信息

        # User ID (偏移48-52字节，hex)
        user_id = struct.unpack('<I', binary_data[48:52])[0]

        # Last log-in Time (偏移16-24字节，8字节FILETIME格式)
        last_login_time_raw = struct.unpack('<Q', binary_data[8:16])[0]
        last_login_time = self.windows_file_time_to_datetime(last_login_time_raw)

        # Last PW change (偏移40-48字节，8字节FILETIME格式)
        last_pw_change_raw = struct.unpack('<Q', binary_data[24:32])[0]
        last_pw_change = self.windows_file_time_to_datetime(last_pw_change_raw)

        # Last failed log-in (偏移56-64字节，8字节FILETIME格式)
        last_failed_login_raw = struct.unpack('<Q', binary_data[40:48])[0]
        last_failed_login = self.windows_file_time_to_datetime(last_failed_login_raw)

        # Log-on Count (偏移64-66字节，2字节整数)
        invalid_pw_count = struct.unpack('<H', binary_data[64:66])[0]

        # Invalid PW Count (偏移66-68字节，2字节整数)
        logon_count = struct.unpack('<H', binary_data[66:68])[0]

        return {
            '上次登录时间': last_login_time,
            '登录次数': logon_count,
            '登录失败次数': invalid_pw_count,
            '上次密码修改时间': last_pw_change,
            '上次登录失败时间': last_failed_login
        }

    # TODO 获取用户信息，待添加更多细节
    def get_users(self) -> dict:
        def hex_key(key):
            _key = hex(key)[2:]
            return _key.upper().zfill(8)
        if self._sam == '':
            return '没有给出SAM注册表文件!'
        if self._software == '':
            return '没有给出SOFTWARE注册表文件!'
        users = []
        reg = Registry.Registry(self._sam)
        root_key = reg.open(r'SAM\Domains\Account\Users')
        sub_keys = root_key.subkey('Names')
        for sub_key in sub_keys.subkeys():
            dic = {}
            username = sub_key.name()
            vtype = sub_key.value('').value_type()
            key = hex_key(vtype)
            user_key = root_key.subkey(key)
            hint = user_key.value('UserPasswordHint').value() if self.check_key('UserPasswordHint', user_key) else b'/'
            expires = user_key.value('AccountExpires').value() if self.check_key('AccountExpires', user_key) else '/'
            f_data = user_key.value('F').value() if self.check_key('F', user_key) else b''
            f_info = self.parse_f_key(f_data)
            software_key = Registry.Registry(self._software).open('Microsoft\Windows NT\CurrentVersion\ProfileList')
            sid = '/'
            for v in software_key.subkeys():
                if v.name().endswith(str(vtype)):
                    sid = v.name()
            dic.update({
                '用户名': username,
                'SID': sid,
                '密码提示': hint.decode('utf8'),
                '过期时间': expires
            })
            dic.update(f_info)
            users.append(dic)
        return users

    def get_recent_name(self, data: bytes) -> bytes:
        _ = b''
        for i in range(0, len(data), 2):
            if data[i:i + 2] == b'\x00\x00':
                break
            _ += data[i:i + 2]
        return _

    def parsePIDL(self, data):
        data = data.split(b"\x04\x00\xef\xbe")
        counter = 0
        path = ''
        for d in data:
            if counter == 0:
                letter = ''
                if d.startswith(
                        b'\x14\x00\x1F\x50\xE0\x4F\xD0\x20\xEA\x3A\x69\x10\xA2\xD8\x08\x00\x2B\x30\x30\x9D\x19\x00\x2F'):
                    letter = d[23:25].decode()
                path += letter
            else:
                format = Struct(
                    'unkuwn' / Bytes(38),
                    'Path' / CString("utf16")
                )
                dd = format.parse(d)
                path += "\\" + str(dd.Path)
            counter = counter + 1
        return path

    def get_recent(self) -> list:
        if self._ntuser == []:
            return '没有给出NTUSER注册表文件!'
        lst = []
        for ntuser in self._ntuser:
            reg = Registry.Registry(ntuser)
            root_key = reg.open(r'Software\Microsoft\Windows\CurrentVersion\Explorer\RecentDocs')
            recent_docs = []
            for sub_key in root_key.subkeys():
                _type = sub_key.name()
                for v in sub_key.values():
                    if v.name() == 'MRUListEx':
                        continue
                    _ = self.get_recent_name(v.value())
                    filename = _.decode('utf-16le')
                    recent_docs.append({
                        '类型': _type,
                        '文件名': filename,
                    })
            lst.append([recent_docs, '最近访问的文档', ntuser])
            root_key = reg.open(r'Software\Microsoft\Windows\CurrentVersion\Explorer\ComDlg32')
            recent_pro = []
            program_key = root_key.subkey('CIDSizeMRU')
            for v in program_key.values():
                if v.name() == 'MRUListEx':
                    continue
                _ = self.get_recent_name(v.value())
                filename = _.decode('utf-16le')
                recent_pro.append({
                    '程序': filename
                })
            program_key = root_key.subkey('LastVisitedPidlMRU')
            for v in program_key.values():
                if v.name() == 'MRUListEx':
                    continue
                _ = self.get_recent_name(v.value())
                filename = _.decode('utf-16le')
                if filename not in [x['程序'] for x in recent_pro]:
                    recent_pro.append({
                        '程序': filename
                    })
            lst.append([recent_pro, '最近打开的程序', ntuser])
            recent_files = []
            file_key = root_key.subkey('OpenSavePidlMRU')
            for sub_key in file_key.subkeys():
                for v in sub_key.values():
                    if v.name() == 'MRUListEx':
                        continue
                    _ = self.parsePIDL(v.value())
                    filename = _
                    recent_files.append({
                        '类型': sub_key.name(),
                        '文件路径': filename
                    })
            # TODO 解决编码问题
            lst.append([recent_files, '最近保存的文件', ntuser])
        return lst

    # TODO 解析最近使用的运行命令，解析开机自启动应用。。。

    def get_system_name(self) -> str:
        if self._system == '':
            return '没有给出SYSTEM注册表文件!'
        reg = Registry.Registry(self._system)
        root_key = reg.open(r'ControlSet001\Control\ComputerName\ComputerName')
        return root_key.value('ComputerName').value()

    def get_login_user(self):
        reg = Registry.Registry(self._software)
        root_key = reg.open(r'Microsoft\Windows\CurrentVersion\Authentication\LogonUI')
        return root_key.value('LastLoggedOnUser').value()

    def get_default_browser(self) -> str:
        if self._ntuser == []:
            return '没有给出NTUSER注册表文件!'
        lst = []
        try:
            for ntuser in self._ntuser:
                reg = Registry.Registry(ntuser)
                root_key = reg.open(r'SOFTWARE\Microsoft\Windows\Shell\Associations\UrlAssociations\http\UserChoice')
                if self.check_key('ProgId', root_key):
                    lst.append([root_key.value('ProgId').value(),ntuser])
                elif self.check_key('Progid', root_key):
                    lst.append([root_key.value('Progid').value(),ntuser])
                else:
                    lst.append(['/',ntuser])
        except:
            lst = ['/', ntuser]
        return lst

    def check_key(self, name: str, key: Registry.RegistryKey) -> bool:
        for v in key.values():
            if v.name() == name:
                return True
        return False

    def get_net_info(self) -> list:
        lst = []
        if self._system == '':
            print_red('没有给出SYSTEM注册表文件!')
            sys.exit(-1)
        reg = Registry.Registry(self._system)
        interface_key = reg.open(r'ControlSet001\Services\Tcpip\Parameters\Interfaces')
        device_key = reg.open(r'ControlSet001\Control\Network\{4D36E972-E325-11CE-BFC1-08002BE10318}')
        keys = interface_key.subkeys()
        _keys = device_key.subkeys()
        for v in keys:
            dic = {}
            for _ in _keys:
                if v.name().upper() == _.name():
                    connection = _.subkey('connection')
                    dic.update(
                        {'名称': connection.value('Name').value() if self.check_key('Name', connection) else '/'})
            dic.update({'IP地址': v.value('IPAddress').value()[0] if self.check_key('IPAddress', v) else '/'})
            dic.update({'子网掩码': v.value('SubnetMask').value()[0] if self.check_key('SubnetMask', v) else '/'})
            dic.update({'网关': v.value('DefaultGateway').value()[0] if self.check_key('DefaultGateway', v) else '/'})
            dic.update(
                {'DHCP网络地址': v.value('DhcpIPAddress').value() if self.check_key('DhcpIPAddress', v) else '/'})
            dic.update({'DHCP网关': v.value('DhcpDefaultGateway').value()[0] if self.check_key('DhcpDefaultGateway',
                                                                                               v) else '/'})
            dic.update(
                {'DHCP子网掩码': v.value('DhcpSubnetMask').value() if self.check_key('DhcpSubnetMask', v) else '/'})
            dic.update({'DHCP服务地址': v.value('DhcpServer').value() if self.check_key('DhcpServer', v) else '/'})
            dic.update({'租赁时间': self.timestamp(v.value('LeaseObtainedTime').value()) if self.check_key(
                'LeaseObtainedTime', v) else '/'})
            dic.update({'过期时间': self.timestamp(v.value('LeaseTerminatesTime').value()) if self.check_key(
                'LeaseTerminatesTime', v) else '/'})
            dic.update({'DHCP服务地址': v.value('DhcpServer').value() if self.check_key('DhcpServer', v) else '/'})
            lst.append(dic)
        return lst

    def get_system_info(self) -> dict:
        system_info = {}
        if self._software == '':
            print_red('没有给出SOFTWARE注册表文件!')
            sys.exit(-1)
        reg = Registry.Registry(self._software)
        root_key = reg.open(r'Microsoft\Windows NT\CurrentVersion')
        system_info.update(
            {'Build信息': root_key.value('BuildLabEx').value() if self.check_key('BuildLabEx', root_key) else '/'})
        system_info.update({'Build版本': root_key.value('CurrentBuildNumber').value() if self.check_key(
            'CurrentBuildNumber', root_key) else '/'})
        system_info.update({'计算机名称': self.get_system_name()})
        system_info.update(
            {'版本信息': root_key.value('EditionID').value() if self.check_key('EditionID', root_key) else '/'})
        system_info.update({'安装时间(本地时区)': self.timestamp(
            root_key.value('InstallDate').value()) if self.check_key('InstallDate', root_key) else '/'})
        system_info.update(
            {'系统名称': root_key.value('ProductName').value() if self.check_key('ProductName', root_key) else '/'})
        system_info.update(
            {'发行ID': root_key.value('ReleaseID').value() if self.check_key('ReleaseID', root_key) else '/'})
        system_info.update(
            {'产品ID': root_key.value('ProductID').value() if self.check_key('ProductID', root_key) else '/'})
        product_key = reg.open(r'Microsoft\Windows NT\CurrentVersion\SoftwareProtectionPlatform')
        system_info.update({'产品密钥备份(非当前密钥)': product_key.value(
            'BackupProductKeyDefault').value() if self.check_key('BackupProductKeyDefault', product_key) else '/'})
        system_info.update({'系统时区': self.get_timezone()})
        system_info.update({'上次正常关机时间(本地时区)': self.get_last_shutdown_time()})
        system_info.update({'上次登录的用户': self.get_login_user()})
        system_info.update(
            {'默认浏览器': self.get_default_browser()[0][0] if self.get_default_browser() != [] else '/'})
        return system_info

    # TODO USB信息
    def get_usb(self) -> list:
        lst = []
        if self._system == '':
            print_red('没有给出SYSTEM注册表文件!')
            sys.exit(-1)
        reg = Registry.Registry(self._system)
        mount_key = reg.open(r'MountedDevices')
        volumes = []
        for v in mount_key.values():
            if v.name().startswith(r'\DosDevices'):
                # _volume = str(v.value().replace(b'\x00', b''))[2:-1]
                _volume = v.value().decode('utf-16le')
                if _volume.__contains__('#'):
                    volume = _volume.split('#')[0]
                    if volume.startswith('{') and volume.endswith('}'):
                        volumes.append(volume)
        for v in mount_key.values():
            if v.name()[10:] in volumes:
                print('===========')
        print(volumes)

    def get_web_location(self) -> list:
        if self._ntuser == []:
            return '没有给出NTUSER注册表文件!'
        lst = []
        for ntuser in self._ntuser:
            locations = []
            reg = Registry.Registry(ntuser)
            root_key = reg.open(r'SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\MountPoints2')
            subkeys = root_key.subkeys()
            for subkey in subkeys:
                location = subkey.name()
                if location.startswith('##'):
                    location = location.replace('#', '\\')
                    locations.append({"地址":location})
            lst.append([locations,'网路位置',ntuser])
        return lst

def analyzeWin(hives: dict):
    winTool = WinTool(hives)
    print_table(winTool.get_system_info(), title='系统信息')
    print_dict(winTool.get_net_info(), winTool.get_net_info()[0].keys(), title='网络信息')
    print_dict(winTool.get_users(), winTool.get_users()[0].keys(), title='用户信息')
    print_dict(winTool.get_nt_hash(),winTool.get_nt_hash()[0].keys(), title='密码信息')
    location_data = winTool.get_web_location()
    for v in location_data:
        if len(v[0]) > 0:
            print_dict(v[0], v[0][0].keys(), title=v[1] + f'[{v[2]}]')
    recent_data = winTool.get_recent()
    for v in recent_data:
        if len(v[0]) > 0:
            print_dict(v[0], v[0][0].keys(), title=v[1] + f'[{v[2]}]')
    # winTool.get_usb()
    print_yellow('[提示]---->最近保存的文件中，如果没有父目录，请在用户的桌面、文档等地方查找文件')
    print_yellow('[提示]---->会努力优化的。。')