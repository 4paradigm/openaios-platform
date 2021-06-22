<!--
 * @Author: liyuying
 * @Date: 2021-04-28 16:20:43
 * @LastEditors: liyuying
 * @LastEditTime: 2021-05-24 10:56:04
 * @Description: file content
-->

# 这里是全局 Model

## model 注册

参考[model 注册](https://umijs.org/zh/guide/with-dva.html#model-%E6%B3%A8%E5%86%8C)

model 分两类，一是全局 model，二是页面 model。全局 model 存于 /src/models/ 目录，所有页面都可引用；页面 model 不能被其他页面所引用。

### 规则如下

- src/models/\*_/_.js 为 global model
- src/pages/**/models/**/\*.js 为 page model
- global model 全量载入，page model 在 production 时按需载入，在 development 时全量载入
- page model 为 page js 所在路径下 models/\*_/_.js 的文件
- page model 会向上查找，比如 page js 为 pages/a/b.js，他的 page model 为 pages/a/b/models/**/\*.js + pages/a/- models/**/\*.js，依次类推
- 约定 model.js 为单文件 model，解决只有一个 model 时不需要建 models 目录的问题，有 model.js 则不去找 models/\*_/_.js

## 如何定义 model

参考[dva 定义 model](https://dvajs.com/guide/getting-started.html#%E5%AE%9A%E4%B9%89-model)

```js
export default {
  namespace: 'products',
  state: [],
  reducers: {
    delete(state, { payload: id }) {
      return state.filter((item) => item.id !== id);
    },
  },
};
```
