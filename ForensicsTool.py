import sys

from PySide6.QtCore import QProcess
from PySide6.QtGui import QIcon
from PySide6.QtWidgets import QApplication, QMainWindow, QWidget, QHBoxLayout, QListWidget, QStackedWidget
from qt_material import apply_stylesheet

from ui.ui import UI

class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("ForensicsTool")
        self.setWindowIcon(QIcon("./icon.ico"))
        self.resize(800, 600)
        main_widget = QWidget()
        self.setCentralWidget(main_widget)

        main_layout = QHBoxLayout()
        main_widget.setLayout(main_layout)

        # 左侧功能列表
        self.sidebar_list = QListWidget()
        # self.sidebar_list.setFont(font)
        self.sidebar_list.currentRowChanged.connect(self.on_list_changed)
        # 右侧功能区
        self.stacked_widget = QStackedWidget()
        ui = UI()
        ui_widgets = ui.init_ui()
        for key, value in ui_widgets.items():
            self.sidebar_list.addItem(key)
            self.stacked_widget.addWidget(value)
        self.sidebar_list.setCurrentRow(0)
        self.sidebar_list.setMaximumWidth(100)
        main_layout.addWidget(self.sidebar_list)
        main_layout.addWidget(self.stacked_widget)

    def on_list_changed(self, current_row):
        self.stacked_widget.setCurrentIndex(current_row)


if __name__ == "__main__":
    app = QApplication(sys.argv)
    apply_stylesheet(app, theme='light_blue.xml')
    window = MainWindow()
    window.show()
    sys.exit(app.exec())
