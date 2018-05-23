// const domain = window.location.origin
const domain = "http://localhost:80"

const config = {
    root: {
        allKeySpace: {
            type: 'GET',
            url: `${domain}/allkeyspace`
        },
        alltablebykeyspace: {
            type: 'GET',
            url: `${domain}/alltablebykeyspace`
        },
        allrowbytable: {
            type: 'GET',
            url: `${domain}/allrowbytable`
        },
    }
}

export default config