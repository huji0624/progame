import Vue from 'vue';
import vant from 'vant';
import 'vant/lib/index.css';
import VConsole from 'vconsole/dist/vconsole.min.js';
import { isDebug } from '../core/config';

export default () => {
  isDebug && new VConsole();
  Vue.use(vant);
};
