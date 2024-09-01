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
	hyperfine --warmup 1 './bin/ant-old RLLR.E-441-C393.25000000' './bin/ant RLLR.E-441-C393.25000000'

ant-rl:
	go build -o bin ./cmd/ant-rl

batch-gen:
	go build -o bin ./cmd/batch-gen

build: ant ant-rl batch-gen
