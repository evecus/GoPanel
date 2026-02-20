<template>
  <div class="login-page">
    <div class="bg-orbs">
      <div class="orb orb1"></div>
      <div class="orb orb2"></div>
      <div class="orb orb3"></div>
    </div>
    <div class="login-card">
      <div class="logo-row">
        <div class="logo-icon">
          <svg width="26" height="26" viewBox="0 0 64 64" fill="none">
            <rect width="64" height="64" rx="14" fill="url(#lg)"/>
            <defs><linearGradient id="lg" x1="0%" y1="0%" x2="100%" y2="100%"><stop offset="0%" stop-color="#6366f1"/><stop offset="100%" stop-color="#06b6d4"/></linearGradient></defs>
            <rect x="10" y="12" width="44" height="28" rx="4" fill="none" stroke="white" stroke-width="3"/>
            <polyline points="16,34 22,28 28,32 34,26 40,30 48,24" fill="none" stroke="#22d3ee" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </div>
        <div>
          <div class="logo-name">GoPanel</div>
          <div class="logo-sub">Server Monitor</div>
        </div>
        <button class="lang-btn" @click="i18n.toggle()">{{ i18n.t('lang') }}</button>
      </div>

      <h2 class="login-title">{{ i18n.locale === 'zh' ? '欢迎回来' : 'Welcome Back' }}</h2>

      <form @submit.prevent="login" class="login-form">
        <div class="field">
          <label>{{ i18n.t('username') }}</label>
          <div class="input-wrap">
            <span class="input-icon">👤</span>
            <input class="input" v-model="form.username" autocomplete="username" required />
          </div>
        </div>
        <div class="field">
          <label>{{ i18n.locale === 'zh' ? '密码' : 'Password' }}</label>
          <div class="input-wrap">
            <span class="input-icon">🔑</span>
            <input class="input" type="password" v-model="form.password" autocomplete="current-password" required />
          </div>
        </div>
        <div class="error-msg" v-if="error">{{ error }}</div>
        <button class="btn btn-primary login-btn" type="submit" :disabled="loading">
          <span v-if="loading" class="animate-spin" style="display:inline-block">⟳</span>
          {{ loading ? (i18n.locale === 'zh' ? '登录中...' : 'Signing in...') : (i18n.locale === 'zh' ? '登录' : 'Sign In') }}
        </button>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import axios from 'axios'
import { useAuthStore } from '../stores/auth.js'
import { useI18n } from '../stores/i18n.js'

const router = useRouter()
const auth = useAuthStore()
const i18n = useI18n()
const loading = ref(false)
const error = ref('')
const form = ref({ username: 'admin', password: '' })

async function login() {
  error.value = ''
  loading.value = true
  try {
    const { data } = await axios.post('/api/login', form.value)
    auth.setToken(data.token, data.username)
    router.push('/dashboard')
  } catch {
    error.value = i18n.locale.value === 'zh' ? '用户名或密码错误' : 'Invalid credentials'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(135deg, #eef2ff 0%, #e0f7fa 50%, #f0fdf4 100%);
  overflow: hidden; position: relative;
}
.bg-orbs { position: absolute; inset: 0; pointer-events: none; }
.orb { position: absolute; border-radius: 50%; filter: blur(60px); opacity: 0.45; animation: float 8s ease-in-out infinite alternate; }
.orb1 { width:320px;height:320px;background:#c7d2fe;top:-80px;left:-80px;animation-delay:0s; }
.orb2 { width:280px;height:280px;background:#a5f3fc;bottom:-60px;right:-60px;animation-delay:-3s; }
.orb3 { width:200px;height:200px;background:#d9f99d;top:40%;left:60%;animation-delay:-5s; }
@keyframes float { from{transform:translate(0,0) scale(1)} to{transform:translate(20px,15px) scale(1.08)} }

.login-card {
  position: relative; z-index: 1;
  background: rgba(255,255,255,0.88); backdrop-filter: blur(24px);
  border: 1px solid rgba(99,102,241,0.18); border-radius: 20px;
  padding: 40px; width: 400px;
  box-shadow: 0 20px 60px rgba(99,102,241,0.15), 0 0 0 1px rgba(255,255,255,0.6);
}
.logo-row { display:flex;align-items:center;gap:12px;margin-bottom:28px; }
.logo-icon { width:48px;height:48px;border-radius:12px;overflow:hidden;flex-shrink:0;box-shadow:0 4px 14px rgba(99,102,241,0.35); }
.logo-name { font-size:20px;font-weight:700;color:#1e1b4b; }
.logo-sub  { font-size:12px;color:#6b7280;margin-top:1px; }
.lang-btn { margin-left:auto;background:rgba(99,102,241,0.08);border:1px solid rgba(99,102,241,0.2);color:#6366f1;padding:5px 12px;border-radius:100px;font-size:12px;font-weight:600;cursor:pointer;transition:all 0.2s; }
.lang-btn:hover { background:rgba(99,102,241,0.15); }

.login-title { font-size:22px;font-weight:700;color:#1e1b4b;margin-bottom:24px; }
.login-form { display:flex;flex-direction:column;gap:16px; }
.field { display:flex;flex-direction:column;gap:6px; }
.field label { font-size:13px;font-weight:500;color:#4f46e5; }
.input-wrap { position:relative; }
.input-icon { position:absolute;left:12px;top:50%;transform:translateY(-50%);font-size:15px;pointer-events:none; }
.input-wrap .input { padding-left:38px; }
.error-msg { background:rgba(244,63,94,0.08);border:1px solid rgba(244,63,94,0.2);color:#f43f5e;padding:10px 14px;border-radius:8px;font-size:13px; }
.login-btn { width:100%;justify-content:center;padding:12px;font-size:14px;margin-top:4px; }

@media (max-width: 480px) {
  .login-card { width: calc(100vw - 32px); padding: 28px 20px; }
}
</style>
