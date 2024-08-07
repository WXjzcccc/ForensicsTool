import tkinter as tk
from tkinter import ttk
from tkinter import messagebox
from tkinterdnd2 import TkinterDnD, DND_FILES
import subprocess
import threading
import json
import re
import os

class CommandLineApp:
    def __init__(self, root):
        self.config = self.get_config()
        self.root = root
        self.root.title("ForensicsTool GUI WorkBench")
        
        self.center_window(800, 600)  # 设置窗口宽度和高度

        self.params_frame = ttk.Frame(root, padding="10")
        self.params_frame.grid(row=0, column=0, sticky=(tk.W, tk.E))
        for i in range(4):
            self.params_frame.grid_columnconfigure(i, weight=1)
        # 模式选择下拉框
        self.mode_label = ttk.Label(self.params_frame, text="模式")
        self.mode_label.grid(row=0, column=0, padx=5, pady=5,sticky='e')
        self.mode_var = tk.StringVar(value="请选择")
        self.mode_combobox = ttk.Combobox(self.params_frame, textvariable=self.mode_var, state="readonly")
        self.mode_combobox['values'] = tuple(self.config['mode'])
        self.mode_combobox.grid(row=0, column=1, padx=5, pady=5,columnspan=3,sticky='w')
        # 任务选择下拉框
        self.type_label = ttk.Label(self.params_frame, text="任务")
        self.type_label.grid(row=1, column=0, padx=5, pady=5,sticky='e')
        self.type_var = tk.StringVar(value="请选择")
        self.type_combobox = ttk.Combobox(self.params_frame, textvariable=self.type_var, state="readonly")
        self.type_combobox['values'] = ('请先选择模式')
        self.type_combobox.grid(row=1, column=1, padx=5, pady=5,columnspan=3,sticky='w')
        # 文件路径输入框
        self.file_label = ttk.Label(self.params_frame, text="file")
        self.file_label.grid(row=2, column=0, padx=5, pady=5,sticky='e')
        self.file_var = tk.StringVar(value='')
        self.file_entry = self.create_entry_with_placeholder(self.file_var, "拖入文件或文件夹")
        self.file_entry.grid(row=2, column=1, padx=5, columnspan=3, pady=5,sticky='w')
        # 批量创建除file外的参数
        self.entrys = self.create_entrys(3)
        # 执行按钮
        self.run_button = ttk.Button(root, text="Run", command=self.execute_command)
        self.run_button.grid(row=1, column=0, pady=10)
        # 输出框
        self.output_text = tk.Text(root, wrap=tk.WORD, state='disabled')
        self.output_text.grid(row=2, column=0, padx=10, pady=10, sticky=(tk.W, tk.E, tk.N, tk.S))

        # 绑定窗口大小改变事件，以动态调整 Combobox 宽度
        self.root.bind('<Configure>', lambda event: self.adjust_height_width())
        # 绑定模式选择事件
        self.mode_combobox.bind("<<ComboboxSelected>>", self.on_mode_selected)
        # 绑定文件拖放事件
        self.file_entry.drop_target_register(DND_FILES)
        self.file_entry.dnd_bind('<<Drop>>', self.drop)

    def create_entrys(self,row) -> dict:
        '''批量创建参数输入框，并返回一个字典entry'''
        params = self.config['params']
        entrys = {}
        col = 0
        for param in params:
            # 每4个元素就换行
            if col == 4:
                row += 1
                col = 0
            label = ttk.Label(self.params_frame, text=param)
            label.grid(row=row, column=col, padx=5, pady=5,sticky='e')
            var = tk.StringVar(value='')
            entry = self.create_entry_with_placeholder(var,params[param])
            entry.grid(row=row, column=col+1, padx=5, pady=5,sticky='w')
            entrys.update({
                param:{
                    'var':var,
                    'entry':entry
                }
            })
            entry = None
            col += 2
        return entrys


    def get_param(self,param:str) -> str:
        '''获取输入的参数'''
        var = self.entrys[param]['var'].get()
        if var == self.config['params'][param]:
            return ''
        else:
            return var


    def create_entry_with_placeholder(self,str_var,placeholder_text):
        '''创建输入框'''
        entry_var = str_var

        def on_focus_in(event, entry, placeholder):
            if entry.get() == placeholder:
                entry.delete(0, tk.END)
                entry.config(foreground='black')

        def on_focus_out(event, entry, placeholder):
            if entry.get() == '':
                entry.insert(0, placeholder)
                entry.config(foreground='grey')

        entry = ttk.Entry(self.params_frame, textvariable=entry_var, foreground='grey')
        entry.insert(0, placeholder_text)
        # entry.grid(row=row, column=column, padx=5, pady=5,sticky=sticky)

        # 绑定 Entry 的 FocusIn 和 FocusOut 事件
        entry.bind("<FocusIn>", lambda event: on_focus_in(event, entry, placeholder_text))
        entry.bind("<FocusOut>", lambda event: on_focus_out(event, entry, placeholder_text))
        return entry

    def drop(self,event):
        '''获取拖放的文件路径并填充到输入框'''
        self.file_var.set(event.data)

    def on_mode_selected(self,event):
        '''在模式选中后，自动更新任务的可选值'''
        self.type_var.set('请选择')
        mode = self.mode_var.get()
        limit_type = self.config['limit'][self.get_number(mode)]
        self.type_combobox['values'] = tuple(self.get_type(limit_type))

    def get_type(self,lst :list) -> list:
        '''根据选中的模式，获取任务的可选值'''
        type_lst = []
        _type = self.config['type']
        for v in _type:
            if int(self.get_number(v)) in lst:
                type_lst.append(v)
        return type_lst

    def get_number(self,s) -> str:
        '''获取模式或任务的序号'''
        pattern = r'【([0-9]{1,2})】'
        ret = re.findall(pattern,s)
        return ret[0]

    def get_config(self) -> json:
        '''读取配置文件为json对象'''
        conf_file = 'config.conf'
        try:
            with open(conf_file,'r',encoding='utf8') as fr:
                data = fr.read()
                json_data = json.loads(data)
        except Exception as e:
            messagebox.showerror('错误',e)
            exit()
        return json_data
    
    def center_window(self, width, height):
        '''自动居中窗口'''
        screen_width = self.root.winfo_screenwidth()
        screen_height = self.root.winfo_screenheight()

        # 计算窗口的 x 和 y 位置
        x = (screen_width - width) // 2
        y = (screen_height - height) // 2

        # 设置窗口大小和位置
        self.root.geometry(f'{width}x{height}+{x}+{y}')

    def execute_command(self):
        '''拼接命令'''
        command = ['python']
        python_file = self.config['filePath']
        if not (os.path.exists(python_file) and os.path.isfile(python_file)):
            messagebox.showerror('错误','未找到ForensicsTool.py，请检查配置文件是否正确！')
        command.append(os.path.abspath(python_file))
        command.append('-m')
        mode = self.get_number(self.mode_var.get())
        command.append(mode)
        if mode != '3':
            command.append('-t')
            command.append(self.get_number(self.type_var.get()))
        file_var = self.file_var.get()
        if file_var != '拖入文件或文件夹':
            command.append('-f')
            command.append(file_var.replace('{','').replace('}',''))
        for entry in self.entrys:
            param = self.get_param(entry)
            if param != '':
                command.append(f'--{entry}')
                command.append(param)
        self.append_text(f"> {' '.join(command)}\n")
        thread = threading.Thread(target=self.run_command, args=(command,))
        thread.start()

    def run_command(self, command :list):
        '''执行命令'''
        # command = ['python',r'D:\Tools\workspace\dev\ForensicsTool\ForensicsTool.py','-m','2','-t','17','-f',r'E:\dd\20240719004\exports\Hawk2.xml','-p','rp6Fz0G8E6p7KzfJS0qg3Agyoek/GoRg12kAptH9VDo=']
        process = subprocess.Popen(command, shell=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
        stdout, stderr = process.communicate()

        if stdout:
            self.append_text(stdout)
        if stderr:
            self.append_text(stderr)

    def append_text(self, text):
        self.output_text.configure(state='normal')
        self.output_text.insert(tk.END, text)
        self.output_text.configure(state='disabled')
        self.output_text.yview(tk.END)

    def adjust_height_width(self):
        '''自动调整宽高'''
        # 获取窗口的当前宽度
        window_width = root.winfo_width()
        window_height = root.winfo_height()
        # 计算 Combobox 宽度（例如，窗口宽度的 80%）
        combobox_width = int(window_width * 0.12)
        # 设置 Combobox 宽度
        self.mode_combobox.config(width=combobox_width)
        self.type_combobox.config(width=combobox_width)
        self.file_entry.config(width=combobox_width+3)
        for entry in self.entrys.keys():
            self.entrys[entry]['entry'].config(width=combobox_width//2-5)
        self.output_text.config(width=combobox_width,height=window_height*0.035)

if __name__ == "__main__":
    root = TkinterDnD.Tk()
    app = CommandLineApp(root)
    root.mainloop()