.DEFAULT_GOAL := build

build: ant ant-rl batch-gen

ant:
	go build -o bin ./cmd/ant

ant-rl:
	go build -o bin ./cmd/ant-rl

batch-gen:
	go build -o bin ./cmd/batch-gen

bench-prep:
	git stash
	make ant
	cp ./bin/ant ./bin/ant-old
	git stash pop
	make ant

#WIDTH = 16384 # default
WIDTH = 1024
bench:
	make bench-prep
	hyperfine --warmup 1 \
		'./bin/ant-old -w $(WIDTH) LRR__0.246000__A7386-B5868__50_000_001' \
		'./bin/ant     -w $(WIDTH) LRR__0.246000__A7386-B5868__50_000_002'

bench-huge:
	make bench-prep
	hyperfine -r 2 \
		'./bin/ant-old -w $(WIDTH) RLL__0.000007__B15160-E10890__500_000_001' \
		'./bin/ant     -w $(WIDTH) RLL__0.000007__B15160-E10890__500_000_002'

STEPS = 50_000_001
compare:
	make bench-prep
	time ./bin/ant-old LRR__0.246000__A7386-B5868__$(STEPS)
	mv        results5/LRR__0.246000__A7386-B5868__$(STEPS).png results5/old.png
	time ./bin/ant     LRR__0.246000__A7386-B5868__$(STEPS)
	mv        results5/LRR__0.246000__A7386-B5868__$(STEPS).png results5/new.png
	open results5/old.png
	open results5/new.png

bench-ant:
	make ant
	hyperfine --warmup 1 './bin/ant LRR__0.246000__A7386-B5868__25_000_002'

prof-ant:
	make ant
	./bin/ant -cpuprofile tmp/ant.prof LRR__0.246000__A7386-B5868__250_000_009
	go tool pprof -http=: -no_browser tmp/ant.prof

prof-ant-mem:
	make ant
	./bin/ant -memprofile tmp/ant-mem.prof RLL__0.000001__B-14917-A-8917__300_000_009
	go tool pprof -http=: -no_browser tmp/ant-mem.prof
