from Registry import Registry
from PrintTool import *
import sys
import datetime

class WinTool:
    _system = ''
    _ntuser = ''
    _sam = ''
    _software = ''

    def __init__(self, hives :dict):
        self._system = hives['SYSTEM'] if 'SYSTEM' in hives.keys() else ''
        self._ntuser = hives['NTUSER'] if 'NTUSER' in hives.keys() else ''
        self._sam = hives['SAM'] if 'SAM' in hives.keys() else ''
        self._software = hives['SOFTWARE'] if 'SOFTWARE' in hives.keys() else ''

    def timestamp(self, time_str :str) -> str:
        dt = datetime.datetime.fromtimestamp(time_str)
        formatted_date = dt.strftime('%Y-%m-%d %H:%M:%S')
        return formatted_date

    def timestamp_hex(self, time_str :str) -> str:
        dt = datetime.datetime.fromtimestamp(int.from_bytes(time_str, byteorder='little') / 10**7 - 11644473600)
        formatted_date = dt.strftime('%Y-%m-%d %H:%M:%S')
        return formatted_date

    def get_timezone(self,) -> str:
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

    def get_system_info(self) -> dict:
        system_info = {}
        if self._software == '':
            print_red('没有给出SOFTWARE注册表文件!')
            sys.exit(-1)
        reg = Registry.Registry(self._software)
        root_key = reg.open(r'Microsoft\Windows NT\CurrentVersion')
        system_info.update({'Build信息':root_key.value('BuildLabEx').value()})
        system_info.update({'版本信息':root_key.value('CurrentBuild').value()})
        system_info.update({'安装时间(本地时区)':self.timestamp(root_key.value('InstallDate').value())})
        system_info.update({'系统名称':root_key.value('ProductName').value()})
        system_info.update({'发行ID':root_key.value('ReleaseID').value()})
        system_info.update({'产品ID':root_key.value('ProductID').value()})
        system_info.update({'系统Build信息':root_key.value('BuildLabEx').value()})
        product_key = reg.open(r'Microsoft\Windows NT\CurrentVersion\SoftwareProtectionPlatform')
        system_info.update({'产品密钥备份(非当前密钥)':product_key.value('BackupProductKeyDefault').value()})
        system_info.update({'系统时区':self.get_timezone()})
        system_info.update({'上次正常关机时间(本地时区)':self.get_last_shutdown_time()})
        return system_info

if __name__ == '__main__':
    winTool = WinTool({'SYSTEM':r"C:\Users\Administrator\Desktop\workspace\source\SYSTEM",
        'NTUSER':r"C:\Users\Administrator\Desktop\workspace\source\NTUSER.DAT",
        'SAM':r"C:\Users\Administrator\Desktop\workspace\source\SAM",
        'SOFTWARE':r"C:\Users\Administrator\Desktop\workspace\source\SOFTWARE",
    })
    print_table(winTool.get_system_info(),title='系统信息')