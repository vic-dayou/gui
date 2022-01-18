### golang 编译app添加图标
1.通过工具或是在线工具生成.ico的图标文件(假定是main.ico)

3.进入到项目的目录(执行go build的地方)

4.创建一个空白文本文件,命名main.rc

5.记录本打开,输入并保存 IDI_ICON1 ICON "main.ico"

6.在项目目录执行下面的命令 windres -o main.syso main.rc ,此时生成了一个main.syso

7.go build -ldfalgs -H="windowsgui" -o cpcn.exe .