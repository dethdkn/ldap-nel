export default createGlobalState(() => {
  const user = ref({ username: '', admin: false })

  const isLoggedIn = computed(() => !!user.value.username)

  const clearSession = async () => {}

  return { user, isLoggedIn, clearSession }
})
