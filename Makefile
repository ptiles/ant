.DEFAULT_GOAL := build

build: ant ant-batch ant-dry ant-rl

ant:
	go build -o bin ./cmd/ant

ant-batch:
	go build -o bin ./cmd/ant-batch

ant-dry:
	go build -o bin ./cmd/ant-dry

ant-rl:
	go build -o bin ./cmd/ant-rl

bench-prep-swiss:
	go1.23.6 build -o bin/ant-1.23.6 ./cmd/ant
	GOEXPERIMENT=noswissmap go1.24.0 build -o bin/ant-1.24.0-map ./cmd/ant
	GOEXPERIMENT=swissmap   go1.24.0 build -o bin/ant-1.24.0-swi ./cmd/ant

bench-swiss:
	make bench-prep-swiss

	hyperfine -i --warmup 1 -r 5 \
		'./bin/ant-1.23.6     -w $(WIDTH) RLL__0.000007__B15160-E10890__500_000_001' \
		'./bin/ant-1.24.0-map -w $(WIDTH) RLL__0.000007__B15160-E10890__500_000_002' \
		'./bin/ant-1.24.0-swi -w $(WIDTH) RLL__0.000007__B15160-E10890__500_000_003' \
		# end

bench-swiss-mem:
	make bench-prep-swiss

	time -lh ./bin/ant-1.23.6     -w $(WIDTH) RLL__0.000007__B15160-E10890__2_000_000_001
	time -lh ./bin/ant-1.24.0-map -w $(WIDTH) RLL__0.000007__B15160-E10890__2_000_000_002
	time -lh ./bin/ant-1.24.0-swi -w $(WIDTH) RLL__0.000007__B15160-E10890__2_000_000_003

prof-ant-swiss:
	make bench-prep-swiss

	./bin/ant-1.24.0-map -cpuprofile tmp/ant.prof RLL__0.000007__B15160-E10890__140_000_009
	go tool pprof -http=: -no_browser tmp/ant.prof

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
	hyperfine -i --warmup 1 \
		'./bin/ant     -w $(WIDTH) RLL__0.000007__B15160-E10890__50_000_002' \
		'./bin/ant-old -w $(WIDTH) RLL__0.000007__B15160-E10890__50_000_001' \
		# end

bench-fast:
	make bench-prep
	hyperfine -i -r 5 \
		'./bin/ant     -w $(WIDTH) RLL__0.000007__B15160-E10890__50_000_002' \
		'./bin/ant-old -w $(WIDTH) RLL__0.000007__B15160-E10890__50_000_001' \
		# end

bench-huge:
	make bench-prep
	hyperfine -i -r 2 \
		'./bin/ant     -w $(WIDTH) RLL__0.000007__B15160-E10890__1_500_000_002' \
		'./bin/ant-old -w $(WIDTH) RLL__0.000007__B15160-E10890__1_500_000_001' \
        # end

bench-mem:
	make bench-prep

	time -lh ./bin/ant     -w $(WIDTH) RLL__0.000007__B15160-E10890__200_000_002
	time -lh ./bin/ant-old -w $(WIDTH) RLL__0.000007__B15160-E10890__200_000_001

STEPS = 50_000_001
compare:
	make bench-prep
	time -lh ./bin/ant-old RLL__0.000007__B15160-E10890__$(STEPS)
	mv            results5/RLL__0.000007__B15160-E10890__$(STEPS).png results5/old.png
	time -lh ./bin/ant     RLL__0.000007__B15160-E10890__$(STEPS)
	mv            results5/RLL__0.000007__B15160-E10890__$(STEPS).png results5/new.png
	open results5/old.png
	open results5/new.png

bench-ant:
	make ant
	hyperfine -i --warmup 1 './bin/ant RLL__0.000007__B15160-E10890__25_000_002'

prof-ant:
	make ant
	./bin/ant -cpuprofile tmp/ant.prof RLL__0.000007__B15160-E10890__140_000_009
	go tool pprof -http=: -no_browser tmp/ant.prof

prof-ant-mem:
	make ant
	./bin/ant -memprofile tmp/ant-mem.prof RLL__0.000007__B15160-E10890__300_000_009
	go tool pprof -http=: -no_browser tmp/ant-mem.prof
