@echo off

del test.exe

go build -o test.exe .

test.exe
