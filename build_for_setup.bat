@echo off
echo Build SecureDNS...

setlocal

echo 32bit windows...
set GOOS=windows
set GOARCH=386

go build -o setup\SecureDNS32.exe

echo done.
echo.

echo 64 bit windows...
set GOOS=windows
set GOARCH=amd64

go build -o setup\SecureDNS64.exe

echo done.

endlocal
