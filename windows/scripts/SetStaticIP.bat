@echo off
C:\Windows\SysWOW64\netsh.exe interface ip set address "Ethernet 2" static 10.20.12.20 255.255.255.240 10.20.12.19
C:\Windows\SysWOW64\netsh.exe interface ip add dns "Ethernet 2" 8.8.8.8
C:\Windows\SysWOW64\netsh.exe interface ip add dns "Ethernet 2" 8.8.4.4 index=2

C:\Windows\SysWOW64\netsh advfirewall firewall set rule group="remote desktop" new enable=Yes
C:\Windows\SysWOW64\netsh advfirewall firewall add rule name="ALL ICMP V4" dir=in action=allow protocol=icmpv4

