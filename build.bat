@echo off
cls
echo off
echo.
echo ----------------------------GoLand 编译----------------------------------
echo. ----------------------------by ohxing-----------------------------------

echo.
echo ----------------------------按"a"键 选择 构建后继续执行---------------------
echo ----------------------------按"b"键 选择 构建后不执行-----------------------
echo ----------------------------按"q"键 选择 退出----------------------------
echo.
SET /P choice=请选择操作项:
IF /I '%Choice:~0,1%'=='a' GOTO a
IF /I '%Choice:~0,1%'=='b' GOTO b
IF /I '%Choice:~0,1%'=='q' exit

SET outpath=.\build\

:a
echo 开始编译
python hyperbole.py build -r

echo 开始运行
cd .\build\
hysteria-windows-amd64.exe panel -c server.yaml -l debug

echo. -----------------------------------by ohxing-----------------------------------
exit

:b
echo 开始编译
python hyperbole.py build -r

echo. -----------------------------------by ohxing-----------------------------------
exit
