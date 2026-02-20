<template>
  <div>
    <div class="toolbar">
      <input class="inp" style="width:200px" v-model="search" :placeholder="`🔍 ${t('search')}...`" />
      <select class="sel" v-model="filterState">
        <option value="">{{ t('all_states') }}</option>
        <option value="active">Active</option>
        <option value="inactive">Inactive</option>
        <option value="failed">Failed</option>
      </select>
      <div style="display:flex;gap:4px">
        <button class="btn btn-sm" :class="sortBy===''?'btn-primary':'btn-ghost'" @click="setSort('')">默认</button>
        <button class="btn btn-sm" :class="sortBy==='memory'?'btn-primary':'btn-ghost'" @click="setSort('memory')">内存</button>
        <button class="btn btn-sm" :class="sortBy==='cpu'?'btn-primary':'btn-ghost'" @click="setSort('cpu')">CPU</button>
        <button v-if="sortBy" class="btn btn-sm btn-ghost" @click="toggleDir" :title="sortDir==='desc'?'从大到小':'从小到大'">{{ sortDir==='desc'?'↓':'↑' }}</button>
      </div>
      <button class="btn btn-ghost btn-sm" @click="load">🔄 {{ t('refresh') }}</button>
      <span style="margin-left:auto;font-size:12px;color:#9ca3af">{{ filtered.length }} 个服务</span>
    </div>

    <div class="card" style="margin-top:14px;padding:0;overflow:hidden">
      <div style="overflow-x:auto">
        <table class="table">
          <thead><tr>
            <th>服务名</th>
            <th>描述</th>
            <th>状态</th>
            <th>内存</th>
            <th>CPU时间</th>
            <th>文件状态</th>
            <th>操作</th>
          </tr></thead>
          <tbody>
            <tr v-for="svc in filtered" :key="svc.unit" style="cursor:pointer" @click="showDetail(svc)">
              <td>
                <span style="font-weight:600;color:#1a1040;font-size:12px">{{ svc.unit }}</span>
                <div v-if="svc.main_pid && svc.main_pid!=='0'" style="font-size:10px;color:#9ca3af">PID {{ svc.main_pid }}</div>
              </td>
              <td style="max-width:180px;overflow:hidden;text-overflow:ellipsis;white-space:nowrap;color:#6b7280;font-size:12px">{{ svc.description }}</td>
              <td>
                <span class="tag" :class="stateTag(svc.active)">{{ svc.active }}</span>
                <span v-if="svc.sub" style="margin-left:4px;font-size:10px;color:#9ca3af">({{ svc.sub }})</span>
              </td>
              <td>
                <span v-if="svc.memory" class="mem-label">{{ svc.memory }}</span>
                <span v-else style="color:#d1d5db;font-size:11px">—</span>
              </td>
              <td>
                <span v-if="svc.cpu_time" style="font-size:11px;color:#7c3aed;font-family:monospace">{{ svc.cpu_time }}</span>
                <span v-else style="color:#d1d5db;font-size:11px">—</span>
              </td>
              <td>
                <span v-if="svc.unit_file_state" class="tag" :class="fileStateTag(svc.unit_file_state)" style="font-size:10px">{{ svc.unit_file_state }}</span>
                <span v-else style="color:#d1d5db;font-size:11px">—</span>
              </td>
              <td @click.stop>
                <div style="display:flex;gap:4px;flex-wrap:wrap">
                  <button class="btn btn-xs btn-cyan"  @click="action(svc,'start')"   v-if="svc.active!=='active'" title="启动">▶</button>
                  <button class="btn btn-xs btn-ghost" @click="action(svc,'stop')"    v-if="svc.active==='active'" title="停止">⏹</button>
                  <button class="btn btn-xs btn-ghost" @click="action(svc,'restart')" title="重启">↺</button>
                  <button class="btn btn-xs btn-ghost" @click="showLogs(svc)"         title="日志">📋</button>
                  <button class="btn btn-xs btn-ghost" @click="openEditor(svc)"       title="编辑配置文件">✏️</button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-if="filtered.length===0 && !loading" class="empty-state">
        <div class="empty-icon">🔧</div>
        <div class="empty-title">未发现系统服务</div>
        <div class="empty-sub">{{ errorMsg || '请确认 systemctl 可用且面板有执行权限' }}</div>
      </div>
      <div v-if="loading" class="empty-state">
        <div style="font-size:32px" class="animate-spin">⚙️</div>
        <div class="empty-sub" style="margin-top:12px">加载中...</div>
      </div>
    </div>

    <!-- 详情 modal -->
    <div class="modal-overlay" v-if="detailModal" @click.self="detailModal=null">
      <div class="modal" style="width:600px;max-width:96vw">
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:20px">
          <div>
            <h3 style="color:#1a1040;font-size:15px;font-weight:700">🔧 {{ detailModal.unit }}</h3>
            <p style="font-size:12px;color:#6b7280;margin-top:2px">{{ detailModal.description }}</p>
          </div>
          <button class="btn btn-ghost btn-sm" @click="detailModal=null">✕</button>
        </div>
        <div class="detail-grid">
          <div class="detail-item">
            <span class="detail-label">运行状态</span>
            <span class="tag" :class="stateTag(detailModal.active)">{{ detailModal.active }} ({{ detailModal.sub }})</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">文件状态</span>
            <span class="tag" :class="fileStateTag(detailModal.unit_file_state)">{{ detailModal.unit_file_state || '—' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">主进程 PID</span>
            <span class="detail-val mono">{{ detailModal.main_pid || '—' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">内存占用</span>
            <span class="detail-val" style="color:#06b6d4;font-weight:600">{{ detailModal.memory || '—' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">CPU 累计时间</span>
            <span class="detail-val" style="color:#7c3aed;font-weight:600">{{ detailModal.cpu_time || '—' }}</span>
          </div>
          <div class="detail-item">
            <span class="detail-label">任务数 (Tasks)</span>
            <span class="detail-val mono">{{ detailModal.tasks || '—' }}</span>
          </div>
          <div class="detail-item" style="grid-column:1/-1">
            <span class="detail-label">启动命令</span>
            <span class="detail-val mono" style="font-size:11px;word-break:break-all">{{ detailModal.exec_start || '—' }}</span>
          </div>
          <div class="detail-item" style="grid-column:1/-1">
            <span class="detail-label">Unit 文件路径</span>
            <span class="detail-val mono" style="font-size:11px;word-break:break-all">{{ detailModal.fragment_path || '—' }}</span>
          </div>
          <div class="detail-item" style="grid-column:1/-1">
            <span class="detail-label">启动时间</span>
            <span class="detail-val mono" style="font-size:11px">{{ detailModal.started_at || '—' }}</span>
          </div>
        </div>
        <div style="display:flex;gap:8px;margin-top:20px;flex-wrap:wrap">
          <button class="btn btn-cyan btn-sm"  @click="action(detailModal,'start');detailModal=null"   v-if="detailModal.active!=='active'">▶ 启动</button>
          <button class="btn btn-ghost btn-sm" @click="action(detailModal,'stop');detailModal=null"    v-if="detailModal.active==='active'">⏹ 停止</button>
          <button class="btn btn-ghost btn-sm" @click="action(detailModal,'restart');detailModal=null">↺ 重启</button>
          <button class="btn btn-ghost btn-sm" @click="showLogs(detailModal);detailModal=null">📋 查看日志</button>
          <button class="btn btn-ghost btn-sm" @click="openEditor(detailModal);detailModal=null">✏️ 编辑配置</button>
          <button class="btn btn-success btn-sm" @click="action(detailModal,'enable');detailModal=null"  v-if="detailModal.unit_file_state!=='enabled'">✓ 开机启动</button>
          <button class="btn btn-danger btn-sm"  @click="action(detailModal,'disable');detailModal=null" v-if="detailModal.unit_file_state==='enabled'">✗ 禁用启动</button>
        </div>
      </div>
    </div>

    <!-- 日志 modal -->
    <div class="modal-overlay" v-if="logModal" @click.self="logModal=null">
      <div class="modal" style="width:760px;max-width:96vw">
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:16px">
          <h3 style="color:#1a1040;font-size:15px">📋 {{ logModal.unit }} 日志</h3>
          <button class="btn btn-ghost btn-sm" @click="logModal=null">✕ 关闭</button>
        </div>
        <pre class="log-box">{{ logContent || '加载中...' }}</pre>
      </div>
    </div>

    <!-- 编辑 Service 文件 modal -->
    <div class="modal-overlay" v-if="editorModal" @click.self="closeEditor">
      <div class="modal" style="width:820px;max-width:96vw">
        <div style="display:flex;align-items:center;justify-content:space-between;margin-bottom:14px">
          <div>
            <h3 style="color:#1a1040;font-size:15px;font-weight:700">✏️ 编辑 {{ editorUnit }}</h3>
            <p style="font-size:11px;color:#9ca3af;margin-top:2px;font-family:monospace">{{ editorPath }}</p>
          </div>
          <div style="display:flex;gap:8px">
            <button class="btn btn-success btn-sm" @click="saveFile" :disabled="saving">{{ saving?'保存中...':'💾 保存并恢复' }}</button>
            <button class="btn btn-ghost btn-sm" @click="closeEditor">✕</button>
          </div>
        </div>
        <div v-if="editorErr" style="background:rgba(244,63,94,0.1);border:1px solid rgba(244,63,94,0.3);border-radius:8px;padding:10px;margin-bottom:12px;font-size:12px;color:#f43f5e">{{ editorErr }}</div>
        <div v-if="saveOk" style="background:rgba(16,185,129,0.1);border:1px solid rgba(16,185,129,0.3);border-radius:8px;padding:10px;margin-bottom:12px;font-size:12px;color:#10b981">✓ 保存成功！已执行 daemon-reload 并恢复原运行状态</div>
        <textarea
          class="editor-box"
          v-model="editorContent"
          spellcheck="false"
          placeholder="加载中..."
        ></textarea>
        <div style="font-size:11px;color:#9ca3af;margin-top:8px">⚠️ 保存后将自动执行 <code>systemctl daemon-reload</code>，若服务原本在运行中将自动重启恢复。</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import axios from 'axios'
import { useI18n } from '../stores/i18n.js'
const i18n = useI18n(); const t = k => i18n.t(k)

const services   = ref([])
const search     = ref('')
const filterState = ref('')
const sortBy     = ref('')
const sortDir    = ref('desc')
const logModal   = ref(null)
const logContent = ref('')
const detailModal = ref(null)
const loading    = ref(false)
const errorMsg   = ref('')

// Editor state
const editorModal   = ref(false)
const editorUnit    = ref('')
const editorPath    = ref('')
const editorContent = ref('')
const editorErr     = ref('')
const saveOk        = ref(false)
const saving        = ref(false)

const filtered = computed(() => services.value.filter(s => {
  const q = search.value.toLowerCase()
  const matchSearch = s.unit?.toLowerCase().includes(q) || s.description?.toLowerCase().includes(q)
  const matchState  = !filterState.value || s.active === filterState.value
  return matchSearch && matchState
}))

function stateTag(s) {
  return s === 'active'   ? 'tag-green'  :
         s === 'failed'   ? 'tag-red'    :
         s === 'inactive' ? 'tag-gray'   : 'tag-yellow'
}
function fileStateTag(s) {
  return s === 'enabled'  ? 'tag-green'  :
         s === 'disabled' ? 'tag-gray'   :
         s === 'masked'   ? 'tag-red'    : 'tag-yellow'
}

function setSort(s) { sortBy.value = s; load() }
function toggleDir() { sortDir.value = sortDir.value==='desc'?'asc':'desc'; load() }

async function load() {
  loading.value = true; errorMsg.value = ''
  try {
    const params = sortBy.value ? `?sort=${sortBy.value}&dir=${sortDir.value}` : ''
    const { data } = await axios.get('/api/services' + params)
    services.value = Array.isArray(data) ? data : []
    if (!services.value.length) errorMsg.value = 'No services returned.'
  } catch(e) {
    services.value = []
    errorMsg.value = e.response?.data?.error || e.message
  } finally {
    loading.value = false
  }
}

async function action(svc, act) {
  try {
    await axios.post(`/api/services/${svc.unit}/${act}`)
    setTimeout(load, 900)
  } catch(e) { alert(e.response?.data?.error || e.message) }
}

function showDetail(svc) { detailModal.value = svc }

async function showLogs(svc) {
  logModal.value = svc; logContent.value = ''
  const { data } = await axios.get(`/api/services/${svc.unit}/logs`)
  logContent.value = data.logs || ''
}

async function openEditor(svc) {
  editorUnit.value = svc.unit
  editorContent.value = ''
  editorPath.value = ''
  editorErr.value = ''
  saveOk.value = false
  editorModal.value = true
  try {
    const { data } = await axios.get(`/api/services/${svc.unit}/file`)
    editorContent.value = data.content || ''
    editorPath.value = data.path || ''
  } catch(e) {
    editorErr.value = '无法读取配置文件: ' + (e.response?.data?.error || e.message)
  }
}

async function saveFile() {
  saving.value = true; editorErr.value = ''; saveOk.value = false
  try {
    await axios.post(`/api/services/${editorUnit.value}/file`, { content: editorContent.value })
    saveOk.value = true
    setTimeout(load, 2000)
  } catch(e) {
    editorErr.value = '保存失败: ' + (e.response?.data?.error || e.message)
  } finally {
    saving.value = false
  }
}

function closeEditor() {
  editorModal.value = false
  editorErr.value = ''
  saveOk.value = false
}

onMounted(load)
</script>

<style scoped>
.toolbar { display:flex;align-items:center;gap:8px;flex-wrap:wrap; }
.inp { background:#f8faff;border:1.5px solid rgba(124,58,237,0.15);color:#1a1040;border-radius:8px;padding:7px 12px;font-size:13px;font-family:inherit;outline:none; }
.inp:focus,.sel:focus { border-color:#06b6d4; }
.sel { background:#f8faff;border:1.5px solid rgba(124,58,237,0.15);color:#1a1040;border-radius:8px;padding:7px 12px;font-size:13px;font-family:inherit;outline:none;cursor:pointer; }
.card { background:#fff;border:1px solid rgba(124,58,237,0.1);border-radius:14px;box-shadow:0 2px 12px rgba(6,182,212,0.06); }
.btn { display:inline-flex;align-items:center;gap:4px;padding:7px 14px;border-radius:8px;font-size:13px;font-weight:500;cursor:pointer;border:none;font-family:inherit;transition:all 0.2s; }
.btn:disabled { opacity:0.6;cursor:not-allowed; }
.btn-sm { padding:5px 11px;font-size:12px; }
.btn-xs { padding:3px 8px;font-size:11px;border-radius:6px;font-weight:600; }
.btn-primary { background:linear-gradient(135deg,#6366f1,#8b5cf6);color:#fff;box-shadow:0 2px 8px rgba(99,102,241,0.3); }
.btn-cyan    { background:linear-gradient(135deg,#06b6d4,#7c3aed);color:#fff;box-shadow:0 2px 6px rgba(6,182,212,0.3); }
.btn-ghost   { background:#fff;color:#6b7280;border:1px solid rgba(124,58,237,0.15); }
.btn-ghost:hover { background:rgba(6,182,212,0.05); }
.btn-success { background:rgba(16,185,129,0.1);color:#10b981;border:1px solid rgba(16,185,129,0.2); }
.btn-danger  { background:rgba(244,63,94,0.08);color:#f43f5e;border:1px solid rgba(244,63,94,0.2); }
.tag { display:inline-flex;align-items:center;padding:2px 8px;border-radius:100px;font-size:11px;font-weight:600; }
.tag-green  { background:rgba(16,185,129,0.1); color:#10b981; }
.tag-red    { background:rgba(244,63,94,0.1);  color:#f43f5e; }
.tag-yellow { background:rgba(245,158,11,0.1); color:#f59e0b; }
.tag-gray   { background:rgba(107,114,128,0.1);color:#6b7280; }
.mem-label { font-size:11px;font-weight:600;color:#06b6d4;font-family:monospace; }
.modal-overlay { position:fixed;inset:0;background:rgba(26,16,64,0.45);backdrop-filter:blur(6px);display:flex;align-items:center;justify-content:center;z-index:1000; }
.modal { background:#fff;border:1px solid rgba(124,58,237,0.15);border-radius:16px;padding:24px;box-shadow:0 20px 60px rgba(124,58,237,0.15);max-height:90vh;overflow-y:auto; }
.log-box { background:#0f172a;color:#e2e8f0;border-radius:10px;padding:16px;font-family:'JetBrains Mono',monospace;font-size:12px;max-height:55vh;overflow-y:auto;white-space:pre-wrap;word-break:break-all; }
.editor-box { width:100%;height:55vh;background:#0f172a;color:#e2e8f0;border:1px solid rgba(124,58,237,0.2);border-radius:10px;padding:14px;font-family:'JetBrains Mono',monospace;font-size:12px;line-height:1.6;resize:vertical;outline:none;box-sizing:border-box;white-space:pre;overflow-wrap:normal;overflow-x:auto; }
.editor-box:focus { border-color:#7c3aed; }
.detail-grid { display:grid;grid-template-columns:1fr 1fr;gap:12px; }
.detail-item { background:#f8faff;border:1px solid rgba(124,58,237,0.08);border-radius:10px;padding:10px 14px;display:flex;flex-direction:column;gap:4px; }
.detail-label { font-size:10px;font-weight:700;text-transform:uppercase;letter-spacing:0.06em;color:#9ca3af; }
.detail-val { font-size:13px;color:#1a1040;font-weight:500; }
.mono { font-family:'JetBrains Mono',monospace; }
.empty-state { text-align:center;padding:50px 20px; }
.empty-icon { font-size:40px;margin-bottom:12px; }
.empty-title { font-size:15px;font-weight:600;color:#1a1040;margin-bottom:6px; }
.empty-sub { font-size:12px;color:#9ca3af; }
.animate-spin { display:inline-block;animation:spin 0.8s linear infinite; }
@keyframes spin { to { transform:rotate(360deg); } }
code { font-family:monospace;background:rgba(99,102,241,0.08);padding:1px 5px;border-radius:4px;font-size:11px; }
</style>
