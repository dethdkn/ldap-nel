import { z } from 'zod'

export const userSchema = z.object({
  username: z.string().min(1, 'Username is required'),
  password: z.string().min(1, 'Password is required'),
  repeatPassword: z.string().min(1, 'Repeat password is required'),
  admin: z.boolean().default(false),
})
  .refine(data => data.password === data.repeatPassword, { message: 'Passwords must match', path: ['repeatPassword'] })

export type User = z.infer<typeof userSchema>
