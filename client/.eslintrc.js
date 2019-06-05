const path = require('path')
module.exports = {
  root: true,
  env: {
    node: true
  },
  extends: ["plugin:vue/essential", 'airbnb-base'],
  // https://github.com/vuejs/eslint-plugin-vue#priority-a-essential-error-prevention
  // consider switching to `plugin:vue/strongly-recommended` or `plugin:vue/recommended` for stricter rules.
  plugins: [
    'vue'
  ],
  settings: {
    'import/resolver': {
      webpack: {
        config: {
          resolve: {
            extensions: ['.js', '.json', '.vue'],
              alias: {
              '@': path.resolve(__dirname, 'src')
            }
          }
        }
      }
    },
    "import/core-modules": ["codemirror"],
  },
  rules: {
    'max-len': ["error", { "code": 200 }],
    // don't require .vue extension when importing
    'import/extensions': ['error', 'always', {
      js: 'never',
      vue: 'never'
    }],
    // disallow reassignment of function parameters
    // disallow parameter object manipulation except for specific exclusions
    'no-param-reassign': ['error', {
      props: true,
      ignorePropertyModificationsFor: [
        'state', // for vuex state
        'acc', // for reduce accumulators
        'e' // for e.returnvalue
      ]
    }],
    // allow optionalDependencies
    'import/no-extraneous-dependencies': ['error', {
      optionalDependencies: ['test/unit/index.js']
    }],
    // allow debugger during development
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'off',

    semi: 0,
    'comma-dangle': 0
  },
  parserOptions: {
    parser: "babel-eslint"
  }
};
