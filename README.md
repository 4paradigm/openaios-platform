[![coverage report](https://gitlab.4paradigm.com/adc/pineapple/badges/feat/fe/coverage.svg)](https://gitlab.4paradigm.com/adc/pineapple/-/commits/feat/fe)

# pineapple

异构算力计算平台，目前可以访问[官网](https://openaios.4paradigm.com)查看

## 组件

- pineapple(core): 处理前端的主要请求，包括浏览，创建环境以及应用
- billing: Pod级别处理计费, 记录用户余额，处理没有费用时的关停
- portal: 前端逻辑
- webhook: 对提交的所有Pod按照annotation以及label进行处理
- webterminal [optional]