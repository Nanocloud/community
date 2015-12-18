mkdir C:\cygwin
Invoke-WebRequest https://community.nanocloud.com/cygwin/setup-x86_64.exe -OutFile C:\cygwin\cygwin-setup-x86_64.exe
C:\cygwin\cygwin-setup-x86_64.exe -a x86_64 -q -R C:\cygwin -P openssh -s http://mirrors.chauf.net/cygwin