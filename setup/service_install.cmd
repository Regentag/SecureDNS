@ECHO OFF
PUSHD "%~dp0"

sc stop "SecureDNS"
sc delete "SecureDNS"
sc create "SecureDNS" binPath= "\"%CD%\SecureDNS.exe\"" start= "auto" DisplayName= "Secure DNS"
sc description "SecureDNS" "DNS proxy service for DOH (DNS over HTTPS) support."
sc start "SecureDNS"

POPD
