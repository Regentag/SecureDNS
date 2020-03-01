[Setup]
AppId={{B83116D9-F507-46E8-B59A-AF714A182295}
AppName=SecureDNS
AppVersion=1.2
AppPublisher=REGENTAG
AppPublisherURL=https://github.com/Regentag/SecureDNS
AppSupportURL=https://github.com/Regentag/SecureDNS
AppUpdatesURL=https://github.com/Regentag/SecureDNS/releases
DefaultDirName={code:GetProgramFiles}\SecureDNS
DefaultGroupName=SecureDNS
DisableProgramGroupPage=yes
InfoBeforeFile=setup_readme.ko_kr.rtf
OutputBaseFilename=securedns_setup
Compression=lzma
SolidCompression=yes
DisableDirPage=yes
AllowUNCPath=False
ShowLanguageDialog=no
AppContact=https://github.com/Regentag/SecureDNS/issues
UninstallDisplaySize=42
UninstallDisplayIcon={uninstallexe}
InfoAfterFile=setup_after.ko_kr.rtf

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Files]
Source: "SecureDNS32.exe"; DestDir: "{app}"; DestName: "SecureDNS.exe"; Flags: ignoreversion 32bit
Source: "SecureDNS64.exe"; DestDir: "{app}"; DestName: "SecureDNS.exe"; Flags: ignoreversion 64bit
Source: "sec-dns.log"; DestDir: "{app}"; Flags: ignoreversion
Source: "service_install.cmd"; DestDir: "{app}"; Flags: ignoreversion
Source: "service_remove.cmd"; DestDir: "{app}"; Flags: ignoreversion

[Run]
Filename: "{app}\service_install.cmd"; WorkingDir: "{app}"; Description: "Install service"

[UninstallRun]
Filename: "{app}\service_remove.cmd"; WorkingDir: "{app}"

[Code]
function GetProgramFiles(Param: string): string;
begin
  if IsWin64 then Result := ExpandConstant('{pf64}')
    else Result := ExpandConstant('{pf32}')
end;
