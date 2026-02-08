//import { Container, Typography , Button } from "@mui/material"
//import { useState } from "react"
import { Routes, Route } from "react-router-dom"
import RegisterPage from "./pages/Register"


function App() {

  /*
  const [responseText , setResponseText] = useState("no response yet");

  
  function handleClick (){
    fetch("/api/health/db")
    .then((res) => {
      if (!res.ok) {
        throw new Error("HTTP error " + res.status);
      }
      return res.json();
    })
    .then((data) => {
      setResponseText(JSON.stringify(data, null, 2));
    })
    .catch((err) => {
      setResponseText("Request failed: " + err.message);
    });
  }
  


  return (
    
    <Container sx={{ mt: 4 }}>
      <Button variant="contained" onClick={handleClick}> Send http request</Button>
      <Typography variant="h4">
        server response: {responseText}
      </Typography>
    </Container>
  )

  */

  return (
    <Routes>
      <Route path="/" element = {<RegisterPage/>} />
    </Routes>
  )
}

export default App
