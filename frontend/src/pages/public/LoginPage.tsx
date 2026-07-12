import { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { useAuth } from '../../context/AuthContext'
import {
  Eye, EyeOff, Mail, Lock, LogIn, ArrowLeft, ShieldCheck, AlertCircle,
  User, Users, ShoppingBag, Package, CreditCard, BarChart3, Brain, Flower2, Loader2,
} from 'lucide-react'

const C = {
  bg:         '#fcf9f3',
  primary:    '#823b18',
  secondary:  '#4a6457',
  gold:       '#cca730',
  surface:    '#f6f3ed',
  onSurface:  '#1c1c18',
  onVariant:  '#54433c',
  outline:    '#87736b',
  outlineVar: '#dac1b8',
  error:      '#ba1a1a',
}

export default function LoginPage() {
  const { login } = useAuth()
  const navigate   = useNavigate()
  const [form, setForm]       = useState({ email: '', password: '' })
  const [error, setError]     = useState('')
  const [loading, setLoading] = useState(false)
  const [showPwd, setShowPwd] = useState(false)
  const [focused, setFocused] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)
    try {
      await login(form.email, form.password)
      navigate('/dashboard')
    } catch {
      setError('Credenciales inválidas. Verifica tu email y contraseña.')
    } finally {
      setLoading(false)
    }
  }

  const inputStyle = (name: string): React.CSSProperties => ({
    width: '100%',
    backgroundColor: C.bg,
    border: `1px solid ${focused === name ? C.primary : C.outlineVar}`,
    borderRadius: 8,
    padding: '13px 16px',
    fontFamily: 'Manrope, sans-serif',
    fontSize: 15,
    color: C.onSurface,
    outline: 'none',
    boxShadow: focused === name ? `0 0 0 3px ${C.primary}18` : 'none',
    transition: 'border-color 0.2s, box-shadow 0.2s',
    boxSizing: 'border-box' as const,
  })

  return (
    <div style={{
      minHeight: '100vh', display: 'flex',
      fontFamily: 'Manrope, sans-serif', color: C.onSurface,
    }}>

      {/* ── PANEL IZQUIERDO: imagen hero ─────────────────────────────────── */}
      <div className="hidden lg:block" style={{ width: '48%', position: 'relative', overflow: 'hidden' }}>
        <img src="/images/hero.png" alt="Vilcashuamán"
          style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
        <div style={{ position: 'absolute', inset: 0, background: 'linear-gradient(135deg, rgba(130,59,24,0.7) 0%, rgba(0,0,0,0.5) 100%)' }} />

        {/* Contenido sobre la imagen */}
        <div style={{ position: 'absolute', inset: 0, padding: 48, display: 'flex', flexDirection: 'column', justifyContent: 'space-between', color: '#fff' }}>
          <Link to="/" style={{ textDecoration: 'none', color: 'inherit' }}>
            <p style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 24, fontWeight: 700, margin: 0 }}>Vilca Suyo</p>
            <p style={{ fontFamily: 'Manrope', fontSize: 12, opacity: 0.7, margin: '2px 0 0', letterSpacing: '0.15em' }}>SABORES ANCESTRALES</p>
          </Link>

          <div>
            <div style={{
              display: 'inline-flex', alignItems: 'center', gap: 8, marginBottom: 24,
              backgroundColor: `${C.gold}30`, border: `1px solid ${C.gold}60`,
              borderRadius: 6, padding: '6px 14px',
              fontFamily: 'Manrope', fontSize: 12, fontWeight: 700, letterSpacing: '0.1em', color: C.gold,
            }}>
              <ShieldCheck size={16} />
              PANEL ADMINISTRATIVO
            </div>
            <h2 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 36, fontWeight: 700, lineHeight: 1.2, marginBottom: 16 }}>
              Gestiona tu restaurante desde aquí
            </h2>
            <p style={{ fontFamily: 'Manrope', fontSize: 15, opacity: 0.8, lineHeight: '24px' }}>
              Clientes, pedidos, inventario, créditos y reportes en un solo lugar.
            </p>

            {/* Módulos disponibles */}
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 10, marginTop: 28 }}>
              {(
                [
                  [Users,'Clientes'],[ShoppingBag,'Pedidos'],
                  [Package,'Inventario'],[CreditCard,'Créditos'],
                  [BarChart3,'Reportes'],[Brain,'IA'],
                ] as [React.ComponentType<{ size?: number; color?: string }>, string][]
              ).map(([Icon, label]) => (
                <div key={label} style={{ display: 'flex', alignItems: 'center', gap: 8, backgroundColor: 'rgba(255,255,255,0.1)', borderRadius: 8, padding: '8px 12px' }}>
                  <Icon color={C.gold} size={18} />
                  <span style={{ fontFamily: 'Manrope', fontSize: 13, fontWeight: 600 }}>{label}</span>
                </div>
              ))}
            </div>
          </div>

          <p style={{ fontFamily: 'Manrope', fontSize: 12, opacity: 0.5 }}>
            © 2026 Vilca Suyo · Vilcashuamán, Ayacucho
          </p>
        </div>
      </div>

      {/* ── PANEL DERECHO: formulario ────────────────────────────────────── */}
      <div style={{
        flex: 1, backgroundColor: C.bg,
        display: 'flex', alignItems: 'center', justifyContent: 'center',
        padding: '40px 24px',
      }}>
        <div style={{ width: '100%', maxWidth: 420 }}>

          {/* Logo mobile */}
          <div className="lg:hidden" style={{ textAlign: 'center', marginBottom: 32 }}>
            <Link to="/" style={{ textDecoration: 'none' }}>
              <p style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 26, fontWeight: 700, color: C.primary, margin: 0 }}>Vilca Suyo</p>
              <p style={{ fontFamily: 'Manrope', fontSize: 11, color: C.outline, letterSpacing: '0.15em', margin: '4px 0 0' }}>ADMINISTRACIÓN</p>
            </Link>
          </div>

          {/* Título */}
          <div style={{ marginBottom: 32 }}>
            <h1 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 30, fontWeight: 700, color: C.onSurface, marginBottom: 8, lineHeight: 1.2 }}>
              Bienvenido de vuelta
            </h1>
            <p style={{ color: C.onVariant, fontSize: 15 }}>
              Ingresa tus credenciales para acceder al panel.
            </p>

            {/* Divider */}
            <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginTop: 20 }}>
              <div style={{ flex: 1, height: 1, backgroundColor: C.outlineVar }} />
              <Flower2 color={C.gold} size={18} />
              <div style={{ flex: 1, height: 1, backgroundColor: C.outlineVar }} />
            </div>
          </div>

          {/* Error */}
          {error && (
            <div style={{
              display: 'flex', alignItems: 'flex-start', gap: 10,
              backgroundColor: '#fff0f0', border: `1px solid ${C.error}30`,
              borderRadius: 8, padding: '12px 16px', marginBottom: 24,
            }}>
              <AlertCircle color={C.error} size={18} style={{ flexShrink: 0 }} />
              <p style={{ fontFamily: 'Manrope', fontSize: 14, color: C.error, margin: 0 }}>{error}</p>
            </div>
          )}

          <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>

            {/* Email */}
            <div>
              <label style={{
                display: 'block', fontFamily: 'Manrope', fontSize: 12, fontWeight: 700,
                letterSpacing: '0.1em', textTransform: 'uppercase' as const, color: C.outline, marginBottom: 8,
              }}>
                Correo electrónico
              </label>
              <div style={{ position: 'relative' }}>
                <Mail size={18} style={{
                  position: 'absolute', left: 14, top: '50%', transform: 'translateY(-50%)',
                  color: focused === 'email' ? C.primary : C.outline, pointerEvents: 'none',
                }} />
                <input required type="email" placeholder="admin@restaurante.com"
                  value={form.email}
                  onChange={e => setForm(p => ({ ...p, email: e.target.value }))}
                  onFocus={() => setFocused('email')}
                  onBlur={() => setFocused(null)}
                  style={{ ...inputStyle('email'), paddingLeft: 44 }} />
              </div>
            </div>

            {/* Contraseña */}
            <div>
              <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 8 }}>
                <label style={{
                  fontFamily: 'Manrope', fontSize: 12, fontWeight: 700,
                  letterSpacing: '0.1em', textTransform: 'uppercase' as const, color: C.outline,
                }}>
                  Contraseña
                </label>
              </div>
              <div style={{ position: 'relative' }}>
                <Lock size={18} style={{
                  position: 'absolute', left: 14, top: '50%', transform: 'translateY(-50%)',
                  color: focused === 'password' ? C.primary : C.outline, pointerEvents: 'none',
                }} />
                <input required type={showPwd ? 'text' : 'password'} placeholder="••••••••"
                  value={form.password}
                  onChange={e => setForm(p => ({ ...p, password: e.target.value }))}
                  onFocus={() => setFocused('password')}
                  onBlur={() => setFocused(null)}
                  style={{ ...inputStyle('password'), paddingLeft: 44, paddingRight: 44 }} />
                <button type="button" onClick={() => setShowPwd(!showPwd)}
                  style={{
                    position: 'absolute', right: 14, top: '50%', transform: 'translateY(-50%)',
                    background: 'none', border: 'none', cursor: 'pointer',
                    color: C.outline, padding: 0, display: 'flex',
                  }}>
                  {showPwd ? <EyeOff size={18} /> : <Eye size={18} />}
                </button>
              </div>
            </div>

            {/* Botón */}
            <button type="submit" disabled={loading}
              style={{
                backgroundColor: C.primary, color: '#fff',
                border: 'none', borderRadius: 8, padding: '15px',
                fontFamily: 'Manrope', fontWeight: 700, fontSize: 15, letterSpacing: '0.03em',
                cursor: loading ? 'not-allowed' : 'pointer',
                opacity: loading ? 0.75 : 1,
                display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 10,
                transition: 'transform 0.2s, opacity 0.2s', marginTop: 4,
              }}
              onMouseEnter={e => { if (!loading) (e.currentTarget as HTMLButtonElement).style.transform = 'scale(1.02)' }}
              onMouseLeave={e => { (e.currentTarget as HTMLButtonElement).style.transform = 'scale(1)' }}>
              {loading ? (
                <>
                  <Loader2 size={18} style={{ animation: 'spin 1s linear infinite' }} />
                  Verificando...
                </>
              ) : (
                <>
                  <LogIn size={18} />
                  Ingresar al Panel
                </>
              )}
            </button>
          </form>

          {/* Credenciales de demo */}
          <div style={{
            marginTop: 28, padding: '16px 20px',
            backgroundColor: C.surface, borderRadius: 10,
            border: `1px solid ${C.outlineVar}`,
          }}>
            <p style={{ fontFamily: 'Manrope', fontSize: 12, fontWeight: 700, color: C.outline, letterSpacing: '0.1em', textTransform: 'uppercase', marginBottom: 10 }}>
              Credenciales de prueba
            </p>
            <div style={{ display: 'flex', flexDirection: 'column', gap: 6 }}>
              {[
                ['Admin', 'admin@restaurante.com', 'password123'],
                ['Empleado', 'carlos@restaurante.com', 'password123'],
              ].map(([rol, email, pwd]) => (
                <button key={rol as string} type="button"
                  onClick={() => setForm({ email: email as string, password: pwd as string })}
                  style={{
                    display: 'flex', alignItems: 'center', justifyContent: 'space-between',
                    backgroundColor: C.bg, border: `1px solid ${C.outlineVar}`,
                    borderRadius: 6, padding: '8px 12px', cursor: 'pointer',
                    transition: 'border-color 0.2s',
                    fontFamily: 'Manrope',
                  }}
                  onMouseEnter={e => (e.currentTarget.style.borderColor = C.primary)}
                  onMouseLeave={e => (e.currentTarget.style.borderColor = C.outlineVar)}>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                    {rol === 'Admin'
                      ? <ShieldCheck size={16} color={C.secondary} />
                      : <User size={16} color={C.secondary} />}
                    <span style={{ fontSize: 13, fontWeight: 600, color: C.onSurface }}>{rol as string}</span>
                  </div>
                  <span style={{ fontSize: 12, color: C.outline }}>{email as string}</span>
                </button>
              ))}
            </div>
          </div>

          {/* Volver */}
          <div style={{ textAlign: 'center', marginTop: 24 }}>
            <Link to="/" style={{
              display: 'inline-flex', alignItems: 'center', gap: 6,
              color: C.outline, textDecoration: 'none',
              fontFamily: 'Manrope', fontSize: 14, transition: 'color 0.2s',
            }}
              onMouseEnter={e => ((e.currentTarget as HTMLAnchorElement).style.color = C.primary)}
              onMouseLeave={e => ((e.currentTarget as HTMLAnchorElement).style.color = C.outline)}>
              <ArrowLeft size={16} />
              Volver a la página del restaurante
            </Link>
          </div>
        </div>
      </div>

      <style>{`
        @keyframes spin {
          from { transform: rotate(0deg); }
          to   { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  )
}
