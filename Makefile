.DEFAULT_GOAL := build

.PHONY: ant ant-rl batch-gen build

ant:
	go build -o bin ./cmd/ant

bench:
	git stash
	make ant
	cp ./bin/ant ./bin/ant-old
	git stash pop
	make ant
	hyperfine --warmup 1 './bin/ant-old RLLRL.E-644+C633.25000001' './bin/ant RLLRL.E-644+C633.25000002'

bench-ant:
	make ant
	hyperfine --warmup 1 './bin/ant RLLR.E-441-C393.25000002'

ant-rl:
	go build -o bin ./cmd/ant-rl

batch-gen:
	go build -o bin ./cmd/batch-gen

build: ant ant-rl batch-gen
