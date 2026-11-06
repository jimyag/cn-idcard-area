# cn-idcard-area

用中国大陆身份证的地址码离线查省、市、区县。码表内置，无第三方 Go 依赖。

最低 Go 版本见 [go.mod](go.mod)。

## Go 包

在项目的 Go 模块中运行：

```sh
go get github.com/jimyag/cn-idcard-area@latest
```

要固定版本，把 `@latest` 换成所需的[发布标签](https://github.com/jimyag/cn-idcard-area/releases)。

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/jimyag/cn-idcard-area"
)

func main() {
	area, ok := idcardarea.Lookup("440305")
	if !ok {
		fmt.Println("未知")
		return
	}
	if err := json.NewEncoder(os.Stdout).Encode(area); err != nil {
		log.Fatal(err)
	}
}
```

```json
{"province":"广东省","city":"深圳市","district":"南山区"}
```

`Lookup` 只用前 6 位查表。旧码仍可查，同一码取启用年份最新的记录。没有对应区划层级的字段为空字符串。

直接传 6 位地址码。支持完整号码的版本也接受 15 位和 18 位输入：只允许 ASCII 数字，18 位的末位也可以是 `X` 或 `x`。查询不校验出生日期或校验位，结果不能用来确定现住址。

输入格式错误或查不到时，`Lookup` 返回 `Area{}, false`。如果所用版本提供 `DataError`，可调用它查看码表加载错误。

## 命令行

```sh
go install github.com/jimyag/cn-idcard-area/cmd/cn-idcard-area@latest
cn-idcard-area 110103
```

```json
{"province":"北京市","city":"","district":"崇文区"}
```

命令只收一个参数，格式同 `Lookup`。查不到时输出 `未知`；参数错误时打印用法并退出 2。

[GitHub Releases](https://github.com/jimyag/cn-idcard-area/releases) 提供 Linux、macOS（darwin）和 Windows 的 amd64、arm64 可执行文件。Linux 和 macOS 下载后用 `chmod +x` 添加执行权限。

推送 `vX.Y.Z` 标签会触发[测试、构建和发布](https://github.com/jimyag/cn-idcard-area/actions/workflows/release.yml)。已有 Release 可手动运行发布工作流，填原标签补附件。

## 数据和许可

代码和文档采用 [MIT](LICENSE)。码表来自 [yescallop/areacodes](https://github.com/yescallop/areacodes) 的 `codes.json`，采用 [CC0](data/LICENSE)。[data/SOURCE](data/SOURCE) 记录来源和最近核验月份。

[码表同步](https://github.com/jimyag/cn-idcard-area/actions/workflows/sync-areacodes.yml)每周一 10:17（UTC+8）运行，也可手动触发。许可一致且测试通过才更新码表。数据不变时，每月只记一次核验，避免[定时任务因长期无活动停用](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#schedule)。
