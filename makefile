ios:
	gomobile bind -v -target ios ./client/ios/hy ./ping
#	 gomobile bind -target=ios -o goPing.xcframework ./ping    打包出goPing
android:
	gomobile bind -v -target android -androidapi 21 ./ping
mtun:
	go build -o mtun
