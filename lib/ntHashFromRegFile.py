# ntHashFromRegFile.exe可能会报毒
# 我是使用nuitka编译的，因为impacket这个库正常安装是肯定会报毒的
# 所以想着编译一下不让他报毒，但是还是会报
# 介意的可以不使用密码哈希的功能，请在WinTool.py中将
# print_dict(winTool.get_nt_hash(),winTool.get_nt_hash()[0].keys(), title='密码信息')
# 注释掉

from __future__ import division
from __future__ import print_function
import logging

from impacket.examples.secretsdump import LocalOperations, SAMHashes

class DumpSecrets:
    def __init__(self, _system, _sam):
        self.__systemHive = _system
        self.__samHive = _sam

    def dump(self):
        try:
            localOperations = LocalOperations(self.__systemHive)
            bootKey = localOperations.getBootKey()
            SAMFileName = self.__samHive
            self.__SAMHashes = SAMHashes(SAMFileName, bootKey, isRemote=False)
            self.__SAMHashes.dump()
        except (Exception, KeyboardInterrupt) as e:
            if logging.getLogger().level == logging.DEBUG:
                import traceback
                traceback.print_exc()


# Process command-line arguments.
if __name__ == '__main__':
    import argparse
    arg_parser = argparse.ArgumentParser()
    arg_parser.add_argument('--system', help='System to dump', required=True)
    arg_parser.add_argument('--sam', help='SAM file to dump', required=True)
    args = arg_parser.parse_args()
    dumper = DumpSecrets(_system=args.system,_sam=args.sam)
    dumper.dump()