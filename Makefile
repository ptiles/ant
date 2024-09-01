.DEFAULT_GOAL := build

.PHONY: ant ant-rl batch-gen build

ant:
	go build -o bin ./cmd/ant

bench-prep:
	git stash
	make ant
	cp ./bin/ant ./bin/ant-old
	git stash pop
	make ant

bench:
	make bench-prep
	hyperfine --warmup 1 './bin/ant-old RLLRL.A-644+E633.50000001' './bin/ant RLLRL.A-644+E633.50000002'

compare:
	make bench-prep
	./bin/ant-old RLLRL.A-644+E633.50000001
	./bin/ant     RLLRL.A-644+E633.50000002
	open results5/RLLRL.A-644+E633.50000001.png
	open results5/RLLRL.A-644+E633.50000002.png

bench-ant:
	make ant
	hyperfine --warmup 1 './bin/ant RLLR.E-441-C393.25000002'

prof-ant:
	make ant
	./bin/ant -cpuprofile tmp/ant.prof RLLR.E-441-C393.250000009
	go tool pprof -http=: -no_browser tmp/ant.prof

ant-rl:
	go build -o bin ./cmd/ant-rl

batch-gen:
	go build -o bin ./cmd/batch-gen

build: ant ant-rl batch-gen
