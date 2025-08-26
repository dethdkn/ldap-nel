import { z } from 'zod'

const passwordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,}$/

export const userSchema = z.object({
  id: z.number().default(0),
  username: z.string().min(1, 'Username is required').trim(),
  password: z.string().min(1, 'Password is required').regex(passwordRegex, 'Password must be at least 8 chars, include upper, lower, number, and symbol'),
  repeatPassword: z.string().min(1, 'Repeat password is required'),
  admin: z.boolean().default(false),
})
  .refine(data => data.password === data.repeatPassword, { message: 'Passwords must match', path: ['repeatPassword'] })

export type User = z.infer<typeof userSchema>

export const updateUserSchema = z.object({
  id: z.number().default(0),
  username: z.string().min(1, 'Username is required').trim(),
  password: z.string().optional().refine(val => !val || passwordRegex.test(val), 'Password must be at least 8 chars, include upper, lower, number, and symbol'),
  repeatPassword: z.string().optional(),
  admin: z.boolean().default(false),
})
  .refine(data => !data.password || data.password === data.repeatPassword, { message: 'Passwords must match', path: ['repeatPassword'] })
