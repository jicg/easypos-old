cd /d %~dp0
nssm install easypos  %cd%\easypos.exe
nssm set easypos Application %cd%\easypos.exe
nssm set easypos AppDirectory  %cd%
nssm set easypos AppParameters web --port 3000
nssm set easypos DisplayName easypos
nssm set easypos Description easypos
nssm set easypos Start SERVICE_AUTO_START
nssm set easypos AppStdout  %cd%\log\service.log
nssm set easypos AppStderr  %cd%\log\service.log
nssm set easypos AppRotateFiles 1
nssm set easypos AppRotateBytes 1048576
nssm start easypos