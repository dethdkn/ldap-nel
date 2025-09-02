import { z } from 'zod'

export const ldapSchema = z.object({
  id: z.number().default(0),
  name: z.string().min(1, 'Name is required').trim(),
  url: z.string().min(1, 'URL is required').trim().transform(val => val.replaceAll('ldap://', '').replaceAll('ldaps://', '')),
  port: z.number().min(1, 'Port is required').max(65535, 'Port must be between 1 and 65535'),
  ssl: z.boolean().default(false),
  base_dn: z.string().min(1, 'Base DN is required').trim(),
  bind_dn: z.string().optional(),
  bind_pass: z.string().optional(),
})

export type Ldap = z.infer<typeof ldapSchema>

export const attributeValueSchema = z.object({
  id: z.number().default(0),
  dn: z.string().min(1, 'DN is required').trim(),
  attribute: z.string().min(1, 'Attribute is required'),
  value: z.string().min(1, 'Value is required'),
})

export type AttributeValue = z.infer<typeof attributeValueSchema>

export const attributeNewValueSchema = attributeValueSchema.extend({
  newValue: z.string().min(1, 'New Value is required'),
})
  .refine(data => data.value !== data.newValue, { message: 'New Value must be different from old value', path: ['newValue'] })

export type AttributeNewValue = z.infer<typeof attributeNewValueSchema>

export const newDnSchema = z.object({
  id: z.number().default(0),
  dn: z.string().min(1, 'DN is required').trim().transform(val => val.endsWith(',') ? val.slice(0, -1) : val),
  attributes: z.array(z.object({
    attribute: z.string().min(1, 'Attribute is required'),
    value: z.string().min(1, 'Value is required'),
  })),
})

export type NewDn = z.infer<typeof newDnSchema>

export const CopyMoveSchema = z.object({
  id: z.number().default(0),
  dn: z.string().min(1, 'Source DN is required').trim(),
  targetDn: z.string().min(1, 'Target DN is required').trim(),
})

export type CopyMove = z.infer<typeof CopyMoveSchema>
