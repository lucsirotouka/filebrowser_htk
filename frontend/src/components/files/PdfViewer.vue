<template>
  <!--
    position:fixed + top:4em fills the area below the header bar.
    @wheel.stop / @touchmove.stop prevent bubbling to #previewer's prevent handler.
  -->
  <div :class="['pdf-viewer', { 'pdf-viewer--inline': props.inline }]" ref="viewerWrap" @touchmove.stop @wheel.stop>
    <vue-pdf-embed
      v-if="containerWidth > 0"
      :source="pdfSource"
      :width="renderWidth"
      :disable-text-layer="true"
      :disable-annotation-layer="true"
      @loaded="onLoaded"
      @loading-failed="onError"
    />

    <div v-if="errorMsg" class="pdf-error">
      <i class="material-icons">error_outline</i>
      <p>{{ errorMsg }}</p>
    </div>

    <!-- Floating zoom toolbar – sticky at bottom so it scrolls with the content -->
    <div class="pdf-zoom-controls">
      <button class="pdf-zoom-btn" @click="zoomOut" :title="t('buttons.zoomOut')">
        <i class="material-icons">remove</i>
      </button>
      <span class="pdf-zoom-label">{{ Math.round(zoom) }}%</span>
      <button class="pdf-zoom-btn" @click="zoomIn" :title="t('buttons.zoomIn')">
        <i class="material-icons">add</i>
      </button>
      <button class="pdf-zoom-btn" @click="resetZoom" :title="t('buttons.resetZoom')">
        <i class="material-icons">fit_screen</i>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from "vue";
import { useI18n } from "vue-i18n";
import VuePdfEmbed from "vue-pdf-embed";
// vue-pdf-embed v2 bundles its own pdfjs-dist; do NOT import pdfjs-dist
// separately to avoid version conflicts.

const props = defineProps<{
  src: string;
  /** When true the viewer is embedded inline (relative), not full-screen fixed */
  inline?: boolean;
}>();
const { t } = useI18n();

const viewerWrap = ref<HTMLElement | null>(null);
const containerWidth = ref(0);
const zoom = ref(100); // percentage  (100 = fit container width)
const errorMsg = ref("");

const pdfSource = computed(() => ({
  url: props.src,
  withCredentials: true,
}));

// Render width responds to zoom; 32 px subtracted for left+right padding
const renderWidth = computed(() =>
  Math.max(100, Math.round((containerWidth.value - 32) * (zoom.value / 100)))
);

const onLoaded = () => { errorMsg.value = ""; };
const onError = (err: unknown) => {
  console.error("PDF load error:", err);
  errorMsg.value = err instanceof Error ? err.message : "Failed to load PDF.";
};

const zoomIn  = () => { zoom.value = Math.min(300, zoom.value + 10); };
const zoomOut = () => { zoom.value = Math.max(20,  zoom.value - 10); };
const resetZoom = () => { zoom.value = 100; };

// Keep containerWidth in sync with the element's actual width
let ro: ResizeObserver | null = null;
onMounted(() => {
  if (!viewerWrap.value) return;
  containerWidth.value = viewerWrap.value.clientWidth;
  ro = new ResizeObserver((entries) => {
    const w = entries[0]?.contentRect.width ?? 0;
    if (w > 0) containerWidth.value = w;
  });
  ro.observe(viewerWrap.value);
});
onBeforeUnmount(() => ro?.disconnect());
</script>

<style scoped>
/*
  position:fixed so the viewer fills exactly the space below the header (4em)
  regardless of any parent overflow settings.  The scroll happens here, not on
  #previewer.
*/
.pdf-viewer {
  position: fixed;
  top: 4em;
  left: 0;
  right: 0;
  bottom: 0;
  overflow-y: auto;
  overflow-x: hidden;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 1em 1em 5em;
  background: #1e1e1e;
  box-sizing: border-box;
  z-index: 1;
}

/* Gap between consecutive PDF pages */
.pdf-viewer :deep(.vue-pdf-embed > div) {
  margin-bottom: 0.75em;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.6);
}

.pdf-error {
  padding: 2em;
  text-align: center;
  color: #ef5350;
}

/* Sticky zoom toolbar – stays visible at bottom-right while scrolling */
.pdf-zoom-controls {
  position: sticky;
  bottom: 1em;
  align-self: flex-end;
  display: flex;
  align-items: center;
  gap: 0.2em;
  background: rgba(40, 40, 40, 0.92);
  border-radius: 24px;
  padding: 0.3em 0.7em;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.5);
  z-index: 10;
  backdrop-filter: blur(4px);
}

.pdf-zoom-label {
  min-width: 3.5em;
  text-align: center;
  font-size: 0.85em;
  color: rgba(255, 255, 255, 0.8);
}

.pdf-zoom-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25em;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: rgba(255, 255, 255, 0.75);
  transition: background 0.15s, color 0.15s;
}
.pdf-zoom-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.15);
}
.pdf-zoom-btn .material-icons {
  font-size: 1.25em;
}

/*
  Inline mode: no fixed positioning; fills its parent container.
  Used when embedding inside Share.vue (sidebar layout).
*/
.pdf-viewer--inline {
  position: relative;
  top: auto;
  left: auto;
  right: auto;
  bottom: auto;
  height: 80vh;
  width: 100%;
  background: #2a2a2a;
  z-index: auto;
}
</style>
