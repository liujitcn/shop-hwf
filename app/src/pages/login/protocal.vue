<script setup lang="ts">
import { useSettingStore } from '@/stores'
const settingStore = useSettingStore()
</script>
<template>
  <view class="container">
    <rich-text :nodes="settingStore.getData('protocal')" />
  </view>
</template>
