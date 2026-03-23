<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { paperService } from '@/services';
import type { Paper } from '@/types/api';

const paper = ref<Paper | null>(null);
const loading = ref<boolean>(true);
const showFullAbstract = ref<boolean>(false);

onMounted(() => {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1];
  const id = (currentPage as unknown as { $page?: { options?: { id?: string } } }).$page?.options?.id;
  
  if (id) {
    loadPaper(id);
  }
});

async function loadPaper(id: string): Promise<void> {
  try {
    loading.value = true;
    paper.value = await paperService.getPaperById(id);
  } catch (error) {
    console.error('加载失败:', error);
  } finally {
    loading.value = false;
  }
}

function openLink(url: string): void {
  uni.navigateTo({
    url: `/pages/webview/webview?url=${encodeURIComponent(url)}`
  });
}

function goBack(): void {
  uni.navigateBack();
}

// 根据 arxiv id 构造链接
function getArxivUrl(id: string): string {
  return `https://arxiv.org/abs/${id}`;
}

function getPdfUrl(id: string): string {
  return `https://arxiv.org/pdf/${id}.pdf`;
}

// 分割作者字符串为数组
function parseAuthors(authorStr: string | null | undefined): string[] {
  if (!authorStr) return [];
  return authorStr.split(',').map(s => s.trim()).filter(Boolean);
}
</script>

<template>
  <view class="page" v-if="paper">
    <!-- Header -->
    <view class="header">
      <text class="back-btn" @click="goBack">‹</text>
      <text class="header-title">论文详情</text>
      <view class="header-placeholder"></view>
    </view>

    <!-- Content -->
    <scroll-view scroll-y class="content">
      <!-- Paper Card -->
      <view class="paper-card">
        <view class="paper-badge">arXiv CS.AI</view>
        <text class="paper-title-cn">{{ paper.title_cn || paper.title }}</text>
        <text class="paper-title-en">{{ paper.title }}</text>
        
        <view class="paper-meta">
          <view class="meta-item">
            <text class="meta-icon">◉</text>
            <text class="meta-text">{{ paper.author ? paper.author.split(',')[0] + ' 等' : '未知作者' }}</text>
          </view>
          <view class="meta-item">
            <text class="meta-icon">◈</text>
            <text class="meta-text">{{ paper.submit_at }}</text>
          </view>
        </view>
      </view>

      <!-- Authors -->
      <view class="section" v-if="parseAuthors(paper.author).length > 0">
        <text class="section-title">作者</text>
        <view class="authors-list">
          <view 
            v-for="(author, idx) in parseAuthors(paper.author)" 
            :key="idx"
            class="author-item"
          >
            <view class="author-avatar">{{ author.charAt(0) }}</view>
            <text class="author-name">{{ author }}</text>
          </view>
        </view>
      </view>

      <!-- Abstract -->
      <view class="section">
        <view class="section-header">
          <text class="section-title">摘要</text>
          <text 
            class="section-action"
            @click="showFullAbstract = !showFullAbstract"
          >
            {{ showFullAbstract ? '收起' : '展开' }}
          </text>
        </view>
        
        <view class="abstract-content">
          <text class="abstract-cn" :class="{ expanded: showFullAbstract }">
            {{ paper.abstract_cn || paper.abstract || '暂无中文摘要' }}
          </text>
          
          <view v-if="showFullAbstract && paper.abstract" class="abstract-en">
            <view class="divider"></view>
            <text class="abstract-en-text">{{ paper.abstract }}</text>
          </view>
        </view>
      </view>

      <!-- Actions -->
      <view class="actions">
        <button class="btn btn--primary" @click="openLink(getArxivUrl(paper.id))">
          <text class="btn-icon">↗</text>
          <text class="btn-text">查看 arXiv</text>
        </button>
        <button class="btn btn--secondary" @click="openLink(getPdfUrl(paper.id))">
          <text class="btn-icon">↓</text>
          <text class="btn-text">下载 PDF</text>
        </button>
      </view>

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

/* ===== Content ===== */
.content {
  flex: 1;
  padding: 24rpx 32rpx;
}

/* ===== Paper Card ===== */
.paper-card {
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  border-radius: 24rpx;
  padding: 40rpx 32rpx;
  margin-bottom: 24rpx;
  position: relative;
  overflow: hidden;
}

.paper-card::before {
  content: '';
  position: absolute;
  top: -50%;
  right: -20%;
  width: 60%;
  height: 200%;
  background: radial-gradient(ellipse, rgba(255,255,255,0.1) 0%, transparent 70%);
}

.paper-badge {
  display: inline-block;
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  font-size: 22rpx;
  font-weight: 600;
  padding: 8rpx 16rpx;
  border-radius: 8rpx;
  margin-bottom: 20rpx;
  position: relative;
  z-index: 1;
}

.paper-title-cn {
  font-size: 36rpx;
  font-weight: 700;
  color: #fff;
  line-height: 1.4;
  display: block;
  margin-bottom: 16rpx;
  position: relative;
  z-index: 1;
}

.paper-title-en {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
  line-height: 1.5;
  display: block;
  position: relative;
  z-index: 1;
}

.paper-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 20rpx;
  margin-top: 28rpx;
  padding-top: 24rpx;
  border-top: 2rpx solid rgba(255, 255, 255, 0.1);
  position: relative;
  z-index: 1;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 10rpx;
}

.meta-icon {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.6);
}

.meta-text {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.85);
}

/* ===== Section ===== */
.section {
  background: #fff;
  border-radius: 20rpx;
  padding: 32rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #1a1a2e;
}

.section-action {
  font-size: 26rpx;
  color: #6366f1;
  font-weight: 500;
}

/* ===== Authors ===== */
.authors-list {
  display: flex;
  flex-wrap: wrap;
  gap: 20rpx;
}

.author-item {
  display: flex;
  align-items: center;
  gap: 12rpx;
  background: #f9fafb;
  padding: 12rpx 20rpx;
  border-radius: 100rpx;
}

.author-avatar {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: #fff;
  font-size: 20rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
}

.author-name {
  font-size: 26rpx;
  color: #4b5563;
  font-weight: 500;
}

/* ===== Abstract ===== */
.abstract-content {
  position: relative;
}

.abstract-cn {
  font-size: 30rpx;
  color: #374151;
  line-height: 1.8;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.abstract-cn.expanded {
  display: block;
  -webkit-line-clamp: unset;
}

.divider {
  height: 2rpx;
  background: #e5e7eb;
  margin: 24rpx 0;
}

.abstract-en {
  margin-top: 24rpx;
}

.abstract-en-text {
  font-size: 26rpx;
  color: #6b7280;
  line-height: 1.7;
  font-style: italic;
}

/* ===== Actions ===== */
.actions {
  display: flex;
  gap: 20rpx;
  margin-top: 8rpx;
}

.btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  padding: 28rpx 0;
  border-radius: 16rpx;
  border: none;
}

.btn--primary {
  background: #1a1a2e;
}

.btn--secondary {
  background: #fff;
  border: 2rpx solid #e5e7eb;
}

.btn-icon {
  font-size: 28rpx;
  color: #6366f1;
}

.btn--primary .btn-icon {
  color: #fff;
}

.btn-text {
  font-size: 28rpx;
  font-weight: 600;
  color: #1a1a2e;
}

.btn--primary .btn-text {
  color: #fff;
}

.bottom-spacer {
  height: 48rpx;
}
</style>
