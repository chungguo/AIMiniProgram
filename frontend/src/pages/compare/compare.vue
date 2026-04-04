<script setup lang="ts">
import { ref, computed } from 'vue';
import { modelService } from '@/services';
import { 
  getNestedValue, 
  formatComparisonValue, 
  getValueClass,
  calculateBestValues,
  isHigherBetterMetric,
  getFamilyName
} from '@/utils/modelHelpers';
import { useLoading, useTabFilter, getPageParams } from '@/composables';
import type { Model, ComparisonCategory } from '@/types/api';

// 使用 useLoading 管理加载状态
const { loading, withLoading } = useLoading();

// 对比数据
const models = ref<Model[]>([]);
const categories = ref<ComparisonCategory[]>([]);

// 使用 useTabFilter 管理分类 Tab
const { activeTab: activeCategory, onTabChange, options: categoryOptions } = useTabFilter<string>({
  options: computed(() => 
    categories.value.map(cat => ({
      label: getCategoryName(cat.key),
      value: cat.key
    }))
  ).value,
  onChange: (value) => {
    // Tab 切换逻辑（如果需要）
    console.log('切换到分类:', value);
  }
});

// 分类名称映射
const categoryNames: Record<string, string> = {
  basic: '基本信息',
  capabilities: '能力特性',
  modalities: '模态支持',
  limits: '限制',
  pricing: '定价'
};

function getCategoryName(key: string): string {
  return categoryNames[key] || key;
}

// 加载对比数据
async function loadCompareData(ids: string[]): Promise<void> {
  await withLoading(async () => {
    const res = await modelService.compareModels(ids);
    
    if (res.success) {
      models.value = res.data.models;
      categories.value = res.data.comparisonCategories;
      // 默认选中第一个分类
      if (categories.value.length > 0 && !activeCategory.value) {
        onTabChange(categories.value[0].key);
      }
    }
  });
}

// 页面加载时获取参数并加载数据
const { ids } = getPageParams(['ids']);
if (ids) {
  const idList = ids.split(',');
  if (idList.length >= 2) {
    loadCompareData(idList);
  }
}

function getModelValue(model: Model, key: string): unknown {
  // 特殊处理 family（可能来自旧 provider 字段）
  if (key === 'family' && !model.family) {
    return model.provider || '-';
  }
  return getNestedValue(model as Record<string, unknown>, key);
}

function getBestClasses(categoryKey: string, itemKey: string): (string | undefined)[] {
  const isHigherBetter = isHigherBetterMetric(itemKey);
  return calculateBestValues(models.value, itemKey, isHigherBetter);
}

function goBack(): void {
  uni.navigateBack();
}
</script>

<template>
  <view class="page">
    <!-- Header -->
    <view class="header">
      <text class="back-btn" @click="goBack">‹ 返回</text>
      <text class="header-title">模型对比</text>
    </view>

    <!-- Loading State -->
    <view v-if="loading" class="loading-container">
      <t-loading theme="spinner" size="80rpx" />
      <text class="loading-text">加载对比数据中...</text>
    </view>

    <template v-else>
      <!-- Models Names -->
      <view class="models-bar">
        <view class="model-tag" v-for="model in models" :key="model.id">
          {{ model.name }}
        </view>
      </view>

      <!-- Category Tabs -->
      <scroll-view scroll-x class="tabs-scroll" show-scrollbar="false">
        <view class="tabs-list">
          <view 
            v-for="cat in categoryOptions" 
            :key="cat.value"
            :class="['tab', { active: activeCategory === cat.value }]"
            @click="onTabChange(cat.value)"
          >
            {{ cat.label }}
          </view>
        </view>
      </scroll-view>

      <!-- Comparison Content -->
      <scroll-view scroll-y class="content">
        <view 
          v-for="category in categories" 
          :key="category.key"
          v-show="activeCategory === category.key"
          class="category-panel"
        >
          <view 
            v-for="item in category.items" 
            :key="item.key"
            class="compare-row"
          >
            <view class="row-label">
              <text class="row-label-name">{{ item.name }}</text>
              <text v-if="item.unit" class="row-label-unit">{{ item.unit }}</text>
            </view>
            
            <view class="row-values">
              <view 
                v-for="(model, idx) in models" 
                :key="model.id"
                :class="[
                  'value-cell', 
                  getValueClass(getModelValue(model, item.key), item.type),
                  getBestClasses(category.key, item.key)[idx]
                ]"
              >
                {{ formatComparisonValue(getModelValue(model, item.key), item.type, item.unit) }}
              </view>
            </view>
          </view>
        </view>
      </scroll-view>
    </template>
  </view>
</template>

<style scoped>
.page {
  height: 100vh;
  background: #f7f8fa;
  display: flex;
  flex-direction: column;
}

/* Loading Container */
.loading-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 24rpx;
}

.loading-text {
  font-size: 28rpx;
  color: #9ca3af;
}

/* Header */
.header {
  background: #1a1a2e;
  padding: 24rpx 32rpx;
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.back-btn {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
  padding: 12rpx 0;
}

.header-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
  flex: 1;
  text-align: center;
  margin-right: 80rpx;
}

/* Models Bar */
.models-bar {
  background: #fff;
  padding: 20rpx 32rpx;
  display: flex;
  gap: 16rpx;
  border-bottom: 2rpx solid #f3f4f6;
  overflow-x: auto;
}

.model-tag {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: #fff;
  font-size: 24rpx;
  font-weight: 600;
  padding: 12rpx 24rpx;
  border-radius: 100rpx;
  white-space: nowrap;
}

/* Tabs */
.tabs-scroll {
  background: #fff;
  border-bottom: 2rpx solid #f3f4f6;
}

.tabs-list {
  display: flex;
  padding: 16rpx 32rpx;
  gap: 12rpx;
}

.tab {
  padding: 16rpx 32rpx;
  font-size: 28rpx;
  color: #6b7280;
  border-radius: 12rpx;
  white-space: nowrap;
  transition: all 0.2s;
}

.tab.active {
  background: #6366f1;
  color: #fff;
  font-weight: 600;
}

/* Content */
.content {
  flex: 1;
  padding: 24rpx 32rpx;
}

.category-panel {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10rpx); }
  to { opacity: 1; transform: translateY(0); }
}

.compare-row {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.row-label {
  margin-bottom: 20rpx;
  padding-bottom: 16rpx;
  border-bottom: 2rpx solid #f3f4f6;
}

.row-label-name {
  font-size: 28rpx;
  font-weight: 600;
  color: #1a1a2e;
}

.row-label-unit {
  font-size: 22rpx;
  color: #9ca3af;
  margin-left: 12rpx;
}

.row-values {
  display: flex;
  gap: 16rpx;
}

.value-cell {
  flex: 1;
  text-align: center;
  padding: 20rpx 16rpx;
  background: #f9fafb;
  border-radius: 12rpx;
  font-size: 28rpx;
  color: #4b5563;
  font-weight: 500;
}

.value-cell.value-yes {
  background: #d1fae5;
  color: #059669;
}

.value-cell.value-no {
  background: #f3f4f6;
  color: #9ca3af;
}

.value-cell.best-value {
  background: #e0e7ff;
  color: #4f46e5;
  font-weight: 700;
}
</style>
