go build  -ldflags="-w -s"
upx easypos.exe
go-msi make --msi easypos.msi --version 0.0.2