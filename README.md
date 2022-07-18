# go-smart
golang 自动化工具

安装
> go install github.com/cocoup/go-smart/tools/gocli@latest

**参考和借鉴了其他第三方库，并做了简单封装和修改**

> 代码自动生成及中间件参考了目前比较火的[go-zero](https://github.com/zeromicro/go-zero)框架
> 
> http模块是基于[gin](http://badi.com)做了简单封装，并集成了部分go-zero的中间件
>
> 数据库模块是基于[gorm](https://gorm.io/zh_CN/)的简单封装
>
> 缓存redis基于[go-redis](https://github.com/go-redis/redis)的简单封装

 目前已完成`http`框架和`mysql`代码的自动生成，后续会完善`Grpc`及其他语言和框架的生成及服务治理相关逻辑的集成。