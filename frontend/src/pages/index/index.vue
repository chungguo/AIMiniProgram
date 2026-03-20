<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { modelService, paperService } from '@/services';
import { getModalityList, getFamilyName, formatPrice } from '@/utils/modelHelpers';
import type { Model, Paper } from '@/types/api';

const featuredModels = ref<Model[]>([]);
const latestPapers = ref<Paper[]>([]);
const loading = ref<boolean>(true);

onMounted(() => {
  loadData();
});

async function loadData(): Promise<void> {
  try {
    const [modelsRes, papersRes] = await Promise.all([
      modelService.getModels({ limit: 4 }),
      paperService.getPapers({ limit: 3 })
    ]);
    
    if (modelsRes.success) featuredModels.value = modelsRes.data;
    if (papersRes.success) latestPapers.value = papersRes.data;
  } catch (error) {
    console.error('加载失败:', error);
  } finally {
    loading.value = false;
  }
}

function goToModels(): void {
  uni.switchTab({ url: '/pages/models/models' });
}

function goToPapers(): void {
  uni.switchTab({ url: '/pages/papers/papers' });
}

function goToModelDetail(id: string): void {
  uni.navigateTo({ url: `/pages/model-detail/model-detail?id=${id}` });
}

function goToPaperDetail(id: string): void {
  uni.navigateTo({ url: `/pages/paper-detail/paper-detail?id=${id}` });
}
</script>

<template>
  <view class="min-h-screen bg-gray-50 pb-30">
    <!-- Header -->
    <view class="bg-gradient-to-br from-indigo-900 to-blue-900 px-32rpx pt-48rpx pb-40rpx">
      <view class="flex-between mb-24rpx">
        <text class="text-36rpx font-bold text-white tracking-wide">AI Hub</text>
        <t-icon name="app" size="64rpx" color="#6366f1" />
      </view>
      <text class="text-28rpx text-white/70 leading-relaxed">探索大模型与前沿论文</text>
    </view>

    <!-- Quick Actions -->
    <view class="px-32rpx py-24rpx flex gap-20rpx">
      <t-cell 
        class="flex-1 rounded-24rpx overflow-hidden"
        :hover="true"
        @click="goToModels"
      >
        <template #leftIcon>
          <view class="w-72rpx h-72rpx rounded-20rpx bg-indigo-100 flex-center mr-20rpx">
            <t-icon name="chart" size="40rpx" color="#6366f1" />
          </view>
        </template>
        <template #title>
          <text class="text-32rpx font-bold text-gray-900">模型对比</text>
        </template>
        <template #description>
          <text class="text-24rpx text-gray-500">主流大模型参数</text>
        </template>
        <template #rightIcon>
          <t-icon name="chevron-right" size="36rpx" color="#9ca3af" />
        </template>
      </t-cell>
      
      <t-cell 
        class="flex-1 rounded-24rpx overflow-hidden"
        :hover="true"
        @click="goToPapers"
      >
        <template #leftIcon>
          <view class="w-72rpx h-72rpx rounded-20rpx bg-blue-100 flex-center mr-20rpx">
            <t-icon name="article" size="40rpx" color="#3b82f6" />
          </view>
        </template>
        <template #title>
          <text class="text-32rpx font-bold text-gray-900">论文研读</text>
        </template>
        <template #description>
          <text class="text-24rpx text-gray-500">中英双语</text>
        </template>
        <template #rightIcon>
          <t-icon name="chevron-right" size="36rpx" color="#9ca3af" />
        </template>
      </t-cell>
    </view>

    <!-- Featured Models -->
    <view class="mt-48rpx px-32rpx">
      <view class="flex-between mb-24rpx">
        <text class="text-36rpx font-bold text-gray-900">热门模型</text>
        <t-link theme="primary" @click="goToModels">
          查看全部 <t-icon name="chevron-right" size="24rpx" />
        </t-link>
      </view>
      
      <scroll-view scroll-x class="-mx-32rpx px-32rpx" show-scrollbar="false">
        <view class="flex gap-20rpx">
          <t-card
            v-for="model in featuredModels"
            :key="model.id"
            class="w-280rpx flex-shrink-0"
            :hover="true"
            @click="goToModelDetail(model.id)"
          >
            <view class="flex-between mb-16rpx">
              <t-tag theme="primary" variant="light" size="small">
                {{ getFamilyName(model) }}
              </t-tag>
              <view class="flex gap-8rpx">
                <text 
                  v-for="cap in getModalityList(model).slice(0, 2)" 
                  :key="cap.name"
                  class="text-24rpx opacity-70"
                >{{ cap.icon }}</text>
              </view>
            </view>
            
            <text class="text-32rpx font-bold text-gray-900 block mb-24rpx">{{ model.name }}</text>
            
            <view class="flex-between pt-20rpx border-t border-gray-100">
              <view class="flex-col">
                <text class="text-28rpx font-bold text-gray-900">{{ formatPrice(model.costInput) }}</text>
                <text class="text-20rpx text-gray-400 mt-4rpx">输入/1M</text>
              </view>
              <view class="text-right">
                <text class="text-28rpx font-bold text-indigo-600">{{ formatPrice(model.costOutput) }}</text>
                <text class="text-20rpx text-gray-400 mt-4rpx block">输出/1M</text>
              </view>
            </view>
          </t-card>
        </view>
      </scroll-view>
    </view>

    <!-- Latest Papers -->
    <view class="mt-48rpx px-32rpx">
      <view class="flex-between mb-24rpx">
        <text class="text-36rpx font-bold text-gray-900">最新论文</text>
        <t-link theme="primary" @click="goToPapers">
          查看全部 <t-icon name="chevron-right" size="24rpx" />
        </t-link>
      </view>
      
      <t-cell-group class="rounded-24rpx overflow-hidden">
        <t-cell
          v-for="paper in latestPapers"
          :key="paper.id"
          :hover="true"
          @click="goToPaperDetail(paper.id)"
        >
          <template #title>
            <text class="text-30rpx font-semibold text-gray-900 leading-relaxed block mb-12rpx">{{ paper.titleCN }}</text>
          </template>
          <template #description>
            <text class="text-24rpx text-gray-500">{{ paper.authors[0] }} · {{ paper.publishDate }}</text>
          </template>
          <template #rightIcon>
            <t-tag theme="primary" variant="light" size="small">{{ paper.keywords[0] }}</t-tag>
          </template>
        </t-cell>
      </t-cell-group>
    </view>
  </view>
</template>

<style scoped>
/* 使用 UnoCSS，无需自定义样式 */
</style>
