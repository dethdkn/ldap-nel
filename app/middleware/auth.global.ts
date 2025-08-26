export default defineNuxtRouteMiddleware(async to => {
  const { setUserSession, clearUserSession } = useUserSession()

  const res = await useRequestFetch()<{ username: string, isAdmin: boolean }>('/server/check-session').catch(() => null)

  if(!res){
    clearUserSession()
    if(['/login', '/setup'].includes(to.path)) return
    return navigateTo('/login')
  }

  setUserSession({ username: res.username, admin: res.isAdmin })
  if(['/login', '/setup'].includes(to.path)) return navigateTo('/')
})
