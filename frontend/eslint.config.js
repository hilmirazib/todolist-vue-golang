import pluginVue from 'eslint-plugin-vue'
import { defineConfigWithVueTs, vueTsConfigs } from '@vue/eslint-config-typescript'
import vuePrettier from '@vue/eslint-config-prettier'

export default defineConfigWithVueTs(
  {
    ignores: ['dist', 'node_modules', '*.config.*'],
  },
  pluginVue.configs['flat/recommended'],
  vueTsConfigs.recommended,
  vuePrettier
)