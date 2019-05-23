<template lang="pug">
  div(
    class="tieme-picker-div")
    el-color-picker(
      class="ml20"
      v-model="theme"
      class="theme-picker"
      popper-class="theme-picker-dropdown")
    //- el-switch(
    //-   class="pb20 ml20"
    //-   active-text="theme dark"
    //-   inactive-text="theme white"
    //-   active-color="#FFFFFF"
    //-   inactive-color="#000000"
    //-   @change="changeTheme"
    //-   v-model="isThemeDark"
    //- )
</template>

<style>
  .tieme-picker-div .ml20 {
    margin-right: 20px;
  }
  .tieme-picker-div .pb20 {
    padding-bottom: 20px;
  }
</style>

<script>
import { version } from 'element-ui/package.json' // element-ui version from node_modules
import Cookies from 'js-cookie'
import concat from 'lodash/concat'

const ORIGINAL_THEME = '#409EFF' // default color

export default {
  data() {
    return {
      chalk: '', // content of theme-chalk css
      theme: ORIGINAL_THEME,
      isThemeDark: false,
    }
  },
  watch: {
    theme(val) {
      const oldVal = this.theme
      if (typeof val !== 'string') return
      const themeCluster = this.getThemeCluster(val.replace('#', ''))
      const originalCluster = this.getThemeCluster(oldVal.replace('#', ''))
      // console.log(themeCluster, originalCluster)
      const getHandler = (variable, id) => (() => {
        const original = this.getThemeCluster(ORIGINAL_THEME.replace('#', ''))
        const newStyle = this.updateStyle(this[variable], original, themeCluster)
        let styleTag = document.getElementById(id)
        if (!styleTag) {
          styleTag = document.createElement('style')
          styleTag.setAttribute('id', id)
          document.head.appendChild(styleTag)
        }
        styleTag.innerText = newStyle
      })
      const chalkHandler = getHandler('chalk', 'chalk-style')
      if (!this.chalk) {
        const url = `https://unpkg.com/element-ui@${version}/lib/theme-chalk/index.css`
        this.getCSSString(url, chalkHandler, 'chalk')
      } else {
        chalkHandler()
      }
      let styles = [].slice.call(document.querySelectorAll('style'))
        .filter((style) => {
          const text = style.innerText

          return new RegExp(oldVal, 'i').test(text) && !/Chalk Variables/.test(text)
        })

      const newStyle = concat([], styles)

      newStyle.forEach((style, index) => {
        const { innerText } = style

        if (typeof innerText !== 'string') return

        newStyle[index].innerText = this.updateStyle(innerText, originalCluster, themeCluster)
      })

      styles = newStyle

      this.$message({
        message: 'The theme is changed!',
        type: 'success'
      })
    }
  },
  created() {
    const isThemeDark = Cookies.get('isThemeDark')

    if (isThemeDark !== undefined) {
      this.isThemeDark = isThemeDark === 'true'
    }
  },
  methods: {
    changeTheme() {
      Cookies.set('isThemeDark', this.isThemeDark)
      this.$router.go(0)
    },
    updateStyle(style, oldCluster, newCluster) {
      const colorOverrides = [] // only capture color overides
      oldCluster.forEach((color, index) => {
        const value = newCluster[index]
        const colorPlain = color.replace(/([()])/g, '\\$1')
        const repl = new RegExp(`(^|})([^{]+{[^{}]+)${colorPlain}\\b([^}]*)(?=})`, 'gi')
        const nestRepl = new RegExp(colorPlain, 'ig') // for greed matching before the 'color'
        let v

        /* eslint no-cond-assign: "error" */

        while ((v = repl.exec(style))) {
          colorOverrides.push(`${v[2].replace(nestRepl, value)}${value}${v[3]}}`) // '}' not captured in the regexp repl to reserve it as locator-boundary
        }
      })
      return colorOverrides.join('')
    },
    getCSSString(url, callback, variable) {
      const xhr = new XMLHttpRequest()
      xhr.onreadystatechange = () => {
        if (xhr.readyState === 4 && xhr.status === 200) {
          this[variable] = xhr.responseText.replace(/@font-face{[^}]+}/, '')
          callback()
        }
      }
      xhr.open('GET', url)
      xhr.send()
    },
    getThemeCluster(theme) {
      const tintColor = (color, tint) => {
        let red = parseInt(color.slice(0, 2), 16)
        let green = parseInt(color.slice(2, 4), 16)
        let blue = parseInt(color.slice(4, 6), 16)

        if (tint === 0) { // when primary color is in its rgb space
          return [red, green, blue].join(',')
        }

        red += Math.round(tint * (255 - red))
        green += Math.round(tint * (255 - green))
        blue += Math.round(tint * (255 - blue))
        red = red.toString(16)
        green = green.toString(16)
        blue = blue.toString(16)

        return `#${red}${green}${blue}`
      }
      const shadeColor = (color, shade) => {
        let red = parseInt(color.slice(0, 2), 16)
        let green = parseInt(color.slice(2, 4), 16)
        let blue = parseInt(color.slice(4, 6), 16)
        red = Math.round((1 - shade) * red)
        green = Math.round((1 - shade) * green)
        blue = Math.round((1 - shade) * blue)
        red = red.toString(16)
        green = green.toString(16)
        blue = blue.toString(16)
        return `#${red}${green}${blue}`
      }
      const clusters = [theme]
      for (let i = 0; i <= 9; i += 1) {
        clusters.push(tintColor(theme, Number((i / 10).toFixed(2))))
      }
      clusters.push(shadeColor(theme, 0.1))
      return clusters
    }
  }
}
</script>

<style>
.theme-picker .el-color-picker__trigger {
  margin-top: 12px;
  height: 26px!important;
  width: 26px!important;
  padding: 2px;
}
.theme-picker-dropdown .el-color-dropdown__link-btn {
  display: none;
}
</style>
