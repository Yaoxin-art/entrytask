# entrytask

## 一、背景及目的
EntryTask项目，通过实现RPC通用框架，使用Go HTTP API，MySQL或Redis操作等能力完成一个用户管理系统，主要功能包含用户（管理员）登录及用户数据的增删查改功能。

## 二、逻辑架构设计
整体分为网关（nginx代）、API层（RPC consumer）、TCP Server层（RPC provider）三层，图片存放于磁盘，前端访问则通过nginx代理，通过配置项确定文件（头像图片）存放的位置，以及配置组装完整访问URL。
![整体架构图](/docs/imgs/framework01.png)

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
接口文档地址：[接口文档](https://confluence.shopee.io/pages/viewpage.action?pageId=597241979 "CF") 

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
    Makefile包含 clean、fmt、vet、cover、test、build、run 功能，使用`make run [web|worker]`启动。
    启动顺序：
        启动mysql服务
        启动redis服务
        启动nginx
        启动web模块
        启动worker模块


## 八、SLA

### coverage
![make cover](/docs/imgs/make_cover.png)

### testcase
![make test](/docs/imgs/make_test.png)

### benchmark
TODO：压测报告
