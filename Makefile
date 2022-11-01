filenames := DataNode/Synth/DATA.txt DataNode/Crematore/DATA.txt DataNode/Grunt/DATA.txt NameNode/DATA.txt

files := $(strip $(foreach f,$(filenames),$(wildcard $(f))))

all: $(filenames)

$(filenames):
	touch $@

clean:
ifneq ($(files),)
	rm -f $(files)
endif

65:
	clear & go run maq65.go
66:
	clear & go run maq66.go
67:
	clear & go run maq67.go
68:
	clear & go run maq68.go