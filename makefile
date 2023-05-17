ios:
	gomobile bind -v -target ios ./client/ios/mtun ./ping

mtun:
	go build -o mtun