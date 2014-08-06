#
#	Makefile for Go Language code
#
# --- CHANGE THESE FOR YOU SOURCE FILE NAME & OUPUT FILENAME (ie minus .go suffix)
SRC=csv2sql.go
OUTNAME=csv2sql
# Go compiler settings
CC=go
CFLAGS=build
#
# To build for Linux 32bit
LIN32=GOOS=linux GOARCH=amd64
# To build for Linux 64bit
LIN64=GOOS=linux GOARCH=amd64
# To build Windows 32 bit version:
WIN32=GOOS=windows GOARCH=386
# To build Windows 64 bit version:
WIN64=GOOS=windows GOARCH=amd64
# To build Mac OS X 32 bit version:
MAC32=GOOS=darwin GOARCH=386
# To build Mac OS X 64 bit version:
MAC64=GOOS=darwin GOARCH=amd64

$(OUTNAME): $(SRC)
	$(LIN64) $(CC) $(CFLAGS) -o $(OUTNAME) $(SRC)

lin32: $(SRC)
	$(LIN32) $(CC) $(CFLAGS) -o $(OUTNAME)-x386 $(SRC)

lin64: $(SRC)
	$(LIN64) $(CC) $(CFLAGS) -o $(OUTNAME) $(SRC)

win32: $(SRC)
	$(WIN32) $(CC) $(CFLAGS) -o $(OUTNAME)-x386.exe $(SRC)

win64: $(SRC)
	$(WIN64) $(CC) $(CFLAGS) -o $(OUTNAME)-x64.exe $(SRC)

mac32: $(SRC)
	$(MAC32) $(CC) $(CFLAGS) -o $(OUTNAME)-mac386 $(SRC)

mac64: $(SRC)
	$(MAC64) $(CC) $(CFLAGS) -o $(OUTNAME)-macx64 $(SRC)

clean:
	rm $(OUTNAME).exe $(OUTNAME)-x64.exe $(OUTNAME)-x386.exe $(OUTNAME) $(OUTNAME)-x386 $(OUTNAME)-macx64 $(OUTNAME)-mac386

all: lin64 win32 win64 mac64
