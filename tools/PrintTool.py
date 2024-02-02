from rich.console import Console

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