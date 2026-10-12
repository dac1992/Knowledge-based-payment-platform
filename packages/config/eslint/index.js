import js from '@eslint/js'
import globals from 'globals'
import tseslint from 'typescript-eslint'
import pluginVue from 'eslint-plugin-vue'
import prettier from 'eslint-config-prettier'

/**
 * 生成各端共用的 ESLint 扁平配置（Flat Config）。
 *
 * @param {{ tsconfigRootDir?: string }} [options]
 * @returns {import('eslint').Linter.Config[]}
 */
export function createEslintConfig(options = {}) {
  return tseslint.config(
    {
      ignores: [
        '**/dist/**',
        '**/node_modules/**',
        '**/unpackage/**',
        '**/coverage/**',
        '**/*.d.ts',
        '**/auto-imports.d.ts',
        '**/components.d.ts',
      ],
    },
    js.configs.recommended,
    ...tseslint.configs.recommended,
    ...pluginVue.configs['flat/recommended'],
    {
      files: ['**/*.vue'],
      languageOptions: {
        parserOptions: {
          parser: tseslint.parser,
          ecmaVersion: 'latest',
          sourceType: 'module',
          tsconfigRootDir: options.tsconfigRootDir,
        },
      },
    },
    {
      files: ['**/*.{js,mjs,cjs,ts,mts,tsx,vue}'],
      languageOptions: {
        ecmaVersion: 'latest',
        sourceType: 'module',
        globals: {
          ...globals.browser,
          ...globals.node,
          ...globals.es2022,
        },
      },
      rules: {
        '@typescript-eslint/no-unused-vars': [
          'warn',
          { argsIgnorePattern: '^_', varsIgnorePattern: '^_' },
        ],
        '@typescript-eslint/no-explicit-any': 'off',
        '@typescript-eslint/consistent-type-imports': 'off',

        // 版式交由 Prettier 负责，避免两条规则互相打架
        'vue/multi-word-component-names': 'off',
        'vue/max-attributes-per-line': 'off',
        'vue/singleline-html-element-content-newline': 'off',
        'vue/html-self-closing': 'off',
        'vue/html-indent': 'off',
        'vue/html-closing-bracket-newline': 'off',
        'vue/attributes-order': 'off',
        'vue/first-attribute-linebreak': 'off',
        'vue/require-default-prop': 'off',
      },
    },
    prettier,
  )
}

export default createEslintConfig
