import { defineStore } from 'pinia'
import { ref } from 'vue'
import axios from 'axios'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem('gp_token') || '')
  const username = ref(localStorage.getItem('gp_user') || 'admin')

  if (token.value) axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`

  function setToken(t, user) {
    token.value = t; username.value = user || 'admin'
    localStorage.setItem('gp_token', t)
    localStorage.setItem('gp_user', username.value)
    axios.defaults.headers.common['Authorization'] = `Bearer ${t}`
  }

  function logout() {
    token.value = ''; username.value = ''
    localStorage.removeItem('gp_token'); localStorage.removeItem('gp_user')
    delete axios.defaults.headers.common['Authorization']
  }

  return { token, username, setToken, logout }
})
