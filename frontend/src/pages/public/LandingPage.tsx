import { useState, useEffect, useRef, useCallback } from 'react'
import { Link } from 'react-router-dom'
import { menuApi } from '../../services/api'
import type { MenuPublico } from '../../types'

// ─── Paleta tomada directamente del Stitch ───────────────────────────────────
const C = {
  bg:           '#fcf9f3',
  primary:      '#823b18',
  primaryCont:  '#a0522d',
  secondary:    '#4a6457',
  gold:         '#cca730',
  surface:      '#f0eee8',
  surfaceHigh:  '#ebe8e2',
  onSurface:    '#1c1c18',
  onVariant:    '#54433c',
  outline:      '#87736b',
  outlineVar:   '#dac1b8',
}

// ─── Hook reveal (IntersectionObserver) ──────────────────────────────────────
function useReveal() {
  const ref = useRef<HTMLDivElement>(null)
  const [active, setActive] = useState(false)
  useEffect(() => {
    const el = ref.current
    if (!el) return
    const obs = new IntersectionObserver(
      ([e]) => { if (e.isIntersecting) { setActive(true); obs.disconnect() } },
      { threshold: 0.12, rootMargin: '0px 0px -40px 0px' }
    )
    obs.observe(el)
    return () => obs.disconnect()
  }, [])
  return { ref, active }
}

function Reveal({ children, delay = 0, className = '' }: {
  children: React.ReactNode; delay?: number; className?: string
}) {
  const { ref, active } = useReveal()
  return (
    <div ref={ref} className={className} style={{
      opacity: active ? 1 : 0,
      transform: active ? 'translateY(0)' : 'translateY(28px)',
      transition: `opacity 0.8s ease-out ${delay}ms, transform 0.8s ease-out ${delay}ms`,
      willChange: 'opacity, transform',
    }}>
      {children}
    </div>
  )
}

// ─── Datos ────────────────────────────────────────────────────────────────────
// Se usa solo si el backend aún no tiene platos cargados en el menú público.
const FALLBACK_MENU = [
  {
    img: '/images/puca-picante.png',
    badge: 'Firma Local',
    nombre: 'Puca Picante Premium',
    desc: 'Chicharrón de cerdo crocante sobre un cremoso de beterraga y maní tostado.',
    precio: 'S/ 64.00',
    delay: 0,
  },
  {
    img: '/images/cuy-chactado.png',
    badge: 'Modern Twist',
    nombre: 'Cuy Chactado de Autor',
    desc: 'Filete de cuy prensado y crocante, acompañado de milhojas de papa andina.',
    precio: 'S/ 78.00',
    delay: 200,
  },
  {
    img: '/images/pachamanca.png',
    badge: 'Tradicional',
    nombre: 'Pachamanca Tradicional',
    desc: 'Cinco tipos de carnes cocidas bajo tierra con piedras volcánicas y chincho aromático.',
    precio: 'S/ 85.00',
    delay: 400,
  },
]

const TESTIMONIOS = [
  {
    stars: 5,
    texto: '"The view of the ruins while dining was spiritual. The Puca Picante is a revelation of flavor. Truly a 5-star experience in the middle of the Andes."',
    nombre: 'Elena Rodriguez',
    origen: 'España',
    initial: 'E',
  },
  {
    stars: 5,
    texto: '"Authentic, high-end, and meaningful. The fusion of Inca history with gourmet food is something I have never seen before. A must-visit!"',
    nombre: 'Jameson Clarke',
    origen: 'Australia',
    initial: 'J',
  },
  {
    stars: 5,
    texto: '"La mejor Pachamanca que he probado. El servicio es impecable y el ambiente te transporta a otra época. Increíble Vilcashuamán."',
    nombre: 'Sofia Méndez',
    origen: 'Lima, Perú',
    initial: 'S',
  },
]

// ─── Componente principal ─────────────────────────────────────────────────────
export default function LandingPage() {
  const [scrollY, setScrollY] = useState(0)
  const [lastScroll, setLastScroll] = useState(0)
  const [navHidden, setNavHidden] = useState(false)
  const [menuOpen, setMenuOpen] = useState(false)
  const [menuItems, setMenuItems] = useState<MenuPublico[]>([])
  const heroImgRef = useRef<HTMLImageElement>(null)

  useEffect(() => {
    menuApi.obtener()
      .then(({ data }) => setMenuItems(data))
      .catch(() => setMenuItems([]))
  }, [])

  const menuCards = menuItems.length > 0
    ? menuItems.map((item, i) => ({
        img: item.imagen_url || '/images/hero.png',
        badge: item.categoria,
        nombre: item.nombre,
        desc: item.descripcion,
        precio: `S/ ${item.precio.toFixed(2)}`,
        delay: i * 200,
      }))
    : FALLBACK_MENU

  const handleScroll = useCallback(() => {
    const cur = window.scrollY
    setScrollY(cur)
    setNavHidden(cur > lastScroll && cur > 200)
    setLastScroll(cur)
    if (heroImgRef.current) {
      heroImgRef.current.style.transform = `translateY(${cur * 0.35}px)`
    }
  }, [lastScroll])

  useEffect(() => {
    window.addEventListener('scroll', handleScroll, { passive: true })
    return () => window.removeEventListener('scroll', handleScroll)
  }, [handleScroll])

  return (
    <div style={{ backgroundColor: C.bg, color: C.onSurface, fontFamily: 'Manrope, sans-serif', overflowX: 'hidden' }}>

      {/* ── NAV ──────────────────────────────────────────────────────────── */}
      <nav style={{
        position: 'fixed', top: 0, left: 0, width: '100%', zIndex: 50,
        display: 'flex', justifyContent: 'center',
        backgroundColor: scrollY > 50 ? 'rgba(252,249,243,0.97)' : 'rgba(252,249,243,0.85)',
        backdropFilter: 'blur(12px)',
        borderBottom: `1px solid ${C.outlineVar}`,
        boxShadow: scrollY > 50 ? '0 2px 16px rgba(0,0,0,0.08)' : 'none',
        transform: navHidden ? 'translateY(-100%)' : 'translateY(0)',
        transition: 'transform 0.3s ease, box-shadow 0.3s ease, background-color 0.3s ease',
        padding: '0 24px',
      }}>
        <div style={{ maxWidth: 1200, width: '100%', display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '14px 0' }}>
          <span style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 22, fontWeight: 700, color: C.primary }}>
            Vilca Suyo
          </span>

          {/* Desktop links */}
          <div className="hidden md:flex" style={{ gap: 28 }}>
            {[['Heritage', '#'], ['Menú', '#menu'], ['Experiencia', '#experiencia'], ['Ubicación', '#ubicacion']].map(([label, href]) => (
              <a key={label} href={href}
                style={{ fontFamily: 'Manrope', fontSize: 13, fontWeight: 600, letterSpacing: '0.05em', color: C.onVariant, textDecoration: 'none' }}
                onMouseEnter={e => (e.currentTarget.style.color = C.primary)}
                onMouseLeave={e => (e.currentTarget.style.color = C.onVariant)}>
                {label}
              </a>
            ))}
          </div>

          <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
            <Link to="/reservar" style={{
              backgroundColor: C.primary, color: '#fff',
              padding: '8px 20px', borderRadius: 8,
              fontFamily: 'Manrope', fontSize: 13, fontWeight: 700, letterSpacing: '0.05em',
              textDecoration: 'none', transition: 'transform 0.2s',
            }}
              onMouseEnter={e => (e.currentTarget.style.transform = 'scale(1.04)')}
              onMouseLeave={e => (e.currentTarget.style.transform = 'scale(1)')}>
              Reservar Mesa
            </Link>
            <Link to="/login" style={{ fontSize: 12, color: C.outline, textDecoration: 'none' }}>Admin</Link>
            <button className="md:hidden" onClick={() => setMenuOpen(!menuOpen)}
              style={{ background: 'none', border: 'none', cursor: 'pointer', color: C.onSurface }}>
              <span className="material-symbols-outlined">{menuOpen ? 'close' : 'menu'}</span>
            </button>
          </div>
        </div>
      </nav>

      {/* Mobile menu */}
      {menuOpen && (
        <div style={{
          position: 'fixed', top: 56, left: 0, width: '100%', zIndex: 49,
          backgroundColor: C.bg, borderBottom: `1px solid ${C.outlineVar}`,
          padding: '16px 24px', display: 'flex', flexDirection: 'column', gap: 16,
        }}>
          {[['Heritage', '#'], ['Menú', '#menu'], ['Experiencia', '#experiencia'], ['Ubicación', '#ubicacion']].map(([l, h]) => (
            <a key={l} href={h} onClick={() => setMenuOpen(false)}
              style={{ color: C.onSurface, textDecoration: 'none', fontWeight: 600, fontSize: 15 }}>{l}</a>
          ))}
        </div>
      )}

      {/* ── HERO ─────────────────────────────────────────────────────────── */}
      <section style={{ position: 'relative', height: '100vh', display: 'flex', alignItems: 'center', justifyContent: 'center', overflow: 'hidden' }}>
        {/* Parallax image */}
        <div style={{ position: 'absolute', inset: 0, overflow: 'hidden' }}>
          <div style={{ position: 'absolute', inset: 0, background: 'rgba(0,0,0,0.45)', zIndex: 1 }} />
          <img ref={heroImgRef} src="/images/hero.png" alt="Complejo Arqueológico de Vilcashuamán"
            style={{ width: '100%', height: '120%', objectFit: 'cover', top: '-10%', position: 'absolute', willChange: 'transform' }} />
        </div>

        {/* Content */}
        <div style={{ position: 'relative', zIndex: 2, textAlign: 'center', color: '#fff', padding: '0 24px', maxWidth: 900 }}>
          <Reveal>
            <h1 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 'clamp(36px, 6vw, 64px)', lineHeight: 1.15, fontWeight: 700, marginBottom: 20 }}>
              Sabores ancestrales en el corazón de Vilcashuamán
            </h1>
          </Reveal>
          <Reveal delay={200}>
            <p style={{ fontSize: 18, lineHeight: '28px', maxWidth: 600, margin: '0 auto 36px', opacity: 0.9, fontFamily: 'Manrope' }}>
              Una experiencia culinaria única donde la historia inca se encuentra con la alta cocina ayacuchana.
            </p>
          </Reveal>
          <Reveal delay={400}>
            <div style={{ display: 'flex', gap: 12, justifyContent: 'center', flexWrap: 'wrap' }}>
              <Link to="/reservar" style={{
                backgroundColor: C.primary, color: '#fff',
                padding: '14px 40px', borderRadius: 8,
                fontFamily: 'Manrope', fontWeight: 700, fontSize: 14, letterSpacing: '0.05em',
                textDecoration: 'none', transition: 'background 0.2s',
              }}>
                Reservar Mesa
              </Link>
              <a href="#menu" style={{
                border: `1px solid ${C.gold}`, color: C.gold,
                padding: '14px 40px', borderRadius: 8,
                fontFamily: 'Manrope', fontWeight: 700, fontSize: 14, letterSpacing: '0.05em',
                textDecoration: 'none', transition: 'all 0.2s',
              }}>
                Ver Menú
              </a>
            </div>
          </Reveal>
        </div>

        {/* Scroll indicator */}
        <div style={{ position: 'absolute', bottom: 24, left: '50%', transform: 'translateX(-50%)', zIndex: 2, animation: 'bounce 2s infinite' }}>
          <span className="material-symbols-outlined" style={{ color: '#fff', fontSize: 40 }}>expand_more</span>
        </div>
      </section>

      {/* ── SOBRE EL LUGAR ───────────────────────────────────────────────── */}
      <section style={{ padding: '80px 24px', position: 'relative', overflow: 'hidden' }}>
        {/* Inca pattern fondo */}
        <div style={{
          position: 'absolute', inset: 0, opacity: 0.05, pointerEvents: 'none',
          backgroundImage: `radial-gradient(${C.primary} 0.5px, transparent 0.5px), radial-gradient(${C.primary} 0.5px, ${C.bg} 0.5px)`,
          backgroundSize: '20px 20px', backgroundPosition: '0 0, 10px 10px',
        }} />
        <div style={{ maxWidth: 1200, margin: '0 auto', display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: 80, alignItems: 'center', position: 'relative', zIndex: 1 }}>
          <Reveal>
            <span style={{ fontFamily: 'Manrope', fontSize: 12, fontWeight: 700, letterSpacing: '0.15em', color: C.primary, textTransform: 'uppercase', display: 'block', marginBottom: 8 }}>
              Nuestra Herencia
            </span>
            <h2 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 'clamp(24px,3vw,32px)', fontWeight: 600, marginBottom: 20, color: C.onSurface, lineHeight: 1.3 }}>
              Un Santuario Gastronómico al Pie del Templo del Sol
            </h2>
            <p style={{ fontSize: 17, lineHeight: '28px', color: C.onVariant, marginBottom: 16 }}>
              Ubicados a solo pasos del Complejo Arqueológico de Vilcashuamán, nuestro restaurante nace como un homenaje a la grandeza del Tahuantinsuyo. Aquí, cada piedra cuenta una historia y cada plato es un viaje en el tiempo.
            </p>
            <p style={{ fontSize: 17, lineHeight: '28px', color: C.onVariant }}>
              Fusionamos técnicas milenarias de cocción en tierra y piedra con la sofisticación de la gastronomía contemporánea, utilizando insumos orgánicos cultivados por comunidades locales de Ayacucho.
            </p>
            <div style={{ display: 'flex', alignItems: 'center', margin: '28px 0', gap: 16 }}>
              <div style={{ flex: 1, height: 1, backgroundColor: C.outlineVar }} />
              <span className="material-symbols-outlined" style={{ color: C.gold, fontSize: 28 }}>filter_vintage</span>
              <div style={{ flex: 1, height: 1, backgroundColor: C.outlineVar }} />
            </div>
          </Reveal>

          <Reveal delay={300}>
            <div style={{ position: 'relative' }}>
              <div style={{ borderRadius: 12, overflow: 'hidden', boxShadow: '0 20px 60px rgba(0,0,0,0.15)', border: '4px solid #fff', aspectRatio: '1/1' }}>
                <img src="/images/hero.png" alt="Vilcashuamán" style={{ width: '100%', height: '100%', objectFit: 'cover' }} />
              </div>
              <div style={{
                position: 'absolute', bottom: -16, left: -16,
                backgroundColor: C.secondary, color: '#fff',
                padding: '20px 24px', borderRadius: 10, boxShadow: '0 8px 24px rgba(0,0,0,0.2)',
              }} className="hidden md:block">
                <p style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 20, fontWeight: 600 }}>Siglo XV</p>
                <p style={{ fontSize: 12, opacity: 0.8, fontFamily: 'Manrope' }}>Fundación del Templo del Sol</p>
              </div>
            </div>
          </Reveal>
        </div>
      </section>

      {/* ── MENÚ DESTACADO ───────────────────────────────────────────────── */}
      <section id="menu" style={{ padding: '80px 24px', backgroundColor: '#f6f3ed' }}>
        <div style={{ maxWidth: 1200, margin: '0 auto' }}>
          <Reveal>
            <div style={{ textAlign: 'center', marginBottom: 60 }}>
              <h2 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 'clamp(24px,3vw,32px)', fontWeight: 600, marginBottom: 12 }}>
                Selecciones de Autor
              </h2>
              <p style={{ fontSize: 16, color: C.onVariant, maxWidth: 520, margin: '0 auto' }}>
                Una curaduría de nuestros platos más emblemáticos, reinterpretados para el paladar moderno sin perder su esencia ancestral.
              </p>
            </div>
          </Reveal>

          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(280px, 1fr))', gap: 24 }}>
            {menuCards.map((card) => (
              <Reveal key={card.nombre} delay={card.delay}>
                <div style={{
                  backgroundColor: C.bg, borderRadius: 12, overflow: 'hidden',
                  boxShadow: '0 2px 8px rgba(0,0,0,0.07)',
                  transition: 'box-shadow 0.4s, transform 0.4s',
                  cursor: 'pointer',
                }}
                  onMouseEnter={e => { (e.currentTarget as HTMLDivElement).style.transform = 'translateY(-8px)'; (e.currentTarget as HTMLDivElement).style.boxShadow = '0 20px 50px rgba(0,0,0,0.15)' }}
                  onMouseLeave={e => { (e.currentTarget as HTMLDivElement).style.transform = 'translateY(0)'; (e.currentTarget as HTMLDivElement).style.boxShadow = '0 2px 8px rgba(0,0,0,0.07)' }}>
                  <div style={{ height: 240, overflow: 'hidden' }}>
                    <img src={card.img} alt={card.nombre}
                      style={{ width: '100%', height: '100%', objectFit: 'cover', transition: 'transform 0.6s' }}
                      onMouseEnter={e => ((e.target as HTMLImageElement).style.transform = 'scale(1.08)')}
                      onMouseLeave={e => ((e.target as HTMLImageElement).style.transform = 'scale(1)')} />
                  </div>
                  <div style={{ padding: '20px 24px', textAlign: 'center' }}>
                    <span style={{
                      backgroundColor: `${C.secondary}18`, color: C.secondary,
                      padding: '3px 10px', borderRadius: 999, fontSize: 12, fontWeight: 600,
                      display: 'inline-block', marginBottom: 10,
                    }}>{card.badge}</span>
                    <h3 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 20, color: C.primary, marginBottom: 8 }}>
                      {card.nombre}
                    </h3>
                    <p style={{ fontSize: 15, color: C.onVariant, lineHeight: '22px', marginBottom: 16 }}>
                      {card.desc}
                    </p>
                    <div style={{ borderTop: `1px solid ${C.outlineVar}`, paddingTop: 14 }}>
                      <span style={{ fontFamily: 'Manrope', fontSize: 15, fontWeight: 700, color: C.gold }}>
                        {card.precio}
                      </span>
                    </div>
                  </div>
                </div>
              </Reveal>
            ))}
          </div>
        </div>
      </section>

      {/* ── EXPERIENCIA TURÍSTICA ─────────────────────────────────────────── */}
      <section id="experiencia" style={{ padding: '80px 24px', backgroundColor: C.bg, overflow: 'hidden' }}>
        <div style={{ maxWidth: 1200, margin: '0 auto', display: 'flex', flexWrap: 'wrap', gap: 80, alignItems: 'center' }}>

          {/* Foto grid */}
          <Reveal className="w-full md:w-[45%]">
            <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 12 }}>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 12 }}>
                <div style={{ borderRadius: 12, overflow: 'hidden', aspectRatio: '3/4' }}>
                  <img src="/images/hero.png" alt="Ushnu Inca"
                    style={{ width: '100%', height: '100%', objectFit: 'cover', transition: 'transform 0.5s' }}
                    onMouseEnter={e => ((e.target as HTMLImageElement).style.transform = 'scale(1.06)')}
                    onMouseLeave={e => ((e.target as HTMLImageElement).style.transform = 'scale(1)')} />
                </div>
                <div style={{ backgroundColor: `${C.primary}0d`, padding: '16px 20px', borderRadius: 12, border: `1px solid ${C.primary}18` }}>
                  <h4 style={{ fontFamily: 'Manrope', fontSize: 13, fontWeight: 700, color: C.primary, letterSpacing: '0.05em', marginBottom: 6 }}>Tours Privados</h4>
                  <p style={{ fontSize: 12, color: C.onVariant }}>Visitas guiadas al Ushnu antes de la cena.</p>
                </div>
              </div>
              <div style={{ display: 'flex', flexDirection: 'column', gap: 12, marginTop: 48 }}>
                <div style={{ backgroundColor: `${C.secondary}0d`, padding: '16px 20px', borderRadius: 12, border: `1px solid ${C.secondary}18` }}>
                  <h4 style={{ fontFamily: 'Manrope', fontSize: 13, fontWeight: 700, color: C.secondary, letterSpacing: '0.05em', marginBottom: 6 }}>Traveler Hub</h4>
                  <p style={{ fontSize: 12, color: C.onVariant }}>Conserjería para viajeros internacionales.</p>
                </div>
                <div style={{ borderRadius: 12, overflow: 'hidden', aspectRatio: '3/4' }}>
                  <img src="/images/puca-picante.png" alt="Puca Picante"
                    style={{ width: '100%', height: '100%', objectFit: 'cover', transition: 'transform 0.5s' }}
                    onMouseEnter={e => ((e.target as HTMLImageElement).style.transform = 'scale(1.06)')}
                    onMouseLeave={e => ((e.target as HTMLImageElement).style.transform = 'scale(1)')} />
                </div>
              </div>
            </div>
          </Reveal>

          {/* Texto */}
          <Reveal delay={300} className="flex-1 min-w-[280px]">
            <span style={{ fontFamily: 'Manrope', fontSize: 12, fontWeight: 700, letterSpacing: '0.15em', color: C.primary, textTransform: 'uppercase', display: 'block', marginBottom: 8 }}>
              Para el Explorador Moderno
            </span>
            <h2 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 'clamp(24px,3vw,32px)', fontWeight: 600, marginBottom: 20, lineHeight: 1.3 }}>
              Más que una cena, un portal a la historia
            </h2>
            <p style={{ fontSize: 17, lineHeight: '28px', color: C.onVariant, marginBottom: 32 }}>
              Diseñamos experiencias completas para viajeros que buscan profundidad. Paquetes que incluyen recorridos privados por el Ushnu con guías expertos, seguidos de una cata de chicha de jora premium y nuestro menú degustación.
            </p>
            {[
              { icon: 'map', titulo: 'Ubicación Privilegiada', sub: 'Vistas directas al Trono del Inca desde nuestras terrazas.' },
              { icon: 'translate', titulo: 'Atención Bilingüe', sub: 'Personal capacitado en español e inglés para viajeros globales.' },
              { icon: 'luggage', titulo: 'Logística de Viaje', sub: 'Coordinación de transporte desde la ciudad de Ayacucho.' },
            ].map(({ icon, titulo, sub }) => (
              <div key={titulo} style={{ display: 'flex', gap: 20, marginBottom: 20 }}>
                <span className="material-symbols-outlined" style={{ color: C.gold, fontSize: 24, marginTop: 2, flexShrink: 0 }}>{icon}</span>
                <div>
                  <h5 style={{ fontFamily: 'Manrope', fontWeight: 700, fontSize: 14, marginBottom: 4 }}>{titulo}</h5>
                  <p style={{ fontSize: 15, color: C.onVariant }}>{sub}</p>
                </div>
              </div>
            ))}
          </Reveal>
        </div>
      </section>

      {/* ── TESTIMONIOS ──────────────────────────────────────────────────── */}
      <section style={{ padding: '80px 24px', backgroundColor: C.surfaceHigh }}>
        <div style={{ maxWidth: 1200, margin: '0 auto' }}>
          <Reveal>
            <h2 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 'clamp(24px,3vw,32px)', fontWeight: 600, textAlign: 'center', marginBottom: 48 }}>
              Ecos de Nuestros Visitantes
            </h2>
          </Reveal>
          <div style={{ display: 'flex', gap: 20, overflowX: 'auto', paddingBottom: 16, scrollSnapType: 'x mandatory' }}>
            {TESTIMONIOS.map((t, i) => (
              <Reveal key={t.nombre} delay={i * 200}>
                <div style={{
                  minWidth: 340, backgroundColor: C.bg, padding: '28px 32px',
                  borderRadius: 12, boxShadow: '0 4px 16px rgba(0,0,0,0.07)',
                  scrollSnapAlign: 'center', transition: 'box-shadow 0.3s',
                }}
                  onMouseEnter={e => ((e.currentTarget as HTMLDivElement).style.boxShadow = '0 12px 40px rgba(0,0,0,0.14)')}
                  onMouseLeave={e => ((e.currentTarget as HTMLDivElement).style.boxShadow = '0 4px 16px rgba(0,0,0,0.07)')}>
                  <div style={{ display: 'flex', gap: 2, marginBottom: 16 }}>
                    {Array(t.stars).fill(0).map((_, j) => (
                      <span key={j} className="material-symbols-outlined" style={{ color: C.gold, fontSize: 20, fontVariationSettings: '"FILL" 1' }}>star</span>
                    ))}
                  </div>
                  <p style={{ fontSize: 17, fontStyle: 'italic', lineHeight: '28px', color: C.onSurface, marginBottom: 24 }}>{t.texto}</p>
                  <div style={{ display: 'flex', alignItems: 'center', gap: 14 }}>
                    <div style={{
                      width: 44, height: 44, borderRadius: '50%',
                      backgroundColor: C.surfaceHigh,
                      display: 'flex', alignItems: 'center', justifyContent: 'center',
                      fontFamily: 'Libre Caslon Text, serif', fontWeight: 700, fontSize: 18, color: C.primary,
                    }}>{t.initial}</div>
                    <div>
                      <p style={{ fontFamily: 'Manrope', fontWeight: 700, fontSize: 14 }}>{t.nombre}</p>
                      <p style={{ fontSize: 12, color: C.onVariant }}>{t.origen}</p>
                    </div>
                  </div>
                </div>
              </Reveal>
            ))}
          </div>
        </div>
      </section>

      {/* ── UBICACIÓN ────────────────────────────────────────────────────── */}
      <section id="ubicacion" style={{ padding: '80px 24px', backgroundColor: C.bg }}>
        <div style={{ maxWidth: 1200, margin: '0 auto', display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(300px, 1fr))', gap: 80, alignItems: 'center' }}>
          <Reveal>
            <span style={{ fontFamily: 'Manrope', fontSize: 12, fontWeight: 700, letterSpacing: '0.15em', color: C.primary, textTransform: 'uppercase', display: 'block', marginBottom: 8 }}>
              ¿Cómo Llegar?
            </span>
            <h2 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 'clamp(22px,3vw,30px)', fontWeight: 600, marginBottom: 28, lineHeight: 1.3 }}>
              Encuéntranos en la Plaza Mayor de Vilcashuamán
            </h2>
            {[
              { icon: 'distance', titulo: 'Desde Ayacucho', sub: '3.5 horas por la carretera vía Cangallo. Ofrecemos traslados privados previa reserva.' },
              { icon: 'location_on', titulo: 'Dirección', sub: 'Jr. Túpac Amaru S/N, Plaza de Armas de Vilcashuamán, Ayacucho.' },
              { icon: 'schedule', titulo: 'Horarios', sub: 'Martes a Domingo, 12:00 pm – 9:00 pm' },
            ].map(({ icon, titulo, sub }) => (
              <div key={titulo} style={{ display: 'flex', gap: 16, marginBottom: 24 }}
                onMouseEnter={e => {
                  const div = (e.currentTarget as HTMLDivElement).querySelector('div') as HTMLDivElement
                  if (div) { div.style.backgroundColor = C.primary; div.style.color = '#fff' }
                }}
                onMouseLeave={e => {
                  const div = (e.currentTarget as HTMLDivElement).querySelector('div') as HTMLDivElement
                  if (div) { div.style.backgroundColor = `${C.primary}12`; div.style.color = C.primary }
                }}>
                <div style={{
                  width: 46, height: 46, borderRadius: '50%', flexShrink: 0,
                  backgroundColor: `${C.primary}12`, display: 'flex', alignItems: 'center', justifyContent: 'center',
                  transition: 'background-color 0.3s, color 0.3s', color: C.primary,
                }}>
                  <span className="material-symbols-outlined" style={{ fontSize: 22 }}>{icon}</span>
                </div>
                <div>
                  <h5 style={{ fontFamily: 'Manrope', fontWeight: 700, fontSize: 14, marginBottom: 4 }}>{titulo}</h5>
                  <p style={{ fontSize: 15, color: C.onVariant, lineHeight: '22px' }}>{sub}</p>
                </div>
              </div>
            ))}
            <a href="https://maps.app.goo.gl/vilcashuaman" target="_blank" rel="noreferrer"
              style={{ display: 'inline-flex', alignItems: 'center', gap: 8, color: C.primary, fontWeight: 700, fontSize: 14, borderBottom: `2px solid ${C.primary}`, paddingBottom: 4, textDecoration: 'none', marginTop: 8 }}>
              Abrir en Google Maps
              <span className="material-symbols-outlined" style={{ fontSize: 18 }}>arrow_forward</span>
            </a>
          </Reveal>

          <Reveal delay={300}>
            <div style={{ height: 400, borderRadius: 12, overflow: 'hidden', boxShadow: '0 4px 24px rgba(0,0,0,0.12)', position: 'relative' }}>
              <img src="/images/hero.png" alt="Ubicación Vilcashuamán"
                style={{ width: '100%', height: '100%', objectFit: 'cover', opacity: 0.85, transition: 'transform 2s' }}
                onMouseEnter={e => ((e.target as HTMLImageElement).style.transform = 'scale(1.04)')}
                onMouseLeave={e => ((e.target as HTMLImageElement).style.transform = 'scale(1)')} />
              <div style={{ position: 'absolute', inset: 0, display: 'flex', alignItems: 'center', justifyContent: 'center', pointerEvents: 'none' }}>
                <div style={{
                  backgroundColor: C.primary, color: '#fff',
                  padding: 16, borderRadius: '50%', boxShadow: '0 8px 32px rgba(0,0,0,0.3)',
                  animation: 'pulse 2s infinite',
                }}>
                  <span className="material-symbols-outlined" style={{ fontSize: 36 }}>location_on</span>
                </div>
              </div>
            </div>
          </Reveal>
        </div>
      </section>

      {/* ── FINAL CTA / RESERVAS ─────────────────────────────────────────── */}
      <section id="reservas" style={{ padding: '80px 24px', position: 'relative', overflow: 'hidden' }}>
        <div style={{ position: 'absolute', inset: 0, backgroundColor: C.primaryCont, opacity: 0.92, zIndex: 0 }} />
        <img src="/images/cuy-chactado.png" alt=""
          style={{ position: 'absolute', inset: 0, width: '100%', height: '100%', objectFit: 'cover', zIndex: -1, opacity: 0.4 }} />

        <div style={{ maxWidth: 1200, margin: '0 auto', position: 'relative', zIndex: 1, textAlign: 'center', color: '#ffe1d6' }}>
          <Reveal>
            <h2 style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 'clamp(32px,5vw,48px)', fontWeight: 700, marginBottom: 16, lineHeight: 1.2 }}>
              Reserva tu experiencia gastronómica en Vilcashuamán
            </h2>
          </Reveal>
          <Reveal delay={200}>
            <p style={{ fontSize: 18, lineHeight: '28px', maxWidth: 620, margin: '0 auto 40px', opacity: 0.9 }}>
              Los cupos son limitados para garantizar la exclusividad y calidad de nuestra atención. Asegura tu mesa hoy.
            </p>
          </Reveal>
          <Reveal delay={400}>
            <Link to="/reservar" style={{
              display: 'inline-block',
              backgroundColor: C.bg, color: C.primary,
              padding: '18px 52px', borderRadius: 10,
              fontFamily: 'Libre Caslon Text, serif', fontSize: 20, fontWeight: 700,
              textDecoration: 'none', boxShadow: '0 12px 40px rgba(0,0,0,0.25)',
              transition: 'transform 0.2s',
            }}
              onMouseEnter={e => ((e.currentTarget as HTMLAnchorElement).style.transform = 'scale(1.04)')}
              onMouseLeave={e => ((e.currentTarget as HTMLAnchorElement).style.transform = 'scale(1)')}>
              Hacer una Reserva
            </Link>
          </Reveal>
          <Reveal delay={500}>
            <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'center', gap: 32, marginTop: 40 }}>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <span className="material-symbols-outlined">call</span>
                <span style={{ fontFamily: 'Manrope', fontWeight: 700, fontSize: 14 }}>+51 945 984 518</span>
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                <span className="material-symbols-outlined">mail</span>
                <span style={{ fontFamily: 'Manrope', fontWeight: 700, fontSize: 14 }}>hola@vilcasuyo.com</span>
              </div>
              <a href="https://wa.me/51945984518" target="_blank" rel="noreferrer"
                style={{ display: 'flex', alignItems: 'center', gap: 8, color: '#a8f0a8', textDecoration: 'none', fontWeight: 700, fontSize: 14, fontFamily: 'Manrope' }}>
                <span className="material-symbols-outlined">chat_bubble</span>
                WhatsApp
              </a>
            </div>
          </Reveal>
        </div>
      </section>

      {/* ── FOOTER ───────────────────────────────────────────────────────── */}
      <footer style={{
        padding: '48px 24px',
        borderTop: `1px solid ${C.outlineVar}`,
        backgroundColor: '#f6f3ed',
        display: 'flex', flexDirection: 'column', alignItems: 'center', gap: 20,
      }}>
        <span style={{ fontFamily: 'Libre Caslon Text, serif', fontSize: 22, fontWeight: 700, color: C.primary }}>
          Vilca Suyo
        </span>
        <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'center', gap: 24 }}>
          {['Heritage', 'Inca Trail', 'Private Dining', 'Sostenibilidad'].map(l => (
            <a key={l} href="#" style={{ color: C.onVariant, textDecoration: 'none', fontSize: 15, transition: 'color 0.2s' }}
              onMouseEnter={e => ((e.currentTarget as HTMLAnchorElement).style.color = C.primary)}
              onMouseLeave={e => ((e.currentTarget as HTMLAnchorElement).style.color = C.onVariant)}>
              {l}
            </a>
          ))}
        </div>
        <div style={{ display: 'flex', gap: 20 }}>
          {['public', 'camera_alt', 'chat_bubble'].map(icon => (
            <a key={icon} href="#"
              style={{ color: C.primary, transition: 'transform 0.2s' }}
              onMouseEnter={e => ((e.currentTarget as HTMLAnchorElement).style.transform = 'scale(1.2)')}
              onMouseLeave={e => ((e.currentTarget as HTMLAnchorElement).style.transform = 'scale(1)')}>
              <span className="material-symbols-outlined">{icon}</span>
            </a>
          ))}
        </div>
        <p style={{ color: C.onVariant, fontSize: 14, textAlign: 'center' }}>
          © 2026 Vilca Suyo Restaurant · Crafted for the Modern Explorer ·{' '}
          <Link to="/login" style={{ color: C.outline, textDecoration: 'none' }}>Admin</Link>
        </p>
      </footer>

      {/* Bounce keyframe */}
      <style>{`
        @keyframes bounce {
          0%, 100% { transform: translateX(-50%) translateY(0); }
          50% { transform: translateX(-50%) translateY(10px); }
        }
        @keyframes pulse {
          0%, 100% { box-shadow: 0 0 0 0 rgba(130,59,24,0.4); }
          50% { box-shadow: 0 0 0 16px rgba(130,59,24,0); }
        }
        ::-webkit-scrollbar { width: 6px; }
        ::-webkit-scrollbar-track { background: #f6f3ed; }
        ::-webkit-scrollbar-thumb { background: ${C.outlineVar}; border-radius: 3px; }
      `}</style>
    </div>
  )
}
