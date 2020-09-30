import { apiUrl } from './core/config';

export default {
  /*
   ** Nuxt rendering mode
   ** See https://nuxtjs.org/api/configuration-mode
   */
  mode: 'spa',
  /*
   ** Nuxt target
   ** See https://nuxtjs.org/api/configuration-target
   */
  target: 'server',
  /*
   ** Headers of the page
   ** See https://nuxtjs.org/api/configuration-head
   */
  head: {
    title: '程序员节游戏',
    meta: [
      { charset: 'utf-8' },
      {
        name: 'viewport',
        content:
          'width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=0, viewport-fit=cover',
      },
      {
        hid: 'description',
        name: 'description',
        content: process.env.npm_package_description || '',
      },
      { httpEquiv: 'X-UA-Compatible', name: 'IE=edge, chrome=1' },
      { name: 'format-detection', content: 'telphone=no, email=no' }, //忽略页面中的数字识别为电话，忽略email识别
      { name: 'apple-mobile-web-app-status-bar-style', content: 'black' }, //苹果工具栏颜色
      { name: 'apple-mobile-web-app-capable', content: 'yes' }, //启用 WebApp 全屏模式，删除苹果默认的工具栏和菜单栏
      { name: 'msapplication-tap-highlight', content: 'no' }, //windows phone 点击无高光
      { name: 'HandheldFriendly', content: 'true' }, //针对手持设备优化
    ],
    link: [{ rel: 'icon', type: 'image/x-icon', href: '/favicon.ico' }],
    script: [
      {
        src: 'https://hm.baidu.com/hm.js?427c3e9ab87e442121e9cb39aca7231f',
      },
      {
        src:
          'https://3gimg.qq.com/lightmap/components/geolocation/geolocation.min.js', //腾讯地图组件
      },
    ],
  },
  /*
   ** Global CSS
   */
  css: [
    '~assets/css/main.css',
    'element-ui/lib/theme-chalk/index.css',
    'normalize.css',
  ],
  /*
   ** Plugins to load before mounting the App
   ** https://nuxtjs.org/guide/plugins
   */
  plugins: [
    // '@/plugins/vant',
    '@/plugins/element-ui',
    '@/plugins/injects',
    '@/plugins/utils',
    '@/plugins/axios',
  ],
  /*
   ** Auto import components
   ** See https://nuxtjs.org/api/configuration-components
   */
  components: true,
  router: {
    base: process.env.NODE_ENV === 'production' ? '/20201024' : '/',
    fallback: true,
  },
  /*
   ** Nuxt.js dev-modules
   */
  buildModules: [],
  /*
   ** Nuxt.js modules
   */
  modules: [
    // Doc: https://axios.nuxtjs.org/usage
    '@nuxtjs/axios',
  ],
  axios: {
    proxy: true,
    credentials: true,
  },
  proxy: {
    '/api': {
      target: apiUrl,
      changeOrigin: true,
      // pathRewrite: {
      //   "^/api": ""
      // }
    },
  },
  /*
   ** Build configuration
   ** See https://nuxtjs.org/api/configuration-build/
   */
  build: {
    transpile: [/^element-ui/],
    postcss: {
      // 添加插件名称作为键，参数作为值
      // 使用npm或yarn安装它们
      plugins: {
        'postcss-pxtorem': false,
        autoprefixer: {},
        // 'postcss-pxtorem': false
      },
      preset: {
        // 更改postcss-preset-env 设置
        autoprefixer: true,
      },
    },
    extend(config, ctx) {
      config.module.rules.push({
        test: /\.(ogg|mp3|wav|mpe?g)$/i,
        loader: 'file-loader',
        options: {
          name: '[path][name].[ext]',
        },
      });
    },
  },
};
