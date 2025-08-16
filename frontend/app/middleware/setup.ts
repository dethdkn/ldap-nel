export default defineNuxtRouteMiddleware(async to => {
  const { empty } = await useRequestFetch()<{ empty: boolean }>('/server/users-empty')

  if(empty && to.path !== '/setup') return navigateTo('/setup')

  if(!empty && to.path === '/setup') return navigateTo('/login')
})
