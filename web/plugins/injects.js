import Vue from 'vue';
// import 'lib-flexible';

import VueAwesomeSwiper from 'vue-awesome-swiper';
import 'swiper/css/swiper.min.css';
import VueClipboard from 'vue-clipboard2';

VueClipboard.config.autoSetContainer = true; // add this line

export default ({ store }) => {
  Vue.use(VueAwesomeSwiper);
  Vue.use(VueClipboard);
};
