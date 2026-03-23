import { readFileSync } from 'fs'
import { resolve } from 'path'

const SNIPPETS_DIR = resolve('./snippets')

export default {
  plugins: [
    {
      name: 'html-snippets',
      
      // dev + build 都会执行
      transformIndexHtml(html) {
        const header = readFileSync(resolve(SNIPPETS_DIR, 'header.html'), 'utf-8')
        return html.replace('<!-- inject:header -->', header)
      },
      
      // dev server 配置：监听 snippets 目录变化
      configureServer(server) {
        server.watcher.add(SNIPPETS_DIR)
      },
      
      // snippets 修改时触发 index.html 热更新
      handleHotUpdate({ file, server }) {
        if (file.startsWith(SNIPPETS_DIR)) {
          // 通知客户端刷新（或者只刷新 HTML）
          server.ws.send({ type: 'full-reload' })
        }
      },
    },
  ],
}