/* eslint-disable no-await-in-loop */
// ! this function is pretty shitty 💩
// ! but it works for now... I need to find a way to change the selected attribute from nuxt ui tree
export default async function handleSearch(treeWrapper: Ref<HTMLElement | undefined>, baseDn: string, fullDn: string){
  const w = treeWrapper.value
  if(!w || !fullDn) return

  const sleep = (ms: number) => new Promise(resolve => setTimeout(resolve, ms))

  const findBtn = (label: string): HTMLButtonElement | null => {
    const buttons = w.querySelectorAll<HTMLButtonElement>('button[role="treeitem"]')
    for(const btn of buttons){
      const span = btn.querySelector<HTMLSpanElement>('span.truncate')
      if(span?.textContent?.trim() === label) return btn
    }
    return null
  }

  const waitExpanded = async (btn: HTMLButtonElement, { timeout = 3000, interval = 50 } = {}) => {
    const start = Date.now()
    while(btn.getAttribute('aria-expanded') !== 'true'){
      if(Date.now() - start > timeout) break
      await sleep(interval)
    }
  }

  const pathParts = fullDn.replaceAll(baseDn ? `,${baseDn}` : '', '').split(',').map(s => s.trim()).filter(Boolean).toReversed()

  if(baseDn){
    const baseBtn = findBtn(baseDn)
    if(baseBtn) if(baseBtn.getAttribute('aria-expanded') !== 'true'){
      baseBtn.click()
      await waitExpanded(baseBtn)
    }
  }

  for(let i = 0; i < pathParts.length; i++){
    const part = pathParts[i]
    const isLast = i === pathParts.length - 1

    let btn = findBtn(part || '')
    if(!btn){
      await sleep(50)
      btn = findBtn(part || '')
      if(!btn) continue
    }

    if(isLast){
      btn.click()
      continue
    }

    if(btn.getAttribute('aria-expanded') !== 'true'){
      btn.click()
      await waitExpanded(btn)
    }
  }
}
