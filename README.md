# cn-idcard-area

在 Go 程序中，用中国大陆身份证号码的前 6 位离线查省、市、区县。JSON 码表随包嵌入，无需联网、加载数据文件或安装第三方 Go 模块。

需要 Go 1.16 或更高版本。

## Go 包

在调用方模块中添加固定版本的依赖：

```sh
go get github.com/jimyag/cn-idcard-area@v0.1.0
```

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

`Lookup(prefix string) (Area, bool)` 采用 Go 的 `value, ok` 约定。传入 6 位 ASCII 数字；格式不符或未收录时返回零值和 `false`。三个 JSON 字段始终存在，没有对应层级时值为 `""`，如北京市的 `city`、东莞市的 `district`。

同一码有多条记录时，取启用年份最新的一条。已撤销的 `110103` 仍返回原崇文区。结果对应地址码的行政区划，不能确定现住址，也不能用来判断身份证是否有效。

## 命令行

```sh
go install github.com/jimyag/cn-idcard-area/cmd/cn-idcard-area@v0.1.0
cn-idcard-area 110103
```

```json
{"province":"北京市","city":"","district":"崇文区"}
```

命令只接受一个六位数字参数。未收录的码输出 `未知`；参数不合法时输出用法并以状态码 2 退出。

## 数据和许可

代码和文档采用 [MIT 许可证](LICENSE)。内置码表来自 [yescallop/areacodes](https://github.com/yescallop/areacodes) 的 `codes.json`，继续采用 [CC0 1.0](data/LICENSE)。保留上游区划树、代码、名称和启用年份，[data/SOURCE](data/SOURCE) 记录数据来源和最近成功检查的月份。

[同步工作流](https://github.com/jimyag/cn-idcard-area/actions/workflows/sync-areacodes.yml)每周一 10:17（UTC+8）检查上游，也可手动运行。工作流用托管 runner 的 `jq` 提取查询数据，核对上游许可证并通过 `go test ./...` 后才写入 `main`。没有变更时每月记录一次成功检查，避免公开仓库因[长期无活动停用定时任务](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#schedule)。
