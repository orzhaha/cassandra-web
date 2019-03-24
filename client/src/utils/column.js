import get from 'lodash/get'
import forEach from 'lodash/forEach'
import includes from 'lodash/includes'
import api from '@/api'
import getType from '@/config/data-type'

const service = api.make('root')

const getTextRect = () => {
  const elId = '__get-text-rect'
  let elNode = document.getElementById(elId)
  if (!elNode) {
    elNode = document.createElement('div')
    elNode.id = elId
    elNode.style.cssText = 'visiblity:hidden; opacity: 0; position: fixed; top: -9999px; left: -9999px'
    document.body.appendChild(elNode)
  }

  return (text, cssText = 'font-size: 14px; padding: 0 10px;') => {
    if (typeof text !== 'string') return {}
    const textElement = `<span style="box-sizing: border-box;${cssText}">${text}</span>`
    elNode.innerHTML = textElement
    const elRect = elNode.getBoundingClientRect()
    return elRect
  }
}

class Column {
  constructor(keyspace, table) {
    this.keyspace = keyspace
    this.table = table
    this.cloumnTextTotalWidth = 0
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

  getCloumnTextTotalWidth() {
    return this.cloumnTextTotalWidth
  }

  async init() {
    const res = await service.request('columns', {
      query: {
        keyspace: this.keyspace,
        table: this.table,
      }
    })

    this.columnData = res.get()

    const gtr = getTextRect()

    forEach(this.columnData, (item, index) => {
      const igtr = gtr(get(item, 'column_name'))
      this.cloumnTextTotalWidth = this.cloumnTextTotalWidth + igtr.width + 6

      if (igtr.width < 80) {
        igtr.width = 80
      }

      this.columnData[index].text_rect = igtr
      this.types[get(item, 'column_name')] = {
        cql: get(item, 'type', ''),
        js: getType(get(item, 'type', '')),
        kind: get(item, 'kind', ''),
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
    if (this.getJSType(key) === 'string') {
      return 'textarea'
    }

    if (this.getJSType(key) === 'object') {
      return 'textarea'
    }

    return 'text'
  }
}

export default Column
