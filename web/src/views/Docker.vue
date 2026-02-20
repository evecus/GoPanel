<template>
  <div>
    <div class="toolbar">
      <input class="inp" style="width:200px" v-model="search" :placeholder="`🔍 ${t('search')}...`" />
      <button class="btn btn-ghost btn-sm" @click="load">🔄 {{ t('refresh') }}</button>
      <span style="margin-left:auto;font-size:12px;color:#9ca3af">{{ filtered.length }} {{ t('container') }}</span>
    </div>

    <div v-if="filtered.length" class="container-grid">
      <div class="ccard" v-for="c in filtered" :key="c.id">
        <div class="cc-head">
          <div class="cc-dot" :class="c.state==='running'?'run':c.state==='paused'?'pause':'stop'"></div>
          <div class="cc-name">{{ c.name }}</div>
          <span class="tag" :class="stateTag(c.state)">{{ c.state }}</span>
        </div>
        <div class="cc-img">🐳 {{ c.image }}</div>
        <div v-if="c.ports" class="cc-ports">
          <span class="tag tag-blue" v-for="p in c.ports?.split(',').slice(0,3)" :key="p">{{ p.trim() }}</span>
        </div>
        <div class="cc-metrics" v-if="c.state==='running'">
          <div class="cm">
            <span style="font-size:10px;color:#9ca3af">CPU</span>
            <div class="mini-bar"><div class="mini-fill" :style="`width:${c.cpu_percent||0}%;background:#6366f1`"></div></div>
            <span style="font-size:11px;color:#6366f1;font-weight:600">{{ c.cpu_percent?.toFixed(1) }}%</span>
          </div>
          <div class="cm">
            <span style="font-size:10px;color:#9ca3af">MEM</span>
            <div class="mini-bar"><div class="mini-fill" :style="`width:${c.mem_percent||0}%;background:#06b6d4`"></div></div>
            <span style="font-size:11px;color:#06b6d4;font-weight:600">{{ c.mem_percent?.toFixed(1) }}%</span>
          </div>
        </div>
        <div class="cc-actions">
          <button class="btn btn-sm btn-cyan" v-if="c.state!=='running'" @click="action(c,'start')">▶ {{ t('start') }}</button>
          <button class="btn btn-sm btn-ghost" v-if="c.state==='running'" @click="action(c,'stop')">⏹ {{ t('stop') }}</button>
          <button class="btn btn-sm btn-ghost" @click="action(c,'restart')">↺ {{ t('restart') }}</button>
          <button class="btn btn-sm btn-ghost" style="margin-left:auto" @click="showLogs(c)">📋 {{ t('logs') }}</button>
        </div>
      </div>
    </div>

    <div class="card empty-state" v-else>
      <div style="font-size:48px;margin-bottom:16px">🐳</div>
      <div style="font-size:16px;font-weight:600;color:#1e1b4b;margin-bottom:6px">{{ t('no_docker') }}</div>
      <div style="font-size:13px;color:#9ca3af">Make sure Docker is running</div>
    </div>

    <!-- Logs modal -->
    <div class="modal-overlay" v-if="logModal" @click.self="logModal=null">
      <div class="modal" style="width:720px;max-width:96vw">
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:16px">
          <h3 style="color:#1e1b4b;font-size:16px">📋 {{ logModal.name }} {{ t('logs') }}</h3>
          <button class="btn btn-ghost btn-sm" @click="logModal=null">✕ {{ t('close') }}</button>
        </div>
        <pre class="log-box">{{ logContent || t('loading') }}</pre>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useI18n } from '../stores/i18n.js'
const i18n = useI18n(); const t = k => i18n.t(k)
const containers = ref([]); const search = ref(''); const logModal = ref(null); const logContent = ref('')
const filtered = computed(() => containers.value.filter(c => c.name?.toLowerCase().includes(search.value.toLowerCase()) || c.image?.toLowerCase().includes(search.value.toLowerCase())))
function stateTag(s) { return s==='running'?'tag tag-green':s==='paused'?'tag tag-yellow':'tag tag-gray' }
async function load() { try { const {data}=await axios.get('/api/docker/containers'); containers.value=data||[] } catch { containers.value=[] } }
async function action(c,act) { try { await axios.post(`/api/docker/containers/${c.id}/${act}`); setTimeout(load,1000) } catch(e) { alert(e.response?.data?.error||e.message) } }
async function showLogs(c) { logModal.value=c; logContent.value=''; const {data}=await axios.get(`/api/docker/containers/${c.id}/logs`); logContent.value=data.logs||'' }
onMounted(load)
</script>
<style scoped>
.toolbar { display:flex;align-items:center;gap:8px;flex-wrap:wrap;margin-bottom:14px; }
.inp { background:#f8faff;border:1.5px solid rgba(99,102,241,0.15);color:#1e1b4b;border-radius:8px;padding:8px 12px;font-size:13px;font-family:inherit;outline:none; }
.inp:focus { border-color:#6366f1; }
.container-grid { display:grid;grid-template-columns:repeat(auto-fill,minmax(300px,1fr));gap:14px; }
.ccard { background:#fff;border:1px solid rgba(99,102,241,0.1);border-radius:14px;padding:16px;box-shadow:0 2px 12px rgba(99,102,241,0.06);transition:transform 0.2s; }
.ccard:hover { transform:translateY(-2px); }
.cc-head { display:flex;align-items:center;gap:8px;margin-bottom:10px; }
.cc-dot  { width:9px;height:9px;border-radius:50%;flex-shrink:0; }
.cc-dot.run   { background:#10b981;box-shadow:0 0 8px rgba(16,185,129,0.5);animation:pulse 2s infinite; }
.cc-dot.pause { background:#f59e0b; }
.cc-dot.stop  { background:#9ca3af; }
.cc-name { font-size:15px;font-weight:700;color:#1e1b4b;flex:1;overflow:hidden;text-overflow:ellipsis; }
.cc-img  { font-size:12px;color:#9ca3af;margin-bottom:8px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap; }
.cc-ports { display:flex;flex-wrap:wrap;gap:4px;margin-bottom:10px; }
.cc-metrics { display:flex;flex-direction:column;gap:6px;margin-bottom:12px;background:rgba(99,102,241,0.04);border-radius:8px;padding:10px; }
.cm { display:flex;align-items:center;gap:6px; }
.mini-bar { flex:1;height:4px;background:#f0f4ff;border-radius:2px;overflow:hidden; }
.mini-fill { height:100%;border-radius:2px;transition:width 0.5s; }
.cc-actions { display:flex;gap:6px;flex-wrap:wrap; }
.card { background:#fff;border:1px solid rgba(99,102,241,0.1);border-radius:14px;padding:18px;box-shadow:0 2px 12px rgba(99,102,241,0.06); }
.empty-state { text-align:center;padding:60px 20px; }
.modal-overlay { position:fixed;inset:0;background:rgba(30,27,75,0.4);backdrop-filter:blur(6px);display:flex;align-items:center;justify-content:center;z-index:1000; }
.modal { background:#fff;border:1px solid rgba(99,102,241,0.15);border-radius:16px;padding:24px;box-shadow:0 20px 60px rgba(99,102,241,0.15); }
.log-box { background:#0f172a;color:#e2e8f0;border-radius:10px;padding:16px;font-family:'JetBrains Mono',monospace;font-size:12px;max-height:50vh;overflow-y:auto;white-space:pre-wrap;word-break:break-all; }
.btn { display:inline-flex;align-items:center;gap:4px;padding:7px 14px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;border:none;font-family:inherit;transition:all 0.2s; }
.btn-sm  { padding:5px 11px;font-size:12px; }
.btn-cyan  { background:linear-gradient(135deg,#06b6d4,#10b981);color:#fff;box-shadow:0 2px 8px rgba(6,182,212,0.3); }
.btn-ghost { background:#fff;color:#6b7280;border:1px solid rgba(99,102,241,0.15); }
.btn-ghost:hover { background:rgba(99,102,241,0.06); }
@keyframes pulse { 0%,100%{opacity:1}50%{opacity:.5} }
@media (max-width:600px) { .container-grid { grid-template-columns:1fr; } }
</style>
