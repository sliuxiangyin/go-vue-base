<script setup>
import { computed } from "vue"
import { MdPreview,MdCatalog  } from 'md-editor-v3'
import 'md-editor-v3/lib/style.css'
const props = defineProps({
  source: { type: String, default: "" },
})

const id = 'preview-only';
const scrollElement = document.documentElement;
</script>

<template>
  <MdPreview :editorId="id" :modelValue="source" />
  <MdCatalog :editorId="id" :scrollElement="scrollElement" />
</template>

<style scoped>
.prose {
  font-size: 14px;
  line-height: 1.6;
}
</style>
