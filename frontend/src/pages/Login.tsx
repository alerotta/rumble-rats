
import { Container, Box, Typography, TextField, Button, Alert } from "@mui/material";
import { useState } from "react";
import {login} from "../api/auth"

export default function LoginPage (){

    const [username, setUsername] = useState("")
    const [password, setPassword] = useState("")
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    const canSubmit =
        username.trim().length > 0 &&
        password.length >= 8 &&
        !loading

    async function onSubmit (e: React.SubmitEvent){
        e.preventDefault()
        setError(null)
        setLoading(true)

        try {
            const user = await login({ username, password })
            console.log("registered:", user)
            } catch (err) {
            setError(err instanceof Error ? err.message : "Registration failed")
            } finally {
            setLoading(false)
        }

    }        
    

    return (
        <Container>
            <Typography>
                Login
            </Typography>

                  {error && (<Alert severity="error" sx={{ mb: 2 }}>{error}</Alert> )}

            <Box component="form" onSubmit={onSubmit} sx={{ display: "grid", gap: 2 }} >

            <TextField
              label="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              autoComplete="username"
              required
              disabled={loading}
            />

            <TextField
              label="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              type="password"
              autoComplete="new-password"
              helperText="Min 8 characters"
              required
              disabled={loading}
            />

            <Button 
            type="submit"
            variant="contained"
            disabled= {!canSubmit} > {loading ? "..." : "Login"} </Button>

            </Box>
        </Container>
    )


}