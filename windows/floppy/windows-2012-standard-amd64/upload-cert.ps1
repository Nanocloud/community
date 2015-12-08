$url = "http://the.earth.li/~sgtatham/putty/latest/x86/pscp.exe"
$output = "c:\windows\system32\pscp.exe"
Invoke-WebRequest -Uri $url -OutFile $output
echo n | C:\Windows\System32\pscp.exe -q -pw "40chanesany" "C:\Users\Administrator\ad2012.cer" antoine@192.168.1.21:/home/antoine/conn