Unicode true
!include "MUI2.nsh"
!include "x64.nsh"

!define APPNAME "MehdiDNSFinder"
!define APPVER "1.0.0"
Name "${APPNAME}"
OutFile "MehdiDNSFinder-Setup.exe"
InstallDir "$PROGRAMFILES64\MehdiDNSFinder"
InstallDirRegKey HKLM "Software\MehdiDNSFinder" "InstallDir"
RequestExecutionLevel admin
SetCompressor /SOLID lzma

!define MUI_ICON "icon.ico"
!define MUI_UNICON "icon.ico"
!define MUI_ABORTWARNING
!define MUI_FINISHPAGE_RUN "$INSTDIR\MehdiDNSFinder.exe"
!define MUI_FINISHPAGE_RUN_TEXT "Run MehdiDNSFinder"
!insertmacro MUI_PAGE_WELCOME
!insertmacro MUI_PAGE_DIRECTORY
!insertmacro MUI_PAGE_INSTFILES
!insertmacro MUI_PAGE_FINISH
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES
!insertmacro MUI_LANGUAGE "Farsi"
!insertmacro MUI_LANGUAGE "English"

Section "Install"
  SetOutPath "$INSTDIR"
  File "MehdiDNSFinder.exe"
  WriteUninstaller "$INSTDIR\Uninstall.exe"
  CreateDirectory "$SMPROGRAMS\${APPNAME}"
  CreateShortcut "$SMPROGRAMS\${APPNAME}\${APPNAME}.lnk" "$INSTDIR\MehdiDNSFinder.exe"
  CreateShortcut "$SMPROGRAMS\${APPNAME}\Uninstall.lnk" "$INSTDIR\Uninstall.exe"
  CreateShortcut "$DESKTOP\${APPNAME}.lnk" "$INSTDIR\MehdiDNSFinder.exe"
  WriteRegStr HKLM "Software\MehdiDNSFinder" "InstallDir" "$INSTDIR"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MehdiDNSFinder" "DisplayName" "${APPNAME}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MehdiDNSFinder" "DisplayVersion" "${APPVER}"
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MehdiDNSFinder" "UninstallString" '"$INSTDIR\Uninstall.exe"'
  WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MehdiDNSFinder" "DisplayIcon" "$INSTDIR\MehdiDNSFinder.exe"
  WriteRegDWORD HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MehdiDNSFinder" "NoModify" 1
SectionEnd

Section "Uninstall"
  Delete "$INSTDIR\MehdiDNSFinder.exe"
  Delete "$INSTDIR\Uninstall.exe"
  RMDir "$INSTDIR"
  Delete "$DESKTOP\${APPNAME}.lnk"
  RMDir /r "$SMPROGRAMS\${APPNAME}"
  DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\MehdiDNSFinder"
  DeleteRegKey HKLM "Software\MehdiDNSFinder"
SectionEnd
