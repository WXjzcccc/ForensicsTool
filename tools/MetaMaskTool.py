from rich.console import Console
from rich.table import Table
import json
import datetime

def print_metamask_table(dic_list: list, title: str = ''):
    console = Console()
    if dic_list:
        table = Table(show_header=True, header_style="bold green", title=f'<{title}>', title_style='yellow')
        
        headers = dic_list[0].keys()
        for header in headers:
            table.add_column(header, style="cyan", justify="center")

        # Add rows to the table
        for dic in dic_list:
            row = [str(dic.get(header, 'N/A')) for header in headers]
            table.add_row(*row)

        table.auto_width = True
        console.print(table)

def analyzeMetaMask(file_path):
    wallets = []
    contacts = []
    transactions = []
    browser_history = []

    with open(file_path, "r", encoding="utf-8") as file:
        data = json.load(file)

        if 'engine' in data:
            json_object = json.loads(data['engine'])
            json_browser = json.loads(data['browser'])

            # Extract wallet information
            for accs in json_object['backgroundState']['AccountTrackerController']['accounts']:
                s = json_object['backgroundState']['PreferencesController']['identities'][accs]['importTime'] / 1000.0
                import_time = datetime.datetime.fromtimestamp(s).strftime('%Y-%m-%d %H:%M:%S.%f')
                balance = int(
                    str(json_object['backgroundState']['AccountTrackerController']['accounts'][accs]['balance']), 16)
                wallet_info = {
                    'Import Time': import_time,
                    'Wallet Name': json_object['backgroundState']['PreferencesController']['identities'][accs]['name'],
                    'Wallet Address': accs,
                    'Balance (ETH)': balance / (10 ** 18)
                }
                wallets.append(wallet_info)

            # Extract contacts information
            for contactId in json_object['backgroundState']['AddressBookController']['addressBook']:
                for contact in json_object['backgroundState']['AddressBookController']['addressBook'][contactId]:
                    contact_info = {
                        'Contact Name':
                            json_object['backgroundState']['AddressBookController']['addressBook'][contactId][contact][
                                'name'],
                        'Wallet Address':
                            json_object['backgroundState']['AddressBookController']['addressBook'][contactId][contact][
                                'address']
                    }
                    contacts.append(contact_info)

            # Extract transactions information
            for transaction in json_object['backgroundState']['TransactionController']['transactions']:
                s = transaction['time'] / 1000.0
                transaction_time = datetime.datetime.fromtimestamp(s).strftime('%Y-%m-%d %H:%M:%S.%f')
                transaction_value = int(str(transaction['txParams']['value']), 16) / (10 ** 18)

                transaction_info = {
                    'Timestamp': transaction_time,
                    'From': transaction['txParams']['from'],
                    'To': transaction['txParams']['to'],
                    'Value (ETH)': transaction_value,
                    'Transaction Hash': transaction.get('hash'),
                    'Error': transaction.get('error')
                }
                transactions.append(transaction_info)

            # Extract browser history information
            for history in json_browser["history"]:
                browser_history_info = {
                    'Name': history["name"],
                    'URL': history["url"]
                }
                browser_history.append(browser_history_info)

    # Print the extracted information using print_metamask_table
    print_metamask_table(wallets, title='Wallets')
    print_metamask_table(contacts, title='Contacts')
    print_metamask_table(transactions, title='Transactions')
    print_metamask_table(browser_history, title='Browser History')
