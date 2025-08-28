import type { TreeItem } from '@nuxt/ui'

function parseDN(dn: string){
  return dn.split(',').map(s => s.trim()).filter(Boolean)
}

const dnEq = (a: string, b: string) => a.toLowerCase() === b.toLowerCase()

function endsWithDN(dn: string, base: string){
  const a = parseDN(dn).toReversed()
  const b = parseDN(base).toReversed()
  if(b.length > a.length) return false

  for(const [i, element] of b.entries()) if(!a[i] || !element || !dnEq(a[i], element)) return false
  return true
}

function getLdapIcon(dn: string){
  if(dn.toLowerCase().startsWith('uid=')) return 'i-lucide-user'
  else if(dn.toLowerCase().startsWith('l=')) return 'i-lucide-building'
  else if(dn.toLowerCase().startsWith('cn=')) return 'i-lucide-component'
}

// eslint-disable-next-line typescript/no-explicit-any
export function buildLdapTree(baseDn: string, dns: string[], onSelect: (e: any)=> void): TreeItem[]{
  const root = {
    label: baseDn,
    fullDn: baseDn.toLowerCase(),
    icon: 'i-lucide-house',
    onSelect,
    children: [] as TreeItem[],
  }

  // eslint-disable-next-line typescript/no-explicit-any
  const index = new Map<string, any>()
  index.set(baseDn.toLowerCase(), root)

  for(const dn of dns){
    if(!endsWithDN(dn, baseDn)) continue

    const parts = parseDN(dn)
    const baseParts = parseDN(baseDn)

    const relative = parts.slice(0, parts.length - baseParts.length).toReversed()

    let parentDn = baseDn
    let parentNode = root

    for(const rdn of relative){
      const currentDn = `${rdn},${parentDn}`
      const key = currentDn.toLowerCase()

      let node = index.get(key)
      if(!node){
        node = { label: rdn, fullDn: key, icon: getLdapIcon(rdn), expandedIcon: 'i-lucide-folder-open', onSelect, children: [] }
        parentNode.children.push(node)
        index.set(key, node)
      }

      parentNode = node
      parentDn = currentDn
    }
  }

  return [root]
}
