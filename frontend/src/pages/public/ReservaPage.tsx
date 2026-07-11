import { useState } from 'react'
import { Link } from 'react-router-dom'
import { reservasApi } from '../../services/api'

const C = {
  bg:          '#fcf9f3',
  primary:     '#823b18',
  secondary:   '#4a6457',
  gold:        '#cca730',
  surface:     '#f6f3ed',
  surfaceHigh: '#ebe8e2',
  onSurface:   '#1c1c18',
  onVariant:   '#54433c',
  outline:     '#87736b',
  outlineVar:  '#dac1b8',
}

const fieldStyle: React.CSSProperties = {
  width: '100%',
  backgroundColor: C.bg,
  border: `1px solid ${C.outlineVar}`,
  borderRadius: 8,
  padding: '12px 16px',
  fontFamily: 'Manrope, sans-serif',
  fontSize: 15,
  color: C.onSurface,
  outline: 'none',
  transition: 'border-color 0.2s',
  boxSizing: 'border-box',
}

const labelStyle: React.CSSProperties = {
  display: 'block',
  fontFamily: 'Manrope, sans-serif',
  fontSize: 12,
  fontWeight: 700,
  letterSpacing: '0.1em',
  textTransform: 'uppercase' as const,
  color: C.outline,
  marginBottom: 8,
}

export default function ReservaPage() {
  const [form, setForm] = useState({
    nombre: '', whatsapp: '', fecha: '', personas: '2', ocasion: '',
  })
  const [enviado, setEnviado] = useState(false)
  const [loading, setLoading] = useState(false)
  const [focused, setFocused] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    try {
      await reservasApi.crear({
        nombre: form.nombre,
        whatsapp: form.whatsapp,
        fecha: form.fecha,
        personas: form.personas,
        ocasion: form.ocasion,
      })
      setEnviado(true)
    } catch (err: any) {
      alert(err.response?.data?.error || 'Error al enviar la reserva, intenta de nuevo')
    } finally {
      setLoading(false)
    }
  }

  const minDate = new Date().toISOString().split('T')[0]

  const inputStyle = (name: string): React.CSSProperties => ({
    ...fieldStyle,
    borderColor: focused === name ? C.primary : C.outlineVar,
    boxShadow: focused === name ? `0 0 0 3px ${C.primary}18` : 'none',
  })

  return (
    <div style={{ minHeight: '100vh', backgroundColor: C.bg, fontFamily: 'Manrope, sans-serif', color: C.onSurface }}>

      {/* ── HEADER ─────────────────────────────────────────────────────── */}
      <header style={{
        position: 'sticky', top: 0, zIndex: 50,
        backgroundColor: 'rgba(252,249,243,0.95)', backdropFilter: 'blur(12px)',
        borderBottom: `1px solid ${C.outlineVar}`,
        padding: '0 24px',
        display: 'flex', justifyContent: 'center',
      }}>
        <div style={{ maxWidth: 1200, width: '100%', display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '14px 0' }}>
          <Link to="/" style={{
            display: 'flex', alignItems: 'center', gap: 8,
            color: C.outline, textDecoration: 'none',
            fontFamily: 'Manrope', fontSize: 13, fontWeight: 600, transition: 'color 0.2s',
          }}
            onMouseEnter={e => ((e.currentTarget as HTMLAnchorElement).style.color = C.primary)}
            onMouseLeave={e => ((e.currentTarget as HTMLAnchorElement).style.color = C.outline)}>
            <span className="material-symbols-outlined" style={{ fontSize: 18 }}>arrow_back</span>
            Volver
          </Link>

          <Link to="/" style={{ textDecoration: 'none', textAlign: 'center' }}>
            <p style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 20, fontWeight: 700, color: C.primary, margin: 0 }}>
              Vilca Suyo
            </p>
            <p style={{ fontFamily: 'Manrope', fontSize: 11, color: C.outline, margin: 0, letterSpacing: '0.1em' }}>
              RESERVAS
            </p>
          </Link>

          <a href="https://wa.me/51945984518" target="_blank" rel="noreferrer"
            style={{
              display: 'flex', alignItems: 'center', gap: 6,
              color: '#2e7d32', textDecoration: 'none',
              fontFamily: 'Manrope', fontSize: 13, fontWeight: 600,
            }}>
            <span className="material-symbols-outlined" style={{ fontSize: 18 }}>chat_bubble</span>
            <span className="hidden sm:inline">WhatsApp</span>
          </a>
        </div>
      </header>

      <div style={{ display: 'flex', minHeight: 'calc(100vh - 57px)' }}>

        {/* ── PANEL IZQUIERDO — imagen + info ────────────────────────────── */}
        <div className="hidden lg:flex" style={{
          width: '42%', position: 'sticky', top: 57, height: 'calc(100vh - 57px)',
          flexDirection: 'column', overflow: 'hidden',
        }}>
          {/* Imagen de fondo */}
          <div style={{ position: 'relative', flex: 1, overflow: 'hidden' }}>
            <img src="/images/hero.png" alt="Vilcashuamán"
              style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
            <div style={{ position: 'absolute', inset: 0, background: 'linear-gradient(to bottom, rgba(130,59,24,0.4) 0%, rgba(0,0,0,0.65) 100%)' }} />

            {/* Contenido sobre imagen */}
            <div style={{ position: 'absolute', inset: 0, padding: 40, display: 'flex', flexDirection: 'column', justifyContent: 'flex-end', color: '#fff' }}>
              {/* Badge urgencia */}
              <div style={{
                display: 'inline-flex', alignItems: 'center', gap: 8,
                backgroundColor: 'rgba(186,26,26,0.85)', border: '1px solid rgba(255,100,100,0.4)',
                borderRadius: 6, padding: '6px 14px',
                fontFamily: 'Manrope', fontSize: 12, fontWeight: 700, letterSpacing: '0.1em',
                marginBottom: 20, alignSelf: 'flex-start',
              }}>
                <span style={{ width: 7, height: 7, borderRadius: '50%', backgroundColor: '#ff6b6b', flexShrink: 0, animation: 'pulse-dot 1.5s infinite' }} />
                SOLO 20 MESAS POR DÍA
              </div>

              <h2 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 32, fontWeight: 700, lineHeight: 1.25, marginBottom: 16 }}>
                Una mesa frente al fin del mundo
              </h2>
              <p style={{ fontFamily: 'Manrope', fontSize: 15, opacity: 0.85, lineHeight: '24px', marginBottom: 28 }}>
                Frente al Ushnu del Imperio Inca, en Vilcashuamán, Ayacucho. A 3,400 m.s.n.m.
              </p>

              {/* Lo que incluye */}
              <div style={{ display: 'flex', flexDirection: 'column', gap: 10 }}>
                {[
                  ['Vista panorámica al Ushnu Inca', 'location_on'],
                  ['Bienvenida con chicha artesanal', 'local_bar'],
                  ['Tour privado con guía experto', 'explore'],
                  ['Sesión fotográfica al atardecer', 'photo_camera'],
                ].map(([texto, icon]) => (
                  <div key={texto as string} style={{ display: 'flex', alignItems: 'center', gap: 10 }}>
                    <span className="material-symbols-outlined" style={{ color: C.gold, fontSize: 18, flexShrink: 0 }}>{icon as string}</span>
                    <span style={{ fontFamily: 'Manrope', fontSize: 14 }}>{texto as string}</span>
                  </div>
                ))}
              </div>

              {/* Garantías */}
              <div style={{ display: 'flex', flexWrap: 'wrap', gap: 10, marginTop: 24 }}>
                {['✓ Confirmación en 2h', '✓ Cancelación gratis 24h', '✓ Sin cobro anticipado'].map(g => (
                  <span key={g} style={{
                    backgroundColor: 'rgba(255,255,255,0.15)', borderRadius: 20,
                    padding: '4px 12px', fontFamily: 'Manrope', fontSize: 12, fontWeight: 600,
                  }}>{g}</span>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* ── PANEL DERECHO — formulario ──────────────────────────────────── */}
        <div style={{ flex: 1, display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '40px 24px', overflowY: 'auto' }}>
          <div style={{ width: '100%', maxWidth: 480 }}>

            {enviado ? (
              /* ── CONFIRMACIÓN ── */
              <div style={{ textAlign: 'center', padding: '40px 0' }}>
                <div style={{
                  width: 64, height: 64, borderRadius: '50%',
                  backgroundColor: `${C.primary}15`, border: `2px solid ${C.primary}40`,
                  display: 'flex', alignItems: 'center', justifyContent: 'center',
                  margin: '0 auto 24px',
                }}>
                  <span className="material-symbols-outlined" style={{ color: C.primary, fontSize: 32 }}>check_circle</span>
                </div>
                <h2 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 28, fontWeight: 700, color: C.primary, marginBottom: 12 }}>
                  ¡Solicitud recibida!
                </h2>
                <p style={{ color: C.onVariant, fontSize: 16, lineHeight: '24px', marginBottom: 8 }}>
                  Hola <strong style={{ color: C.onSurface }}>{form.nombre}</strong>, confirmaremos tu reserva
                  en menos de 2 horas por WhatsApp.
                </p>
                <p style={{ color: C.outline, fontSize: 14, marginBottom: 32 }}>
                  {form.fecha} · {form.personas} {form.personas === '1' ? 'persona' : 'personas'}
                </p>

                {/* Garantías */}
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr 1fr', gap: 10, marginBottom: 32 }}>
                  {[['check','Confirmación','en 2h'],['event_available','Cancelación','gratis 24h'],['payments','Sin cobro','anticipado']].map(([icon,t,s]) => (
                    <div key={t as string} style={{ border: `1px solid ${C.outlineVar}`, borderRadius: 8, padding: '12px 8px', textAlign: 'center' }}>
                      <span className="material-symbols-outlined" style={{ color: C.gold, fontSize: 20 }}>{icon as string}</span>
                      <p style={{ fontFamily: 'Manrope', fontSize: 12, fontWeight: 700, color: C.onSurface, marginTop: 4 }}>{t as string}</p>
                      <p style={{ fontFamily: 'Manrope', fontSize: 11, color: C.outline }}>{s as string}</p>
                    </div>
                  ))}
                </div>

                <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
                  <a href={`https://wa.me/51945984518?text=${encodeURIComponent(
                    `Hola, soy ${form.nombre} y quiero confirmar mi reserva para el ${form.fecha}, para ${form.personas} ${form.personas === '1' ? 'persona' : 'personas'}${form.ocasion ? `, con motivo de: ${form.ocasion}` : ''}. ¡Gracias!`
                  )}`}
                    target="_blank" rel="noreferrer"
                    style={{
                      display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 10,
                      backgroundColor: '#2e7d32', color: '#fff', borderRadius: 8,
                      padding: '14px', fontFamily: 'Manrope', fontWeight: 700, fontSize: 15,
                      textDecoration: 'none', transition: 'opacity 0.2s',
                    }}>
                    <span className="material-symbols-outlined">chat_bubble</span>
                    Confirmar por WhatsApp
                  </a>
                  <Link to="/" style={{
                    display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 8,
                    border: `1px solid ${C.outlineVar}`, color: C.onVariant, borderRadius: 8,
                    padding: '14px', fontFamily: 'Manrope', fontWeight: 600, fontSize: 15,
                    textDecoration: 'none',
                  }}>
                    Volver al inicio
                  </Link>
                </div>
              </div>

            ) : (
              /* ── FORMULARIO ── */
              <>
                {/* Header */}
                <div style={{ marginBottom: 36 }}>
                  <span style={{ fontFamily: 'Manrope', fontSize: 12, fontWeight: 700, letterSpacing: '0.15em', color: C.primary, textTransform: 'uppercase', display: 'block', marginBottom: 8 }}>
                    Vilca Suyo · Reservas
                  </span>
                  <h1 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 32, fontWeight: 700, color: C.onSurface, marginBottom: 10, lineHeight: 1.2 }}>
                    Asegura tu Mesa
                  </h1>
                  <p style={{ color: C.onVariant, fontSize: 15, lineHeight: '22px' }}>
                    Confirmaremos disponibilidad en menos de 2 horas.
                  </p>

                  {/* Divider ornamental */}
                  <div style={{ display: 'flex', alignItems: 'center', gap: 12, marginTop: 20 }}>
                    <div style={{ flex: 1, height: 1, backgroundColor: C.outlineVar }} />
                    <span className="material-symbols-outlined" style={{ color: C.gold, fontSize: 20 }}>filter_vintage</span>
                    <div style={{ flex: 1, height: 1, backgroundColor: C.outlineVar }} />
                  </div>
                </div>

                <form onSubmit={handleSubmit} style={{ display: 'flex', flexDirection: 'column', gap: 20 }}>

                  {/* Nombre */}
                  <div>
                    <label style={labelStyle}>Nombre completo *</label>
                    <input required type="text" placeholder="¿Cómo te llamamos?"
                      value={form.nombre}
                      onChange={e => setForm(p => ({ ...p, nombre: e.target.value }))}
                      onFocus={() => setFocused('nombre')}
                      onBlur={() => setFocused(null)}
                      style={inputStyle('nombre')} />
                  </div>

                  {/* WhatsApp */}
                  <div>
                    <label style={labelStyle}>WhatsApp *</label>
                    <div style={{ display: 'flex' }}>
                      <div style={{
                        backgroundColor: C.surface, border: `1px solid ${C.outlineVar}`,
                        borderRight: 'none', borderRadius: '8px 0 0 8px',
                        padding: '12px 14px', fontFamily: 'Manrope', fontSize: 15, fontWeight: 700,
                        color: C.onVariant, display: 'flex', alignItems: 'center', whiteSpace: 'nowrap',
                        borderColor: focused === 'whatsapp' ? C.primary : C.outlineVar,
                      }}>+51</div>
                      <input required type="tel" placeholder="945 984 518"
                        value={form.whatsapp}
                        onChange={e => setForm(p => ({ ...p, whatsapp: e.target.value }))}
                        onFocus={() => setFocused('whatsapp')}
                        onBlur={() => setFocused(null)}
                        style={{ ...inputStyle('whatsapp'), borderRadius: '0 8px 8px 0', flex: 1 }} />
                    </div>
                    <p style={{ fontFamily: 'Manrope', fontSize: 12, color: C.outline, marginTop: 6 }}>
                      Enviamos la confirmación a este número
                    </p>
                  </div>

                  {/* Fecha + Personas */}
                  <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 14 }}>
                    <div>
                      <label style={labelStyle}>Fecha *</label>
                      <input required type="date" min={minDate}
                        value={form.fecha}
                        onChange={e => setForm(p => ({ ...p, fecha: e.target.value }))}
                        onFocus={() => setFocused('fecha')}
                        onBlur={() => setFocused(null)}
                        style={inputStyle('fecha')} />
                    </div>
                    <div>
                      <label style={labelStyle}>Personas *</label>
                      <select required value={form.personas}
                        onChange={e => setForm(p => ({ ...p, personas: e.target.value }))}
                        onFocus={() => setFocused('personas')}
                        onBlur={() => setFocused(null)}
                        style={{ ...inputStyle('personas'), cursor: 'pointer', appearance: 'auto' }}>
                        {['1','2','3','4','5','6','7','8','9','10+'].map(n => (
                          <option key={n} value={n}>{n} {n === '1' ? 'persona' : 'personas'}</option>
                        ))}
                      </select>
                    </div>
                  </div>

                  {/* Ocasión */}
                  <div>
                    <label style={labelStyle}>
                      Ocasión{' '}
                      <span style={{ fontWeight: 400, textTransform: 'none', letterSpacing: 0, color: C.outline }}>(opcional)</span>
                    </label>
                    <select value={form.ocasion}
                      onChange={e => setForm(p => ({ ...p, ocasion: e.target.value }))}
                      onFocus={() => setFocused('ocasion')}
                      onBlur={() => setFocused(null)}
                      style={{ ...inputStyle('ocasion'), cursor: 'pointer', appearance: 'auto', color: form.ocasion ? C.onSurface : C.outline }}>
                      <option value="">Sin ocasión especial</option>
                      <option value="cumpleaños">🎂 Cumpleaños</option>
                      <option value="aniversario">💍 Aniversario</option>
                      <option value="negocios">💼 Reunión de negocios</option>
                      <option value="grupo">👥 Grupo turístico</option>
                      <option value="luna">🌙 Luna de miel</option>
                    </select>
                  </div>

                  {/* CTAs */}
                  <div style={{ display: 'flex', flexDirection: 'column', gap: 12, marginTop: 8 }}>
                    <button type="submit" disabled={loading}
                      style={{
                        backgroundColor: C.primary, color: '#fff',
                        border: 'none', borderRadius: 8, padding: '16px',
                        fontFamily: 'Manrope', fontWeight: 700, fontSize: 15, letterSpacing: '0.05em',
                        cursor: loading ? 'not-allowed' : 'pointer', opacity: loading ? 0.7 : 1,
                        transition: 'transform 0.2s, opacity 0.2s',
                        display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 8,
                      }}
                      onMouseEnter={e => { if (!loading) (e.currentTarget as HTMLButtonElement).style.transform = 'scale(1.02)' }}
                      onMouseLeave={e => { (e.currentTarget as HTMLButtonElement).style.transform = 'scale(1)' }}>
                      {loading ? (
                        <>
                          <span className="material-symbols-outlined" style={{ fontSize: 18, animation: 'spin 1s linear infinite' }}>progress_activity</span>
                          Enviando...
                        </>
                      ) : (
                        <>
                          <span className="material-symbols-outlined" style={{ fontSize: 18 }}>event_available</span>
                          Confirmar Reserva
                        </>
                      )}
                    </button>

                    <a href={`https://wa.me/51945984518?text=${encodeURIComponent(
                      `Hola, quiero reservar una mesa en Vilca Suyo${form.nombre ? ` — soy ${form.nombre}` : ''}${form.fecha ? `, para el ${form.fecha}` : ''}${form.personas ? `, ${form.personas} ${form.personas === '1' ? 'persona' : 'personas'}` : ''}${form.ocasion ? `, con motivo de: ${form.ocasion}` : ''}.`
                    )}`}
                      target="_blank" rel="noreferrer"
                      style={{
                        display: 'flex', alignItems: 'center', justifyContent: 'center', gap: 10,
                        border: `1px solid #4caf5050`, backgroundColor: '#f1f8f1',
                        color: '#2e7d32', borderRadius: 8, padding: '14px',
                        fontFamily: 'Manrope', fontWeight: 700, fontSize: 14,
                        textDecoration: 'none', transition: 'background 0.2s',
                      }}>
                      <span className="material-symbols-outlined" style={{ fontSize: 20 }}>chat_bubble</span>
                      Reservar por WhatsApp
                    </a>
                  </div>

                  {/* Nota */}
                  <p style={{ textAlign: 'center', color: C.outline, fontSize: 12, lineHeight: '18px', borderTop: `1px solid ${C.outlineVar}`, paddingTop: 16 }}>
                    Al enviar aceptas nuestros términos. Sin cobro anticipado — solo pagas cuando visitas.
                  </p>
                </form>
              </>
            )}
          </div>
        </div>
      </div>

      <style>{`
        @keyframes pulse-dot {
          0%, 100% { opacity: 1; transform: scale(1); }
          50% { opacity: 0.6; transform: scale(1.3); }
        }
        @keyframes spin {
          from { transform: rotate(0deg); }
          to { transform: rotate(360deg); }
        }
        input[type="date"]::-webkit-calendar-picker-indicator { cursor: pointer; opacity: 0.6; }
      `}</style>
    </div>
  )
}
