import { NavLink, Outlet, useNavigate } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import {
  LayoutDashboard, Users, ShoppingBag, Package,
  CreditCard, Brain, LogOut, Menu, X, UtensilsCrossed, FileBarChart, CalendarClock
} from 'lucide-react'
import { useState } from 'react'

const C = {
  bg:         '#fcf9f3',
  primary:    '#823b18',
  secondary:  '#4a6457',
  gold:       '#cca730',
  surface:    '#f6f3ed',
  surfaceHigh:'#ebe8e2',
  onSurface:  '#1c1c18',
  onVariant:  '#54433c',
  outline:    '#87736b',
  outlineVar: '#dac1b8',
}

const navItems = [
  { to: '/dashboard', label: 'Inicio', icon: LayoutDashboard, exact: true },
  { to: '/dashboard/clientes', label: 'Clientes', icon: Users },
  { to: '/dashboard/reservas', label: 'Reservas', icon: CalendarClock },
  { to: '/dashboard/pedidos', label: 'Pedidos', icon: ShoppingBag },
  { to: '/dashboard/inventario', label: 'Inventario', icon: Package },
  { to: '/dashboard/creditos', label: 'Créditos', icon: CreditCard, adminOnly: true },
  { to: '/dashboard/ia', label: 'IA / Alertas', icon: Brain, adminOnly: true },
  { to: '/dashboard/menu', label: 'Menú público', icon: UtensilsCrossed, adminOnly: true },
  { to: '/dashboard/reportes', label: 'Reportes', icon: FileBarChart, adminOnly: true },
]

export default function DashboardLayout() {
  const { user, logout, isAdmin } = useAuth()
  const navigate = useNavigate()
  const [open, setOpen] = useState(false)

  const handleLogout = () => { logout(); navigate('/login') }
  const visibleItems = navItems.filter(i => !i.adminOnly || isAdmin)

  return (
    <div style={{ display: 'flex', height: '100vh', backgroundColor: C.bg, fontFamily: 'Manrope, sans-serif', color: C.onSurface }}>

      {/* SIDEBAR */}
      <aside id="sidebar" style={{
        position: 'fixed', insetInlineStart: 0, insetBlock: 0, zIndex: 50, width: 280,
        backgroundColor: C.bg, display: 'flex', flexDirection: 'column',
        boxShadow: '0 4px 20px rgba(0,0,0,0.08)',
        borderRight: `1px solid ${C.outlineVar}`,
        transform: open ? 'translateX(0)' : 'translateX(-100%)',
        transition: 'transform 0.3s ease',
      }}>

        {/* Header */}
        <div style={{
          padding: '24px 20px', borderBottom: `1px solid ${C.outlineVar}`,
          display: 'flex', flexDirection: 'column', gap: 8,
        }}>
          <div>
            <p style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 22, fontWeight: 700, color: C.primary, margin: 0 }}>
              Vilca Suyo
            </p>
            <p style={{ fontFamily: 'Manrope', fontSize: 11, color: C.outline, letterSpacing: '0.1em', margin: '2px 0 0' }}>
              ADMINISTRACIÓN
            </p>
          </div>
          <div style={{
            display: 'flex', alignItems: 'center', gap: 10,
            backgroundColor: `${C.primary}10`, borderRadius: 8, padding: '8px 12px',
            borderLeft: `3px solid ${C.primary}`,
          }}>
            <div style={{
              width: 36, height: 36, borderRadius: '50%', backgroundColor: `${C.primary}20`,
              display: 'flex', alignItems: 'center', justifyContent: 'center',
              fontFamily: 'Libre Caslon Text, serif', fontWeight: 700, color: C.primary, fontSize: 14,
            }}>
              {user?.nombre?.[0]}
            </div>
            <div style={{ flex: 1, minWidth: 0 }}>
              <p style={{ fontFamily: 'Manrope', fontWeight: 600, fontSize: 13, color: C.onSurface, margin: 0 }}>
                {user?.nombre}
              </p>
              <p style={{ fontFamily: 'Manrope', fontSize: 11, color: C.outline, margin: '2px 0 0' }}>
                {user?.rol === 'admin' ? 'Administrador' : 'Empleado'}
              </p>
            </div>
          </div>
        </div>

        {/* Navegación */}
        <nav style={{ flex: 1, padding: '12px 8px', overflowY: 'auto', display: 'flex', flexDirection: 'column', gap: 4 }}>
          {visibleItems.map(({ to, label, icon: Icon }) => (
            <NavLink
              key={to}
              to={to}
              end
              onClick={() => setOpen(false)}
              style={({ isActive }) => ({
                display: 'flex', alignItems: 'center', gap: 12, padding: '10px 14px',
                borderRadius: 10, textDecoration: 'none', transition: 'all 0.2s',
                backgroundColor: isActive ? `${C.primary}15` : 'transparent',
                borderLeft: `3px solid ${isActive ? C.primary : 'transparent'}`,
                color: isActive ? C.primary : C.onVariant,
                fontWeight: isActive ? 600 : 500,
                fontSize: 14,
              })}
            >
              <Icon size={18} />
              {label}
              {label === 'IA / Alertas' && <span style={{ fontFamily: 'Manrope', fontSize: 10, backgroundColor: C.gold, color: C.onSurface, borderRadius: 4, padding: '2px 6px', fontWeight: 700, marginLeft: 'auto' }}>PRO</span>}
            </NavLink>
          ))}
        </nav>

        {/* Logout */}
        <button onClick={handleLogout}
          style={{
            display: 'flex', alignItems: 'center', gap: 12, padding: '12px 14px', margin: '8px',
            backgroundColor: 'transparent', border: `1px solid ${C.outlineVar}`, borderRadius: 10,
            fontFamily: 'Manrope', fontSize: 14, fontWeight: 600, color: '#ba1a1a',
            cursor: 'pointer', transition: 'all 0.2s',
          }}
          onMouseEnter={e => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = '#ffebee'; (e.currentTarget as HTMLButtonElement).style.borderColor = '#ba1a1a' }}
          onMouseLeave={e => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = 'transparent'; (e.currentTarget as HTMLButtonElement).style.borderColor = C.outlineVar }}>
          <LogOut size={18} /> Cerrar sesión
        </button>
      </aside>

      {/* Overlay móvil */}
      {open && (
        <div
          style={{
            position: 'fixed', inset: 0, backgroundColor: 'rgba(0,0,0,0.4)', zIndex: 40, display: 'lg:hidden',
          }}
          onClick={() => setOpen(false)}
        />
      )}

      {/* CONTENIDO PRINCIPAL */}
      <div id="dashboard-content" style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden', transition: 'margin-left 0.3s ease' }}>

        {/* Header móvil */}
        <header style={{
          display: 'none', gridColumn: '1 / -1', backgroundColor: C.surface, borderBottom: `1px solid ${C.outlineVar}`,
          padding: '12px 16px', alignItems: 'center', gap: 12,
        }} className="lg:hidden">
          <button onClick={() => setOpen(!open)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: C.outline, padding: 0 }}>
            {open ? <X size={22} /> : <Menu size={22} />}
          </button>
          <p style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 18, fontWeight: 700, color: C.primary, margin: 0 }}>
            Vilca Suyo
          </p>
        </header>

        {/* Main content */}
        <main style={{
          flex: 1, overflowY: 'auto', padding: '32px 24px', backgroundColor: C.bg,
        }}>
          <Outlet />
        </main>
      </div>

      <style>{`
        @media (min-width: 1024px) {
          #sidebar { transform: translateX(0) !important; }
          #dashboard-content { margin-left: 280px; }
          header { display: none !important; }
        }
        @media (max-width: 1023.98px) {
          header { display: flex; }
        }
        ::-webkit-scrollbar { width: 8px; }
        ::-webkit-scrollbar-track { background: transparent; }
        ::-webkit-scrollbar-thumb { background: ${C.outlineVar}; border-radius: 4px; }
      `}</style>
    </div>
  )
}
