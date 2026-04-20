# 如何添加 jigsaw：
# 
# Android:
# 1. 在同一目录下 clone mtun，hysteria 和 outline-apps-v1-17-0 项目
#		+ Android 不能使用更高版本的 outline, 实测网速会大大下降，原因不明
# 2. 运行 `ln -s ../outline-apps-v1-17-0 ../outline-apps` 创建目录链接
#		+ Windows 命令: powershell -Command "New-Item -ItemType SymbolicLink -Path '../outline-apps' -Target '../outline-apps-v1-17-0'"
# 3. 进入 mtun 目录，`make android` 打包，将包含 jigsaw 协议。
#
# ios:
# 1. 在同一目录下 clone mtun，hysteria 和 outline-apps-client-ios-v1-19-0-rc-1 项目
# 2. 运行 `ln -s ../outline-apps-client-ios-v1-19-0-rc-1 ../outline-apps` 创建目录链接
# 3. 进入 mtun 目录，`make ios` 打包，将包含 jigsaw 协议。

OUTLINE_DIR=../outline-apps

# 导出环境变量到所有规则  环境变量用于解决下面这个url的issue
# https://github.com/golang/go/issues/71827#issuecomment-2669425491
export GODEBUG=gotypesalias=0
export CGO_CFLAGS=-fstack-protector-strong
export MACOSX_DEPLOYMENT_TARGET=12.0

export JAVA_TOOL_OPTIONS = -Dfile.encoding=utf-8

ios:
	gomobile bind -v -target ios ./client/ios/hy ./ping ${OUTLINE_DIR}/client/go/outline/platerrors ${OUTLINE_DIR}/client/go/outline/tun2socks ${OUTLINE_DIR}/client/go/outline
#	 gomobile bind -target=ios -o goPing.xcframework ./ping    打包出goPing
android:
	gomobile bind -v -target android -androidapi 21 ./client/ios/hy ./ping ${OUTLINE_DIR}/client/go/outline/platerrors ${OUTLINE_DIR}/client/go/outline/tun2socks ${OUTLINE_DIR}/client/go/outline
mtun:
	go build -o mtun
