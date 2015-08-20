RAGEL_FLAGS = -Z -G2

all: parser.go

%.go: ragel/%.rl
	ragel $(RAGEL_FLAGS) $< -o $@
