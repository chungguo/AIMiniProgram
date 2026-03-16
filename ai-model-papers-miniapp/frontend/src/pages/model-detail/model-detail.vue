<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { modelService } from '@/services';
import type { Model } from '@/types/api';

const model = ref<Model | null>(null);
const loading = ref<boolean>(true);
const activeTab = ref<'overview' | 'pricing' | 'benchmarks'>('overview');

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

function getCapabilities(capabilities: Record<string, boolean>): string[] {
  const list: string[] = [];
  if (capabilities.text) list.push('文本');
  if (capabilities.image) list.push('图像');
  if (capabilities.audio) list.push('音频');
  if (capabilities.video) list.push('视频');
  if (capabilities.file) list.push('文件');
  return list;
}

function goBack(): void {
  uni.navigateBack();
}
</script>

<template>
  <view class="page" v-if="model">
    <!-- Header - 固定顶部 -->
    <view class="header">
      <text class="back-btn" @click="goBack">‹</text>
      <text class="header-title">模型详情</text>
      <view class="header-placeholder"></view>
    </view>

    <!-- Hero Card -->
    <view class="hero-card">
      <view class="hero-badge">{{ model.provider }}</view>
      <text class="hero-name">{{ model.name }}</text>
      <text class="hero-desc">{{ model.description }}</text>
      
      <view class="hero-chips">
        <text class="chip">{{ model.architecture }}</text>
        <text class="chip">{{ model.releaseDate }}</text>
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
        :class="['tab', { active: activeTab === 'benchmarks' }]"
        @click="activeTab = 'benchmarks'"
      >性能</view>
    </view>

    <!-- Tab Content -->
    <scroll-view scroll-y class="content">
      <!-- Overview Tab -->
      <view v-show="activeTab === 'overview'" class="tab-panel">
        <!-- Capabilities -->
        <view class="section">
          <text class="section-title">能力支持</text>
          <view class="caps-grid">
            <view 
              v-for="cap in getCapabilities(model.capabilities)" 
              :key="cap"
              class="cap-item"
            >
              <text class="cap-name">{{ cap }}</text>
              <text class="cap-check">✓</text>
            </view>
          </view>
        </view>

        <!-- Specs -->
        <view class="section">
          <text class="section-title">技术规格</text>
          <view class="specs-list">
            <view class="spec-item">
              <text class="spec-label">上下文窗口</text>
              <text class="spec-value">{{ model.contextWindow >= 1000 ? (model.contextWindow / 1000) + 'K' : model.contextWindow }} tokens</text>
            </view>
            <view class="spec-item">
              <text class="spec-label">最大输出</text>
              <text class="spec-value">{{ model.maxTokens }} tokens</text>
            </view>
            <view class="spec-item">
              <text class="spec-label">生成速度</text>
              <text class="spec-value">{{ model.speed.tokensPerSecond }} tokens/s</text>
            </view>
          </view>
        </view>

        <!-- Features -->
        <view class="section">
          <text class="section-title">特性</text>
          <view class="features-list">
            <text 
              v-for="(feature, idx) in model.features" 
              :key="idx"
              class="feature-item"
            >
              {{ feature }}
            </text>
          </view>
        </view>
      </view>

      <!-- Pricing Tab -->
      <view v-show="activeTab === 'pricing'" class="tab-panel">
        <view class="pricing-card">
          <view class="pricing-row">
            <view class="pricing-info">
              <text class="pricing-label">输入价格</text>
              <text class="pricing-desc">每百万 tokens</text>
            </view>
            <view class="pricing-amount">
              <text class="pricing-currency">$</text>
              <text class="pricing-number">{{ model.pricing.inputPrice }}</text>
            </view>
          </view>
          <view class="pricing-divider"></view>
          <view class="pricing-row">
            <view class="pricing-info">
              <text class="pricing-label">输出价格</text>
              <text class="pricing-desc">每百万 tokens</text>
            </view>
            <view class="pricing-amount">
              <text class="pricing-currency">$</text>
              <text class="pricing-number">{{ model.pricing.outputPrice }}</text>
            </view>
          </view>
        </view>

        <view class="pricing-note">
          <text class="note-title">定价说明</text>
          <text class="note-text">价格可能因使用量、地区等因素有所变动，请以官方最新信息为准。</text>
        </view>
      </view>

      <!-- Benchmarks Tab -->
      <view v-show="activeTab === 'benchmarks'" class="tab-panel">
        <view class="benchmarks-list">
          <view v-if="model.quality.mmlu" class="benchmark-item">
            <view class="benchmark-header">
              <text class="benchmark-name">MMLU</text>
              <text class="benchmark-score">{{ model.quality.mmlu }}%</text>
            </view>
            <view class="benchmark-bar">
              <view class="benchmark-fill" :style="{ width: model.quality.mmlu + '%' }"></view>
            </view>
            <text class="benchmark-desc">多学科知识理解</text>
          </view>

          <view v-if="model.quality.humanEval" class="benchmark-item">
            <view class="benchmark-header">
              <text class="benchmark-name">HumanEval</text>
              <text class="benchmark-score">{{ model.quality.humanEval }}%</text>
            </view>
            <view class="benchmark-bar">
              <view class="benchmark-fill" :style="{ width: model.quality.humanEval + '%' }"></view>
            </view>
            <text class="benchmark-desc">代码生成能力</text>
          </view>

          <view v-if="model.quality.mtBench" class="benchmark-item">
            <view class="benchmark-header">
              <text class="benchmark-name">MT-Bench</text>
              <text class="benchmark-score">{{ model.quality.mtBench }}</text>
            </view>
            <view class="benchmark-bar">
              <view class="benchmark-fill" :style="{ width: (model.quality.mtBench / 10 * 100) + '%' }"></view>
            </view>
            <text class="benchmark-desc">多轮对话质量</text>
          </view>

          <view v-if="model.quality.math" class="benchmark-item">
            <view class="benchmark-header">
              <text class="benchmark-name">Math</text>
              <text class="benchmark-score">{{ model.quality.math }}%</text>
            </view>
            <view class="benchmark-bar">
              <view class="benchmark-fill" :style="{ width: model.quality.math + '%' }"></view>
            </view>
            <text class="benchmark-desc">数学推理能力</text>
          </view>
        </view>
      </view>

      <!-- Bottom Spacer -->
      <view class="bottom-spacer"></view>
    </scroll-view>
  </view>
</template>

<style scoped>
.page {
  height: 100vh;
  background: #f7f8fa;
  display: flex;
  flex-direction: column;
}

/* ===== Header ===== */
.header {
  background: #1a1a2e;
  padding: 24rpx 32rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.back-btn {
  width: 64rpx;
  height: 64rpx;
  font-size: 40rpx;
  color: #fff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.header-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #fff;
}

.header-placeholder {
  width: 64rpx;
}

/* ===== Hero Card ===== */
.hero-card {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  margin: 24rpx 32rpx;
  border-radius: 24rpx;
  padding: 40rpx 32rpx;
  position: relative;
  overflow: hidden;
}

.hero-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -20%;
  width: 60%;
  height: 200%;
  background: radial-gradient(ellipse, rgba(255,255,255,0.1) 0%, transparent 70%);
}

.hero-badge {
  display: inline-block;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  font-size: 22rpx;
  font-weight: 600;
  padding: 8rpx 16rpx;
  border-radius: 8rpx;
  margin-bottom: 16rpx;
}

.hero-name {
  font-size: 48rpx;
  font-weight: 700;
  color: #fff;
  display: block;
  margin-bottom: 16rpx;
  position: relative;
  z-index: 1;
}

.hero-desc {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.85);
  line-height: 1.5;
  display: block;
  margin-bottom: 24rpx;
  position: relative;
  z-index: 1;
}

.hero-chips {
  display: flex;
  gap: 12rpx;
  position: relative;
  z-index: 1;
}

.chip {
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
  font-size: 22rpx;
  padding: 10rpx 20rpx;
  border-radius: 8rpx;
}

/* ===== Tabs ===== */
.tabs {
  display: flex;
  background: #fff;
  margin: 0 32rpx 24rpx;
  border-radius: 16rpx;
  padding: 8rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.tab {
  flex: 1;
  text-align: center;
  padding: 20rpx 0;
  font-size: 28rpx;
  color: #6b7280;
  border-radius: 12rpx;
  transition: all 0.2s;
}

.tab.active {
  background: #6366f1;
  color: #fff;
  font-weight: 600;
}

/* ===== Content ===== */
.content {
  flex: 1;
  padding: 0 32rpx;
}

.tab-panel {
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

.section {
  background: #fff;
  border-radius: 20rpx;
  padding: 32rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.section-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #1a1a2e;
  margin-bottom: 24rpx;
  display: block;
}

/* ===== Capabilities ===== */
.caps-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16rpx;
}

.cap-item {
  background: #f9fafb;
  border-radius: 16rpx;
  padding: 24rpx 16rpx;
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.cap-name {
  font-size: 26rpx;
  color: #4b5563;
}

.cap-check {
  font-size: 24rpx;
  color: #10b981;
  font-weight: 700;
}

/* ===== Specs ===== */
.specs-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.spec-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 2rpx solid #f3f4f6;
}

.spec-item:last-child {
  border-bottom: none;
}

.spec-label {
  font-size: 28rpx;
  color: #6b7280;
}

.spec-value {
  font-size: 28rpx;
  font-weight: 600;
  color: #1a1a2e;
}

/* ===== Features ===== */
.features-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.feature-item {
  background: #f9fafb;
  padding: 20rpx 24rpx;
  border-radius: 12rpx;
  font-size: 28rpx;
  color: #4b5563;
  position: relative;
  padding-left: 48rpx;
}

.feature-item::before {
  content: '•';
  position: absolute;
  left: 20rpx;
  color: #6366f1;
  font-weight: 700;
}

/* ===== Pricing ===== */
.pricing-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 32rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.pricing-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.pricing-info {
  display: flex;
  flex-direction: column;
}

.pricing-label {
  font-size: 30rpx;
  font-weight: 600;
  color: #1a1a2e;
}

.pricing-desc {
  font-size: 24rpx;
  color: #9ca3af;
  margin-top: 4rpx;
}

.pricing-amount {
  display: flex;
  align-items: flex-start;
}

.pricing-currency {
  font-size: 28rpx;
  font-weight: 600;
  color: #6366f1;
  margin-top: 8rpx;
}

.pricing-number {
  font-size: 56rpx;
  font-weight: 700;
  color: #6366f1;
}

.pricing-divider {
  height: 2rpx;
  background: #f3f4f6;
  margin: 28rpx 0;
}

.pricing-note {
  background: #fef3c7;
  border-radius: 16rpx;
  padding: 24rpx;
}

.note-title {
  font-size: 26rpx;
  font-weight: 600;
  color: #92400e;
  display: block;
  margin-bottom: 8rpx;
}

.note-text {
  font-size: 24rpx;
  color: #b45309;
  line-height: 1.5;
}

/* ===== Benchmarks ===== */
.benchmarks-list {
  display: flex;
  flex-direction: column;
  gap: 24rpx;
}

.benchmark-item {
  background: #fff;
  border-radius: 20rpx;
  padding: 32rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.benchmark-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.benchmark-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #1a1a2e;
}

.benchmark-score {
  font-size: 36rpx;
  font-weight: 700;
  color: #6366f1;
}

.benchmark-bar {
  height: 12rpx;
  background: #f3f4f6;
  border-radius: 6rpx;
  overflow: hidden;
  margin-bottom: 12rpx;
}

.benchmark-fill {
  height: 100%;
  background: linear-gradient(90deg, #6366f1 0%, #8b5cf6 100%);
  border-radius: 6rpx;
  transition: width 0.6s ease;
}

.benchmark-desc {
  font-size: 24rpx;
  color: #9ca3af;
}

.bottom-spacer {
  height: 48rpx;
}
</style>
