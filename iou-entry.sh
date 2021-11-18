#!/bin/sh
set -e

echo "数据挂载目录： [$MNT_DIR]"
echo "仓库存放目录： [$REPOS_DIR]"
echo "任务文件目录： [$CRON_FILE_DIR]"
echo "当前执行目录： [$PWD]"
echo "判断数据存放目录是已存在或者需要创建"
export DD_DATA_DIR="$MNT_DIR/dd_data"
# 该变量需要传递给ddbot同步脚本仓库使用
export SCRIPTS_REPO_BASE_DIR="$REPOS_DIR/dd_scripts"
export DD_CRON_FILE_PATH="$CRON_FILE_DIR/dd_scripts_cron.sh"

if [ -d $DD_DATA_DIR ]; then
    echo "数据存放目录[$DD_DATA_DIR]已存在，请自行检查配置文件是否正确完整..."
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
        echo "# DATA_LOGS_DIR 为临时脚本存放目录" >>"$DD_DATA_DIR/env.sh"
        echo "export DATA_LOGS_DIR=\"$DD_DATA_DIR/logs\"" >>"$DD_DATA_DIR/env.sh"
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
    echo "# DATA_LOGS_DIR 为临时脚本存放目录" >>"$DD_DATA_DIR/env.sh"
    echo "export DATA_LOGS_DIR=\"$DD_DATA_DIR/logs\"" >>"$DD_DATA_DIR/env.sh"
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
echo "目前只构建三个平台（amd64,arm64,arm）的ddbot，其他架构平台暂未发现使用者，如果有欢迎上报，并且只知道arch为x86_64(amd64)，aarch64(arm64)所以其他的就归到arm上"
if [ "$(arch)" == "x86_64" ]; then
    echo "当前容器运行于amd64平台"
    cp $PWD/ddbot/ddbot-amd64 /usr/local/bin/ddbot
elif [ "$(arch)" == "aarch64" ]; then
    echo "当前容器运行于arm64平台"
    cp $PWD/ddbot/ddbot-arm64 /usr/local/bin/ddbot
else
    echo "当前容器运行于arm平台"
    cp $PWD/ddbot/ddbot-arm /usr/local/bin/ddbot
fi
chmod +x /usr/local/bin/ddbot

echo "[dd_scripts]开始同步仓库..."
sleep 2
ddbot -up syncRepo | sed -e "s|^|[->exec ddbot sync repo] |"
echo "[dd_scripts]仓库同步完成..."

echo "[$DD_CRON_FILE_PATH]任务处理开始..."

#排要扫面的文件crontab的文件名
excludeFile="smiek_jd_zdjr.js,JS_USER_AGENTS.js,JD_DailyBonus.js,JDJRValidator"
logDir="$DD_DATA_DIR/logs"
# 查找指定目录下脚本内的定时任务配置信息
findDirCronFile() {
    if [ $1 ]; then
        findDir=$SCRIPTS_REPO_BASE_DIR/$1
    else
        findDir=$SCRIPTS_REPO_BASE_DIR
    fi
    echo "  开始查找$findDir目录下脚本文件内的crontab任务定义..."
    for scriptFile in $(ls -l $findDir | grep "^-" | awk '{print $9}' | tr "\n" " "); do
        cron=$(sed -n "s/.*crontab=[\"\|']\(.*\)[\"\|'].*/\1/p" "$findDir/$scriptFile")
        if [ "$cron" != "" ] && [ "$(echo $excludeFile | grep "$scriptFile")" == "" ]; then
            cronName=$(sed -n "s/.*new Env([\"\|']\(.*\)[\"\|']).*/\1/p" "$findDir/$scriptFile")
            echo "      #$cronName($findDir/$scriptFile)"
            echo "      $cron node $findDir/$scriptFile >> $logDir/$(echo $scriptFile | sed "s/\.js/.log/g") 2>&1 &"
        fi
    done
}

# 循环查找dd_scripts仓库目录下的脚本文件夹
cd $SCRIPTS_REPO_BASE_DIR
echo "" >$DD_CRON_FILE_PATH
for scriptDir in $(ls -l $SCRIPTS_REPO_BASE_DIR | grep "^d" | grep "dd" | awk '{print $9}' | tr "\n" " "); do
    findDirCronFile $scriptDir
done

findDirCronFile

echo "[$DD_CRON_FILE_PATH]任务处理完成..."
