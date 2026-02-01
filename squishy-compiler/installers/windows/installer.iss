#define appName       "squishy-compiler"
#define appVersion    "1.0.0"
#define appPublisher  "Cod2rDude"
#define appURL        "https://github.com/Cod2rDude/squishy"
#define appExeName    "squishy-compiler.exe"

[Setup]
AppId={{019C1ABF-C746-7C93-8638-4ADA309A9165}
AppName={#appName}
AppVersion={#appVersion}
AppPublisher={#appPublisher}
AppPublisherURL={#appURL}
AppSupportURL={#appURL}
AppUpdatesURL={#appURL}

ArchitecturesInstallIn64BitMode=x64 arm64

DefaultDirName={autopf}\{#appName}
DefaultGroupName={#appName}

WizardStyle=modern
DisableProgramGroupPage=no
LicenseFile=..\..\..\LICENSE
DisableWelcomePage=no

Compression=lzma2/ultra64
SolidCompression=yes
OutputBaseFilename=Squishy_Compiler_Setup_v{#appVersion}
OutputDir=dist

ChangesEnvironment=yes
PrivilegesRequired=admin

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "envPath"; Description: "Add to system PATH environment variable"; GroupDescription: "Additional icons"; Flags: checkedonce

[Files]
Source: "..\..\bin\squishy-compiler-windows-amd64.exe"; DestDir: "{app}"; DestName: "squishy.exe"; Check: IsX64; Flags: ignoreversion
Source: "..\..\bin\squishy-compiler-windows-arm64.exe"; DestDir: "{app}"; DestName: "squishy.exe"; Check: IsArm64; Flags: ignoreversion

[Icons]
Name: "{group}\{#appName}"; Filename: "{app}\{#appExeName}"
Name: "{group}\Uninstall {#appName}"; Filename: "{uninstallexe}"

[Registry]
Root: HKCR; Subkey: ".squishy"; ValueType: string; ValueName: ""; ValueData: "SquishyFile"; Flags: uninsdeletevalue
Root: HKCR; Subkey: ".sqy"; ValueType: string; ValueName: ""; ValueData: "SquishyFile"; Flags: uninsdeletevalue
Root: HKLM; Subkey: "SYSTEM\CurrentControlSet\Control\Session Manager\Environment"; \
    ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}"; \
    Tasks: envPath; Check: NeedsAddPath(ExpandConstant('{app}'))
    
[Code]
function NeedsAddPath(Param: string): boolean;
var
  OrigPath: string;
begin
  if not RegQueryStringValue(HKLM,
    'SYSTEM\CurrentControlSet\Control\Session Manager\Environment',
    'Path', OrigPath)
  then begin
    Result := True;
    exit;
  end;
  Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;

[Run]
Filename: "https://github.com/Cod2rDude/squishy"; Description: "Visit Squishy Documentation for more info"; Flags: postinstall shellexec nowait