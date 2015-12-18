C:\cygwin\bin\bash -l -c "C:/install-cygwin-sshd.sh"
netsh advfirewall firewall add rule name="sshd" dir=in action=allow program="%SystemDrive%\cygwin\usr\sbin\sshd.exe" enable=yes
netsh advfirewall firewall add rule name="ssh" dir=in action=allow protocol=TCP localport=22
net start sshd