import Vue from 'vue';
import Router from 'vue-router';
import Home from '@/components/home';
import Keyspace from '@/components/keyspace';
import Table from '@/components/table';
import Row from '@/components/row';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/keyspace',
      name: 'keyspace',
      component: Keyspace,
      children: [
        {
          path: ':keyspace/table',
          name: 'table',
          component: Table,
        },
        {
          path: ':table/row',
          name: 'row',
          component: Row,
        }
      ]
    },
  ],
});
