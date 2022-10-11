# kxapp-common
各个项目经常使用的共享的代码

1.设置go环境变量
GONOPROXY=github.com/kxapp-com/*


2.增加 github ssh 访问功能
3. 修改如下文件添加配置
C:\Users\admin\.gitconfig

[user]
	name = lilili87222
	email = li_li_li87222@163.com
[url "ssh://git@github.com/"]
	insteadOf = https://github.com/