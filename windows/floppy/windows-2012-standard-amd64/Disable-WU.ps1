New-Item HKLM:\SOFTWARE\Policies\Microsoft\Windows -Name WindowsUpdate
New-Item HKLM:\SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate -Name AU
New-ItemProperty HKLM:\SOFTWARE\Policies\Microsoft\Windows\WindowsUpdate\AU -Name NoAutoUpdate -Value 1