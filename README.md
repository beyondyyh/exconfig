# 项目名称

Exconfig 提供基于 [`consul KV`](https://www.consul.io/api-docs/kv) 配置中心实现异步更新、获取key、数据格式化等操作。

## 快速开始

### 使用Exconfig

**创建一个exconfig**
> `func New(cfg *Config) (manifest *Manifest, err error)`

- config介绍
    - ConsulServerAddr: consul服务端http地址
    - Datacenter: 数据中心
    - KeyPrefix: key路径，支持监听指定路径下的所有key
    - DiscoverySpan: 服务发现间隔时间，默认1min
    - Logger: 日志打印类，用法参考 [`gdp logger`](https://gitlab.weibo.cn/gdp/gdp/wikis/模块列表/Logger#newwriter)

- 使用默认config示例：
    ```go
    ecfg, err := exconfig.New(exconfig.DefaultConfig())
    ```

- 使用自定义config示例：
    ```go
	config := &exconfig.Config{
		ConsulServerAddr: "http://consul-dev.im.weibo.cn:8500",
		Datacenter:       "kylin_dev",
        KeyPrefix:        "mp_service/release/manifest",
        Logger:           gdpLogger.GetWriter("ral"),
	}
	ecfg, err := exconfig.New(config)
    ```

**生成 consul-api 配置**
> `func (m *Manifest) GenerateEnv() []string {`

**获取 key 原始数据**
> `func (m *Manifest) Acquire(key string) (*api.KVPair, error) {`

**数据类型转换**
> 更多用法参考源码 [`reply.go`](https://gitlab.weibo.cn/gdp/exconfig/blob/master/reply.go)

- Int: func Int(reply *api.KVPair, err error) (int, error) {
- Int64: func Int64(reply *api.KVPair, err error) (int64, error) {
- Bool: func Bool(reply *api.KVPair, err error) (bool, error) {
- Sets: func Sets(reply *api.KVPair, err error, separator string) (set.Set, error) {
- Json: func Json(reply *api.KVPair, err error, v interface{}) error {

**示例**

```go
// 只实例化一次放入内存
var ecfg *exconfig.Manifest

func init() {
	config := &exconfig.Config{
		ConsulServerAddr: "http://consul-dev.im.weibo.cn:8500",
		Datacenter:       "kylin_dev",
		KeyPrefix:        "mp_service/release/manifest",
	}

	var err error
    ecfg, err = exconfig.New(config)
    // 出错时记录日志并关闭服务发现
	if err != nil {
        log.Fatal(err)
		ecfg.Close()
    }
}

func main() {
    // 获取 hello/foo 并转换为string
    foo, err := exconfig.String(ecfg.Acquire("hello/foo"))
	log.Printf("foo:%s, err:%v\n", foo, err)

    // 获取 enable 并转换为bool
    enable, err := exconfig.Bool(ecfg.Acquire("enable"))
    log.Printf("enable:%v, err:%+v\n", enable, err)

    // 获取whitelist 转换为集合
    reply, err := ecfg.Acquire("whitelist")
    users, err := exconfig.Sets(reply, err, "\n")
    // 判断uid 2017223391 是否在白名单内
    isViper := users.Contains("2017223391")
    log.Printf("users:%+v, err:%+v isViper:%t\n", users.Elements(), err, isViper)
}
```
