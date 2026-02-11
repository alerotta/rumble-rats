
import { Routes, Route } from "react-router-dom"
import RegisterPage from "./pages/Register"
import LoginPage from "./pages/Login"
import MenuPage from "./pages/Menu"


function App() {

  return (
    <Routes>
      <Route path="/" element = {<MenuPage/>}/>
      <Route path="/register" element = {<RegisterPage/>}/>
      <Route path="/login" element= {<LoginPage/>}/>
    </Routes>
  )
}

export default App
