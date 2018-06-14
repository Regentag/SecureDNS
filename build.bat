rsrc -manifest comctrl6.manifest -ico icon.ico -o rsrc.syso
go build -ldflags="-s -w -H windowsgui"
