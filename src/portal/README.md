# pineapple

> AIOS 社区版

## 一、Development

```bash
# 安装依赖包
npm ci

# 启动本地环境
npm run start

# 构建生产包
npm run build

# 测试
npm run test

# 格式化
npm run lint
```

[store 即 models 开发参考](./src/models/README.md)

## 二、前端技术选型

前端技术选型
TypeScript + UMI@3 + React Hooks + DVA + AntD@4 + LESS + Jest + cess-ui

简单介绍

- TypeScript: 主编程语言
- react hooks: 不用再写 class
- roadhog 是基于 webpack 的封装工具，目的是简化 webpack 的配置
- umi 可以简单地理解为 roadhog + 路由(高级玩法，如权限路由)
- dva 目前是纯粹的数据流，简单理解为 Redux + react-router
- less: CSS 预处理器，推荐使用 less，为跟 **antd** 保持一致
- AntD: React 组件库
- Jest: 单元测试库

环境要求：
node@^10

- umijs: <https://umijs.org/>
- dvajs: <https://dvajs.com/>
- react hooks: <https://reactjs.org/docs/hooks-intro.html>
- antd: <https://ant.design/docs/react/introduce-cn>
