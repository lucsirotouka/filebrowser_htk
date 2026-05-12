<template>
  <div id="previewer">
    <header-bar>
      <action icon="arrow_back" :label="$t('buttons.back')" @action="goBack" />
      <title>{{ filename }}</title>
      <action
        icon="file_download"
        :label="$t('buttons.download')"
        @action="download"
      />
    </header-bar>

    <div v-if="!src" class="message" style="padding-top: 6em; text-align: center">
      <i class="material-icons" style="font-size: 3em">error_outline</i>
      <p>No PDF source provided.</p>
    </div>
    <Suspense v-else>
      <PdfViewer :src="src" />
      <template #fallback>
        <div class="loading delayed">
          <div class="spinner">
            <div class="bounce1"></div>
            <div class="bounce2"></div>
            <div class="bounce3"></div>
          </div>
        </div>
      </template>
    </Suspense>
  </div>
</template>

<script setup lang="ts">
import { computed, defineAsyncComponent } from "vue";
import { useRoute, useRouter } from "vue-router";
import HeaderBar from "@/components/header/HeaderBar.vue";
import Action from "@/components/header/Action.vue";

const PdfViewer = defineAsyncComponent(
  () => import("@/components/files/PdfViewer.vue")
);

const route = useRoute();
const router = useRouter();

// The PDF URL is passed as query ?src=...
const src = computed(() => {
  const s = route.query.src;
  return typeof s === "string" ? s : "";
});

// Extract filename from URL for the title bar
const filename = computed(() => {
  if (!src.value) return "PDF";
  try {
    const url = new URL(src.value, window.location.origin);
    const parts = url.pathname.split("/");
    return decodeURIComponent(parts[parts.length - 1] || "PDF");
  } catch {
    return "PDF";
  }
});

const goBack = () => router.back();

const download = () => {
  if (src.value) {
    // Replace inline=true with inline=false (or add dl param) to trigger download
    const url = new URL(src.value, window.location.origin);
    url.searchParams.set("inline", "false");
    window.open(url.toString());
  }
};
</script>
