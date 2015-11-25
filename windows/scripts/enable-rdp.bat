::
:: http://blogs.technet.com/b/bruce_adamczak/archive/2013/02/12/windows-2012-core-survival-guide-remote-desktop.aspx -> in powershell
::
:: enable RDP:http://technet.microsoft.com/en-us/library/cc782195(v=ws.10).aspx
%SystemRoot%\System32\reg.exe ADD "HKLM\SYSTEM\CurrentControlSet\Control\Terminal Server" /v fDenyTSConnections /t REG_DWORD /d 0 /f
:: Allow "insecure" connections
%SystemRoot%\System32\reg.exe ADD "HKLM\System\CurrentControlSet\Control\Terminal Server\WinStations\RDP-Tcp" /v UserAuthentication /t REG_DWORD /d 0 /f

:: Firewall
%SystemRoot%\SysWOW64\netsh advfirewall firewall set rule group="remote desktop" new enable=Yes
%SystemRoot%\SysWOW64\netsh advfirewall firewall add rule name="ALL ICMP V4" dir=in action=allow protocol=icmpv4

