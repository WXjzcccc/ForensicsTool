from rich.console import Console
from rich.table import Table

def print_green(s :str):
    console = Console()
    console.print(s, style="green")

def print_green_key(s :str,key :str):
    console = Console()
    console.print(s, style="green", end='')
    console.print(f'<{key}>', style="underline")

def print_red(s :str):
    console = Console()
    console.print(s, style="red")

def print_red_key(s :str,key :str):
    console = Console()
    console.print(s, style="red", end='')
    console.print(f'<{key}>', style="underline")

def print_yellow(s :str):
    console = Console()
    console.print(s, style="yellow")

def print_yellow_key(s :str,key :str):
    console = Console()
    print(s, end='')
    console.print(f'<{key}>', style="underline")

def print_dict(dics :list ,head :list ,title :str=''):
    console = Console()
    table = Table(show_header=True, header_style="bold green", title=f'<{title}>', title_style='yellow')
    for h in head:
        table.add_column(h, style="cyan" ,justify="center")
    for dic in dics:
        table.add_row(*[str(dic[column_name]) for column_name in dic.keys()])
    table.auto_width = True
    console.print(table)

def print_table(dic :dict ,title :str=''):
    console = Console()
    table = Table(show_header=True, header_style="bold green", title=f'<{title}>', title_style='yellow')
    table.add_column('名称', style="cyan" ,justify="center")
    table.add_column('值', style="cyan" ,justify="center")
    for k,v in dic.items():
        table.add_row(k,v)
    table.auto_width = True
    console.print(table)