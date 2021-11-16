#!/bin/sh
set -e

echo "数据挂载目录： [$MNT_DIR]"
echo "仓库存放目录： [$REPOS_DIR]"
echo "任务文件目录： [$CRON_FILE_PATH]"
echo "当前执行目录： [$PWD]"
echo "判断数据存放目录是已存在或者需要创建"
export DD_DATA_DIR="$MNT_DIR/dd_data"
# 该变量需要传递给ddbot同步脚本仓库使用
export SCRIPTS_REPO_BASE_DIR="$REPOS_DIR/dd_scripts"

if [ -d $DD_DATA_DIR ]; then
    echo "[$DD_DATA_DIR]数据存放目录，请自行检查配置文件是否正确完整..."
    if [ ! -d "$DD_DATA_DIR/custom_scripts" ]; then
        mkdir -p $DD_DATA_DIR/custom_scripts
    fi
    if [ ! -d "$DD_DATA_DIR/logs" ]; then
        mkdir -p $DD_DATA_DIR/logs
    fi
    if [ ! -f "$DD_DATA_DIR/env.sh" ]; then
        echo "#↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓仓库&bot所需环境变量↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓#" >>"$DD_DATA_DIR/env.sh"
        echo "# REPO_BASE_DIR 为脚本仓库目录，非本仓库目录" >>"$DD_DATA_DIR/env.sh"
        echo "export REPO_BASE_DIR=\"$SCRIPTS_REPO_BASE_DIR\"" >>"$DD_DATA_DIR/env.sh"
        echo "# DATA_BASE_DIR 为数据存放的根目录" >>"$DD_DATA_DIR/env.sh"
        echo "export DATA_BASE_DIR=\"$DD_DATA_DIR\"" >>"$DD_DATA_DIR/env.sh"
        echo "# CUSTOM_SCRIPTS_DIR 为临时脚本存放目录" >>"$DD_DATA_DIR/env.sh"
        echo "export CUSTOM_SCRIPTS_DIR=\"$DD_DATA_DIR/custom_scripts\"" >>"$DD_DATA_DIR/env.sh"
        echo "# ENV_FILE_PATH 为env.sh文件路径，如果自己调整过，此变量请更新" >>"$DD_DATA_DIR/env.sh"
        echo "export ENV_FILE_PATH=\"$DD_DATA_DIR/env.sh\"" >>"$DD_DATA_DIR/env.sh"
        echo "# WSKEY_FILE_PATH 为wskey文件路径，如果自己调整过，此变量请更新" >>"$DD_DATA_DIR/env.sh"
        echo "export WSKEY_FILE_PATH=\"$DD_DATA_DIR/cookies_wskey.list\"" >>"$DD_DATA_DIR/env.sh"
        echo "# DDCK_FILE_PATH 为cookies文件路径，如果自己调整过，此变量请更新" >>"$DD_DATA_DIR/env.sh"
        echo "export DDCK_FILE_PATH=\"$DD_DATA_DIR/cookies.list\"" >>"$DD_DATA_DIR/env.sh"
        echo "# REPLY_KEYBOARD_FILE_PATH 为reply_keyboar快捷回复按钮配置文件路径，如果自己调整过，此变量请更新" >>"$DD_DATA_DIR/env.sh"
        echo "export REPLY_KEYBOARD_FILE_PATH=\"$DD_DATA_DIR/reply_keyboard.list\"" >>"$DD_DATA_DIR/env.sh"
        echo "#↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑仓库&bot所需环境变量↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑#" >>"$DD_DATA_DIR/env.sh"
    fi
else
    echo "[$DD_DATA_DIR]数据存放目录不存在，开始初始化数据存放目录和配置文件..."
    mkdir -p $DD_DATA_DIR
    cd $DD_DATA_DIR
    mkdir -p $DD_DATA_DIR/custom_scripts
    mkdir -p $DD_DATA_DIR/logs
    echo "" >"$DD_DATA_DIR/cookies_wskey.list"
    echo "" >"$DD_DATA_DIR/cookies.list"
    echo "" >"$DD_DATA_DIR/env.sh"
    echo "" >"$DD_DATA_DIR/reply_keyboard.list"
    echo "#↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓仓库&bot所需环境变量↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓#" >>"$DD_DATA_DIR/env.sh"
    echo "# REPO_BASE_DIR 为脚本仓库目录，非本仓库目录" >>"$DD_DATA_DIR/env.sh"
    echo "export REPO_BASE_DIR=\"$SCRIPTS_REPO_BASE_DIR\"" >>"$DD_DATA_DIR/env.sh"
    echo "# DATA_BASE_DIR 为数据存放的根目录" >>"$DD_DATA_DIR/env.sh"
    echo "export DATA_BASE_DIR=\"$DD_DATA_DIR\"" >>"$DD_DATA_DIR/env.sh"
    echo "# CUSTOM_SCRIPTS_DIR 为临时脚本存放目录" >>"$DD_DATA_DIR/env.sh"
    echo "export CUSTOM_SCRIPTS_DIR=\"$DD_DATA_DIR/custom_scripts\"" >>"$DD_DATA_DIR/env.sh"
    echo "# ENV_FILE_PATH 为env.sh文件路径，如果自己调整过，此变量请更新" >>"$DD_DATA_DIR/env.sh"
    echo "export ENV_FILE_PATH=\"$DD_DATA_DIR/env.sh\"" >>"$DD_DATA_DIR/env.sh"
    echo "# WSKEY_FILE_PATH 为wskey文件路径，如果自己调整过，此变量请更新" >>"$DD_DATA_DIR/env.sh"
    echo "export WSKEY_FILE_PATH=\"$DD_DATA_DIR/cookies_wskey.list\"" >>"$DD_DATA_DIR/env.sh"
    echo "# DDCK_FILE_PATH 为cookies文件路径，如果自己调整过，此变量请更新" >>"$DD_DATA_DIR/env.sh"
    echo "export DDCK_FILE_PATH=\"$DD_DATA_DIR/cookies.list\"" >>"$DD_DATA_DIR/env.sh"
    echo "# REPLY_KEYBOARD_FILE_PATH 为reply_keyboar快捷回复按钮配置文件路径，如果自己调整过，此变量请更新" >>"$DD_DATA_DIR/env.sh"
    echo "export REPLY_KEYBOARD_FILE_PATH=\"$DD_DATA_DIR/reply_keyboard.list\"" >>"$DD_DATA_DIR/env.sh"
    echo "#↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑仓库&bot所需环境变量↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑#" >>"$DD_DATA_DIR/env.sh"
    echo "[$DD_DATA_DIR]数据存放目录初始化完成..."
fi

# 判断平台架构使用对应平台版本的ddbot
echo "目前只构建三个平台（and64,arm64,arm）的ddbot，其他架构平台暂未发现使用者，如果有欢迎上报，并且只知道arch为x86_64(amd64)，aarch64(arm64)所以其他的就归到arm上"
if [ "$(arch)" == "x86_64" ]; then
    echo "amd64"
    cp $PWD/ddbot/ddbot-amd64 /usr/local/bin/ddbot
elif [ "$(arch)" == "aarch64" ]; then
    echo "arm64"
    cp $PWD/ddbot/ddbot-arm64 /usr/local/bin/ddbot
else
    echo "arm"
    cp $PWD/ddbot/ddbot-arm /usr/local/bin/ddbot
fi
chmod +x /usr/local/bin/ddbot

echo "开始同步仓库dd_scripts..."
ddbot -up syncRepo
echo "dd_scripts仓库完成..."
