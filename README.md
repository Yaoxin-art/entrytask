# entrytask

## 一、背景及目的
EntryTask项目，通过实现RPC通用框架，使用Go HTTP API，MySQL或Redis操作等能力完成一个用户管理系统，主要功能包含用户（管理员）登录及用户数据的增删查改功能。

## 二、逻辑架构设计
整体分为网关（nginx代）、API层（RPC consumer）、TCP Server层（RPC provider）三层，图片存放于磁盘，前端访问则通过nginx代理，通过配置项确定文件（头像图片）存放的位置，以及配置组装完整访问URL。
![整体架构图](/docs/imgs/framework01.png)

项目结构：
![项目结构](/docs/imgs/project_struct.png)

## 三、核心逻辑详细设计

### 1、登录
![登录流程图](/docs/imgs/login.png)

### 2、已登录用户查看用户信息
![已登录用户查看用户信息流程图](/docs/imgs/info.png)

### 3、更新昵称
![更新昵称流程图](/docs/imgs/update_nick.png)

### 4、更新头像
![更新头像流程图](/docs/imgs/update_profile.png)

## 四、接口设计
接口文档地址：

## 五、存储设计
存储使用到redis和mysql，其中mysql单表存1kw用户信息（表结构简单，内容较少），redis用户缓存用户信息（使用hash数据类型，不包含密码）及登录用户的token与用户名映射关系（string数据类型）。

## 六、外部依赖与限制
静态图片存放于磁盘，并通过nginx做代理访问，nginx需要的主要配置见 [nginx配置依赖](/docs/nginx_conf.md)

## 七、部署方案与环境要求
配置，固定写入到代码常量使用：
    1、数据库连接配置；
    2、Redis连接配置；
    3、上传图片保存路径；
    4、图片访问URL前缀；

部署方案：
    Makefile包含 clean、fmt、vet、cover、test、build、run 等功能，使用`make run[Rpc|Web]`启动。
    启动顺序：
        启动mysql服务
        启动redis服务
        启动nginx
        启动web模块
        启动worker模块


## 八、SLA

### 目标

- 数据库必须有超过1000w的用户数据
- 结果必须正确
- 每个请求都要包含RPC调用以及Mysql/Redis访问
- 性能要求
    - 200固定用户并发，HTTP API QPS大于3000
    - 200随机用户并发，HTTP API QPS大于1000
    - 2000固定用户并发，HTTP API QPS大于1500
    - 2000随机用户并发，HTTP API QPS大于800

### 基准测试方法
选择go自带benchmark测试模块，大致使用方式如下：
```go
func BenchmarkInfoRandom(b *testing.B) {
	// 初始化httpclient 和 user列表（取20000用户用于随机）
	initForBenchmark(false)
	defer destroyHttpClients()

	b.ResetTimer()
	parallelism := clientSize
	b.SetParallelism(parallelism)
	fmt.Printf("parallelism:%d \n", parallelism)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var err error
			id := rand.Intn(clientSize)
			client := clients[id]
			uid := rand.Intn(userSize)
			// 固定用户时 各client各使用一个用户，即 uid=id
			user := users[uid]
			requestUrl := httpServerAddr + "/user/find?username=" + user.Username
			req, err := http.NewRequest(http.MethodGet, requestUrl, nil)
			if err != nil {
				b.Error(err)
				continue
			}
			res, err := client.Do(req)
			if err != nil {
				b.Error(err)
				continue
			}
			if res.StatusCode != http.StatusOK {
				resBody, _ := ioutil.ReadAll(res.Body)
				b.Error(string(resBody))
				err := res.Body.Close()
				if err != nil {
					b.Error(err)
				}
				continue
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				logrus.Errorf("get body err:%v", err)
			}
			err = res.Body.Close()
			if err != nil {
				logrus.Errorf("close body err:%v", err)
			}
			logrus.Infof("info response:%v", string(body[:]))
		}
	})
}
```

golang的benchmark测试用例会默认执行超过1秒，并统计执行情况。
该测试方式，固定用户下请求都命中redis，随机用户则基本穿透到mysql。
用户数据在数据库中的信息概要为：
![用户数据预览](/docs/imgs/data_info.png)
正常http请求返回结果为：
````text
INFO[0003] info response:{"code":1,"msg":"success","data":{"id":18502,"username":"zero60f6807d","nickname":"zero No.60f6807d","profile":"http://127.0.0.1/profile/default.jpg"}}
````

通过Makefile编写基准测试命令 `make benchInfoFix` 、 `make benchInfoRandom` 执行两种测试：
```makefile
benchInfoFix:
	go test -v ./cmd/web/router -test.bench InfoFix 

benchInfoRandom:
	go test -v ./cmd/web/router -test.bench InfoRandom 
```

#### make benchInfoFix 执行结果
并发度200的情况下(清空redis后执行)：
```
goos: darwin
goarch: amd64
pkg: git.garena.com/zhenrong.zeng/entrytask/cmd/web/router
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkInfoFix
benchmark info fix parallelism:200 
benchmark info fix parallelism:200 
benchmark info fix parallelism:200 
benchmark info fix parallelism:200 
BenchmarkInfoFix-12    	  482190	    150885 ns/op
PASS
ok  	git.garena.com/zhenrong.zeng/entrytask/cmd/web/router	74.960s
```

并发度2000的情况下(清空redis后执行)：
```
goos: darwin
goarch: amd64
pkg: git.garena.com/zhenrong.zeng/entrytask/cmd/web/router
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkInfoFix
benchmark info fix parallelism:2000 
benchmark info fix parallelism:2000 
benchmark info fix parallelism:2000 
benchmark info fix parallelism:2000 
BenchmarkInfoFix-12    	  461310	    148644 ns/op
PASS
ok  	git.garena.com/zhenrong.zeng/entrytask/cmd/web/router	70.690s
```

#### make benchInfoRandom 执行结果

并发度200的情况下(清空redis后执行)：
```
goos: darwin
goarch: amd64
pkg: git.garena.com/zhenrong.zeng/entrytask/cmd/web/router
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkInfoRandom
benchmark info random parallelism:200 
benchmark info random parallelism:200 
benchmark info random parallelism:200 
benchmark info random parallelism:200 
benchmark info random parallelism:200 
BenchmarkInfoRandom-12    	  477351	    160185 ns/op
PASS
ok  	git.garena.com/zhenrong.zeng/entrytask/cmd/web/router	122.942s
```

并发度2000的情况下(清空redis后执行)：
```
goos: darwin
goarch: amd64
pkg: git.garena.com/zhenrong.zeng/entrytask/cmd/web/router
cpu: Intel(R) Core(TM) i7-9750H CPU @ 2.60GHz
BenchmarkInfoRandom
benchmark info random parallelism:2000 
benchmark info random parallelism:2000 
benchmark info random parallelism:2000 
benchmark info random parallelism:2000 
benchmark info random parallelism:2000 
BenchmarkInfoRandom-12    	  484503	    171512 ns/op
PASS
ok  	git.garena.com/zhenrong.zeng/entrytask/cmd/web/router	142.591s
```

**结果分析：**


| 测试项 | 固定200用户并发 | 随机200用户并发 | 固定2000用户并发 | 随机2000用户并发 |
| --- | --- | --- | --- | --- |
| 平均耗时 | 150885 ns/op | 160185 ns/op | 144715 ns/op | 171512 ns/op |
| 平均QPS | 482190(times) / 74.960(s) = 6432.63 (t/s) | 477351(times) / 122.942(s) = 3638.72(t/s) | 501991(times) / 74.602(s) = 6728.92(t/s) | 484503(times) / 142.591(s) = 3397.85(t/s) |
| 目标QPS | 3000 | 1000 | 1500 | 800 |


以上测试结果表明项目达到该要求。

实现样式如：
登录页面：
![登录页面](/docs/imgs/home_login.png)


修改昵称页面：
![修改昵称页面](/docs/imgs/home_nick.png)


修改头像页面：
![修改头像页面](/docs/imgs/home_profile.png)
