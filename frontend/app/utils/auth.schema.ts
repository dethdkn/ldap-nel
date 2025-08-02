import { z } from 'zod'

export const authSchema = z.object({
  username: z.string().min(1, 'Username is required'),
  password: z.string().min(1, 'Password is required'),
})

export type Auth = z.infer<typeof authSchema>
