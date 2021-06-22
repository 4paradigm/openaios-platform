<!--
 * @Author: liyuying
 * @Date: 2021-05-11 13:52:32
 * @LastEditors: liyuying
 * @LastEditTime: 2021-05-19 13:59:12
 * @Description: file content
-->

# openApi

## 一、生成 SDK（参考文档：https://openapi-generator.tech/docs/installation）

```bash
# 安装依赖包（首次执行）
npm install @openapitools/openapi-generator-cli -g

# 生成SDK
cd src/portal/src/openApi

npx @openapitools/openapi-generator-cli generate -i ../../../../doc/api/main.yaml -g typescript-axios -o api


# 导出api实例
在openApi文件夹下创建index.ts，并实例化SDK中导出的api

```
