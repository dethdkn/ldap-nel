export default createGlobalState(() => {
  const updateModal = ref(false)
  const deleteModal = ref(false)

  const deleteState = ref<AttributeValue>({ id: 0, dn: '', attribute: '', value: '' })

  function openDeleteModal(ldap: number, dn: string, attr: string, val: string){
    deleteState.value = { id: ldap, dn, attribute: attr, value: val }
    deleteModal.value = true
  }

  whenever(() => !deleteModal.value, () => deleteState.value = { id: 0, dn: '', attribute: '', value: '' })

  return { updateModal, deleteModal, deleteState, openDeleteModal }
})
