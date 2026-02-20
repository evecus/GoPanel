<template>
  <div>
    <div style="display:flex;align-items:center;gap:10px;margin-bottom:16px">
      <button class="btn btn-ghost btn-sm" @click="load">🔄 {{ t('refresh') }}</button>
      <span style="font-size:12px;color:#9ca3af">{{ t('crontab_desc') }}</span>
    </div>
    <div class="card" v-if="jobs.length">
      <table class="table">
        <thead><tr><th>{{ t('source') }}</th><th>{{ t('schedule') }}</th><th>{{ t('command') }}</th></tr></thead>
        <tbody>
          <tr v-for="(job,i) in jobs" :key="i">
            <td>
              <span class="tag" :class="job.source==='[user]'?'tag-indigo':'tag-violet'">{{ job.source }}</span>
            </td>
            <td><code class="cron-expr" :title="cronDesc(job.schedule)">{{ job.schedule }}</code></td>
            <td class="mono cmd-cell">{{ job.command }}</td>
          </tr>
        </tbody>
      </table>
    </div>
    <div class="card empty-state" v-else>
      <div class="empty-icon">🕐</div>
      <div class="empty-title">{{ t('no_crontab') }}</div>
      <div class="empty-sub">{{ t('no_crontab_sub') }}</div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { useI18n } from '../stores/i18n.js'
const i18n = useI18n(); const t = k => i18n.t(k)
const jobs = ref([])
async function load() {
  const {data} = await axios.get('/api/crontab')
  jobs.value = (data||[]).map(line => {
    const parts = line.match(/^(\[.*?\])\s+(.+)$/)
    if (parts) {
      const rest = parts[2].trim().split(/\s+/)
      const schedule = rest.slice(0,5).join(' ')
      const command = rest.slice(5).join(' ')
      return { source: parts[1], schedule, command }
    }
    const rest = line.trim().split(/\s+/)
    return { source: '[system]', schedule: rest.slice(0,5).join(' '), command: rest.slice(5).join(' ') }
  })
}
function cronDesc(expr) {
  const parts = expr.split(' ')
  if (parts.length < 5) return expr
  const [min,hour,dom,mon,dow] = parts
  if (min==='*'&&hour==='*') return 'Every minute'
  if (dom==='*'&&mon==='*'&&dow==='*') return `Daily at ${hour}:${min}`
  return expr
}
onMounted(load)
</script>
<style scoped>
.card { background:#fff;border:1px solid rgba(99,102,241,0.1);border-radius:14px;padding:0;box-shadow:0 2px 12px rgba(99,102,241,0.06);overflow:hidden; }
.cron-expr { font-family:'JetBrains Mono',monospace;font-size:12px;background:rgba(99,102,241,0.08);color:#4f46e5;padding:2px 8px;border-radius:4px;cursor:help; }
.cmd-cell { font-family:'JetBrains Mono',monospace;font-size:12px;max-width:300px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap; }
.empty-state { text-align:center;padding:60px 20px; }
.empty-icon  { font-size:48px;margin-bottom:16px; }
.empty-title { font-size:16px;font-weight:600;color:#1e1b4b;margin-bottom:6px; }
.empty-sub   { font-size:13px;color:#9ca3af; }
</style>
