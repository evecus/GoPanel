<template>
  <div>
    <div class="toolbar">
      <input class="input" style="width:200px" v-model="search" :placeholder="`🔍 ${t('search')}...`" />
      <div style="display:flex;gap:4px">
        <button class="btn btn-sm" :class="sortBy==='cpu'?'btn-primary':'btn-ghost'" @click="setSortBy('cpu')">CPU</button>
        <button class="btn btn-sm" :class="sortBy==='mem'?'btn-primary':'btn-ghost'" @click="setSortBy('mem')">MEM</button>
        <button class="btn btn-sm btn-ghost" @click="toggleDir" :title="sortDir==='desc'?'从大到小':'从小到大'">{{ sortDir==='desc'?'↓ 从大':'↑ 从小' }}</button>
      </div>
      <button class="btn btn-ghost btn-sm" @click="load">🔄 {{ t('refresh') }}</button>
      <span style="margin-left:auto;font-size:12px;color:#9ca3af">{{ filtered.length }} / {{ processes.length }}</span>
    </div>
    <div class="card" style="margin-top:14px;padding:0;overflow:hidden">
      <div style="overflow-x:auto;max-height:640px;overflow-y:auto">
        <table class="table">
          <thead>
            <tr>
              <th>PID</th><th>{{ t('process_name') }}</th><th>{{ t('username') }}</th>
              <th>CPU%</th><th>MEM%</th><th>RSS</th><th>{{ t('status') }}</th><th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="p in filtered" :key="p.pid">
              <td class="mono" style="color:#6366f1">{{ p.pid }}</td>
              <td style="font-weight:600;color:#1e1b4b">{{ p.name }}</td>
              <td style="color:#9ca3af">{{ p.username }}</td>
              <td>
                <div style="display:flex;align-items:center;gap:6px">
                  <div style="width:40px;height:4px;background:#f0f4ff;border-radius:2px;overflow:hidden">
                    <div style="height:100%;border-radius:2px;transition:width 0.3s" :style="`width:${Math.min(p.cpu_percent,100)}%;background:${p.cpu_percent>50?'#f43f5e':p.cpu_percent>20?'#f59e0b':'#6366f1'}`"></div>
                  </div>
                  <span :style="`color:${p.cpu_percent>50?'#f43f5e':p.cpu_percent>20?'#f59e0b':'#6366f1'}`">{{ p.cpu_percent?.toFixed(1) }}%</span>
                </div>
              </td>
              <td :style="`color:${p.mem_percent>20?'#f43f5e':p.mem_percent>10?'#f59e0b':'#10b981'}`">{{ p.mem_percent?.toFixed(1) }}%</td>
              <td class="mono">{{ fmtMem(p.mem_rss) }}</td>
              <td><span class="tag" :class="stTag(p.status)">{{ p.status }}</span></td>
              <td>
                <button class="btn btn-danger btn-xs" @click="confirmKill(p)">✕</button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    <div class="modal-overlay" v-if="killTarget" @click.self="killTarget=null">
      <div class="modal">
        <h3 style="color:#1e1b4b;margin-bottom:10px">{{ t('confirm_kill') }}</h3>
        <p style="color:#4b5563;font-size:14px;margin-bottom:20px">{{ t('confirm_kill_desc') }} <strong style="color:#f43f5e">{{ killTarget.name }}</strong> (PID {{ killTarget.pid }})?</p>
        <div style="display:flex;gap:8px;justify-content:flex-end">
          <button class="btn btn-ghost" @click="killTarget=null">{{ t('cancel') }}</button>
          <button class="btn btn-danger" @click="doKill">{{ t('kill') }}</button>
        </div>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useI18n } from '../stores/i18n.js'
const i18n = useI18n(); const t = k => i18n.t(k)
const processes = ref([]); const search = ref(''); const sortBy = ref('cpu'); const sortDir = ref('desc'); const killTarget = ref(null)
const filtered = computed(() => processes.value.filter(p => p.name?.toLowerCase().includes(search.value.toLowerCase())))
function fmtMem(b) { if(!b) return '0'; if(b<1048576) return (b/1024).toFixed(0)+'KB'; return (b/1048576).toFixed(1)+'MB' }
function stTag(s) { if(s==='R'||s==='running') return 'tag tag-green'; if(s==='Z'||s==='zombie') return 'tag tag-red'; return 'tag tag-gray' }
function setSortBy(s) { sortBy.value=s; load() }
function toggleDir() { sortDir.value = sortDir.value==='desc'?'asc':'desc'; load() }
async function load() { const {data}=await axios.get(`/api/processes?sort=${sortBy.value}&dir=${sortDir.value}&limit=100`); processes.value=data||[] }
function confirmKill(p) { killTarget.value=p }
async function doKill() {
  try { await axios.delete(`/api/processes/${killTarget.value.pid}`); killTarget.value=null; load() }
  catch(e) { alert('Failed: '+(e.response?.data?.error||e.message)) }
}
onMounted(load)
</script>
<style scoped>
.toolbar { display:flex;align-items:center;gap:8px;flex-wrap:wrap; }
.input { background:#f8faff;border:1.5px solid rgba(99,102,241,0.15);color:#1e1b4b;border-radius:8px;padding:8px 12px;font-size:13px;font-family:inherit;outline:none; }
.input:focus { border-color:#6366f1; }
.card { background:#fff;border:1px solid rgba(99,102,241,0.1);border-radius:14px;box-shadow:0 2px 12px rgba(99,102,241,0.06); }
.btn { display:inline-flex;align-items:center;gap:4px;padding:7px 14px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;border:none;font-family:inherit;transition:all 0.2s; }
.btn-sm { padding:5px 11px;font-size:12px; }
.btn-xs { padding:3px 8px;font-size:11px;border-radius:6px; }
.btn-primary { background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;box-shadow:0 2px 8px rgba(99,102,241,0.3); }
.btn-ghost   { background:#fff;color:#6b7280;border:1px solid rgba(99,102,241,0.15); }
.btn-ghost:hover { background:rgba(99,102,241,0.06); }
.btn-danger  { background:rgba(244,63,94,0.08);color:#f43f5e;border:1px solid rgba(244,63,94,0.2); }
.btn-danger:hover { background:rgba(244,63,94,0.15); }
.modal-overlay { position:fixed;inset:0;background:rgba(30,27,75,0.4);backdrop-filter:blur(6px);display:flex;align-items:center;justify-content:center;z-index:1000; }
.modal { background:#fff;border:1px solid rgba(99,102,241,0.15);border-radius:16px;padding:28px;min-width:380px;box-shadow:0 20px 60px rgba(99,102,241,0.15); }
</style>
