<template>
  <div class="network-page">
    <div class="cards-row">
      <div class="stat-card" style="--c:#6366f1">
        <div class="sc-icon">🔗</div>
        <div class="sc-val">{{ network.connections }}</div>
        <div class="sc-lbl">{{ t('connections') }}</div>
      </div>
      <div class="stat-card" style="--c:#10b981">
        <div class="sc-icon">↑</div>
        <div class="sc-val">{{ fmtSpeed(totalUp) }}</div>
        <div class="sc-lbl">{{ t('upload') }}</div>
      </div>
      <div class="stat-card" style="--c:#06b6d4">
        <div class="sc-icon">↓</div>
        <div class="sc-val">{{ fmtSpeed(totalDown) }}</div>
        <div class="sc-lbl">{{ t('download') }}</div>
      </div>
      <div class="stat-card" style="--c:#f59e0b">
        <div class="sc-icon">📡</div>
        <div class="sc-val">{{ (network.interfaces||[]).length }}</div>
        <div class="sc-lbl">{{ t('interfaces') }}</div>
      </div>
    </div>

    <div class="iface-list">
      <div class="card iface-card" v-for="iface in network.interfaces" :key="iface.name">
        <div class="iface-head">
          <div class="iface-dot" :class="iface.speed_up||iface.speed_down?'active':'idle'"></div>
          <span class="iface-name">{{ iface.name }}</span>
          <div class="iface-addrs">
            <span class="tag tag-blue" v-for="addr in iface.addrs?.slice(0,2)" :key="addr">{{ addr }}</span>
          </div>
        </div>
        <div class="iface-stats">
          <div class="is">
            <div class="is-icon up">↑</div>
            <div>
              <div class="is-speed">{{ fmtSpeed(iface.speed_up) }}</div>
              <div class="is-total">{{ t('total_sent') }}: {{ fmt(iface.bytes_sent) }}</div>
            </div>
          </div>
          <div class="is">
            <div class="is-icon dn">↓</div>
            <div>
              <div class="is-speed">{{ fmtSpeed(iface.speed_down) }}</div>
              <div class="is-total">{{ t('total_recv') }}: {{ fmt(iface.bytes_recv) }}</div>
            </div>
          </div>
          <div class="is">
            <div class="is-icon pk">📦</div>
            <div>
              <div class="is-speed">{{ (iface.packets_sent||0) + (iface.packets_recv||0) }}</div>
              <div class="is-total">{{ t('packets') }}</div>
            </div>
          </div>
        </div>
        <div class="speed-bars">
          <div style="display:flex;align-items:center;gap:8px;margin-bottom:4px">
            <span style="font-size:11px;color:#10b981;width:16px">↑</span>
            <div class="bar" style="flex:1"><div class="bf" :style="`width:${Math.min((iface.speed_up||0)/1048576*20,100)}%;background:#10b981`"></div></div>
          </div>
          <div style="display:flex;align-items:center;gap:8px">
            <span style="font-size:11px;color:#6366f1;width:16px">↓</span>
            <div class="bar" style="flex:1"><div class="bf" :style="`width:${Math.min((iface.speed_down||0)/1048576*20,100)}%;background:#6366f1`"></div></div>
          </div>
        </div>
      </div>
    </div>

    <div class="card">
      <div class="ch"><span class="ct">{{ t('network_history') }}</span></div>
      <div ref="chartEl" style="height:220px"></div>
    </div>
  </div>
</template>
<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import axios from 'axios'
import * as echarts from 'echarts'
import { useI18n } from '../stores/i18n.js'
const i18n = useI18n(); const t = k => i18n.t(k)
const network = ref({ interfaces:[], connections:0, total_sent:0, total_recv:0 })
const chartEl = ref(null); let chart = null, wsh = null, speedHistory = ref([])
const totalUp = computed(() => (network.value.interfaces||[]).reduce((s,i)=>s+(i.speed_up||0),0))
const totalDown = computed(() => (network.value.interfaces||[]).reduce((s,i)=>s+(i.speed_down||0),0))
function fmt(b) { if(!b) return '0 B'; const u=['B','KB','MB','GB','TB'],i=Math.min(Math.floor(Math.log(b)/Math.log(1024)),4); return (b/Math.pow(1024,i)).toFixed(1)+' '+u[i] }
function fmtSpeed(b) { if(!b) return '0 B/s'; if(b<1024) return b+'B/s'; if(b<1048576) return (b/1024).toFixed(1)+'KB/s'; return (b/1048576).toFixed(1)+'MB/s' }
async function load() { const {data} = await axios.get('/api/network'); network.value = data }
function updateChart() {
  if (!chart) return
  const now = new Date(); const label = `${String(now.getHours()).padStart(2,'0')}:${String(now.getMinutes()).padStart(2,'0')}:${String(now.getSeconds()).padStart(2,'0')}`
  speedHistory.value.push({ t: label, up: totalUp.value/1024, down: totalDown.value/1024 })
  if (speedHistory.value.length > 60) speedHistory.value.shift()
  chart.setOption({ xAxis:{data:speedHistory.value.map(p=>p.t)}, series:[{data:speedHistory.value.map(p=>p.up.toFixed(1))},{data:speedHistory.value.map(p=>p.down.toFixed(1))}] })
}
function initChart() {
  chart = echarts.init(chartEl.value, null, {renderer:'svg'})
  chart.setOption({
    backgroundColor:'transparent',
    tooltip:{trigger:'axis',backgroundColor:'rgba(255,255,255,0.95)',borderColor:'rgba(99,102,241,0.2)',textStyle:{color:'#1e1b4b',fontSize:12},formatter:params=>`${params[0].axisValue}<br>↑ ${params[0].data} KB/s<br>↓ ${params[1].data} KB/s`},
    legend:{data:['Upload','Download'],right:0,textStyle:{color:'#6b7280',fontSize:12}},
    grid:{left:56,right:16,top:32,bottom:20},
    xAxis:{type:'category',data:[],axisLabel:{color:'#9ca3af',fontSize:10},axisLine:{lineStyle:{color:'rgba(99,102,241,0.1)'}},splitLine:{show:false}},
    yAxis:{type:'value',axisLabel:{color:'#9ca3af',fontSize:11,formatter:'{value}KB/s'},splitLine:{lineStyle:{color:'rgba(99,102,241,0.06)'}}},
    series:[
      {name:'Upload',type:'line',smooth:true,data:[],showSymbol:false,lineStyle:{color:'#10b981',width:2},areaStyle:{color:{type:'linear',x:0,y:0,x2:0,y2:1,colorStops:[{offset:0,color:'rgba(16,185,129,0.2)'},{offset:1,color:'rgba(16,185,129,0)'}]}}},
      {name:'Download',type:'line',smooth:true,data:[],showSymbol:false,lineStyle:{color:'#6366f1',width:2},areaStyle:{color:{type:'linear',x:0,y:0,x2:0,y2:1,colorStops:[{offset:0,color:'rgba(99,102,241,0.18)'},{offset:1,color:'rgba(99,102,241,0)'}]}}},
    ]
  })
}
onMounted(async () => {
  await load(); initChart()
  wsh = e => { const {event,data}=e.detail; if(event==='metrics') { network.value=data.network; updateChart() } }
  window.addEventListener('ws-msg', wsh)
})
onUnmounted(() => { window.removeEventListener('ws-msg', wsh); chart?.dispose() })
</script>
<style scoped>
.network-page { display:flex;flex-direction:column;gap:14px; }
.cards-row { display:grid;grid-template-columns:repeat(4,1fr);gap:14px; }
.stat-card { background:#fff;border:1px solid rgba(99,102,241,0.1);border-radius:14px;padding:20px;text-align:center;box-shadow:0 2px 12px rgba(99,102,241,0.06);position:relative;overflow:hidden;transition:transform 0.2s; }
.stat-card:hover { transform:translateY(-2px); }
.stat-card::after { content:'';position:absolute;bottom:0;left:0;right:0;height:3px;background:var(--c); }
.sc-icon { font-size:28px;margin-bottom:8px; }
.sc-val  { font-size:24px;font-weight:700;color:var(--c);font-family:monospace; }
.sc-lbl  { font-size:12px;color:#9ca3af;margin-top:4px; }
.iface-list { display:grid;grid-template-columns:repeat(auto-fill,minmax(300px,1fr));gap:14px; }
.card { background:#fff;border:1px solid rgba(99,102,241,0.1);border-radius:14px;padding:18px;box-shadow:0 2px 12px rgba(99,102,241,0.06); }
.ch { display:flex;align-items:center;justify-content:space-between;margin-bottom:14px; }
.ct { font-size:14px;font-weight:600;color:#1e1b4b; }
.iface-card { }
.iface-head { display:flex;align-items:center;gap:8px;flex-wrap:wrap;margin-bottom:14px; }
.iface-dot  { width:8px;height:8px;border-radius:50%;flex-shrink:0; }
.iface-dot.active { background:#10b981;box-shadow:0 0 8px rgba(16,185,129,0.5);animation:pulse 2s infinite; }
.iface-dot.idle   { background:#9ca3af; }
.iface-name { font-size:15px;font-weight:700;color:#1e1b4b;flex-shrink:0; }
.iface-addrs { display:flex;gap:4px;flex-wrap:wrap;flex:1; }
.iface-stats { display:grid;grid-template-columns:1fr 1fr 1fr;gap:8px;margin-bottom:12px; }
.is { display:flex;align-items:center;gap:8px; }
.is-icon { width:28px;height:28px;border-radius:8px;display:flex;align-items:center;justify-content:center;font-size:13px;font-weight:700;flex-shrink:0; }
.is-icon.up { background:rgba(16,185,129,0.1);color:#10b981; }
.is-icon.dn { background:rgba(99,102,241,0.1);color:#6366f1; }
.is-icon.pk { background:rgba(245,158,11,0.1); }
.is-speed { font-size:13px;font-weight:600;color:#1e1b4b;font-family:monospace; }
.is-total { font-size:11px;color:#9ca3af; }
.bar { height:5px;background:#f0f4ff;border-radius:3px;overflow:hidden; }
.bf  { height:100%;border-radius:3px;transition:width 0.5s; }
.speed-bars { margin-top:8px; }
@media (max-width:768px) { .cards-row { grid-template-columns:1fr 1fr; } .iface-list { grid-template-columns:1fr; } .iface-stats { grid-template-columns:1fr 1fr; } }
@keyframes pulse { 0%,100%{opacity:1}50%{opacity:.5} }
</style>
