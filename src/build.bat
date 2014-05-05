@echo off
if exist src.exe (
del src.exe -f
)
set PWD=%cd%\..
echo Welcome to use GoMvc!
echo GoMvc  is a lightweight web framework ,QQ Group: 184572648
set GOPATH=%GOPATH%;%PWD%
echo Building ...
go build .
if exist src.exe (
 echo Build succeed!!!!
)
pause