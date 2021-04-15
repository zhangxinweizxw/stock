all:build
	
build:
	gox -osarch="linux/amd64" -output ./bin/stock

clean:
	@rm -rf bin
	 
test:
	go test ./go/... -race

out:
	rm stock
 -rf
	mkdir -p stock
/conf.d
	cp -a conf.d/stock
_dev.xml stock
/conf.d/
	cp -a bin/stock
 stock
/
	tar -zcf stock
.tar.gz stock
/
