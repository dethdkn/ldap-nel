export default async () => {
  const { data: avaiableLdaps, refresh } = await useFetch<string[]>('/server/saved-ldaps')

  const selectedLdap = ref<string>()

  async function refreshLdaps(){
    await refresh()
    selectedLdap.value = undefined
  }

  return { avaiableLdaps, refreshLdaps, selectedLdap }
}
