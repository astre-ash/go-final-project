include .env
export

run:
	go run ./cmd/todo/main.go

test-db:
	go test -v -run ^TestDB$$ ./tests

test-nextdate:
	 go test -v -run ^TestNextDate$$ ./tests

test-addtask:
	go test -v -run ^TestAddTask$$ ./tests

test-tasks:
	go test -v -run ^TestTasks$$ ./tests

test-edit:
	go test -run ^TestEditTask$$ ./tests	

test-done:
	go test -run ^TestDone$$ ./tests

test-del:
	go test -run ^TestDelTask$$ ./tests	

test-all:
	go test -v ./...

run-auth:
	TODO_PASSWORD=123 go run ./cmd/todo/main.go	

