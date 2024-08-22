SRC_FILE=hello.s
OPTS=-nostdlib

# SRC_FILE=hello-libc.s
# OPTS=

OBJ_FILE=a.out
ARCH=x86_64

.PHONY: run
run:
	docker compose up # && docker compose down -v

.PHONY: exec
exec:
	./${OBJ_FILE}

.PHONY: asm
asm:
	gcc ${OPTS} ${SRC_FILE}

.PHONY: show
show:
	readelf --syms ${OBJ_FILE}

.PHONY: header
header:
	readelf --file-header ${OBJ_FILE}