#!/bin/bash

export ENV_UNIX_SOCK_PATH=./test.listen.sock

test -S "${ENV_UNIX_SOCK_PATH}" || exec env spat="${ENV_UNIX_SOCK_PATH}" sh -c '
	echo the path ${spat} is not a socket.
	exit 1
'

./cmd/postreqs2usock/postreqs2usock
