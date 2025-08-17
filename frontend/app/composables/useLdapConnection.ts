export default createGlobalState(async () => {
  const { data: avaiableLdaps, refresh } = await useFetch<{ id: number, name: string }[]>('/server/saved-ldaps')

  const selectedLdap = ref<number>()

  async function refreshLdaps(){
    await refresh()
    selectedLdap.value = undefined
  }

  return { avaiableLdaps, refreshLdaps, selectedLdap }
})
