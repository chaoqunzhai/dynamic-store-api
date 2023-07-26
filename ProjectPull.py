import os
import sys

PROJECT = os.environ.get('PROJECT')
print("PROJECT",PROJECT)
if PROJECT == "home":
    BASE_PATH = "/Users/zhaichaoqun/workspace/goProjects/src/"
    WEAPP_PATH = "/Users/zhaichaoqun/workspace/vueProjects/"
else:
    BASE_PATH = "/Users/zhaichaoqun/workspace/goProjects/src/"
    WEAPP_PATH = "/Users/zhaichaoqun/workspace/wechat/"
projectList = [
    {
        "name": "universe",
        "desc": "宇宙帮",
        "list": [
            {
                "name": "universe-app",
                "desc": "宇宙帮-客户端小程序项目",
                "url": "https://github.com/chaoqunzhai/universe-app",
                "path": os.path.join(WEAPP_PATH, "universe-app"),
            },
            {
                "name": "universe-leader-app",
                "desc": "宇宙帮-帮主户端小程序项目",
                "url": "https://github.com/chaoqunzhai/universe-leader-app",
                "path": os.path.join(WEAPP_PATH, "universe-leader-app"),
            },
            {
                "name": "universe",
                "desc": "宇宙帮-核心API",
                "url": "https://github.com/chaoqunzhai/universe",
                "path": os.path.join(BASE_PATH, "universe"),
            },
            {
                "name": "vue-universe",
                "desc": "宇宙帮-后端PC管理页面",
                "url": "https://github.com/chaoqunzhai/vue-universe",
                "path": "/Users/zhaichaoqun/workspace/vueProjects/vue-universe",
            },
        ]
    },
    {
        "name": "dynamic",
        "desc": "动创云订货软件",
        "list": [
            {
                "name": "dynamic-app",
                "desc": "动创云-客户端订货项目",
                "url": "https://github.com/chaoqunzhai/dynamic-app",
                "path": os.path.join(WEAPP_PATH, "dynamic-app"),
            },
            {
                "name": "dynamic-store-api",
                "desc": "动创云-超管大B后端API",
                "url": "https://github.com/chaoqunzhai/dynamic-store-api",
                "path": os.path.join(BASE_PATH, "dynamic-store-api"),
            },
            {
                "name": "dynamic-weapp-api",
                "desc": "动创云-小程序后端API",
                "url": "https://github.com/chaoqunzhai/dynamic-weapp-api",
                "path": os.path.join(BASE_PATH, "dynamic-weapp-api"),
            },
            {
                "name": "dynamic-web",
                "desc": "动创云-大B管理页面",
                "url": "https://github.com/chaoqunzhai/dynamic-web",
                "path": "/Users/zhaichaoqun/workspace/vueProjects/dynamic-web",
            },
        ]
    },
]

def pullUniverse():
    for ject in projectList:
        if ject.get("name") != "universe": continue

        print("开始下载[", ject.get("name"), "] 描述", ject.get("desc"))
        for row in ject.get("list"):
            print("=================>>子项目[", row.get("name"), "] 描述", row.get("desc"))
            os.chdir(row.get("path"))
            gitPull = "git pull"
            os.system(gitPull)


def pullDynamic():
    for ject in projectList:
        if ject.get("name") != "dynamic": continue
        print("开始下载[", ject.get("name"), "] 描述", ject.get("desc"))
        for row in ject.get("list"):
            print("=================>>子项目[", row.get("name"), "] 描述", row.get("desc"))
            os.chdir(row.get("path"))
            gitPull = "git pull"
            os.system(gitPull)


def commitUniverse(commit):
    for ject in projectList:
        if ject.get("name") != "universe": continue
        print("开始提交[", ject.get("name"), "] 描述", ject.get("desc"))
        for row in ject.get("list"):
            print("=================>>子项目[", row.get("name"), "] 描述", row.get("desc"))
            os.chdir(row.get("path"))
            gitPush = "git add . && git commit -m " + "'" + commit + "'"
            os.system(gitPush)
            os.system("git push origin main")


def commitDynamic(commit):
    for ject in projectList:
        if ject.get("name") != "dynamic": continue
        print("开始提交[", ject.get("name"), "] 描述", ject.get("desc"))
        for row in ject.get("list"):
            print("=================>>子项目[", row.get("name"), "] 描述", row.get("desc"))
            os.chdir(row.get("path"))
            gitPush = "git add . && git commit -m " + "'" + commit + "'"
            os.system(gitPush)
            os.system("git push origin main")


if __name__ == '__main__':
    if len(sys.argv) != 4:
        raise "参数缺少"

    projectName = sys.argv[1]
    action = sys.argv[2]
    commit = sys.argv[3]

    if action == "pull":
        if projectName == "all":
            pullUniverse()
            pullDynamic()
        elif projectName == "u":
            pullUniverse()
        elif projectName == "d":
            pullDynamic()
    else:
        action = "push"
        if projectName == "all":
            commitUniverse(commit)
            commitDynamic(commit)
        elif projectName == "u":
            commitUniverse(commit)
        elif projectName == "d":
            commitDynamic(commit)
