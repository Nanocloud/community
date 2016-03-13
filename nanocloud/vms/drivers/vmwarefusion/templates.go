package vmwarefusion

const vmxTemplate = `.encoding = "UTF-8"
config.version = "8"
displayName = "{{.Name}}"
ethernet0.present = "TRUE"
ethernet0.connectionType = "nat"
ethernet0.virtualDev = "e1000e"
ethernet0.wakeOnPcktRcv = "FALSE"
ethernet0.addressType = "generated"
ethernet0.linkStatePropagation.enable = "TRUE"
pciBridge0.present = "TRUE"
pciBridge4.present = "TRUE"
pciBridge4.virtualDev = "pcieRootPort"
pciBridge4.functions = "8"
pciBridge5.present = "TRUE"
pciBridge5.virtualDev = "pcieRootPort"
pciBridge5.functions = "8"
pciBridge6.present = "TRUE"
pciBridge6.virtualDev = "pcieRootPort"
pciBridge6.functions = "8"
pciBridge7.present = "TRUE"
pciBridge7.virtualDev = "pcieRootPort"
pciBridge7.functions = "8"
pciBridge0.pciSlotNumber = "17"
pciBridge4.pciSlotNumber = "21"
pciBridge5.pciSlotNumber = "22"
pciBridge6.pciSlotNumber = "23"
pciBridge7.pciSlotNumber = "24"
scsi0.pciSlotNumber = "160"
usb.pciSlotNumber = "-1"
ethernet0.pciSlotNumber = "192"
sound.pciSlotNumber = "-1"
vmci0.pciSlotNumber = "35"
sata0.pciSlotNumber = "36"
guestOS = "windows8srv-64"
hpet0.present = "TRUE"
{{ if .WindowsInstallISO }}
sata0.present = "TRUE"
sata0:1.present = "TRUE"
sata0:1.fileName = "{{.WindowsInstallISO}}"
sata0:1.deviceType = "cdrom-image"
{{end}}
{{ if .NanocloudInstallISO }}
sata0:2.present = "TRUE"
sata0:2.fileName = "{{.NanocloudInstallISO}}"
sata0:2.deviceType = "cdrom-image"
{{ end }}
vmci0.present = "TRUE"
mem.hotadd = "TRUE"
memsize = "{{.RAM}}"
powerType.powerOff = "soft"
powerType.powerOn = "soft"
powerType.reset = "soft"
powerType.suspend = "soft"
scsi0.present = "TRUE"
scsi0.virtualDev = "lsisas1068"
scsi0:0.fileName = "{{.WMDKHardDrive}}"
scsi0:0.present = "TRUE"
tools.synctime = "TRUE"
virtualHW.productCompatibility = "hosted"
virtualHW.version = "12"
msg.autoanswer = "TRUE"
uuid.action = "create"
numvcpus = "{{.CPU}}"
hgfs.mapRootShare = "FALSE"
hgfs.linkRootShare = "FALSE"
numa.autosize.vcpu.maxPerVirtualNode = "2"
numa.autosize.cookie = "20001"
uuid.bios = "56 4d 36 d7 33 af 02 3c-6b 8a 6f 0c 68 eb 49 3b"
uuid.location = "56 4d 36 d7 33 af 02 3c-6b 8a 6f 0c 68 eb 49 3b"
scsi0:0.redo = ""
scsi0.sasWWID = "50 05 05 67 33 af 02 30"
vmci0.id = "1760250171"
vm.genid = "-1033732349146259583"
vm.genidX = "-5789930034112383728"
monitor.phys_bits_used = "42"
vmotion.checkpointFBSize = "100663296"
cleanShutdown = "TRUE"
softPowerOff = "FALSE"
floppy0.present = "FALSE"
acpi.smbiosVersion2.7 = "FALSE"
acpi.mouseVMW0003 = "FALSE"
vcpu.hotadd = "TRUE"
vmotion.checkpointSVGAPrimarySize = "100663296"
bios.bootOrder = "CDROM"`

const autounattend = `<?xml version="1.0" encoding="utf-8"?>
<unattend xmlns="urn:schemas-microsoft-com:unattend">
	<settings pass="windowsPE">
		<component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
			<DiskConfiguration>
				<Disk wcm:action="add">
					<CreatePartitions>
						<CreatePartition wcm:action="add">
							<Type>Primary</Type>
							<Order>1</Order>
							<Size>350</Size>
						</CreatePartition>
						<CreatePartition wcm:action="add">
							<Order>2</Order>
							<Type>Primary</Type>
							<Extend>true</Extend>
						</CreatePartition>
					</CreatePartitions>
					<ModifyPartitions>
						<ModifyPartition wcm:action="add">
							<Active>true</Active>
							<Format>NTFS</Format>
							<Label>boot</Label>
							<Order>1</Order>
							<PartitionID>1</PartitionID>
						</ModifyPartition>
						<ModifyPartition wcm:action="add">
							<Format>NTFS</Format>
							<Label>Windows 2012 R2</Label>
							<Letter>C</Letter>
							<Order>2</Order>
							<PartitionID>2</PartitionID>
						</ModifyPartition>
					</ModifyPartitions>
					<DiskID>0</DiskID>
					<WillWipeDisk>true</WillWipeDisk>
				</Disk>
			</DiskConfiguration>
			<ImageInstall>
				<OSImage>
					<InstallFrom>
						<MetaData wcm:action="add">
							<Key>/IMAGE/NAME </Key>
							<Value>Windows Server 2012 R2 SERVERSTANDARD</Value>
						</MetaData>
					</InstallFrom>
					<InstallTo>
						<DiskID>0</DiskID>
						<PartitionID>2</PartitionID>
					</InstallTo>
				</OSImage>
			</ImageInstall>

			<UserData>
				<!-- Product Key from http://technet.microsoft.com/en-us/library/jj612867.aspx -->
				<ProductKey>
					<!-- Do not uncomment the Key element if you are using trial ISOs -->
					<!-- You must uncomment the Key element (and optionally insert your own key) if you are using retail or volume license ISOs -->
					<!-- <Key></Key> -->
					<WillShowUI>OnError</WillShowUI>
				</ProductKey>
				<AcceptEula>true</AcceptEula>
				<FullName>Admin</FullName>
				<Organization>Admin</Organization>
			</UserData>
		</component>

		<component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-International-Core-WinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
			<SetupUILanguage>
				<UILanguage>en-US</UILanguage>
			</SetupUILanguage>
			<InputLocale>en-US</InputLocale>
			<SystemLocale>en-US</SystemLocale>
			<UILanguage>en-US</UILanguage>
			<UILanguageFallback>en-US</UILanguageFallback>
			<UserLocale>en-US</UserLocale>
		</component>

	</settings>
	<settings pass="specialize">

		<component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
			<OEMInformation>
				<HelpCustomized>false</HelpCustomized>
			</OEMInformation>
			<ComputerName>{{.Hostname}}</ComputerName>
			<TimeZone>{{.TimeZone}}</TimeZone>
			<RegisteredOwner/>
		</component>

		<!-- disable server manager auto-start -->
		<component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-ServerManager-SvrMgrNc" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
			<DoNotOpenServerManagerAtLogon>true</DoNotOpenServerManagerAtLogon>
		</component>

		<!-- disable annoying IE security  -->
		<component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-IE-ESC" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
			<IEHardenAdmin>false</IEHardenAdmin>
			<IEHardenUser>false</IEHardenUser>
		</component>
	</settings>

	<settings pass="oobeSystem">
		<component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
			<AutoLogon>
				<Password>
					<Value>{{.Password}}</Value>
					<PlainText>true</PlainText>
				</Password>
				<Enabled>true</Enabled>
				<Username>Administrator</Username>
			</AutoLogon>
			<FirstLogonCommands>
				<SynchronousCommand wcm:action="add">
					<CommandLine>E:\plaza.exe</CommandLine>
					<Description>Install Plaza</Description>
					<Order>1</Order>
					<RequiresUserInput>true</RequiresUserInput>
				</SynchronousCommand>
			</FirstLogonCommands>
			<OOBE>
				<HideEULAPage>true</HideEULAPage>
				<HideLocalAccountScreen>true</HideLocalAccountScreen>
				<HideOEMRegistrationScreen>true</HideOEMRegistrationScreen>
				<HideOnlineAccountScreens>true</HideOnlineAccountScreens>
				<HideWirelessSetupInOOBE>true</HideWirelessSetupInOOBE>
				<NetworkLocation>Home</NetworkLocation>
				<ProtectYourPC>1</ProtectYourPC>
			</OOBE>
			<UserAccounts>
				<AdministratorPassword>
					<Value>{{.Password}}</Value>
					<PlainText>true</PlainText>
				</AdministratorPassword>
			</UserAccounts>
			<RegisteredOwner/>
		</component>
	</settings>
</unattend>`
