export default createGlobalState(() => {
  const toast = useToast()

  const addModal = ref(false)
  const updateModal = ref(false)
  const jpegPhotoModal = ref(false)
  const passwordModal = ref(false)
  const deleteModal = ref(false)
  const copyModal = ref(false)
  const moveModal = ref(false)
  const deleteDnModal = ref(false)

  const addState = ref<AttributeValue>({ id: 0, dn: '', attribute: '', value: '' })
  const updateState = ref<AttributeNewValue>({ id: 0, dn: '', attribute: '', value: '', newValue: '' })
  const jpegPhotoState = ref<AttributeNewValue>({ id: 0, dn: '', attribute: '', value: '', newValue: '' })
  const passwordState = ref<AttributeNewValue>({ id: 0, dn: '', attribute: '', value: '', newValue: '' })
  const deleteState = ref<AttributeValue>({ id: 0, dn: '', attribute: '', value: '' })
  const copyState = ref<CopyMove>({ id: 0, dn: '', targetDn: '' })
  const moveState = ref<CopyMove>({ id: 0, dn: '', targetDn: '' })
  const deleteDnState = ref({ id: 0, dn: '', smart: false })

  function openAddModal(ldap: number, dn: string){
    if(!dn) return toast.add({ title: 'Please select a DN first', icon: 'i-lucide-shield-alert', color: 'error' })
    addState.value = { id: ldap, dn, attribute: '', value: '' }
    addModal.value = true
  }
  function openUpdateModal(ldap: number, dn: string, attr: string, val: string){
    if(!dn) return toast.add({ title: 'Please select a DN first', icon: 'i-lucide-shield-alert', color: 'error' })
    if(attr === 'jpegPhoto') return openJpegPhotoModal(ldap, dn, attr, val)
    if(attr === 'userPassword') return openPasswordModal(ldap, dn, attr, val)
    updateState.value = { id: ldap, dn, attribute: attr, value: val, newValue: val }
    updateModal.value = true
  }
  function openJpegPhotoModal(ldap: number, dn: string, attr: string, val: string){
    if(!dn) return toast.add({ title: 'Please select a DN first', icon: 'i-lucide-shield-alert', color: 'error' })
    jpegPhotoState.value = { id: ldap, dn, attribute: attr, value: val, newValue: val }
    jpegPhotoModal.value = true
  }
  function openPasswordModal(ldap: number, dn: string, attr: string, val: string){
    if(!dn) return toast.add({ title: 'Please select a DN first', icon: 'i-lucide-shield-alert', color: 'error' })
    passwordState.value = { id: ldap, dn, attribute: attr, value: val, newValue: val }
    passwordModal.value = true
  }
  function openDeleteModal(ldap: number, dn: string, attr: string, val: string){
    if(!dn) return toast.add({ title: 'Please select a DN first', icon: 'i-lucide-shield-alert', color: 'error' })
    deleteState.value = { id: ldap, dn, attribute: attr, value: val }
    deleteModal.value = true
  }
  function openCopyModal(ldap: number, dn: string){
    if(!dn) return toast.add({ title: 'Please select a DN first', icon: 'i-lucide-shield-alert', color: 'error' })
    copyState.value = { id: ldap, dn, targetDn: '' }
    copyModal.value = true
  }
  function openMoveModal(ldap: number, dn: string){
    if(!dn) return toast.add({ title: 'Please select a DN first', icon: 'i-lucide-shield-alert', color: 'error' })
    moveState.value = { id: ldap, dn, targetDn: '' }
    moveModal.value = true
  }
  function openDeleteDnModal(ldap: number, dn: string){
    if(!dn) return toast.add({ title: 'Please select a DN first', icon: 'i-lucide-shield-alert', color: 'error' })
    deleteDnState.value = { id: ldap, dn, smart: false }
    deleteDnModal.value = true
  }

  whenever(() => !addModal.value, () => addState.value = { id: 0, dn: '', attribute: '', value: '' })
  whenever(() => !updateModal.value, () => updateState.value = { id: 0, dn: '', attribute: '', value: '', newValue: '' })
  whenever(() => !jpegPhotoModal.value, () => jpegPhotoState.value = { id: 0, dn: '', attribute: '', value: '', newValue: '' })
  whenever(() => !passwordModal.value, () => passwordState.value = { id: 0, dn: '', attribute: '', value: '', newValue: '' })
  whenever(() => !deleteModal.value, () => deleteState.value = { id: 0, dn: '', attribute: '', value: '' })
  whenever(() => !copyModal.value, () => copyState.value = { id: 0, dn: '', targetDn: '' })
  whenever(() => !moveModal.value, () => moveState.value = { id: 0, dn: '', targetDn: '' })
  whenever(() => !deleteDnModal.value, () => deleteDnState.value = { id: 0, dn: '', smart: false })

  return { addModal, addState, openAddModal, updateModal, updateState, openUpdateModal, jpegPhotoModal, jpegPhotoState, passwordModal, passwordState, deleteModal, deleteState, openDeleteModal, copyModal, copyState, openCopyModal, moveModal, moveState, openMoveModal, deleteDnModal, deleteDnState, openDeleteDnModal }
})
