import { BrowserRouter, Navigate,Route, Routes } from 'react-router-dom' 
import Index from './Pages/index'
import './App.css'

function App() {

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Index />} />
        <Route path="/index" element={<Index />} />
        <Route path="/interfaz" element={<Index />} />
        <Route path="/reportes" element={<Index />} />
        <Route path="*" element={<Navigate to="/" replace={true} />} exact={true} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
