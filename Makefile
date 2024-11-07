all:
	go build -buildmode=c-shared -o in_p1aic.so .

fast:
	go build in_gdummy.go

clean:
	rm -rf *.so *.h *~