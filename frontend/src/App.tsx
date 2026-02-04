import { Container, Typography } from "@mui/material"
import { useEffect } from "react"

function App() {

  useEffect(() => {
    fetch("/api/health")
      .then((response) => {
        console.log("Raw response:", response);

        if (!response.ok) {
          throw new Error(`HTTP error ${response.status}`);
        }

        return response.json();
      })
      .then((data) => {
        console.log("Parsed JSON:", data);
      })
      .catch((error) => {
        console.error("Request failed:", error);
      });
  }, []);

  return (
    <Container sx={{ mt: 4 }}>
      <Typography variant="h4">
        Open DevTools → Console
      </Typography>
    </Container>
  )
}

export default App
