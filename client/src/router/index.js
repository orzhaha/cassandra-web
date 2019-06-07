import Vue from 'vue'
import Router from 'vue-router'
import Home from '@/components/home'
import Main from '@/components/main'
import TableList from '@/components/table-list'
import Table from '@/components/table'
import Content from '@/components/table/content'
import Rows from '@/components/table/rows'
import RowsToken from '@/components/table/rowstoken'
import Columns from '@/components/table/columns'
import Definition from '@/components/table/definition'
import Export from '@/components/table/export'
import Import from '@/components/table/import'
import Truncate from '@/components/table/truncate'
import Find from '@/components/table/find'
import Query from '@/components/query'
import HostInfo from '@/components/hostinfo'

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
          path: 'hostinfo',
          name: 'hostinfo',
          component: HostInfo,
        },

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
              path: 'rowstoken/:pagesize',
              name: 'rowstoken',
              component: RowsToken,
            },

            {
              path: 'find',
              name: 'find',
              component: Find,
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
              path: 'truncate',
              name: 'truncate',
              component: Truncate,
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
