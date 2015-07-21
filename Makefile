RAGEL_FLAGS = -Z -G2

parser.go:
	ragel $(RAGEL_FLAGS) ragel/parser.rl -o $@
