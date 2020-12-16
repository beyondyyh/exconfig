# 项目名称

Exconfig 提供基于 [`consul KV`](https://www.consul.io/api-docs/kv) 配置中心实现异步更新、获取key、数据格式化等操作。

## 快速开始

### 使用Exconfig

**创建一个exconfig**
> `func New(cfg *Config, opts ...Option) (manifest *Manifest, err error) {`

- config介绍
    - ConsulServerAddr: consul服务端http地址
    - Datacenter: 数据中心
    - KeyPrefix: key路径，支持监听指定路径下的所有key
- opts可选参数
    - `WithSpan(d time.Duration)`: 服务发现间隔时间，默认1min
    - `WithLogger(l hclog.Logger)`: 日志打印类，用法参考 [hclog](https://github.com/hashicorp/go-hclog)

- 使用默认config示例：
    ```go
    ecfg, err := exconfig.New(exconfig.DefaultConfig())
    ```

- 使用自定义config示例：
    ```go
	ecfg, err := exconfig.New(
		&exconfig.Config{
			ConsulServerAddr: "http://consul-dev.im.weibo.cn:8500/",
			Datacenter:       "kylin_dev",
			KeyPrefix:        "mp_service/release/manifest",
		},
		exconfig.WithSpan(10*time.Second),
		exconfig.WithLogger(hclog.New(&hclog.LoggerOptions{
			Name:       "exconfig-example",
			JSONFormat: true,
			Color:      hclog.AutoColor,
		})),
	)
	if err != nil {
		log.Fatal(err)
		ecfg.Close()
	}
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

参考 [example/main.go](example/main.go)