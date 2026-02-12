import { Typography, Button, Stack } from "@mui/material";
import { useNavigate } from "react-router-dom";
import { useState } from "react"
import { apiFetch } from "../api/apiFetch";
import {refresh} from "../api/auth"

export default function MenuPage (){

    const navigate = useNavigate()
    const [msg, setMsg] = useState("")


    async function check() {
        setMsg("checking...")
            try {
                const data = await apiFetch<{ username: string }>("/api/me")
                setMsg(`✅ logged in as: ${data.username}`)
            }
            catch (e) {
                if (e instanceof Error) {
                    setMsg(e.message)
                } else {
                    setMsg("Unknown error")
                }
        }
    }

    async function onrefresh(){
        setMsg("refreshing...")
            try {
                await refresh()
                setMsg("done")
            }
            catch (err) {
            setMsg(err instanceof Error ? err.message : "refresh failed")
            }
            }

    


    return(
        <>
        <Typography > this is the menu page</Typography>
        <Stack direction="column">
            <Button variant="contained" onClick={() => navigate("/login")}> login </Button>
            <Button variant="contained" onClick={() => navigate("/register")}> create new account </Button>
            <Button variant="contained" onClick={check}>Check auth</Button>
            <Button variant="contained" onClick={onrefresh}>Refresh</Button>
        </Stack>
        <Typography > {msg}</Typography>
        </>
    )
}