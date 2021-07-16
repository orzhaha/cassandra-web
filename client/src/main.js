// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import '@/plugins/element-ui';
import '@/plugins/vue-codemirror'
import Vue from 'vue';
import App from './App';
import router from './router';
import api from '@/api'

Vue.config.productionTip = false;

const service = api.make('root')


async function init() {
  try {
    const res = await service.request('readonly')
    const ret = res.get()

    Vue.prototype.$readonly = ret.readonly

    /* eslint-disable no-new */
    new Vue({
      el: '#app',
      router,
      components: { App },
      template: '<App/>',
      render: h => h(App),
    });
  } catch (error) {
    // eslint-disable-next-line
    console.log(error)
  }
}

init();
