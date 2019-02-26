import get from 'lodash/get'
import forEach from 'lodash/forEach'
import keys from 'lodash/keys'
import last from 'lodash/last'
import filter from 'lodash/filter'

const MAPPING = {
  ascii: 'string',
  bigint: 'number', // NOTE: int64
  blob: 'string', // NOTE: bytes
  boolean: 'boolean',
  counter: 'number', // NOTE: int64
  date: 'string',
  decimal: 'number',
  double: 'number', // NOTE: int64
  float: 'number',
  frozen: 'object', // NOTE: blob
  inet: 'string',
  int: 'number',
  list: 'object', // NOTE: array
  map: 'object',
  set: 'object',
  smallint: 'number',
  text: 'string',
  time: 'number', // NOTE: int64
  timestamp: 'any', // NOTE: string or number
  timeuuid: 'any', // NOTE: uuid
  tinyint: 'number',
  tuple: 'object',
  uuid: 'any', // NOTE: uuid
  varchar: 'string',
  varint: 'number'
}

let regExpRule = ''
const objectItems = filter(keys(MAPPING), key => MAPPING[key] === 'object')

forEach(objectItems, (key) => {
  regExpRule += `^(${key})`

  if (key !== last(objectItems)) {
    regExpRule += '|'
  }
})

const objectRegExp = new RegExp(regExpRule, 'g')

export default function getType(cqlType) {
  if (get(MAPPING, cqlType)) {
    return get(MAPPING, cqlType)
  }

  if (cqlType.match(objectRegExp)) {
    return 'object'
  }

  return 'string'
}
