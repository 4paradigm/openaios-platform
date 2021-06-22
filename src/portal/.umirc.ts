/*
 * @Author: liyuying
 * @Date: 2021-04-23 11:55:28
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-11 14:16:03
 * @Description: file content
 */
import { defineConfig } from 'umi';
const proxyTarget = 'https://develop.pineapple-test.com/';
import { routes } from './src/router/routes';
const isProduction = process.env.NODE_ENV === 'production';

export default defineConfig({
  publicPath: isProduction ? '/' : '/',
  title: 'AIOS 社区版',
  links: [
    // href的图片你可以放在public里面，直接./图片名.png 就可以了
    { rel: 'icon', href: './favicon.png' },
  ],
  routes: routes,
  dynamicImport: {
    loading: '@/components/Loading',
  },
  dva: { skipModelValidate: true },
  hash: true,
  proxy: {
    '/api': {
      target: proxyTarget,
      changeOrigin: true,
      pathRewrite: {
        '^': '',
      },
      secure: false,
    },
    '/web-terminal': {
      target: proxyTarget,
      changeOrigin: true,
      pathRewrite: {
        '^': '',
      },
      secure: false,
    },
    '/terminal': {
      target: proxyTarget,
      changeOrigin: true,
      pathRewrite: {
        '^': '',
      },
      secure: false,
    },
    '/api/auth/pineapple': {
      target: proxyTarget,
      changeOrigin: true,
      pathRewrite: {
        '^': 'auth',
      },
      secure: false,
    },
  },
  // dll: false,
  // disableCSSModules: false,
  lessLoader: {
    modifyVars: {
      // 或者可以通过 less 文件覆盖（文件路径为绝对路径）,注入全局less变量
      hack: `true; @import "~@/assets/styles/global.less"`,
    },
  },
  chainWebpack(config) {
    config.module
      .rule('customfont')
      .test(/\.otf$/)
      .use('file-loader')
      .loader('file-loader');
  },
  copy: [{ from: 'node_modules/libarchive.js/dist', to: './' }],
});
