@ECHO OFF
PUSHD "%~dp0"

echo This script should be run with administrator privileges.
echo Right click - run as administrator.
echo Press any key if you're running it as administrator.
pause

sc stop "SecureDNS"
sc delete "SecureDNS"
sc create "SecureDNS" binPath= "\"%CD%\SecureDNS.exe\"" start= "auto" DisplayName= "Secure DNS"
sc description "SecureDNS" "DNS proxy service for DOH (DNS over HTTPS) support."
sc start "SecureDNS"
pause

POPD
