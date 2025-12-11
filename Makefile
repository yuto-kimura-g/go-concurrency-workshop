.PHONY: help gen w1 w2 w3 w4 s1 s2 s3 s4

# Default target
help:
	@echo "Go Concurrency Workshop - Available Commands"
	@echo ""
	@echo "Log Generation:"
	@echo "  make gen          Generate log files"
	@echo ""
	@echo "Workshop Phases:"
	@echo "  make w1           Run workshop phase 1"
	@echo "  make w2           Run workshop phase 2"
	@echo "  make w3           Run workshop phase 3"
	@echo "  make w4           Run workshop phase 4"
	@echo ""
	@echo "Solution Phases:"
	@echo "  make s1           Run solution phase 1"
	@echo "  make s2           Run solution phase 2"
	@echo "  make s3           Run solution phase 3"
	@echo "  make s4           Run solution phase 4"

# Log Generation
gen:
	go run cmd/loggen/main.go

# Workshop Phases
w1:
	go run ./workshop/phase1/main.go

w2:
	go run ./workshop/phase2/main.go

w3:
	go run ./workshop/phase3/main.go

w4:
	go run ./workshop/phase4/main.go

# Solution Phases
s1:
	go run ./solutions/phase1/main.go

s2:
	go run ./solutions/phase2/main.go

s3:
	go run ./solutions/phase3/main.go

s4:
	go run ./solutions/phase4/main.go
