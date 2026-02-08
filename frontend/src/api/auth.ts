export type RegisterRequest = {
    username : string
    email: string 
    password: string
}

export type RegisterResponse = {
  id: string
  name: string
  email: string
  // later you might add: accessToken?: string, refreshToken?: string
}

export async function register(payload: RegisterRequest): Promise <RegisterResponse>{
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
    return res.json()
}