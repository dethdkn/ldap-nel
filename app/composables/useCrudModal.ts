export default createGlobalState(() => {
  const toast = useToast()

  const addModal = ref(false)
  const updateModal = ref(false)
  const deleteModal = ref(false)

  const addState = ref<AttributeValue>({ id: 0, dn: '', attribute: '', value: '' })
  const updateState = ref<AttributeNewValue>({ id: 0, dn: '', attribute: '', value: '', newValue: '' })
  const deleteState = ref<AttributeValue>({ id: 0, dn: '', attribute: '', value: '' })

  function openAddModal(ldap: number, dn: string){
    if(!dn) return toast.add({ title: 'Please select a DN first', icon: 'i-lucide-shield-alert', color: 'error' })
    addState.value = { id: ldap, dn, attribute: '', value: '' }
    addModal.value = true
  }
  function openUpdateModal(ldap: number, dn: string, attr: string, val: string){
    updateState.value = { id: ldap, dn, attribute: attr, value: val, newValue: val }
    updateModal.value = true
  }
  function openDeleteModal(ldap: number, dn: string, attr: string, val: string){
    deleteState.value = { id: ldap, dn, attribute: attr, value: val }
    deleteModal.value = true
  }

  whenever(() => !addModal.value, () => addState.value = { id: 0, dn: '', attribute: '', value: '' })
  whenever(() => !updateModal.value, () => updateState.value = { id: 0, dn: '', attribute: '', value: '', newValue: '' })
  whenever(() => !deleteModal.value, () => deleteState.value = { id: 0, dn: '', attribute: '', value: '' })

  return { addModal, addState, openAddModal, updateModal, updateState, openUpdateModal, deleteModal, deleteState, openDeleteModal }
})
