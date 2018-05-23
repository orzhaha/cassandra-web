import ServiceProxy from 'service-proxy'
import config from './config'

const sp = new ServiceProxy(config, { cors: false })

export default sp
