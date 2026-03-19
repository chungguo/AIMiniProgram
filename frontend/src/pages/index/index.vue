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
  <view class="page">
    <!-- Header -->
    <view class="header">
      <view class="header-top">
        <text class="logo">AI Hub</text>
        <view class="header-icon">◈</view>
      </view>
      <text class="header-title">探索大模型与前沿论文</text>
    </view>

    <!-- Quick Actions -->
    <view class="quick-actions">
      <view class="action-card action-card--primary" @click="goToModels">
        <view class="action-icon">
          <text class="action-icon-text">◈</text>
        </view>
        <view class="action-content">
          <text class="action-title">模型对比</text>
          <text class="action-desc">主流大模型参数</text>
        </view>
        <text class="action-arrow">›</text>
      </view>
      
      <view class="action-card" @click="goToPapers">
        <view class="action-icon action-icon--secondary">
          <text class="action-icon-text">◉</text>
        </view>
        <view class="action-content">
          <text class="action-title">论文研读</text>
          <text class="action-desc">中英双语</text>
        </view>
        <text class="action-arrow">›</text>
      </view>
    </view>

    <!-- Featured Models -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">热门模型</text>
        <text class="section-more" @click="goToModels">查看全部 ›</text>
      </view>
      
      <scroll-view scroll-x class="models-scroll" show-scrollbar="false">
        <view class="models-list">
          <view 
            v-for="model in featuredModels" 
            :key="model.id"
            class="model-card"
            @click="goToModelDetail(model.id)"
          >
            <view class="model-header">
              <text class="model-family">{{ getFamilyName(model) }}</text>
              <view class="model-caps">
                <text 
                  v-for="cap in getModalityList(model).slice(0, 2)" 
                  :key="cap.name"
                  class="model-cap"
                >{{ cap.icon }}</text>
              </view>
            </view>
            <text class="model-name">{{ model.name }}</text>
            <view class="model-footer">
              <view class="model-stat">
                <text class="model-stat-value">{{ formatPrice(model.costInput) }}</text>
                <text class="model-stat-label">输入/1M</text>
              </view>
              <view class="model-price">
                <text class="model-price-value">{{ formatPrice(model.costOutput) }}</text>
                <text class="model-price-label">输出/1M</text>
              </view>
            </view>
          </view>
        </view>
      </scroll-view>
    </view>

    <!-- Latest Papers -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">最新论文</text>
        <text class="section-more" @click="goToPapers">查看全部 ›</text>
      </view>
      
      <view class="papers-list">
        <view 
          v-for="paper in latestPapers" 
          :key="paper.id"
          class="paper-item"
          @click="goToPaperDetail(paper.id)"
        >
          <view class="paper-main">
            <text class="paper-title">{{ paper.titleCN }}</text>
            <text class="paper-meta">{{ paper.authors[0] }} · {{ paper.publishDate }}</text>
          </view>
          <view class="paper-tag">{{ paper.keywords[0] }}</view>
        </view>
      </view>
    </view>
  </view>
</template>

<style scoped>
.page {
  min-height: 100vh;
  background: #f7f8fa;
  padding-bottom: 120rpx;
}

/* Header */
.header {
  background: linear-gradient(165deg, #1a1a2e 0%, #16213e 100%);
  padding: 48rpx 32rpx 40rpx;
}

.header-top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.logo {
  font-size: 36rpx;
  font-weight: 700;
  color: #fff;
  letter-spacing: 0.05em;
}

.header-icon {
  width: 64rpx;
  height: 64rpx;
  background: rgba(99, 102, 241, 0.2);
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32rpx;
  color: #6366f1;
}

.header-title {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.5;
}

/* Quick Actions */
.quick-actions {
  padding: 24rpx 32rpx;
  display: flex;
  gap: 20rpx;
}

.action-card {
  flex: 1;
  background: #fff;
  border-radius: 24rpx;
  padding: 32rpx 24rpx;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
  position: relative;
}

.action-card--primary {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
}

.action-icon {
  width: 72rpx;
  height: 72rpx;
  background: rgba(99, 102, 241, 0.1);
  border-radius: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20rpx;
}

.action-card--primary .action-icon {
  background: rgba(255, 255, 255, 0.2);
}

.action-icon-text {
  font-size: 36rpx;
  color: #6366f1;
}

.action-card--primary .action-icon-text {
  color: #fff;
}

.action-content {
  flex: 1;
}

.action-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #1a1a2e;
  display: block;
  margin-bottom: 8rpx;
}

.action-card--primary .action-title {
  color: #fff;
}

.action-desc {
  font-size: 24rpx;
  color: #6b7280;
}

.action-card--primary .action-desc {
  color: rgba(255, 255, 255, 0.8);
}

.action-arrow {
  position: absolute;
  top: 32rpx;
  right: 24rpx;
  font-size: 36rpx;
  color: #9ca3af;
  font-weight: 300;
}

.action-card--primary .action-arrow {
  color: rgba(255, 255, 255, 0.6);
}

/* Section */
.section {
  margin-top: 48rpx;
  padding: 0 32rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.section-title {
  font-size: 36rpx;
  font-weight: 700;
  color: #1a1a2e;
}

.section-more {
  font-size: 26rpx;
  color: #6366f1;
  font-weight: 500;
}

/* Models Scroll */
.models-scroll {
  margin: 0 -32rpx;
  padding: 0 32rpx;
}

.models-list {
  display: flex;
  gap: 20rpx;
}

.model-card {
  flex-shrink: 0;
  width: 280rpx;
  background: #fff;
  border-radius: 24rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.model-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.model-family {
  font-size: 22rpx;
  font-weight: 600;
  color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
  padding: 6rpx 14rpx;
  border-radius: 8rpx;
}

.model-caps {
  display: flex;
  gap: 8rpx;
}

.model-cap {
  font-size: 24rpx;
  opacity: 0.7;
}

.model-name {
  font-size: 32rpx;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 24rpx;
  display: block;
}

.model-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 20rpx;
  border-top: 2rpx solid #f3f4f6;
}

.model-stat {
  display: flex;
  flex-direction: column;
}

.model-stat-value {
  font-size: 28rpx;
  font-weight: 700;
  color: #1a1a2e;
}

.model-stat-label {
  font-size: 20rpx;
  color: #9ca3af;
  margin-top: 4rpx;
}

.model-price {
  text-align: right;
}

.model-price-value {
  font-size: 28rpx;
  font-weight: 700;
  color: #6366f1;
}

.model-price-label {
  font-size: 20rpx;
  color: #9ca3af;
}

/* Papers List */
.papers-list {
  background: #fff;
  border-radius: 24rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.paper-item {
  padding: 28rpx 24rpx;
  border-bottom: 2rpx solid #f3f4f6;
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.paper-item:last-child {
  border-bottom: none;
}

.paper-item:active {
  background: #f9fafb;
}

.paper-main {
  flex: 1;
  padding-right: 20rpx;
}

.paper-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #1a1a2e;
  line-height: 1.4;
  display: block;
  margin-bottom: 12rpx;
}

.paper-meta {
  font-size: 24rpx;
  color: #6b7280;
}

.paper-tag {
  flex-shrink: 0;
  font-size: 22rpx;
  font-weight: 500;
  color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
  padding: 10rpx 18rpx;
  border-radius: 8rpx;
}
</style>
