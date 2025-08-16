export default createGlobalState(() => {
  const userInfo = ref({ username: '', admin: false })
  const sessionToken = useCookie('session')

  const isLoggedIn = computed(() => !!userInfo.value.username)
  const user = computed(() => ({
    username: userInfo.value.username,
    admin: userInfo.value.admin,
    token: sessionToken.value,
  }))

  const setUserSession = ({ username, admin, token }: { username: string, admin: boolean, token?: string }) => {
    userInfo.value.username = username
    userInfo.value.admin = admin

    if(token) sessionToken.value = token
  }

  const clearUserSession = () => {
    userInfo.value = { username: '', admin: false }
    sessionToken.value = undefined
  }

  return { user, isLoggedIn, setUserSession, clearUserSession }
})
