export default createGlobalState(() => {
  const updateModal = ref(false)
  const deleteModal = ref(false)

  const updateState = ref<AttributeNewValue>({ id: 0, dn: '', attribute: '', value: '', newValue: '' })
  const deleteState = ref<AttributeValue>({ id: 0, dn: '', attribute: '', value: '' })

  function openUpdateModal(ldap: number, dn: string, attr: string, val: string){
    updateState.value = { id: ldap, dn, attribute: attr, value: val, newValue: val }
    updateModal.value = true
  }
  function openDeleteModal(ldap: number, dn: string, attr: string, val: string){
    deleteState.value = { id: ldap, dn, attribute: attr, value: val }
    deleteModal.value = true
  }

  whenever(() => !updateModal.value, () => updateState.value = { id: 0, dn: '', attribute: '', value: '', newValue: '' })
  whenever(() => !deleteModal.value, () => deleteState.value = { id: 0, dn: '', attribute: '', value: '' })

  return { updateModal, updateState, openUpdateModal, deleteModal, deleteState, openDeleteModal }
})
