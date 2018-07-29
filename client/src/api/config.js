const domain = window.location.origin
// const domain = 'http://localhost:8083'

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
  }
}

export default config
