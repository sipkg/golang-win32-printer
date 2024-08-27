# golang-win32-printer
使用golang 语言封装win32 print API 支持打印图片,字符串，文件
# 开发目的
用于在windows 系统下 进行系统打印功能
最初想法是想为web 打印增加打印插件 以及获取系统信息,目前web打印功能太弱.

已封装函数
. DeleteDC
. CreateDC
. TextOut
. StartDoc
. EndDoc
. StartPage
. EndPage
. OpenPrinter
. ClosePrinter
. StartDocPrinter
. CloseDocPrinter
. StartPagePrinter
. ClosePagePrinter

功能正在完善中....欢迎参与改进
