ifeq ($(OS),Windows_NT)
	BIN_EXT = .exe
else
endif

binary:
	go build

default: binary
	./fake$(BIN_EXT)

simple: binary
	./fake$(BIN_EXT) \
		--students 10 \
		--judges 4 \
		--rooms 1 \
		--divisions 5 \
		--events 3 \
		--judge-talent 1
