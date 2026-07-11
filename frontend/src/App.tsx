import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import PrivateRoute from './components/layout/PrivateRoute'
import DashboardLayout from './components/layout/DashboardLayout'

import LandingPage       from './pages/public/LandingPage'
import LoginPage         from './pages/public/LoginPage'
import ReservaPage       from './pages/public/ReservaPage'
import DashboardHome     from './pages/dashboard/DashboardHome'
import ClientesPage      from './pages/dashboard/ClientesPage'
import ReservasPage      from './pages/dashboard/ReservasPage'
import PedidosPage       from './pages/dashboard/PedidosPage'
import InventarioPage    from './pages/dashboard/InventarioPage'
import CreditosPage      from './pages/dashboard/CreditosPage'
import IAPage            from './pages/dashboard/IAPage'
import MenuPage          from './pages/dashboard/MenuPage'
import ReportesPage      from './pages/dashboard/ReportesPage'

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          {/* Públicas */}
          <Route path="/"         element={<LandingPage />} />
          <Route path="/login"    element={<LoginPage />} />
          <Route path="/reservar" element={<ReservaPage />} />

          {/* Dashboard privado */}
          <Route path="/dashboard" element={
            <PrivateRoute>
              <DashboardLayout />
            </PrivateRoute>
          }>
            <Route index           element={<DashboardHome />} />
            <Route path="clientes" element={<ClientesPage />} />
            <Route path="reservas" element={<ReservasPage />} />
            <Route path="pedidos"  element={<PedidosPage />} />
            <Route path="inventario" element={<InventarioPage />} />
            <Route path="creditos"   element={<CreditosPage />} />
            <Route path="ia"         element={<IAPage />} />
            <Route path="menu"       element={<MenuPage />} />
            <Route path="reportes"   element={<ReportesPage />} />
          </Route>

          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  )
}
