const webpack = require('webpack')

// NOTE: key 為編譯過後實際檔案路徑與名稱
const vendors = {
  'vendor/element-theme-dark': [
    'element-theme-dark'
  ],
}

// NOTE: 取 vendor 內的所有套件，以供實作
let allModules = Object.keys(vendors).reduce((modules, fileName) => {
  return modules.concat(vendors[fileName])
}, [])

let lastFileName = null

module.exports = (
  Object.keys(vendors).map((fileName, i, allFilesName) => {
    return new webpack.optimize.CommonsChunkPlugin({
      name: fileName,
      minChunks: (module) => {
        if (i === 0 && lastFileName == null) {
          lastFileName = fileName
        }

        if (lastFileName !== fileName) {
          const numberOfDependenciesOnPreviousVendorFile = vendors[allFilesName[i - 1]].length

          allModules = allModules.slice(numberOfDependenciesOnPreviousVendorFile)

          lastFileName = fileName
        }

        return (
          module.resource && allModules.some(dependency => module.resource.indexOf(`node_modules/${dependency}/`) > -1)
        )
      }
    })
  })
)