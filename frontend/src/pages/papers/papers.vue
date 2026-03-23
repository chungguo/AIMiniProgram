<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { paperService } from '@/services';
import type { Paper } from '@/types/api';

const papers = ref<Paper[]>([]);
const loading = ref<boolean>(false);
const currentPage = ref<number>(1);
const hasMore = ref<boolean>(true);

const PAGE_SIZE = 10;

onMounted(() => {
  loadPapers();
});

async function loadPapers(refresh: boolean = false): Promise<void> {
  if (loading.value) return;
  
  if (refresh) {
    currentPage.value = 1;
    hasMore.value = true;
    papers.value = [];
  }
  
  try {
    loading.value = true;
    const res = await paperService.getPapers({
      page: currentPage.value,
      limit: PAGE_SIZE
    });
    
    if (refresh) papers.value = res.data;
    else papers.value.push(...res.data);
    // 判断是否还有更多：当前页返回数据量等于 pageSize 且总数大于已加载数量
    hasMore.value = res.data.length === PAGE_SIZE && papers.value.length < res.pagination.total;
  } catch (error) {
    console.error('加载论文失败:', error);
  } finally {
    loading.value = false;
  }
}

function loadMore(): void {
  if (!hasMore.value || loading.value) return;
  currentPage.value++;
  loadPapers();
}

function goToDetail(id: string): void {
  uni.navigateTo({ url: `/pages/paper-detail/paper-detail?id=${id}` });
}
</script>

<template>
  <view class="page">
    <!-- Header -->
    <view class="header">
      <text class="header-title">论文库</text>
      <text class="header-count">{{ papers.length }} 篇论文</text>
    </view>

    <!-- Papers List -->
    <scroll-view 
      scroll-y 
      class="content"
      @scrolltolower="loadMore"
      refresher-enabled
      @refresherrefresh="loadPapers(true)"
    >
      <view class="papers-list">
        <view 
          v-for="paper in papers" 
          :key="paper.id"
          class="paper-card"
          @click="goToDetail(paper.id)"
        >
          <view class="paper-header">
            <text class="paper-category">arXiv CS.AI</text>
            <text class="paper-date">{{ paper.submit_at }}</text>
          </view>
          
          <text class="paper-title">{{ paper.title_cn || paper.title }}</text>
          <text class="paper-abstract">{{ (paper.abstract_cn || paper.abstract || '').slice(0, 80) }}...</text>
          
          <view class="paper-footer">
            <view class="paper-authors">
              <text class="author-avatar" v-for="(author, idx) in (paper.author || '').split(',').map(s => s.trim()).filter(Boolean).slice(0, 2)" :key="idx">
                {{ author.charAt(0) }}
              </text>
              <text v-if="(paper.author || '').split(',').length > 2" class="author-more">+{{ (paper.author || '').split(',').length - 2 }}</text>
            </view>
            <view class="paper-meta">
              <text class="paper-arrow">›</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Load More -->
      <view class="load-more">
        <text v-if="loading" class="load-text">加载中...</text>
        <text v-else-if="!hasMore" class="load-text load-text--end">已经到底了</text>
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
  background: linear-gradient(165deg, #1a1a2e 0%, #16213e 100%);
  padding: 32rpx;
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

/* ===== Content ===== */
.content {
  flex: 1;
  padding: 24rpx 32rpx;
}

.papers-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.paper-card {
  background: #fff;
  border-radius: 20rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.paper-card:active {
  transform: scale(0.98);
}

.paper-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.paper-category {
  font-size: 22rpx;
  font-weight: 600;
  color: #6366f1;
  background: rgba(99, 102, 241, 0.1);
  padding: 6rpx 14rpx;
  border-radius: 8rpx;
}

.paper-date {
  font-size: 22rpx;
  color: #9ca3af;
}

.paper-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #1a1a2e;
  line-height: 1.4;
  display: block;
  margin-bottom: 12rpx;
}

.paper-abstract {
  font-size: 26rpx;
  color: #6b7280;
  line-height: 1.6;
  margin-bottom: 20rpx;
}

.paper-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.paper-authors {
  display: flex;
  align-items: center;
}

.author-avatar {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  color: #fff;
  font-size: 22rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: -12rpx;
  border: 2rpx solid #fff;
}

.author-more {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: #f3f4f6;
  color: #6b7280;
  font-size: 20rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-left: 8rpx;
}

.paper-meta {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.paper-arrow {
  font-size: 32rpx;
  color: #d1d5db;
  font-weight: 300;
}

/* ===== Load More ===== */
.load-more {
  text-align: center;
  padding: 40rpx 0;
}

.load-text {
  font-size: 26rpx;
  color: #9ca3af;
}

.load-text--end {
  color: #d1d5db;
}

.bottom-spacer {
  height: 32rpx;
}
</style>
