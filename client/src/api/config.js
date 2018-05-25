// const domain = window.location.origin
const domain = 'http://localhost:8081'

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
  }
}

export default config
