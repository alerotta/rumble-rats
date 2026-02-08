
import { Container, Box, Typography, TextField, Button, Alert } from "@mui/material";
import { useState } from "react";
import {register} from "../api/auth"

export default function RegisterPage (){

    const [username, setUsername] = useState("")
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [loading, setLoading] = useState(false)
    const [error, setError] = useState<string | null>(null)

    const canSubmit =
        username.trim().length > 0 &&
        email.trim().length > 0 &&
        password.length >= 8 &&
        !loading

    async function onSubmit (e: React.SubmitEvent){
        e.preventDefault()
        setError(null)
        setLoading(true)

        try {
            const user = await register({ username, email, password })
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
                Create Account
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
              label="Email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              autoComplete="email"
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
            disabled= {!canSubmit} > {loading ? "Creating..." : "Create User"} </Button>

            </Box>
        </Container>
    )


}