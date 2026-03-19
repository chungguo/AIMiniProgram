<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { modelService } from '@/services';
import { getModalityList, getFamilyName, formatNumber, formatPrice } from '@/utils/modelHelpers';
import type { Model } from '@/types/api';

const model = ref<Model | null>(null);
const loading = ref<boolean>(true);
const activeTab = ref<'overview' | 'pricing' | 'capabilities'>('overview');

onMounted(() => {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1];
  const id = (currentPage as unknown as { $page?: { options?: { id?: string } } }).$page?.options?.id;
  
  if (id) {
    loadModel(id);
  }
});

async function loadModel(id: string): Promise<void> {
  try {
    loading.value = true;
    const res = await modelService.getModelById(id);
    if (res.success) model.value = res.data;
  } catch (error) {
    console.error('加载失败:', error);
  } finally {
    loading.value = false;
  }
}

// 特性标签
const featureTags = computed(() => {
  if (!model.value) return [];
  const tags: { icon: string; label: string; active: boolean }[] = [];
  tags.push({ icon: '🧠', label: '推理', active: model.value.reasoning });
  tags.push({ icon: '🔧', label: '工具调用', active: model.value.toolCall });
  tags.push({ icon: '📎', label: '附件', active: model.value.attachment });
  tags.push({ icon: '📋', label: '结构化输出', active: model.value.structuredOutput });
  tags.push({ icon: '🌡️', label: '温度调节', active: model.value.temperature });
  tags.push({ icon: '🔓', label: '开源权重', active: model.value.openWeights });
  return tags;
});

// 定价项目
const pricingItems = computed(() => {
  if (!model.value) return [];
  const items: { label: string; value: number }[] = [];
  if (model.value.costInput > 0) items.push({ label: '输入', value: model.value.costInput });
  if (model.value.costOutput > 0) items.push({ label: '输出', value: model.value.costOutput });
  if (model.value.costReasoning > 0) items.push({ label: '推理', value: model.value.costReasoning });
  if (model.value.costCacheRead > 0) items.push({ label: '缓存读取', value: model.value.costCacheRead });
  if (model.value.costCacheWrite > 0) items.push({ label: '缓存写入', value: model.value.costCacheWrite });
  if (model.value.costInputAudio > 0) items.push({ label: '音频输入', value: model.value.costInputAudio });
  if (model.value.costOutputAudio > 0) items.push({ label: '音频输出', value: model.value.costOutputAudio });
  return items;
});

function goBack(): void {
  uni.navigateBack();
}
</script>

<template>
  <view class="page" v-if="model">
    <!-- Header -->
    <view class="header">
      <text class="back-btn" @click="goBack">‹</text>
      <text class="header-title">模型详情</text>
      <view class="header-placeholder"></view>
    </view>

    <!-- Hero Card -->
    <view class="hero-card">
      <view class="hero-badge">{{ getFamilyName(model) }}</view>
      <text class="hero-name">{{ model.name }}</text>
      <text class="hero-desc">{{ model.description || '暂无描述' }}</text>
      
      <view class="hero-chips">
        <text v-if="model.architecture" class="chip">{{ model.architecture }}</text>
        <text v-if="model.releaseDate" class="chip">{{ model.releaseDate }}</text>
        <text v-if="model.knowledge" class="chip">知识截止: {{ model.knowledge }}</text>
      </view>
    </view>

    <!-- Tab Navigation -->
    <view class="tabs">
      <view 
        :class="['tab', { active: activeTab === 'overview' }]"
        @click="activeTab = 'overview'"
      >概览</view>
      <view 
        :class="['tab', { active: activeTab === 'pricing' }]"
        @click="activeTab = 'pricing'"
      >定价</view>
      <view 
        :class="['tab', { active: activeTab === 'capabilities' }]"
        @click="activeTab = 'capabilities'"
      >能力</view>
    </view>

    <!-- Tab Content -->
    <scroll-view scroll-y class="content">
      <!-- Overview Tab -->
      <view v-show="activeTab === 'overview'" class="tab-panel">
        <!-- Limits -->
        <view class="section">
          <text class="section-title">限制</text>
          <view class="stats-grid">
            <view class="stat-card">
              <text class="stat-value">{{ formatNumber(model.limitContext) }}</text>
              <text class="stat-label">上下文窗口</text>
            </view>
            <view class="stat-card">
              <text class="stat-value">{{ formatNumber(model.limitInput) }}</text>
              <text class="stat-label">最大输入</text>
            </view>
            <view class="stat-card">
              <text class="stat-value">{{ formatNumber(model.limitOutput) }}</text>
              <text class="stat-label">最大输出</text>
            </view>
          </view>
        </view>

        <!-- Modalities -->
        <view class="section">
          <text class="section-title">模态支持</text>
          <view class="modality-section">
            <view class="modality-group">
              <text class="modality-label">输入</text>
              <view class="modality-list">
                <text 
                  v-for="mod in getModalityList(model)" 
                  :key="mod.name"
                  class="modality-tag"
                >{{ mod.icon }} {{ mod.name }}</text>
              </view>
            </view>
            <view class="modality-group">
              <text class="modality-label">输出</text>
              <view class="modality-list">
                <text 
                  v-for="mod in model.modalitiesOutput" 
                  :key="mod"
                  class="modality-tag"
                >{{ mod }}</text>
              </view>
            </view>
          </view>
        </view>

        <!-- Features -->
        <view class="section">
          <text class="section-title">特性</text>
          <view class="feature-grid">
            <view 
              v-for="tag in featureTags" 
              :key="tag.label"
              :class="['feature-item', { active: tag.active }]"
            >
              <text class="feature-icon">{{ tag.icon }}</text>
              <text class="feature-label">{{ tag.label }}</text>
              <text class="feature-status">{{ tag.active ? '✓' : '✗' }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Pricing Tab -->
      <view v-show="activeTab === 'pricing'" class="tab-panel">
        <view class="section">
          <text class="section-title">价格（每百万 tokens）</text>
          <view class="pricing-list">
            <view 
              v-for="item in pricingItems" 
              :key="item.label"
              class="pricing-item"
            >
              <text class="pricing-label">{{ item.label }}</text>
              <text class="pricing-value">{{ formatPrice(item.value) }}</text>
            </view>
            <view v-if="pricingItems.length === 0" class="pricing-empty">
              <text>暂无定价信息</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Capabilities Tab -->
      <view v-show="activeTab === 'capabilities'" class="tab-panel">
        <view class="section">
          <text class="section-title">详细信息</text>
          <view class="info-list">
            <view class="info-item">
              <text class="info-label">ID</text>
              <text class="info-value">{{ model.id }}</text>
            </view>
            <view class="info-item">
              <text class="info-label">家族</text>
              <text class="info-value">{{ model.family }}</text>
            </view>
            <view class="info-item">
              <text class="info-label">发布日期</text>
              <text class="info-value">{{ model.releaseDate || '-' }}</text>
            </view>
            <view class="info-item">
              <text class="info-label">知识截止日期</text>
              <text class="info-value">{{ model.knowledge || '-' }}</text>
            </view>
            <view class="info-item">
              <text class="info-label">最后更新</text>
              <text class="info-value">{{ model.lastUpdated || '-' }}</text>
            </view>
            <view class="info-item">
              <text class="info-label">开源权重</text>
              <text class="info-value">{{ model.openWeights ? '是' : '否' }}</text>
            </view>
            <view v-if="model.interleavedField" class="info-item">
              <text class="info-label">推理字段</text>
              <text class="info-value">{{ model.interleavedField }}</text>
            </view>
          </view>
        </view>
      </view>
    </scroll-view>
  </view>

  <!-- Loading -->
  <view v-else-if="loading" class="loading">
    <text class="loading-text">加载中...</text>
  </view>

  <!-- Empty -->
  <view v-else class="empty">
    <text class="empty-text">模型不存在</text>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background: #f7f8fa;
  display: flex;
  flex-direction: column;
}

/* Header */
.header {
  background: linear-gradient(165deg, #1a1a2e 0%, #16213e 100%);
  padding: 48rpx 32rpx 32rpx;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.back-btn {
  font-size: 48rpx;
  color: #fff;
  font-weight: 300;
  width: 60rpx;
}

.header-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
}

.header-placeholder {
  width: 60rpx;
}

/* Hero Card */
.hero-card {
  background: #fff;
  margin: -20rpx 32rpx 0;
  border-radius: 24rpx;
  padding: 32rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
  position: relative;
  z-index: 10;
}

.hero-badge {
  display: inline-block;
  background: rgba(99, 102, 241, 0.1);
  color: #6366f1;
  font-size: 24rpx;
  font-weight: 600;
  padding: 8rpx 20rpx;
  border-radius: 8rpx;
  margin-bottom: 16rpx;
}

.hero-name {
  font-size: 40rpx;
  font-weight: 700;
  color: #1a1a2e;
  display: block;
  margin-bottom: 16rpx;
}

.hero-desc {
  font-size: 26rpx;
  color: #6b7280;
  line-height: 1.5;
  display: block;
  margin-bottom: 20rpx;
}

.hero-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
}

.chip {
  background: #f3f4f6;
  color: #4b5563;
  font-size: 22rpx;
  padding: 8rpx 16rpx;
  border-radius: 8rpx;
}

/* Tabs */
.tabs {
  display: flex;
  margin: 32rpx 32rpx 0;
  background: #fff;
  border-radius: 16rpx;
  padding: 8rpx;
}

.tab {
  flex: 1;
  text-align: center;
  padding: 20rpx 0;
  font-size: 28rpx;
  color: #6b7280;
  border-radius: 12rpx;
}

.tab.active {
  background: #1a1a2e;
  color: #fff;
  font-weight: 600;
}

/* Content */
.content {
  flex: 1;
  margin-top: 24rpx;
}

.tab-panel {
  padding: 0 32rpx 40rpx;
}

/* Section */
.section {
  margin-top: 32rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 20rpx;
  display: block;
}

/* Stats Grid */
.stats-grid {
  display: flex;
  gap: 16rpx;
}

.stat-card {
  flex: 1;
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  text-align: center;
}

.stat-value {
  font-size: 32rpx;
  font-weight: 700;
  color: #6366f1;
  display: block;
}

.stat-label {
  font-size: 22rpx;
  color: #9ca3af;
  margin-top: 8rpx;
  display: block;
}

/* Modality */
.modality-section {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
}

.modality-group {
  margin-bottom: 20rpx;
}

.modality-group:last-child {
  margin-bottom: 0;
}

.modality-label {
  font-size: 24rpx;
  color: #9ca3af;
  margin-bottom: 12rpx;
  display: block;
}

.modality-list {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
}

.modality-tag {
  background: rgba(99, 102, 241, 0.1);
  color: #6366f1;
  font-size: 24rpx;
  padding: 10rpx 20rpx;
  border-radius: 8rpx;
}

/* Feature Grid */
.feature-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16rpx;
}

.feature-item {
  background: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  display: flex;
  align-items: center;
  gap: 16rpx;
  opacity: 0.6;
}

.feature-item.active {
  opacity: 1;
}

.feature-icon {
  font-size: 32rpx;
}

.feature-label {
  flex: 1;
  font-size: 26rpx;
  color: #1a1a2e;
}

.feature-status {
  font-size: 24rpx;
  color: #10b981;
}

.feature-item:not(.active) .feature-status {
  color: #9ca3af;
}

/* Pricing */
.pricing-list {
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
}

.pricing-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx;
  border-bottom: 2rpx solid #f3f4f6;
}

.pricing-item:last-child {
  border-bottom: none;
}

.pricing-label {
  font-size: 28rpx;
  color: #4b5563;
}

.pricing-value {
  font-size: 32rpx;
  font-weight: 700;
  color: #6366f1;
}

.pricing-empty {
  padding: 48rpx;
  text-align: center;
  color: #9ca3af;
}

/* Info List */
.info-list {
  background: #fff;
  border-radius: 16rpx;
  overflow: hidden;
}

.info-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx;
  border-bottom: 2rpx solid #f3f4f6;
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 26rpx;
  color: #6b7280;
}

.info-value {
  font-size: 26rpx;
  color: #1a1a2e;
  font-weight: 500;
}

/* Loading & Empty */
.loading, .empty {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.loading-text, .empty-text {
  font-size: 28rpx;
  color: #9ca3af;
}
</style>
