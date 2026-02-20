<template>
  <div class="layout" :class="{ 'sidebar-open': sidebarOpen }">
    <!-- Mobile overlay -->
    <div class="mobile-overlay" v-if="sidebarOpen" @click="sidebarOpen=false"></div>

    <!-- Sidebar -->
    <aside class="sidebar">
      <div class="sidebar-inner">
        <div class="sidebar-logo">
          <div class="logo-box">
            <svg width="38" height="38" viewBox="0 0 64 64" fill="none" xmlns="http://www.w3.org/2000/svg">
              <defs>
                <linearGradient id="sl-grad" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#06b6d4"/>
                  <stop offset="100%" stop-color="#7c3aed"/>
                </linearGradient>
              </defs>
              <rect width="64" height="64" rx="14" fill="url(#sl-grad)"/>
              <rect x="10" y="12" width="44" height="28" rx="4" fill="none" stroke="white" stroke-width="3"/>
              <polyline points="16,34 22,28 28,32 34,26 40,30 48,24" fill="none" stroke="#22d3ee" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
              <rect x="22" y="44" width="20" height="3" rx="1.5" fill="white" opacity="0.7"/>
              <rect x="28" y="40" width="8" height="4" rx="1" fill="white" opacity="0.5"/>
            </svg>
          </div>
          <div class="logo-text">
            <span class="logo-name">GoPanel</span>
            <span class="logo-v">v1.0</span>
          </div>
          <button class="close-btn" @click="sidebarOpen=false">✕</button>
        </div>

        <nav class="nav">
          <div class="nav-section" v-for="section in navSections" :key="section.label">
            <div class="nav-label">{{ section.label }}</div>
            <router-link
              v-for="item in section.items" :key="item.path"
              :to="item.path" class="nav-item"
              :class="{ active: $route.path === item.path }"
              @click="sidebarOpen=false"
            >
              <span class="nav-icon" :style="`background:${item.color}`">{{ item.emoji }}</span>
              <span class="nav-name">{{ i18n.t(item.key) }}</span>
            </router-link>
          </div>
        </nav>

        <div class="sidebar-bottom">
          <div class="ws-pill" :class="wsConnected ? 'on' : 'off'">
            <span class="ws-dot"></span>
            {{ wsConnected ? i18n.t('realtime') : i18n.t('disconnected') }}
          </div>
          <div class="user-row">
            <div class="user-avatar">{{ auth.username?.[0]?.toUpperCase() }}</div>
            <span class="user-name">{{ auth.username }}</span>
            <button class="logout-btn" @click="logout" title="Logout">⏻</button>
          </div>
        </div>
      </div>
    </aside>

    <!-- Main -->
    <div class="main-wrap">
      <header class="topbar">
        <button class="hamburger" @click="sidebarOpen=true">☰</button>
        <div class="topbar-title">
          <span class="grad-text">{{ i18n.t(currentKey) }}</span>
        </div>
        <div class="topbar-actions">
          <button class="icon-btn" @click="i18n.toggle()" title="Switch Language">🌐</button>
          <div class="clock">{{ currentTime }}</div>
        </div>
      </header>
      <main class="content">
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth.js'
import { useI18n } from '../stores/i18n.js'
import dayjs from 'dayjs'

const route = useRoute(), router = useRouter()
const auth = useAuthStore(), i18n = useI18n()
const wsConnected = ref(false), currentTime = ref(''), sidebarOpen = ref(false)
let ws = null, timer = null

const navSections = computed(() => [
  {
    label: i18n.locale.value === 'zh' ? '监控' : 'Monitor',
    items: [
      { path:'/dashboard', key:'dashboard', emoji:'📊', color:'linear-gradient(135deg,#06b6d4,#7c3aed)' },
    ]
  },
  {
    label: i18n.locale.value === 'zh' ? '管理' : 'Manage',
    items: [
      { path:'/processes', key:'processes', emoji:'⚙️', color:'linear-gradient(135deg,#f59e0b,#f97316)' },
      { path:'/docker',    key:'docker',    emoji:'🐳', color:'linear-gradient(135deg,#06b6d4,#10b981)' },
      { path:'/services',  key:'services',  emoji:'🔧', color:'linear-gradient(135deg,#10b981,#06b6d4)' },
    ]
  },
  {
    label: i18n.locale.value === 'zh' ? '系统' : 'System',
    items: [
      { path:'/settings',  key:'settings',  emoji:'⚡', color:'linear-gradient(135deg,#f43f5e,#ec4899)' },
    ]
  },
])

const routeKeyMap = {'/dashboard':'dashboard','/processes':'processes','/docker':'docker','/services':'services','/settings':'settings'}
const currentKey = computed(() => routeKeyMap[route.path] || 'dashboard')

function connectWS() {
  const proto = location.protocol === 'https:' ? 'wss' : 'ws'
  ws = new WebSocket(`${proto}://${location.host}/api/ws`)
  ws.onopen  = () => { wsConnected.value = true }
  ws.onclose = () => { wsConnected.value = false; setTimeout(connectWS, 3000) }
  ws.onmessage = (e) => window.dispatchEvent(new CustomEvent('ws-msg', { detail: JSON.parse(e.data) }))
}

function logout() { auth.logout(); router.push('/login') }

onMounted(() => {
  connectWS()
  timer = setInterval(() => { currentTime.value = dayjs().format('HH:mm:ss') }, 1000)
  currentTime.value = dayjs().format('HH:mm:ss')
})
onUnmounted(() => { ws?.close(); clearInterval(timer) })
</script>

<style scoped>
.layout { display:flex; height:100vh; overflow:hidden; }

.sidebar {
  width: 220px; flex-shrink: 0;
  background: #fff;
  border-right: 1px solid rgba(6,182,212,0.1);
  box-shadow: 4px 0 24px rgba(6,182,212,0.06);
  display: flex; flex-direction: column;
  transition: transform 0.25s ease;
  z-index: 100;
}
.sidebar-inner { display:flex;flex-direction:column;height:100%;overflow:hidden; }
.sidebar-logo {
  display:flex;align-items:center;gap:10px;
  padding:18px 16px 14px;
  border-bottom:1px solid rgba(124,58,237,0.08);
}
.logo-box { width:38px;height:38px;border-radius:10px;overflow:hidden;flex-shrink:0;box-shadow:0 3px 10px rgba(6,182,212,0.3); }
.logo-name { font-size:15px;font-weight:700;color:#1e1b4b; }
.logo-v { font-size:10px;color:#7c3aed;background:rgba(6,182,212,0.1);padding:1px 6px;border-radius:4px;margin-left:4px; }
.logo-text { display:flex;align-items:center;flex:1; }
.close-btn { display:none;background:none;border:none;color:#6b7280;cursor:pointer;font-size:16px;padding:4px; }

.nav { flex:1;overflow-y:auto;padding:12px 10px; }
.nav-section { margin-bottom:16px; }
.nav-label { font-size:10px;font-weight:700;text-transform:uppercase;letter-spacing:0.08em;color:#9ca3af;padding:0 8px;margin-bottom:4px; }
.nav-item {
  display:flex;align-items:center;gap:10px;
  padding:9px 10px;border-radius:10px;
  font-size:13px;font-weight:500;color:#4b5563;
  text-decoration:none;cursor:pointer;margin-bottom:2px;
  transition:all 0.15s;
}
.nav-item:hover { background:rgba(6,182,212,0.06);color:#1e1b4b; }
.nav-item.active {
  background:linear-gradient(135deg,rgba(6,182,212,0.12),rgba(124,58,237,0.08));
  color:#0891b2;font-weight:600;
  border:1px solid rgba(124,58,237,0.15);
}
.nav-icon { width:28px;height:28px;border-radius:7px;display:flex;align-items:center;justify-content:center;font-size:14px;flex-shrink:0;box-shadow:0 2px 6px rgba(0,0,0,0.1); }
.nav-name { flex:1; }

.sidebar-bottom { padding:12px 10px;border-top:1px solid rgba(124,58,237,0.08); }
.ws-pill { display:flex;align-items:center;gap:6px;font-size:11px;font-weight:500;padding:5px 10px;border-radius:100px;margin-bottom:8px;width:fit-content; }
.ws-pill.on  { background:rgba(16,185,129,0.1);color:#10b981; }
.ws-pill.off { background:rgba(244,63,94,0.1); color:#f43f5e; }
.ws-dot { width:6px;height:6px;border-radius:50%;background:currentColor; }
.ws-pill.on .ws-dot { animation:pulse 2s infinite; }

.user-row { display:flex;align-items:center;gap:8px;padding:4px 2px; }
.user-avatar { width:28px;height:28px;border-radius:50%;background:linear-gradient(135deg,#06b6d4,#7c3aed);color:white;display:flex;align-items:center;justify-content:center;font-size:12px;font-weight:700;flex-shrink:0; }
.user-name { font-size:13px;font-weight:500;color:#1e1b4b;flex:1;overflow:hidden;text-overflow:ellipsis; }
.logout-btn { background:none;border:none;color:#9ca3af;cursor:pointer;font-size:16px;padding:2px;transition:color 0.2s; }
.logout-btn:hover { color:#f43f5e; }

/* Main */
.main-wrap { flex:1;display:flex;flex-direction:column;overflow:hidden;min-width:0; }
.topbar {
  display:flex;align-items:center;gap:12px;
  padding:14px 24px;
  background:#fff;border-bottom:1px solid rgba(124,58,237,0.08);
  box-shadow:0 2px 12px rgba(6,182,212,0.06);
  flex-shrink:0;
}
.hamburger { display:none;background:none;border:1px solid var(--border);border-radius:8px;padding:6px 10px;cursor:pointer;font-size:16px;color:var(--text3); }
.topbar-title { font-size:18px;font-weight:700;flex:1; }
.topbar-actions { display:flex;align-items:center;gap:10px; }
.icon-btn { background:rgba(124,58,237,0.08);border:1px solid rgba(124,58,237,0.15);border-radius:8px;width:34px;height:34px;cursor:pointer;font-size:16px;display:flex;align-items:center;justify-content:center;transition:all 0.2s; }
.icon-btn:hover { background:rgba(124,58,237,0.15); }
.clock { font-family:'JetBrains Mono',monospace;font-size:13px;color:#7c3aed;background:rgba(124,58,237,0.08);border:1px solid rgba(124,58,237,0.15);padding:6px 12px;border-radius:8px;font-weight:500; }

.content { flex:1;overflow-y:auto;padding:20px 24px; }

.mobile-overlay { display:none;position:fixed;inset:0;background:rgba(30,27,75,0.4);z-index:99; }

/* Mobile */
@media (max-width: 768px) {
  .sidebar { position:fixed;top:0;left:0;height:100vh;transform:translateX(-100%);z-index:100; }
  .layout.sidebar-open .sidebar { transform:translateX(0); }
  .layout.sidebar-open .mobile-overlay { display:block; }
  .hamburger { display:flex; }
  .close-btn { display:block; }
  .content { padding:14px 16px; }
  .topbar { padding:12px 16px; }
}
</style>
