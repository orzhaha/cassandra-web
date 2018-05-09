/*
* NOTE: element components
* https://github.com/ElemeFE/element/blob/master/components.json
*/
import Vue from 'vue';
import { Loading } from 'element-ui';
import '@/style/theme/index.css';

Vue.use(Loading.directive);

Vue.prototype.$loading = Loading.service;
