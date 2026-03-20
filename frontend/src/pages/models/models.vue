<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { modelService } from '@/services';
import { getFamilyName, formatNumber, formatPrice } from '@/utils/modelHelpers';
import type { Model } from '@/types/api';

const models = ref<Model[]>([]);
const families = ref<string[]>([]);
const loading = ref<boolean>(true);
const selectedFamily = ref<string>('');
const searchKeyword = ref<string>('');
const compareList = ref<string[]>([]);

const filteredModels = computed<Model[]>(() => {
  let result = models.value;
  if (selectedFamily.value) {
    result = result.filter((m: Model) => m.family === selectedFamily.value);
  }
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter((m: Model) => 
      m.name.toLowerCase().includes(keyword) ||
      getFamilyName(m).toLowerCase().includes(keyword)
    );
  }
  return result;
});

onMounted(() => {
  loadData();
});

async function loadData(): Promise<void> {
  try {
    const [modelsRes, familiesRes] = await Promise.all([
      modelService.getModels({ limit: 100 }),
      modelService.getFamilies()
    ]);
    if (modelsRes.success) models.value = modelsRes.data;
    if (familiesRes.success) families.value = familiesRes.data;
  } catch (error) {
    console.error('加载失败:', error);
  } finally {
    loading.value = false;
  }
}

function toggleCompare(id: string): void {
  const index = compareList.value.indexOf(id);
  if (index > -1) {
    compareList.value.splice(index, 1);
  } else if (compareList.value.length < 3) {
    compareList.value.push(id);
  } else {
    uni.showToast({ title: '最多选择3个', icon: 'none' });
  }
}

function goToCompare(): void {
  if (compareList.value.length < 2) {
    uni.showToast({ title: '至少选择2个', icon: 'none' });
    return;
  }
  uni.navigateTo({ url: `/pages/compare/compare?ids=${compareList.value.join(',')}` });
}

function goToDetail(id: string): void {
  uni.navigateTo({ url: `/pages/model-detail/model-detail?id=${id}` });
}
</script>

<template>
  <view class="min-h-screen bg-gray-50">
    <!-- Header -->
    <view class="bg-gradient-to-br from-indigo-900 to-blue-900 px-32rpx pt-48rpx pb-32rpx">
      <text class="text-40rpx font-bold text-white">大模型库</text>
      <text class="text-26rpx text-white/60 ml-16rpx">{{ models.length }} 个模型</text>
    </view>

    <!-- Search -->
    <view class="-mt-24rpx mx-32rpx">
      <t-search
        v-model="searchKeyword"
        placeholder="搜索模型名称..."
        shape="round"
        :clearable="true"
      />
    </view>

    <!-- Family Filter -->
    <scroll-view scroll-x class="mt-24rpx px-32rpx" show-scrollbar="false">
      <view class="flex gap-16rpx py-8rpx">
        <t-tag
          :theme="selectedFamily === '' ? 'primary' : 'default'"
          shape="round"
          @click="selectedFamily = ''"
        >全部</t-tag>
        <t-tag
          v-for="family in families"
          :key="family"
          :theme="selectedFamily === family ? 'primary' : 'default'"
          shape="round"
          @click="selectedFamily = family"
        >{{ family }}</t-tag>
      </view>
    </scroll-view>

    <!-- Compare Bar -->
    <view v-if="compareList.length > 0" class="mx-32rpx mt-24rpx p-24rpx bg-indigo-900 rounded-16rpx flex-between">
      <text class="text-28rpx text-white">已选 {{ compareList.length }} 个模型</text>
      <t-button theme="primary" size="small" @click="goToCompare">开始对比</t-button>
    </view>

    <!-- Models List -->
    <view class="px-32rpx py-24rpx">
      <t-cell-group class="rounded-20rpx overflow-hidden">
        <t-cell
          v-for="model in filteredModels"
          :key="model.id"
          :hover="true"
          @click="goToDetail(model.id)"
        >
          <template #leftIcon>
            <t-checkbox
              :checked="compareList.includes(model.id)"
              @change="toggleCompare(model.id)"
              @click.stop
            />
          </template>
          <template #title>
            <view class="flex items-center gap-16rpx">
              <text class="text-32rpx font-bold text-gray-900">{{ model.name }}</text>
              <t-tag theme="primary" variant="light" size="small">{{ getFamilyName(model) }}</t-tag>
            </view>
          </template>
          <template #description>
            <text class="text-24rpx text-gray-500">{{ formatNumber(model.limitContext, 'tokens') }}</text>
          </template>
          <template #rightIcon>
            <view class="flex gap-32rpx text-right">
              <view class="flex-col">
                <text class="text-32rpx font-bold text-gray-900">{{ formatPrice(model.costInput) }}</text>
                <text class="text-20rpx text-gray-400 mt-4rpx">输入/1M</text>
              </view>
              <view class="flex-col">
                <text class="text-32rpx font-bold text-indigo-600">{{ formatPrice(model.costOutput) }}</text>
                <text class="text-20rpx text-gray-400 mt-4rpx">输出/1M</text>
              </view>
            </view>
          </template>
        </t-cell>
      </t-cell-group>
    </view>

    <!-- Empty -->
    <view v-if="filteredModels.length === 0 && !loading" class="py-120rpx text-center">
      <t-icon name="search" size="80rpx" color="#d1d5db" />
      <text class="text-28rpx text-gray-400 mt-24rpx block">未找到匹配的模型</text>
    </view>
  </view>
</template>

<style scoped>
/* 使用 UnoCSS */
</style>
