#This Makefile supports the building of this project on Windows and Linux
# Using 'make' will build both versions for ARMV6 and ARMV7 architecture
# Using 'make armv6l' or 'make armv7l' will only build the selected Version
# 'make Clean' will removes all generated Files and Folders by this script!
#!!! BUILD Output will be generated in Folders 'armv6l' or 'armv7', if they exist they will be deleted first!!!!


################################################################################
#On default Build both versions for ARMV6 and for ARMV7 by performing methode armv6l and armv7l 
all: armv6l armv7l

#WINDOWS
################################################################################
ifeq ($(OS),Windows_NT) #IF Windows

#Build instruction for ARMV6 on Windows
# - Setting env.-variables
# - Creating folder armv6l if not exist
# - Let Go build executable and store it in folder armv6l (-w No DWARF debugging information, -s No generation of the Go symbol table)
armv6l: export GOARCH = arm
armv6l: export GOARM = 6
armv6l: export GOOS = linux
armv6l: pkged.go
	go get -v ./...
	if not exist armv6l mkdir armv6l
	go build -v -o ./armv6l/berrymse -ldflags="-w -s"

#Build instruction for ARMV7 on Windows
# - Setting env.-variables
# - Create folder armv7l if not exist
# - Let Go build executable and store it in folder armv7l (-w No DWARF debugging information, -s No generation of the Go symbol table)
armv7l: export GOARCH = arm
armv7l: export GOARM = 7
armv7l: export GOOS = linux
armv7l: pkged.go
	go get -v ./...
	if not exist armv7l mkdir armv7l
	go build -v -o ./armv7l/berrymse -ldflags="-w -s"


#Clean up for Windows
# - Check if folder or file exist and delete it with content without asking for premissions
clean:
	if exist armv6l rmdir /Q /S armv6l
	if exist armv7l rmdir /Q /S armv7l
	if exist pkged.go del /Q pkged.go




#LINUX
################################################################################
else #If Linux

#Build instruction for ARMV6 on Linux
# - Creating folder armv6l
# - Setting env.-variables and let Go build executable and store it in folder armv6l (-w No DWARF debugging information, -s No generation of the Go symbol table)
armv6l: pkged.go
	mkdir -p armv6l
	GOARCH=arm GOARM=7 GOOS=linux go build -v -o ./armv6l/berrymse -ldflags="-w -s"

#Build instruction for ARMV7  on Linux
# - Creating folder armv7l
# - Setting env.-variables
# - Setting env.-variables and let Go build executable and store it in folder armv7l (-w No DWARF debugging information, -s No generation of the Go symbol table)
armv7l: pkged.go
	mkdir -p armv7l
	GOARCH=arm GOARM=7 GOOS=linux go build -v -o ./armv7l/berrymse -ldflags="-w -s"


#Clean up for Linux
# - Delete all created files and folders
clean:
	rm -rf armv6l armv7l pkged.go


endif



################################################################################

#Methode pkger compresses the content of web/static
pkged.go: web/**
	pkger -o ./cmd/berryMSE

#Register Methodes avoid a conflict with a file or folder with the same name
.PHONY: armv6l armv7l
.PHONY: pkged.go
.PHONY: clean