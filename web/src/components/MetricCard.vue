<template>
  <div class="metric-card" :style="`--accent:${accent}`">
    <div class="mc-header">
      <div class="mc-emoji">{{ emoji }}</div>
      <div class="mc-ring-wrap">
        <svg viewBox="0 0 40 40" class="mc-ring">
          <circle cx="20" cy="20" r="16" class="ring-bg"/>
          <circle cx="20" cy="20" r="16" class="ring-val" :style="ringStyle"/>
        </svg>
        <span class="ring-pct">{{ Math.round(clampedVal) }}</span>
      </div>
    </div>
    <div class="mc-label">{{ title }}</div>
    <div class="mc-bar"><div class="mc-fill" :style="`width:${clampedVal}%`"></div></div>
    <div class="mc-sub">{{ sub }}</div>
    <div class="mc-accent-line"></div>
  </div>
</template>
<script setup>
import { computed } from 'vue'
const props = defineProps({ title:String, value:{type:Number,default:0}, sub:String, emoji:String, accent:{type:String,default:'#6366f1'} })
const clampedVal = computed(() => Math.min(Math.max(props.value||0, 0), 100))
const C = 2 * Math.PI * 16
const ringStyle = computed(() => ({ strokeDasharray: C, strokeDashoffset: C * (1 - clampedVal.value/100), stroke: props.accent }))
</script>
<style scoped>
.metric-card {
  background:#fff;border:1px solid rgba(99,102,241,0.12);
  border-radius:14px;padding:18px;
  box-shadow:0 2px 16px rgba(99,102,241,0.08);
  position:relative;overflow:hidden;
  transition:transform 0.2s,box-shadow 0.2s;
}
.metric-card:hover { transform:translateY(-2px);box-shadow:0 8px 28px rgba(99,102,241,0.14); }
.mc-accent-line { position:absolute;top:0;left:0;right:0;height:3px;background:var(--accent);border-radius:14px 14px 0 0; }
.mc-header { display:flex;align-items:center;justify-content:space-between;margin-bottom:10px; }
.mc-emoji { font-size:24px; }
.mc-ring-wrap { position:relative;width:44px;height:44px; }
.mc-ring { width:100%;height:100%;transform:rotate(-90deg); }
.ring-bg  { fill:none;stroke:#f0f4ff;stroke-width:4; }
.ring-val { fill:none;stroke-width:4;stroke-linecap:round;transition:stroke-dashoffset 0.6s ease; }
.ring-pct { position:absolute;inset:0;display:flex;align-items:center;justify-content:center;font-size:10px;font-weight:700;color:#1e1b4b; }
.mc-label { font-size:13px;font-weight:500;color:#6b7280;margin-bottom:6px; }
.mc-bar { height:5px;background:#f0f4ff;border-radius:3px;overflow:hidden;margin-bottom:6px; }
.mc-fill { height:100%;border-radius:3px;background:var(--accent);transition:width 0.6s ease; }
.mc-sub { font-size:11px;color:#9ca3af;overflow:hidden;text-overflow:ellipsis;white-space:nowrap; }
</style>
