$url = "http://the.earth.li/~sgtatham/putty/latest/x86/pscp.exe"
$output = "c:\windows\system32\pscp.exe"
Invoke-WebRequest -Uri $url -OutFile $output
echo n | C:\Windows\System32\pscp.exe -q -pw PASSWORD "C:\Users\Administrator\ad2012.cer" USER@HOST:PATH
