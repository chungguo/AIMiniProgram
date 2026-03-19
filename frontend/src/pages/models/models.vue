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

    <!-- Family Filter -->
    <scroll-view scroll-x class="filter-scroll" show-scrollbar="false">
      <view class="filter-list">
        <view 
          :class="['filter-item', { active: selectedFamily === '' }]"
          @click="selectedFamily = ''"
        >全部</view>
        <view 
          v-for="family in families" 
          :key="family"
          :class="['filter-item', { active: selectedFamily === family }]"
          @click="selectedFamily = family"
        >{{ family }}</view>
      </view>
    </scroll-view>

    <!-- Compare Bar -->
    <view v-if="compareList.length > 0" class="compare-bar">
      <text class="compare-text">已选 {{ compareList.length }} 个模型</text>
      <button class="compare-btn" @click="goToCompare">开始对比</button>
    </view>

    <!-- Models List -->
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
            <text class="model-family">{{ getFamilyName(model) }}</text>
            <text class="model-context">{{ formatNumber(model.limitContext, 'tokens') }}</text>
          </view>
        </view>
        <view class="model-stats">
          <view class="stat">
            <text class="stat-value">{{ formatPrice(model.costInput) }}</text>
            <text class="stat-label">输入/1M</text>
          </view>
          <view class="stat">
            <text class="stat-value stat-value--price">{{ formatPrice(model.costOutput) }}</text>
            <text class="stat-label">输出/1M</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Empty -->
    <view v-if="filteredModels.length === 0 && !loading" class="empty">
      <text class="empty-icon">⌕</text>
      <text class="empty-text">未找到匹配的模型</text>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background: #f7f8fa;
}

/* Header */
.header {
  background: linear-gradient(165deg, #1a1a2e 0%, #16213e 100%);
  padding: 48rpx 32rpx 32rpx;
}

.header-title {
  font-size: 40rpx;
  font-weight: 700;
  color: #fff;
}

.header-count {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.6);
  margin-left: 16rpx;
}

/* Search */
.search-box {
  margin: -24rpx 32rpx 0;
  background: #fff;
  border-radius: 16rpx;
  padding: 20rpx 24rpx;
  display: flex;
  align-items: center;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
  position: relative;
  z-index: 10;
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
}

.search-placeholder {
  color: #9ca3af;
}

/* Filter */
.filter-scroll {
  margin-top: 24rpx;
  padding: 0 32rpx;
}

.filter-list {
  display: flex;
  gap: 16rpx;
  padding: 8rpx 0;
}

.filter-item {
  padding: 16rpx 32rpx;
  background: #fff;
  border-radius: 32rpx;
  font-size: 26rpx;
  color: #6b7280;
  white-space: nowrap;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.filter-item.active {
  background: #6366f1;
  color: #fff;
}

/* Compare Bar */
.compare-bar {
  margin: 24rpx 32rpx;
  background: #1a1a2e;
  border-radius: 16rpx;
  padding: 24rpx 32rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.compare-text {
  font-size: 28rpx;
  color: #fff;
}

.compare-btn {
  margin: 0;
  padding: 16rpx 32rpx;
  background: #6366f1;
  color: #fff;
  font-size: 26rpx;
  border-radius: 12rpx;
  line-height: 1;
}

/* Models List */
.models-list {
  padding: 0 32rpx 32rpx;
}

.model-item {
  background: #fff;
  border-radius: 20rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.model-info {
  flex: 1;
}

.model-row {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.model-row--secondary {
  margin-top: 12rpx;
}

.model-name {
  font-size: 32rpx;
  font-weight: 700;
  color: #1a1a2e;
}

.model-family {
  font-size: 24rpx;
  color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
  padding: 6rpx 16rpx;
  border-radius: 8rpx;
}

.model-context {
  font-size: 24rpx;
  color: #6b7280;
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

.compare-check {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  border: 2rpx solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 16rpx;
  font-size: 24rpx;
  color: #fff;
}

.compare-check.checked {
  background: #6366f1;
  border-color: #6366f1;
}

/* Empty */
.empty {
  padding: 120rpx 32rpx;
  text-align: center;
}

.empty-icon {
  font-size: 80rpx;
  color: #d1d5db;
}

.empty-text {
  font-size: 28rpx;
  color: #9ca3af;
  margin-top: 24rpx;
  display: block;
}
</style>
