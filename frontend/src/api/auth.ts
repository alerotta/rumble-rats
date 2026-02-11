import { setSession } from "./session"
import type { RegisterRequest, LoginRequest, AuthResponse } from "./types"


export async function register(payload: RegisterRequest): Promise <AuthResponse>{
    const res = await fetch ("/api/auth/register" , {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
    })

    if (!res.ok){
        let message = "registration failed"
        try {
            const data = await res.json()
            message = data?.message || data?.error || message
        }
        catch{
            //ignore
        }
        throw new Error(message)
    }
    const data = (await res.json()) as AuthResponse
    setSession(data)
    return data
}

export async function login(payload: LoginRequest): Promise <AuthResponse>{
    const res = await fetch ("/api/auth/login" , {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
    })

    if (!res.ok){
        let message = "login failed"
        try {
            const data = await res.json()
            message = data?.message || data?.error || message
        }
        catch{
            //ignore
        }
        throw new Error(message)
    }
    const data = (await res.json()) as AuthResponse
    setSession(data)
    return data
}

