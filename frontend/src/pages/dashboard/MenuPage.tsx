import { useEffect, useState } from 'react';
import { menuApi } from '../../services/api';
import type { MenuPublico } from '../../types';
import { Plus, X, EyeOff, Eye, Trash2, UtensilsCrossed, Pencil } from 'lucide-react';

const C = {
  bg: '#fefaf5',
  surface: '#ffffff',
  primary: '#b4532b',
  primaryDark: '#8b3a1a',
  secondary: '#3f6b5c',
  gold: '#d4af37',
  textPrimary: '#2c2c2a',
  textSecondary: '#6b5e55',
  borderLight: '#efe5dc',
  alert: '#e53e3e',
};

const formInicial = { categoria: '', nombre: '', descripcion: '', precio: '', imagen_url: '', orden: '' };

export default function MenuPage() {
  const [items, setItems] = useState<MenuPublico[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [editandoId, setEditandoId] = useState<number | null>(null);
  const [form, setForm] = useState(formInicial);
  const [guardando, setGuardando] = useState(false);

  const cargar = async () => {
    setLoading(true);
    try {
      const { data } = await menuApi.listarAdmin();
      setItems(data);
    } catch (error) {
      console.error('Error al cargar el menú:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    cargar();
  }, []);

  const abrirCrear = () => {
    setEditandoId(null);
    setForm(formInicial);
    setShowForm(true);
  };

  const abrirEditar = (item: MenuPublico) => {
    setEditandoId(item.id);
    setForm({
      categoria: item.categoria,
      nombre: item.nombre,
      descripcion: item.descripcion,
      precio: String(item.precio),
      imagen_url: item.imagen_url,
      orden: String(item.orden),
    });
    setShowForm(true);
  };

  const handleGuardar = async (e: React.FormEvent) => {
    e.preventDefault();
    setGuardando(true);
    const payload = {
      categoria: form.categoria,
      nombre: form.nombre,
      descripcion: form.descripcion,
      precio: Number(form.precio) || 0,
      imagen_url: form.imagen_url,
      orden: Number(form.orden) || 0,
    };
    try {
      if (editandoId !== null) {
        await menuApi.actualizar(editandoId, payload);
      } else {
        await menuApi.crear(payload);
      }
      setShowForm(false);
      setEditandoId(null);
      setForm(formInicial);
      cargar();
    } catch (err: any) {
      alert(err.response?.data?.error || `Error al ${editandoId !== null ? 'actualizar' : 'crear'} el item del menú`);
    } finally {
      setGuardando(false);
    }
  };

  const toggleDisponible = async (item: MenuPublico) => {
    try {
      await menuApi.actualizar(item.id, { disponible: !item.disponible });
      cargar();
    } catch (err: any) {
      alert(err.response?.data?.error || 'Error al actualizar el item');
    }
  };

  const eliminar = async (item: MenuPublico) => {
    if (!confirm(`¿Eliminar "${item.nombre}" del menú?`)) return;
    try {
      await menuApi.eliminar(item.id);
      cargar();
    } catch (err: any) {
      alert(err.response?.data?.error || 'Error al eliminar el item');
    }
  };

  return (
    <div style={{ maxWidth: '1400px', margin: '0 auto', padding: '24px 20px', backgroundColor: C.bg, minHeight: '100vh', fontFamily: '"Inter", system-ui, sans-serif' }}>

      <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'space-between', alignItems: 'center', marginBottom: '32px', gap: '16px' }}>
        <div>
          <h1 style={{ fontSize: 'clamp(1.75rem, 5vw, 2rem)', fontWeight: 700, color: C.textPrimary, margin: '0 0 6px 0', letterSpacing: '-0.02em' }}>
            Menú público
          </h1>
          <p style={{ fontSize: '0.9rem', color: C.textSecondary, margin: 0 }}>
            {items.length} platos • visible en la página de inicio
          </p>
        </div>
        <button
          onClick={abrirCrear}
          style={{
            display: 'flex', alignItems: 'center', gap: '8px', backgroundColor: C.primary, color: '#fff',
            border: 'none', padding: '10px 20px', borderRadius: '40px', fontSize: '0.9rem', fontWeight: 500,
            cursor: 'pointer', transition: 'all 0.2s ease', boxShadow: '0 2px 6px rgba(0,0,0,0.05)',
          }}
          onMouseEnter={(e) => { e.currentTarget.style.backgroundColor = C.primaryDark; }}
          onMouseLeave={(e) => { e.currentTarget.style.backgroundColor = C.primary; }}
        >
          <Plus size={18} /> Nuevo plato
        </button>
      </div>

      {showForm && (
        <div
          style={{ position: 'fixed', inset: 0, backgroundColor: 'rgba(0,0,0,0.5)', backdropFilter: 'blur(4px)', zIndex: 1000, display: 'flex', alignItems: 'center', justifyContent: 'center', padding: '20px' }}
          onClick={(e) => { if (e.target === e.currentTarget) setShowForm(false); }}
        >
          <div style={{ backgroundColor: C.surface, borderRadius: '28px', maxWidth: '500px', width: '100%', boxShadow: '0 25px 40px -12px rgba(0,0,0,0.3)', overflow: 'hidden' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', padding: '20px 24px', borderBottom: `1px solid ${C.borderLight}` }}>
              <h2 style={{ fontSize: '1.4rem', fontWeight: 600, color: C.textPrimary, margin: 0 }}>{editandoId !== null ? 'Editar plato' : 'Nuevo plato'}</h2>
              <button onClick={() => setShowForm(false)} style={{ background: 'none', border: 'none', cursor: 'pointer', color: C.textSecondary, padding: '4px', display: 'flex' }}>
                <X size={20} />
              </button>
            </div>
            <form onSubmit={handleGuardar} style={{ padding: '24px' }}>
              <div style={{ display: 'grid', gap: '20px' }}>
                <div>
                  <label style={labelStyle}>Categoría *</label>
                  <input required value={form.categoria} onChange={(e) => setForm((p) => ({ ...p, categoria: e.target.value }))} style={inputStyle} placeholder="Entradas, Fondos, Postres..." />
                </div>
                <div>
                  <label style={labelStyle}>Nombre *</label>
                  <input required value={form.nombre} onChange={(e) => setForm((p) => ({ ...p, nombre: e.target.value }))} style={inputStyle} />
                </div>
                <div>
                  <label style={labelStyle}>Descripción</label>
                  <input value={form.descripcion} onChange={(e) => setForm((p) => ({ ...p, descripcion: e.target.value }))} style={inputStyle} />
                </div>
                <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '16px' }}>
                  <div>
                    <label style={labelStyle}>Precio (S/)</label>
                    <input type="number" step="0.01" min="0" value={form.precio} onChange={(e) => setForm((p) => ({ ...p, precio: e.target.value }))} style={inputStyle} />
                  </div>
                  <div>
                    <label style={labelStyle}>Orden</label>
                    <input type="number" value={form.orden} onChange={(e) => setForm((p) => ({ ...p, orden: e.target.value }))} style={inputStyle} placeholder="0" />
                  </div>
                </div>
                <div>
                  <label style={labelStyle}>URL de imagen</label>
                  <input value={form.imagen_url} onChange={(e) => setForm((p) => ({ ...p, imagen_url: e.target.value }))} style={inputStyle} placeholder="/images/plato.png" />
                </div>
              </div>
              <div style={{ display: 'flex', gap: '12px', marginTop: '32px' }}>
                <button type="button" onClick={() => setShowForm(false)} style={{ flex: 1, padding: '12px', borderRadius: '40px', border: `1px solid ${C.borderLight}`, backgroundColor: 'transparent', color: C.textSecondary, fontWeight: 500, cursor: 'pointer' }}>
                  Cancelar
                </button>
                <button type="submit" disabled={guardando} style={{ flex: 1, padding: '12px', borderRadius: '40px', border: 'none', backgroundColor: C.primary, color: '#fff', fontWeight: 500, cursor: 'pointer', opacity: guardando ? 0.7 : 1 }}>
                  {guardando ? 'Guardando...' : editandoId !== null ? 'Guardar cambios' : 'Crear plato'}
                </button>
              </div>
            </form>
          </div>
        </div>
      )}

      {loading ? (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {[...Array(4)].map((_, i) => (
            <div key={i} style={{ height: '80px', backgroundColor: C.surface, borderRadius: '20px', border: `1px solid ${C.borderLight}` }} />
          ))}
        </div>
      ) : items.length === 0 ? (
        <div style={{ textAlign: 'center', padding: '48px 24px', backgroundColor: C.surface, borderRadius: '28px', border: `1px solid ${C.borderLight}` }}>
          <UtensilsCrossed size={48} style={{ color: C.textSecondary, opacity: 0.5, marginBottom: '16px' }} />
          <p style={{ color: C.textSecondary, fontSize: '1rem', margin: 0 }}>Aún no hay platos en el menú público</p>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {items.map((item) => (
            <div key={item.id} style={{ backgroundColor: C.surface, borderRadius: '20px', border: `1px solid ${C.borderLight}`, padding: '16px 20px', display: 'flex', flexWrap: 'wrap', justifyContent: 'space-between', alignItems: 'center', gap: '12px', opacity: item.disponible ? 1 : 0.6 }}>
              <div style={{ flex: '2', minWidth: '180px' }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                  <p style={{ fontWeight: 600, color: C.textPrimary, margin: 0, fontSize: '1rem' }}>{item.nombre}</p>
                  <span style={{ fontSize: '0.7rem', fontWeight: 500, color: C.secondary, backgroundColor: `${C.secondary}15`, padding: '2px 10px', borderRadius: '40px' }}>
                    {item.categoria}
                  </span>
                </div>
                {item.descripcion && <p style={{ fontSize: '0.8rem', color: C.textSecondary, margin: '4px 0 0' }}>{item.descripcion}</p>}
              </div>
              <div style={{ display: 'flex', alignItems: 'center', gap: '16px' }}>
                <span style={{ fontWeight: 700, color: C.gold, fontSize: '0.95rem' }}>S/ {item.precio.toFixed(2)}</span>
                <button onClick={() => abrirEditar(item)} title="Editar"
                  style={{ background: 'none', border: `1px solid ${C.borderLight}`, borderRadius: '50%', width: 36, height: 36, display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', color: C.primary }}>
                  <Pencil size={16} />
                </button>
                <button onClick={() => toggleDisponible(item)} title={item.disponible ? 'Marcar como agotado' : 'Marcar como disponible'}
                  style={{ background: 'none', border: `1px solid ${C.borderLight}`, borderRadius: '50%', width: 36, height: 36, display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', color: item.disponible ? C.secondary : C.textSecondary }}>
                  {item.disponible ? <Eye size={16} /> : <EyeOff size={16} />}
                </button>
                <button onClick={() => eliminar(item)} title="Eliminar"
                  style={{ background: 'none', border: `1px solid ${C.borderLight}`, borderRadius: '50%', width: 36, height: 36, display: 'flex', alignItems: 'center', justifyContent: 'center', cursor: 'pointer', color: C.alert }}>
                  <Trash2 size={16} />
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}

const labelStyle: React.CSSProperties = {
  display: 'block', fontSize: '0.75rem', fontWeight: 500, color: C.textSecondary,
  marginBottom: '6px', textTransform: 'uppercase', letterSpacing: '0.5px',
};

const inputStyle: React.CSSProperties = {
  width: '100%', padding: '10px 14px', border: `1px solid ${C.borderLight}`, borderRadius: '16px',
  fontSize: '0.9rem', backgroundColor: C.surface, color: C.textPrimary, outline: 'none',
};
