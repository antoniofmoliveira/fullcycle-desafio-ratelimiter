witherror:
	go run cmd/main.go --qt-tokens 100000 --time-frame-seconds 1 --simulate-slow-requests --seed-for-simulate-slow-requests 100 --simulate-errors  --seed-for-simulate-errors 100

withouterror:
	go run cmd/main.go --qt-tokens 100000 --time-frame-seconds 1