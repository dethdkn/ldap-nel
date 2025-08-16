import { z } from 'zod'

export const userSchema = z.object({
  id: z.number().default(0),
  username: z.string().min(1, 'Username is required').trim(),
  password: z.string().min(1, 'Password is required').regex(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,}$/, 'Password must be at least 8 chars, include upper, lower, number, and symbol'),
  repeatPassword: z.string().min(1, 'Repeat password is required'),
  admin: z.boolean().default(false),
})
  .refine(data => data.password === data.repeatPassword, { message: 'Passwords must match', path: ['repeatPassword'] })

export type User = z.infer<typeof userSchema>

export const updateUserSchema = userSchema.extend({
  password: userSchema.shape.password.optional(),
  repeatPassword: userSchema.shape.repeatPassword.optional(),
})
  .refine(data => !data.password || data.password === data.repeatPassword, { message: 'Passwords must match', path: ['repeatPassword'] })
