import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/home'
import Main from '@/components/main'
import TableList from '@/components/table-list'
import Table from '@/components/table'
import Content from '@/components/table/content'
import Rows from '@/components/table/rows'
import Columns from '@/components/table/columns'
import Definition from '@/components/table/definition'
import Export from '@/components/table/export'
import Import from '@/components/table/import'
import Query from '@/components/query'

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/main',
      name: 'Main',
      component: Main,
      children: [
        {
          path: 'query/:query',
          name: 'query',
          component: Query,
        },
        {
          path: 'keyspace/:keyspace',
          name: 'table-list',
          component: TableList
        },

        {
          path: 'keyspace/:keyspace/:table',
          component: Table,
          children: [
            {
              path: '/',
              name: 'content',
              component: Content,
            },

            {
              path: 'rows/:page/:pagesize',
              name: 'rows',
              component: Rows,
            },

            {
              path: 'columns',
              name: 'columns',
              component: Columns,
            },

            {
              path: 'definition',
              name: 'definition',
              component: Definition,
            },

            {
              path: 'export',
              name: 'export',
              component: Export,
            },

            {
              path: 'import',
              name: 'import',
              component: Import,
            },

            {
              path: '*',
              redirect: { name: 'content' }
            }
          ],
        }
      ]
    },
    {
      path: '*',
      redirect: { name: 'home' }
    },
  ],
});
