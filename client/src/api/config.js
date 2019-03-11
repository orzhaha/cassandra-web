export const domain = window.location.origin
// export const domain = 'http://localhost:8083'

const config = {
  root: {
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
    query: {
      type: 'POST',
      url: `${domain}/query`
    },
    save: {
      type: 'POST',
      url: `${domain}/save`
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
    import: {
      type: 'POST',
      url: `${domain}/import`
    }
  }
}

export default config
