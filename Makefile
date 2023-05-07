#
#	Makefile for Go Language code
#
SRC=csv2sql.go 
OUTNAME=bin/csv2sql
# Go compiler settings
CC=go
CFLAGS=build -gcflags=all=-dwarf=false -ldflags="-s -w" -trimpath
RFLAGS=run
#
# To build for Linux 32bit ARM7
AARCH32=GOOS=linux GOARCH=arm
# To build for Linux 64bit ARM64
AARCH64=GOOS=linux GOARCH=arm64
# To build for Linux 32bit
LIN32=GOOS=linux GOARCH=386
# To build for Linux 64bit
LIN64=GOOS=linux GOARCH=amd64
# To build Windows 32 bit version:
WIN32=GOOS=windows GOARCH=386
# To build Windows 64 bit version:
WIN64=GOOS=windows GOARCH=amd64
# To build Windows amr64 bit version:
WINARM64=GOOS=windows GOARCH=arm64
# To build macOS 64 bit version:
MAC64=GOOS=darwin GOARCH=amd64
# To build macOS M1 or M2 arm64 version:
MACARM64=GOOS=darwin GOARCH=arm64
# To build FreeBSD 64 bit version:
FREE64=GOOS=freebsd GOARCH=amd64

aarch32: $(SRC)
	$(AARCH32) $(CC) $(CFLAGS) -o $(OUTNAME)-aarch32
aarch64: $(SRC)
	$(AARCH64) $(CC) $(CFLAGS) -o $(OUTNAME)-aarch64
lin32: $(SRC)
	$(LIN32) $(CC) $(CFLAGS) -o $(OUTNAME)-linux-x86
lin64: $(SRC)
	$(LIN64) $(CC) $(CFLAGS) -o $(OUTNAME)-linux-x64
win32: $(SRC)
	$(WIN32) $(CC) $(CFLAGS) -o $(OUTNAME)-windows-x86.exe
win64: $(SRC)
	$(WIN64) $(CC) $(CFLAGS) -o $(OUTNAME)-windows-x64.exe
winarm64: $(SRC)
	$(WINARM64) $(CC) $(CFLAGS) -o $(OUTNAME)-windows-arm64.exe
mac64: $(SRC)
	$(MAC64) $(CC) $(CFLAGS) -o $(OUTNAME)-mac-x64
macarm64: $(SRC)
	$(MACARM64) $(CC) $(CFLAGS) -o $(OUTNAME)-mac-arm64
free64: $(SRC)
	$(FREE64) $(CC) $(CFLAGS) -o $(OUTNAME)-freebsd64
run: $(SRC)
	$(CC) $(RFLAGS) $(SRC)

clean:
	rm $(OUTNAME)-aarch32 $(OUTNAME)-aarch64 $(OUTNAME)-linux-x86 $(OUTNAME)-linux-x64 $(OUTNAME)-windows-x64.exe $(OUTNAME)-windows-x86.exe $(OUTNAME)-windows-arm64.exe $(OUTNAME)-mac-x64 $(OUTNAME)-freebsd64 $(OUTNAME)-mac-arm64

all: aarch32 aarch64 lin32 lin64 win32 win64 winarm64 mac64 free64 macarm64
