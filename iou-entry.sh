#!/bin/sh
set -e

echo "数据挂载目录： [$MNT_DIR]"
echo "仓库存放目录： [$REPOS_DIR]"
echo "任务文件目录： [$CRON_FILE_DIR]"
echo "当前执行目录： [$PWD]"
echo "判断数据存放目录是已存在或者需要创建"
export DD_DATA_DIR="$MNT_DIR/dd_data"
export IOU_DD_DIR=$(pwd)
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

source "$DD_DATA_DIR/env.sh"

# 去iou-dd仓库目录处理相关配置
cd $IOU_DD_DIR
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

echo -e
echo "[dd_scripts]开始同步仓库..."
sleep 2
ddbot -up syncRepo | sed -e "s|^|[->exec ddbot sync repo] |"
echo "[dd_scripts]仓库同步完成..."

cd $SCRIPTS_REPO_BASE_DIR

echo "检查package.json依赖是否有更新..."
if [ ! -d node_modules ]; then
    echo -e "检测到首次部署, 运行 npm install..."
    cat package.json >old_package.json
    npm install --loglevel error --prefix
else
    if [ "$(cat old_package.json)" != "$(cat package.json)" ]; then
        echo -e "检测到package.json有变化，运行 npm install..."
        cat package.json >old_package.json
        npm install --loglevel error --prefix
    else
        echo -e "检测到package.json无变化，跳过...\n"
    fi
fi

echo -e
echo "[$DD_CRON_FILE_PATH] 定时任务处理逻辑顺序如下："
echo "[$DD_CRON_FILE_PATH] 1：循环查找对应子目录脚本内的 crontab 配置"
echo "[$DD_CRON_FILE_PATH] 2：查找对应的仓库根目录脚本内的 crontab 配置"
echo "[$DD_CRON_FILE_PATH] 3：处理旧版的任务配置 docker/crontab_list.sh 取出来1，2两步没处理到的任务"
echo "[$DD_CRON_FILE_PATH] 4：执行配置的 mod shell，可能会增加的任务，调用mod shell会默认把任务要写入的文件路径作为参数传进去，时候使用自行决定"
echo "[$DD_CRON_FILE_PATH] 5：排除 env.sh 配置的排除任务"
echo "[$DD_CRON_FILE_PATH] 6：追加配置的自定crontab_list.sh"
echo "[$DD_CRON_FILE_PATH] 7：替换 node 执行命令为 ddnode"
echo "[$DD_CRON_FILE_PATH] 8：得到本仓库最终的定时任务配置文件"
echo "[$DD_CRON_FILE_PATH] "
echo "[$DD_CRON_FILE_PATH] 任务处理开始..."

#排除要扫描的文件crontab的文件名
logDir="$DD_DATA_DIR/logs"
CRONFILES="xxx.js"
# 查找指定目录下脚本内的定时任务配置信息
findDirCronFile() {
    if [ $1 ]; then
        findDir=$SCRIPTS_REPO_BASE_DIR/$1
    else
        findDir=$SCRIPTS_REPO_BASE_DIR
    fi
    echo "[$DD_CRON_FILE_PATH]   开始查找$findDir目录下脚本文件内的crontab任务定义..."
    for scriptFile in $(ls -l $findDir | grep "^-" | awk '{print $9}' | tr "\n" " "); do
        cron=$(sed -n "s/.*crontab=[\"\|']\(.*\)[\"\|'].*/\1/p" "$findDir/$scriptFile")
        if [ "$cron" != "" ]; then
            cronName=$(sed -n "s/.*new Env([\"\|']\(.*\)[\"\|']).*/\1/p" "$findDir/$scriptFile")
            # echo "      #$cronName($findDir/$scriptFile)"
            # echo "      $cron node $findDir/$scriptFile >> $logDir/$(echo $scriptFile | sed "s/\.js/.log/g") 2>&1 &"
            echo "#$cronName($findDir/$scriptFile)" >>$DD_CRON_FILE_PATH
            echo "$cron node $findDir/$scriptFile >>$logDir/$(echo $scriptFile | sed "s/\.js/.log/g") 2>&1 &" >>$DD_CRON_FILE_PATH
            echo "" >>$DD_CRON_FILE_PATH
            CRONFILES="$CRONFILES\|$scriptFile"
        fi
    done
}

# 循环查找dd_scripts仓库目录下的脚本文件夹
cd $SCRIPTS_REPO_BASE_DIR

echo "#↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓ [$SCRIPTS_REPO_BASE_DIR] 仓库任务列表 ↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓↓#" >$DD_CRON_FILE_PATH
echo "" >>$DD_CRON_FILE_PATH
for scriptDir in $(ls -l $SCRIPTS_REPO_BASE_DIR | grep "^d" | grep "dd" | awk '{print $9}' | tr "\n" " "); do
    findDirCronFile $scriptDir
done

echo "[$DD_CRON_FILE_PATH]   "
findDirCronFile

# echo "[$DD_CRON_FILE_PATH]   "
# echo "[$DD_CRON_FILE_PATH]   处理兼容旧版脚本任务配置文件"
# # echo $CRONFILES
# cat $SCRIPTS_REPO_BASE_DIR/docker/crontab_list.sh | grep -v "$CRONFILES" >>$DD_CRON_FILE_PATH

echo "[$DD_CRON_FILE_PATH]   "
echo "[$DD_CRON_FILE_PATH]   处理 mod shell"
if [ $CUSTOM_SHELL_FILE ]; then
    echo "" >"$SCRIPTS_REPO_BASE_DIR/shell_mod.sh"
    if expr "$CUSTOM_SHELL_FILE" : 'http.*' &>/dev/null; then
        echo "[$DD_CRON_FILE_PATH]   [CUSTOM_SHELL_FILE] 自定义shell脚本为远程脚本，开始下在自定义远程脚本 $CUSTOM_SHELL_FILE ..."
        wget -O "$SCRIPTS_REPO_BASE_DIR/shell_mod.sh" "$CUSTOM_SHELL_FILE" | sed -e "s|^|[$DD_CRON_FILE_PATH]   [CUSTOM_SHELL_FILE] |"
        echo "[$DD_CRON_FILE_PATH]   [CUSTOM_SHELL_FILE] 下载完成..."
        echo "[$DD_CRON_FILE_PATH]   [CUSTOM_SHELL_FILE] 自定义 mod shell 开始执行..."
        echo "#远程自定义shell脚本追加定时任务" >>$DD_CRON_FILE_PATH
        sh "$SCRIPTS_REPO_BASE_DIR/shell_mod.sh" | sed -e "s|^|[$DD_CRON_FILE_PATH]   [CUSTOM_SHELL_FILE->exec] |"
        echo "[$DD_CRON_FILE_PATH]   [CUSTOM_SHELL_FILE] 自定义 mod shell 执行结束..."
    else
        if [ ! -f "$CUSTOM_SHELL_FILE" ]; then
            echo "[$DD_CRON_FILE_PATH]   [CUSTOM_SHELL_FILE] 自定义shell脚本为docker挂载脚本文件，但是指定挂载文件 $CUSTOM_SHELL_FILE 不存在，跳过执行..."
        else
            cat "$CUSTOM_SHELL_FILE" >"$SCRIPTS_REPO_BASE_DIR/shell_mod.sh"
            echo "[$DD_CRON_FILE_PATH]   [CUSTOM_SHELL_FILE] 自定义 mod shell 开始执行..."
            echo "#挂载自定义shell脚本追加定时任务" >>$DD_CRON_FILE_PATH
            sh "$SCRIPTS_REPO_BASE_DIR/shell_mod.sh" | sed -e "s|^|[$DD_CRON_FILE_PATH]   [CUSTOM_SHELL_FILE->exec] |"
            echo "[$DD_CRON_FILE_PATH]   [CUSTOM_SHELL_FILE] 自定义 mod shell 执行结束..."
        fi
    fi
fi

#根据EXCLUDE_CRON配置的关键字剔除相关任务 EXCLUDE_CRON="cfd,joy"
echo "[$DD_CRON_FILE_PATH] "
echo "[$DD_CRON_FILE_PATH]   处理 EXCLUDE_CRON 配置的关键字剔除相关任务..."
if [ $EXCLUDE_CRON ]; then
    for kw in $(echo $EXCLUDE_CRON | tr "," " "); do
        matchCron=$(cat $DD_CRON_FILE_PATH | grep -v "$kw")
        if [ -z "$matchCron" ]; then
            echo "[$DD_CRON_FILE_PATH]   [EXCLUDE_CRON] 关键词 $kw 未匹配到任务"
            echo "[$DD_CRON_FILE_PATH]   [EXCLUDE_CRON] "
        else
            echo "[$DD_CRON_FILE_PATH]   [EXCLUDE_CRON] 根据关键词 $kw 剔除的任务..."
            echo "$matchCron" | sed -e "s|^|[$DD_CRON_FILE_PATH]   [EXCLUDE_CRON] |"
            echo "[$DD_CRON_FILE_PATH]   [EXCLUDE_CRON] "
            sed -i '/'"$kw"'/d' $DD_CRON_FILE_PATH
        fi
    done
fi

echo "[$DD_CRON_FILE_PATH] "
echo "[$DD_CRON_FILE_PATH]   处理 CUSTOM_LIST_FILE 配置的追加自定义任务..."
if [ "$CUSTOM_LIST_FILE" ]; then
    echo "[$DD_CRON_FILE_PATH]   [CUSTOM_LIST_FILE] 配置了自定义任务文件：$CUSTOM_LIST_FILE ..."
    if [ -f "$CUSTOM_LIST_FILE" ]; then
        echo "" >>$CUSTOM_LIST_FILE
        cat "$CUSTOM_LIST_FILE" >>$DD_CRON_FILE_PATH

    else
        echo "[$DD_CRON_FILE_PATH]   [CUSTOM_LIST_FILE] 配置的自定义任务文件：$CUSTOM_LIST_FILE ，但是文件不存在或者路径错误，跳过..."
    fi
fi

echo "[$DD_CRON_FILE_PATH]   "
echo "[$DD_CRON_FILE_PATH]   处理替换 node 执行命令为 ddnode，增加 ts 输出日志时间"
sed -i "s/ node / ddnode /g" $DD_CRON_FILE_PATH
sed -i "s/\(\| ts\| |ts\|| ts\)//g" $DD_CRON_FILE_PATH
sed -i "/ddBot/!s/>>/\|sed -e \"s|^|\$(date +'%Y-%m-%d %H:%M:%S') | \" >>/g" $DD_CRON_FILE_PATH
sed -i "/\(>&1 &\|> &1 &\)/!s/>&1/>\&1 \&/g" $DD_CRON_FILE_PATH
sed -i "s|/data/logs|$DD_DATA_DIR/logs|g" $DD_CRON_FILE_PATH
sed -i "s|/scripts/|$SCRIPTS_REPO_BASE_DIR/|g" $DD_CRON_FILE_PATH

echo "" >>$DD_CRON_FILE_PATH
echo "#↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑ [$SCRIPTS_REPO_BASE_DIR] 仓库任务列表 ↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑↑#" >>$DD_CRON_FILE_PATH
echo "[$DD_CRON_FILE_PATH] "
echo "[$DD_CRON_FILE_PATH] 任务处理完成..."
echo -e

echo "增加 submitShareCode.sh 脚本，用于清理日志，调用ddbot提交互助码到助力池..."
sed -i "3 i 30 23 * * * sleep \$((RANDOM % 400)); sh $SCRIPTS_REPO_BASE_DIR/submitShareCode.sh |sed -e \"s|^|\$(date +'%Y-%m-%d %H:%M:%S') | \" >>$logDir/logs/submitCode.log 2>&1 &" $DD_CRON_FILE_PATH
sed -i "4 i  " $DD_CRON_FILE_PATH
(
    cat <<EOF
#!/bin/sh
set -e
curr_dt=\$(date +'%Y-%m-%d' | awk '{print \$1}')
echo "清除非当日(\$curr_dt)产生的日志，准备提交互助码码到助力池"
for dd_log in \$(ls "$DD_DATA_DIR/logs/" | grep "^.*log\$" | grep -v "sharecode"); do
      sed -i "/^\$curr_dt.*/!d" "$DD_DATA_DIR/logs/\$dd_log"
done
# 开始上传
ddbot -up commitShareCode
EOF
) >"$SCRIPTS_REPO_BASE_DIR/submitShareCode.sh"

echo "增加 ddnode 命令脚本..."
echo "" >/usr/local/bin/ddnode
chmod +x /usr/local/bin/ddnode
(
    cat <<EOF
#!/bin/sh
set -e

if [ -f $DD_DATA_DIR/env.sh ]; then
    source $DD_DATA_DIR/env.sh
fi

first=\$1
cmd=\$*
# 判断命令是否需要执行混淆后的js脚本
if [ -n "\$(echo \$cmd | grep "_hx.js")" ]; then
    if [ \$DEFAULT_EXEC_HX_SCRIPT ] && [ \$DEFAULT_EXEC_HX_SCRIPT == "Y" ]; then
        echo "[ddnode] 配置了 DEFAULT_EXEC_HX_SCRIPT=Y，混淆脚本执行命令继续..."
    else
        echo '[ddnode] 执行的为混淆脚本，退出执行。如需启用请配置 [export DEFAULT_EXEC_HX_SCRIPT="Y"]'
        exit 0
    fi
fi

# 指令交给node后台执行
if [ "\$1" == "conc" ]; then
    for job in \$(cat \$DDCK_FILE_PATH | grep -v "#" | paste -s -d ' '); do
        {
            export JD_COOKIE=\$job && node \${cmd/\$1/}
        } &
    done
elif [ -n "\$(echo \$first | sed -n "/^[0-9]\+\$/p")" ]; then
    echo "[ddnode] 指定 ck idx=\${first} 执行 [node\${cmd/\$1/}] 命令"
    echo ""
    export JD_COOKIE=\$(cat \$DDCK_FILE_PATH | grep -v "#\|^$" | sed -n "\${first}p") && node \${cmd/\$1/}
elif [ -n "\$(cat \$DDCK_FILE_PATH | grep "pt_pin=\$first")" ]; then
    echo "[ddnode] 指定 ck pin=\$first 执行 [node\${cmd/\$1/}]"
    export JD_COOKIE=\$(cat \$DDCK_FILE_PATH | grep "pt_pin=\$first") && node \${cmd/\$1/}
else
    echo "[ddnode] 执行 [node \${cmd}] 命令"
    echo ""
    export JD_COOKIE=\$(cat \$DDCK_FILE_PATH | grep -v "#\|^$" | paste -s -d '&') && node \$cmd
fi
EOF
) >/usr/local/bin/ddnode
