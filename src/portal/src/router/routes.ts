/*
 * @Author: liyuying
 * @Date: 2021-05-24 17:56:37
 * @LastEditors: liyuying
 * @LastEditTime: 2021-06-09 17:31:21
 * @Description: file content
 */
export const routes = [
  // ⚠注意：component是相对于pages文件夹的路径
  {
    path: '/',
    component: '../layouts/index',
    routes: [
      // 有侧边栏页面
      {
        path: './',
        component: '../layouts/WithAsideNav',
        routes: [
          // 首页
          { path: 'home', exact: true, component: './Home' },
          { path: 'home/message/:id', exact: true, component: './Home/MessageDetail' },
          // 开发环境
          { path: 'devEnvironment', component: './DevEnvironment' },
          { path: 'devEnvironment/create', component: './DevEnvironment/Create' },
          // 应用管理
          { path: 'application_instance', component: './Application/Instance' },
          { path: 'application_instance/create/:name', component: './Application/Instance/Create' },
          { path: 'application_instance/:name', component: './Application/Instance/Detail' },
          { path: 'application_chart', component: './Application/Chart' },
          { path: 'application_chart/create', component: './Application/Chart/Create' },
          // 文件管理
          { path: 'file', component: './File' },
          // 镜像管理
          { path: 'private_mirror', component: './Mirror/Private' },
          { path: 'public_mirror', component: './Mirror/Public' },
          {
            path: '500',
            component: './500',
          },
          { component: './404' },
        ],
      },
    ],
  },
  {
    path: '*',
    component: './redirect',
  },
];
