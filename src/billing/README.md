# 计费系统 Billing server

该模块主要承接pineapple的计费功能，其中包含有用户的账户管理，算力规格的管理，以及用户任务的计费。下面简单介绍每个子目录的相关信息：

## conf
该模块主要是读取相关的环境变量，目前主要是mongodb相关。

## handler
实现了对外的接口（主要是对webserver以及webhook），接口文档是`doc/api/billing.yaml`
注意：instance这个接口已经被废弃，但是还没有删除。

## utils
实现了一些工具包，方便handler使用

- billing-utils.go：更改用户账户信息的相关函数，主要被`handler/account.go`使用。

- computeunit-utils.go：查询computeunit以及group的相关函数，主要被`handler/computeunit.go`使用

- k8s-utils.go：k8s客户端，主要被heartbeat使用

## heartbeat
用于计费，使用go rountine执行，每分钟执行一次，对用户进行扣费，在log中输出所有无法被计费的pod。
当用户的余额不足时，执行用户的callback方法，向webserver发请求，杀死该用户的所有任务。

## 一些设想
通过go channel的方式实现一个消息队列，使得一些请求作为一个独立的模块向同一个模块发送信号，实现加钱或者扣钱操作。
可以参考 https://segmentfault.com/a/1190000024518618 这篇博客。
