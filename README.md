# SecureDNS
DOH service for local PC w/ Cloudflare.
Cloudflare의 DOH 서비스를 사용하여 로컬 PC에 안전한 DNS 연결을 제공합니다.

 * https://developers.cloudflare.com/1.1.1.1/dns-over-https/

현 버전은 DNS 서버 기능만 제공하며, PC의 DNS 설정은 수동으로 변경해 주어야 합니다.
네트워크 어댑터 설정에서 IPv4의 DNS 주소를 "127.0.0.1"로 변경하여 DOH를 사용할 수 있습니다.
프로그램이 실행되어 있는 동안에만 작동하므로 종료 후에는 반드시 DNS 주소를 본래대로 되돌려 주어야 합니다.

# Todo
 1. PC의 DNS 설정을 자동으로 변경하는 기능
 1. UI 개선
 1. DNS Wireformat을 대신하여 JSON 적용 (Google DOH 지원)

# 기타
 * 아이콘 출처 : http://www.iconarchive.com/show/ios7-icons-by-icons8/Network-Security-Checked-icon.html
