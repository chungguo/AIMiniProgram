<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { modelService } from '@/services';
import type { Model, Provider } from '@/types/api';

const models = ref<Model[]>([]);
const providers = ref<Provider[]>([]);
const loading = ref<boolean>(true);
const selectedProvider = ref<string>('');
const searchKeyword = ref<string>('');
const compareList = ref<string[]>([]);

const filteredModels = computed<Model[]>(() => {
  let result = models.value;
  if (selectedProvider.value) {
    result = result.filter((m: Model) => m.providerId === selectedProvider.value);
  }
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase();
    result = result.filter((m: Model) => 
      m.name.toLowerCase().includes(keyword) ||
      m.provider.toLowerCase().includes(keyword)
    );
  }
  return result;
});

onMounted(() => {
  loadData();
});

async function loadData(): Promise<void> {
  try {
    const [modelsRes, providersRes] = await Promise.all([
      modelService.getModels(),
      modelService.getProviders()
    ]);
    if (modelsRes.success) models.value = modelsRes.data;
    if (providersRes.success) providers.value = providersRes.data;
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
  <view class="page">
    <!-- Header -->
    <view class="header">
      <text class="header-title">大模型库</text>
      <text class="header-count">{{ models.length }} 个模型</text>
    </view>

    <!-- Search -->
    <view class="search-box">
      <text class="search-icon">⌕</text>
      <input 
        v-model="searchKeyword"
        class="search-input"
        placeholder="搜索模型名称..."
        placeholder-class="search-placeholder"
      />
    </view>

    <!-- Provider Filter -->
    <scroll-view scroll-x class="filter-scroll" show-scrollbar="false">
      <view class="filter-list">
        <view 
          :class="['filter-item', { active: selectedProvider === '' }]"
          @click="selectedProvider = ''"
        >全部</view>
        <view 
          v-for="provider in providers" 
          :key="provider.id"
          :class="['filter-item', { active: selectedProvider === provider.id }]"
          @click="selectedProvider = provider.id"
        >{{ provider.name }}</view>
      </view>
    </scroll-view>

    <!-- Compare Bar -->
    <view v-if="compareList.length > 0" class="compare-bar">
      <text class="compare-text">已选 {{ compareList.length }} 个模型</text>
      <button class="compare-btn" @click="goToCompare">开始对比</button>
    </view>

    <!-- Models List - 移动卡片列表 -->
    <view class="models-list">
      <view 
        v-for="model in filteredModels" 
        :key="model.id"
        class="model-item"
        @click="goToDetail(model.id)"
      >
        <view class="model-info">
          <view class="model-row">
            <text class="model-name">{{ model.name }}</text>
            <view 
              :class="['compare-check', { checked: compareList.includes(model.id) }]"
              @click.stop="toggleCompare(model.id)"
            >
              <text v-if="compareList.includes(model.id)">✓</text>
            </view>
          </view>
          <view class="model-row model-row--secondary">
            <text class="model-provider">{{ model.provider }}</text>
            <text class="model-context">{{ model.contextWindow >= 1000 ? (model.contextWindow / 1000) + 'K' : model.contextWindow }} 上下文</text>
          </view>
        </view>
        <view class="model-stats">
          <view class="stat">
            <text class="stat-value">{{ model.quality.mmlu }}</text>
            <text class="stat-label">MMLU</text>
          </view>
          <view class="stat">
            <text class="stat-value stat-value--price">${{ model.pricing.inputPrice }}</text>
            <text class="stat-label">/1M</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Empty -->
    <view v-if="filteredModels.length === 0 && !loading" class="empty">
      <text class="empty-icon">⌕</text>
      <text class="empty-text">没有找到模型</text>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 120rpx;
}

/* ===== Header ===== */
.header {
  padding: 32rpx;
  background: linear-gradient(165deg, #1a1a2e 0%, #16213e 100%);
}

.header-title {
  font-size: 40rpx;
  font-weight: 700;
  color: #fff;
  display: block;
  margin-bottom: 8rpx;
}

.header-count {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.6);
}

/* ===== Search ===== */
.search-box {
  margin: 24rpx 32rpx;
  background: #fff;
  border-radius: 16rpx;
  padding: 20rpx 24rpx;
  display: flex;
  align-items: center;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.search-icon {
  font-size: 32rpx;
  color: #9ca3af;
  margin-right: 16rpx;
}

.search-input {
  flex: 1;
  font-size: 30rpx;
  color: #1a1a2e;
  height: 48rpx;
}

.search-placeholder {
  color: #9ca3af;
}

/* ===== Filter ===== */
.filter-scroll {
  margin-bottom: 24rpx;
  padding: 0 32rpx;
}

.filter-list {
  display: flex;
  gap: 16rpx;
}

.filter-item {
  padding: 16rpx 28rpx;
  background: #fff;
  border-radius: 100rpx;
  font-size: 26rpx;
  color: #4b5563;
  white-space: nowrap;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.filter-item.active {
  background: #6366f1;
  color: #fff;
}

/* ===== Compare Bar ===== */
.compare-bar {
  margin: 0 32rpx 24rpx;
  background: #1a1a2e;
  border-radius: 16rpx;
  padding: 20rpx 24rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.compare-text {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
}

.compare-btn {
  background: #6366f1;
  color: #fff;
  font-size: 26rpx;
  font-weight: 600;
  padding: 16rpx 32rpx;
  border-radius: 12rpx;
  border: none;
  margin: 0;
}

/* ===== Models List ===== */
.models-list {
  padding: 0 32rpx;
}

.model-item {
  background: #fff;
  border-radius: 20rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.model-item:active {
  transform: scale(0.98);
}

.model-info {
  flex: 1;
}

.model-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12rpx;
}

.model-row--secondary {
  margin-bottom: 0;
}

.model-name {
  font-size: 34rpx;
  font-weight: 700;
  color: #1a1a2e;
}

.compare-check {
  width: 48rpx;
  height: 48rpx;
  border: 2rpx solid #e5e7eb;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24rpx;
  color: #fff;
  margin-left: 16rpx;
}

.compare-check.checked {
  background: #6366f1;
  border-color: #6366f1;
}

.model-provider {
  font-size: 26rpx;
  color: #6366f1;
  font-weight: 500;
}

.model-context {
  font-size: 24rpx;
  color: #9ca3af;
}

.model-stats {
  display: flex;
  gap: 32rpx;
  text-align: right;
}

.stat {
  display: flex;
  flex-direction: column;
}

.stat-value {
  font-size: 32rpx;
  font-weight: 700;
  color: #1a1a2e;
}

.stat-value--price {
  color: #6366f1;
}

.stat-label {
  font-size: 20rpx;
  color: #9ca3af;
  margin-top: 4rpx;
}

/* ===== Empty ===== */
.empty {
  text-align: center;
  padding: 120rpx 32rpx;
}

.empty-icon {
  font-size: 64rpx;
  color: #d1d5db;
  display: block;
  margin-bottom: 20rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #6b7280;
}
</style>
