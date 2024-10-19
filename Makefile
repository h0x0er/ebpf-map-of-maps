bpf-build:
	@bash build.sh generate

build:
	@bash build.sh all

run: build
	@chmod +x ./maps
	@sudo ./maps

log:
	@bash build.sh logk

symbols:
	@cat /proc/kallsyms

clean:
	@rm map*
