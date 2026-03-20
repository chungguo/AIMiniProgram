import { defineConfig } from 'unocss';
import presetWeapp from 'unocss-preset-weapp';

export default defineConfig({
  presets: [
    presetWeapp()
  ],
  rules: [
    ['rpx', { unit: 'rpx' }]
  ],
  shortcuts: {
    'flex-center': 'flex items-center justify-center',
    'flex-between': 'flex items-center justify-between',
    'flex-col': 'flex flex-col',
    'text-ellipsis': 'overflow-hidden whitespace-nowrap text-ellipsis',
    'card': 'bg-white rounded-16rpx p-24rpx shadow-sm'
  }
});
