# cn-idcard-area

用中国大陆身份证号码的前 6 位离线查行政区划。数据编译进 Go 二进制，不依赖第三方 Go 模块，查询时不需要联网或单独携带 CSV。

需要 Go 1.22 或更高版本。在仓库根目录运行：

```console
$ go run . 440305
广东省 / 深圳市 / 南山区
$ go run . 110103
北京市 / 崇文区
$ go run . 999999
未知
```

只接受一个参数，必须是 6 位 ASCII 数字。未收录的 6 位码输出 `未知`；参数格式错误时输出用法，编译后的 CLI 退出码为 `2`。直辖市、省直管县和不设区的市按实际层级输出，省略不存在的层级。

同一码有多条记录时，取启用年份最新的一条。历史码保留该记录中的地名，例如 `110103` 输出原崇文区，不追踪为东城区。查询结果是地址码对应的区划，不能确定持证人的现住址，也不能用来判断身份证是否有效。

内置数据来自 [yescallop/areacodes](https://github.com/yescallop/areacodes)，按 [CC0 1.0](data/LICENSE) 提供。[data/SOURCE](data/SOURCE) 记录了当前数据对应的上游提交和最近成功检查的月份。

[同步工作流](https://github.com/jimyag/cn-idcard-area/actions/workflows/sync-areacodes.yml)每周一 10:17（UTC+8）检查上游，也可手动运行。上游 CSV 变更且通过 `go test ./...` 后，工作流将数据和来源记录一起提交到 `main`。上游许可变更或测试失败时，保留已有数据。即使码表没有变更，每月也会记录一次成功检查，避免公开仓库因[长期无活动停用定时任务](https://docs.github.com/en/actions/reference/workflows-and-actions/events-that-trigger-workflows#schedule)。
