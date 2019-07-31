module.exports = {
  outputDir: 'dist',
  assetsDir: 'static',
  pages: {
    index: {
      // page 的入口
      entry: 'src/main.js',
      // 模板來源
      template: 'src/index.pug',
      // 在 dist/index.html 的輸出
      filename: 'index.html',
      // 當使用 title 選項時，
      // template 中的 title 標籤需要是 <title><%= htmlWebpackPlugin.options.title %></title>
      title: 'Index Page',
      // 在這個頁面中包含的塊，默認情況下會包含
      // 提取出來的通用 chunk 和 vendor chunk。
      chunks: ['chunk-vendors', 'chunk-common', 'index']
    },
  },
}