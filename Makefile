.DEFAULT_GOAL := build

build: ant ant-batch ant-dry

ant:
	go build -o bin ./cmd/ant

ant-batch:
	go build -o bin ./cmd/ant-batch

ant-dry:
	go build -o bin ./cmd/ant-dry

bench-prep-swiss:
	go build -o bin ./cmd/ant
	go1.23.8 build -o bin/ant-1.23.8 ./cmd/ant
	GOEXPERIMENT=noswissmap go1.24.2 build -o bin/ant-1.24.2-old ./cmd/ant
	GOEXPERIMENT=swissmap   go1.24.2 build -o bin/ant-1.24.2-swi ./cmd/ant

#RECT = ''
#RECT = '(-423424,-430080)-(423424,375296)/128'
#RECT = '(-423424,-430080)-(423424,375296)/1024'
RECT = '(-423424,-430080)-(423424,375296)/2048'
ANT = RLL__7e-06__B15160-E10890

bench-swiss:
	make bench-prep-swiss

	hyperfine -i --warmup 1 -r 5 \
		"./bin/ant            -r $(RECT) $(ANT)__500_000_000" \
		"./bin/ant-1.23.8     -r $(RECT) $(ANT)__500_000_001" \
		"./bin/ant-1.24.2-old -r $(RECT) $(ANT)__500_000_002" \
		"./bin/ant-1.24.2-swi -r $(RECT) $(ANT)__500_000_003" \
		# end

bench-swiss-mem:
	make bench-prep-swiss

	time -lh ./bin/ant-1.23.8     -r $(RECT) $(ANT)__2_000_000_001
	time -lh ./bin/ant-1.24.2-old -r $(RECT) $(ANT)__2_000_000_002
	time -lh ./bin/ant-1.24.2-swi -r $(RECT) $(ANT)__2_000_000_003

prof-ant-swiss:
	make bench-prep-swiss

	./bin/ant-1.24.2-swi -cpuprofile tmp/ant-swi.prof $(ANT)__140_000_009
	go tool pprof -http=: -no_browser tmp/ant-swi.prof

bench-prep:
	git stash
	make ant
	cp ./bin/ant ./bin/ant-old
	git stash pop
	make ant

bench:
	make bench-prep
	hyperfine -i --warmup 1 \
		"./bin/ant     -r $(RECT) $(ANT)__50_000_002" \
		"./bin/ant-old -r $(RECT) $(ANT)__50_000_001" \
		# end

bench-fast:
	make bench-prep
	hyperfine -i -r 5 \
		"./bin/ant     -r $(RECT) $(ANT)__50_000_002" \
		"./bin/ant-old -r $(RECT) $(ANT)__50_000_001" \
		# end

bench-large:
	make bench-prep
	hyperfine -i -r 2 \
		"./bin/ant     -r $(RECT) $(ANT)__250_000_002" \
		"./bin/ant-old -r $(RECT) $(ANT)__250_000_001" \
		# end

bench-huge:
	make bench-prep
	hyperfine -i -r 2 \
		"./bin/ant     -r $(RECT) $(ANT)__2_500_000_002" \
		"./bin/ant-old -r $(RECT) $(ANT)__2_500_000_001" \
		# end

bench-prep-dry:
	git stash
	make ant-dry
	cp ./bin/ant-dry ./bin/ant-dry-old
	git stash pop
	make ant-dry

bench-dry:
	make bench-prep-dry
	hyperfine -i --warmup 1 \
		"./bin/ant-dry     $(ANT)__50_000_002" \
		"./bin/ant-dry-old $(ANT)__50_000_001" \
		# end

bench-dry-large:
	make bench-prep-dry
	hyperfine -i --warmup 1 \
		"./bin/ant-dry     $(ANT)__250_000_002" \
		"./bin/ant-dry-old $(ANT)__250_000_001" \
		# end

bench-mem:
	make bench-prep

	time -lh ./bin/ant     -r $(RECT) $(ANT)__2_000_000_002
	time -lh ./bin/ant-old -r $(RECT) $(ANT)__2_000_000_001

STEPS = 50_000_001
compare:
	make bench-prep
	time -lh ./bin/ant-old $(ANT)__$(STEPS)
	mv            results5/$(ANT)__$(STEPS).png results5/old.png
	time -lh ./bin/ant     $(ANT)__$(STEPS)
	mv            results5/$(ANT)__$(STEPS).png results5/new.png
	open results5/old.png
	open results5/new.png

bench-ant:
	make ant
	hyperfine -i --warmup 1 "./bin/ant $(ANT)__25_000_002"

prof-ant:
	make ant
	./bin/ant -cpuprofile tmp/ant.prof -r $(RECT) $(ANT)__140_000_009
	go tool pprof -http=: -no_browser tmp/ant.prof

prof-ant-mem:
	make ant
	./bin/ant -memprofile tmp/ant-mem.prof $(ANT)__300_000_009
	go tool pprof -http=: -no_browser tmp/ant-mem.prof

test:
	go test ./pgrid
