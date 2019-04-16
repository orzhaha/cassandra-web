export const domain = window.location.origin
// export const domain = 'http://localhost:8083'

const config = {
  root: {
    find: {
      type: 'POST',
      url: `${domain}/find`
    },
    query: {
      type: 'POST',
      url: `${domain}/query`
    },
    save: {
      type: 'POST',
      url: `${domain}/save`
    },
    delete: {
      type: 'POST',
      url: `${domain}/delete`
    },
    import: {
      type: 'POST',
      url: `${domain}/import`
    },
    rowtoken: {
      type: 'POST',
      url: `${domain}/rowtoken`
    },
    keyspace: {
      type: 'GET',
      url: `${domain}/keyspace`
    },
    table: {
      type: 'GET',
      url: `${domain}/table`
    },
    row: {
      type: 'GET',
      url: `${domain}/row`
    },
    columns: {
      type: 'GET',
      url: `${domain}/columns`
    },
    describe: {
      type: 'GET',
      url: `${domain}/describe`
    },
    export: {
      type: 'GET',
      url: `${domain}/export`
    },
    hostinfo: {
      type: 'GET',
      url: `${domain}/hostinfo`
    }
  }
}

export default config
