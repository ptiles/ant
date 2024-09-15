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

#WIDTH = 16384 # default
WIDTH = 1024
bench:
	make bench-prep
	hyperfine --warmup 1 \
		'./bin/ant-old -w $(WIDTH) LRR__0.246000__A7386-B5868__50_000_001' \
		'./bin/ant     -w $(WIDTH) LRR__0.246000__A7386-B5868__50_000_002'

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

ant-rl:
	go build -o bin ./cmd/ant-rl

batch-gen:
	go build -o bin ./cmd/batch-gen

build: ant ant-rl batch-gen
