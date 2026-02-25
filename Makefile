# 如何添加 jigsaw：
# 
# Android:
# 1. 在同一目录下 clone mtun，hysteria 和 outline-apps-v1-17-0 项目
#		+ Android 不能使用更高版本的 outline, 实测网速会大大下降，原因不明
# 2. 运行 `go work init ./mtun ./outline-apps-v1-17-0` 生成 go.work 文件，
# 3. 进入 mtun 目录，make android 打包，将包含 jigsaw 协议。
#
# ios:
# 1. 在同一目录下 clone mtun，hysteria 和 outline-apps-client-macos-v1-19-5 项目
# 2. 运行 `go work init ./mtun ./outline-apps-client-macos-v1-19-5` 生成 go.work 文件，
# 3. 进入 mtun 目录，make ios 打包，将包含 jigsaw 协议。

OUTLINE_DIR_ANDROID=../outline-apps-v1-17-0/client/go/outline
OUTLINE_DIR_IOS=../outline-apps-client-macos-v1-19-5/client/go/outline

# 导出环境变量到所有规则  环境变量用于解决下面这个url的issue
# https://github.com/golang/go/issues/71827#issuecomment-2669425491
export GODEBUG=gotypesalias=0
export CGO_CFLAGS=-fstack-protector-strong
export MACOSX_DEPLOYMENT_TARGET=12.0

ios:
	gomobile bind -v -target ios ./client/ios/hy ./ping ${OUTLINE_DIR_IOS}/platerrors ${OUTLINE_DIR_IOS}/tun2socks ${OUTLINE_DIR_IOS}
#	 gomobile bind -target=ios -o goPing.xcframework ./ping    打包出goPing
android:
	gomobile bind -v -target android -androidapi 21 ./client/ios/hy ./ping ${OUTLINE_DIR_ANDROID}/platerrors ${OUTLINE_DIR_ANDROID}/tun2socks ${OUTLINE_DIR_ANDROID}
mtun:
	go build -o mtun
