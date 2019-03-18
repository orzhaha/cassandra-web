import get from 'lodash/get'
import forEach from 'lodash/forEach'
import includes from 'lodash/includes'
import api from '@/api'
import getType from '@/config/data-type'

const service = api.make('root')


class Column {
  constructor(keyspace, table) {
    this.keyspace = keyspace
    this.table = table
    this.partitionKey = []
    this.clusteringKey = []
    this.columnData = []
    this.types = {}
  }

  getKeyspace() {
    return this.keyspace
  }

  getTable() {
    return this.table
  }

  getColumnData() {
    return this.columnData
  }

  async init() {
    const res = await service.request('columns', {
      query: {
        keyspace: this.keyspace,
        table: this.table,
      }
    })

    this.columnData = res.get()

    forEach(this.columnData, (item) => {
      if (get(item, 'column_name')) {
        this.types[item.column_name] = {
          cql: get(item, 'type', ''),
          js: getType(get(item, 'type', '')),
          kind: get(item, 'kind', ''),
        }
      }

      if (get(item, 'kind') === 'partition_key') {
        this.partitionKey.push(get(item, 'column_name'))
      }

      if (get(item, 'kind') === 'clustering') {
        this.clusteringKey.push(get(item, 'column_name'))
      }
    })
  }

  isPartitionKey(key) {
    return includes(this.partitionKey, key)
  }

  isClusteringKey(key) {
    return includes(this.clusteringKey, key)
  }

  getCqlType(key) {
    return get(this.types, `${key}.cql`, '')
  }

  getJSType(key) {
    return get(this.types, `${key}.js`, 'string')
  }

  inputType(key) {
    if (this.getJSType(key) === 'number') {
      return ''
    }

    if (this.getJSType(key) === 'string') {
      return 'textarea'
    }

    return null
  }
}

export default Column
