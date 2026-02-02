# 如何添加 jigsaw：
# 1. 在同一目录下 clone mtun，hysteria 和 outline-apps-v1-17-0 项目，
# 2. 运行 `go work init ./mtun ./outline-apps-v1-17-0` 生成 go.work 文件，
# 3. 进入 mtun 目录，make ios/android 打包，将包含 jigsaw 协议。

OUTLINE_DIR=../outline-apps-v1-17-0/client/go/outline

ios:
	gomobile bind -v -target ios ./client/ios/hy ./ping ${OUTLINE_DIR}/platerrors ${OUTLINE_DIR}/tun2socks ${OUTLINE_DIR}
#	 gomobile bind -target=ios -o goPing.xcframework ./ping    打包出goPing
android:
	gomobile bind -v -target android -androidapi 21 ./client/ios/hy ./ping ${OUTLINE_DIR}/platerrors ${OUTLINE_DIR}/tun2socks ${OUTLINE_DIR}
mtun:
	go build -o mtun
