@echo off
echo Build SecureDNS...

setlocal

echo 64 bit windows...
set GOOS=windows
set GOARCH=amd64

go build -o setup\SecureDNS.exe

echo done.

endlocal
