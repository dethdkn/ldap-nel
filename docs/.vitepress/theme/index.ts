import { inject } from '@vercel/analytics'
import theme from 'vitepress/theme'
import './extend.css'

inject()

export default {
  ...theme,
}
