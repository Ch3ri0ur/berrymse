#Check if Windows or Linux
ifeq ($(OS),Windows_NT)
   export RM = del /Q
   export FixPath = $(subst /,\,$1)
   export WinMode = 1
else
	export RM = rm -f
	export FixPath = $1
	export WinMode = 0
endif

all:

	@echo $(WinMode)
	@echo $(value WinMode)

all: armv6l armv7l

ifeq ($(WinMode),1)
armv6l: pkged.go
	mkdir armv6l
	set GOARCH=arm
	set GOARM=6
	set GOOS=linux
	go build -v -o ./armv6l/berrymse
armv7l: pkged.go
	mkdir armv7l
	set GOARCH=arm
	set GOARM=7
	set GOOS=linux
	go build -v -o ./armv7l/berrymse
else
armv6l: pkged.go
	mkdir -p armv6l
	GOARCH=arm GOARM=7 GOOS=linux go build -v -o ./armv6l/berrymse -ldflags="-w -s"
armv7l: pkged.go
	mkdir -p armv7l
	GOARCH=arm GOARM=7 GOOS=linux go build -v -o ./armv7l/berrymse -ldflags="-w -s"
endif

pkged.go: web/**
	pkger

clean:
	$(RM) $(call FixPath,objs/*)

#missing didnt knew how to convert
#clean:
#	rm -rf armv6l armv7l pkged.go

#.PHONY: armv6l armv7l
#.PHONY: clean