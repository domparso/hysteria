@echo off
cls
echo off
echo.
echo ----------------------------GoLand ����----------------------------------
echo. ----------------------------by ohxing-----------------------------------

echo.
echo ----------------------------��"a"�� ѡ�� ���������ִ��---------------------
echo ----------------------------��"b"�� ѡ�� ������ִ��-----------------------
echo ----------------------------��"q"�� ѡ�� �˳�----------------------------
echo.
SET /P choice=��ѡ�������:
IF /I '%Choice:~0,1%'=='a' GOTO a
IF /I '%Choice:~0,1%'=='b' GOTO b
IF /I '%Choice:~0,1%'=='q' exit

SET outpath=.\build\

:a
echo ��ʼ����
python hyperbole.py build -r

echo ��ʼ����
cd .\build\
hysteria-windows-amd64.exe panel -c server.yaml -l debug

echo. -----------------------------------by ohxing-----------------------------------
exit

:b
echo ��ʼ����
python hyperbole.py build -r

echo. -----------------------------------by ohxing-----------------------------------
exit
