import { useEffect, useState } from 'react';
import { reportesApi } from '../../services/api';
import type { Cliente, ReporteMovimientos, ReporteVentas } from '../../types';
import { FileDown, FileSpreadsheet, TrendingUp, Users, Package } from 'lucide-react';

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
  alertBg: '#fff5f5',
};

function descargarArchivo(blob: Blob, filename: string) {
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = filename;
  a.click();
  window.URL.revokeObjectURL(url);
}

type Tab = 'ventas' | 'deudores' | 'inventario';

const tabs: { id: Tab; label: string; icon: typeof TrendingUp }[] = [
  { id: 'ventas', label: 'Ventas', icon: TrendingUp },
  { id: 'deudores', label: 'Clientes con deuda', icon: Users },
  { id: 'inventario', label: 'Movimientos de inventario', icon: Package },
];

export default function ReportesPage() {
  const [tab, setTab] = useState<Tab>('ventas');

  return (
    <div style={{ maxWidth: '1400px', margin: '0 auto', padding: '24px 20px', backgroundColor: C.bg, minHeight: '100vh', fontFamily: '"Inter", system-ui, sans-serif' }}>
      <div style={{ marginBottom: '28px' }}>
        <h1 style={{ fontSize: 'clamp(1.75rem, 5vw, 2rem)', fontWeight: 700, color: C.textPrimary, margin: '0 0 6px 0', letterSpacing: '-0.02em' }}>
          Reportes
        </h1>
        <p style={{ fontSize: '0.9rem', color: C.textSecondary, margin: 0 }}>Ventas, deuda de clientes y movimientos de inventario</p>
      </div>

      <div style={{ display: 'flex', gap: '8px', marginBottom: '28px', flexWrap: 'wrap' }}>
        {tabs.map(({ id, label, icon: Icon }) => (
          <button key={id} onClick={() => setTab(id)}
            style={{
              display: 'flex', alignItems: 'center', gap: '8px', padding: '10px 18px', borderRadius: '40px',
              border: `1px solid ${tab === id ? C.primary : C.borderLight}`,
              backgroundColor: tab === id ? C.primary : C.surface,
              color: tab === id ? '#fff' : C.textSecondary,
              fontSize: '0.85rem', fontWeight: 600, cursor: 'pointer', transition: 'all 0.2s',
            }}>
            <Icon size={16} /> {label}
          </button>
        ))}
      </div>

      {tab === 'ventas' && <ReporteVentasTab />}
      {tab === 'deudores' && <ReporteDeudoresTab />}
      {tab === 'inventario' && <ReporteInventarioTab />}
    </div>
  );
}

function BotonesDescarga({ onPDF, onExcel }: { onPDF: () => void; onExcel: () => void }) {
  return (
    <div style={{ display: 'flex', gap: '10px' }}>
      <button onClick={onPDF} style={botonSecundario}>
        <FileDown size={16} /> PDF
      </button>
      <button onClick={onExcel} style={botonSecundario}>
        <FileSpreadsheet size={16} /> Excel
      </button>
    </div>
  );
}

// ─── Ventas ───────────────────────────────────────────────────────────────────

function ReporteVentasTab() {
  const [periodo, setPeriodo] = useState<'diario' | 'semanal' | 'mensual'>('diario');
  const [fecha, setFecha] = useState('');
  const [reporte, setReporte] = useState<ReporteVentas | null>(null);
  const [loading, setLoading] = useState(true);

  const cargar = async () => {
    setLoading(true);
    try {
      const { data } = await reportesApi.ventas(periodo, fecha || undefined);
      setReporte(data);
    } catch (error) {
      console.error('Error al cargar el reporte de ventas:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { cargar(); }, [periodo, fecha]);

  const descargar = async (formato: 'pdf' | 'excel') => {
    try {
      const { data } = await reportesApi.ventasArchivo(periodo, formato, fecha || undefined);
      descargarArchivo(data, `reporte-ventas.${formato === 'pdf' ? 'pdf' : 'xlsx'}`);
    } catch (err: any) {
      alert('Error al descargar el reporte');
    }
  };

  return (
    <div style={{ backgroundColor: C.surface, borderRadius: '20px', border: `1px solid ${C.borderLight}`, padding: '24px' }}>
      <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'space-between', alignItems: 'center', gap: '16px', marginBottom: '20px' }}>
        <div style={{ display: 'flex', gap: '10px', flexWrap: 'wrap' }}>
          <select value={periodo} onChange={(e) => setPeriodo(e.target.value as any)} style={selectStyle}>
            <option value="diario">Diario</option>
            <option value="semanal">Semanal</option>
            <option value="mensual">Mensual</option>
          </select>
          <input type="date" value={fecha} onChange={(e) => setFecha(e.target.value)} style={selectStyle} />
        </div>
        <BotonesDescarga onPDF={() => descargar('pdf')} onExcel={() => descargar('excel')} />
      </div>

      {loading ? (
        <p style={{ color: C.textSecondary }}>Cargando...</p>
      ) : !reporte ? (
        <p style={{ color: C.textSecondary }}>No se pudo cargar el reporte</p>
      ) : (
        <>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '16px', marginBottom: '24px' }}>
            <StatTile label="Total ventas" valor={`S/ ${reporte.total_ventas.toFixed(2)}`} color={C.gold} />
            <StatTile label="Total pedidos" valor={String(reporte.total_pedidos)} color={C.secondary} />
            <StatTile label="Periodo" valor={`${reporte.desde} a ${reporte.hasta}`} color={C.primary} />
          </div>

          {reporte.por_dia.length === 0 ? (
            <p style={{ color: C.textSecondary, textAlign: 'center', padding: '24px' }}>Sin ventas en este periodo</p>
          ) : (
            <div style={{ overflowX: 'auto' }}>
              <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.9rem' }}>
                <thead>
                  <tr style={{ borderBottom: `2px solid ${C.borderLight}`, textAlign: 'left' }}>
                    <th style={thStyle}>Fecha</th>
                    <th style={thStyle}>Pedidos</th>
                    <th style={thStyle}>Total (S/)</th>
                  </tr>
                </thead>
                <tbody>
                  {reporte.por_dia.map((d) => (
                    <tr key={d.fecha} style={{ borderBottom: `1px solid ${C.borderLight}` }}>
                      <td style={tdStyle}>{d.fecha}</td>
                      <td style={tdStyle}>{d.cantidad_pedidos}</td>
                      <td style={tdStyle}>{d.total.toFixed(2)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}
    </div>
  );
}

// ─── Deudores ─────────────────────────────────────────────────────────────────

function ReporteDeudoresTab() {
  const [clientes, setClientes] = useState<Cliente[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    (async () => {
      setLoading(true);
      try {
        const { data } = await reportesApi.deudores();
        setClientes(data.deudores);
      } catch (error) {
        console.error('Error al cargar deudores:', error);
      } finally {
        setLoading(false);
      }
    })();
  }, []);

  const descargar = async (formato: 'pdf' | 'excel') => {
    try {
      const { data } = await reportesApi.deudoresArchivo(formato);
      descargarArchivo(data, `reporte-deudores.${formato === 'pdf' ? 'pdf' : 'xlsx'}`);
    } catch {
      alert('Error al descargar el reporte');
    }
  };

  const totalDeuda = clientes.reduce((acc, c) => acc + c.deuda_total, 0);

  return (
    <div style={{ backgroundColor: C.surface, borderRadius: '20px', border: `1px solid ${C.borderLight}`, padding: '24px' }}>
      <div style={{ display: 'flex', justifyContent: 'flex-end', marginBottom: '20px' }}>
        <BotonesDescarga onPDF={() => descargar('pdf')} onExcel={() => descargar('excel')} />
      </div>

      {loading ? (
        <p style={{ color: C.textSecondary }}>Cargando...</p>
      ) : clientes.length === 0 ? (
        <p style={{ color: C.textSecondary, textAlign: 'center', padding: '24px' }}>No hay clientes con deuda pendiente</p>
      ) : (
        <>
          <div style={{ marginBottom: '20px' }}>
            <StatTile label="Deuda total" valor={`S/ ${totalDeuda.toFixed(2)}`} color={C.alert} />
          </div>
          <div style={{ overflowX: 'auto' }}>
            <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.9rem' }}>
              <thead>
                <tr style={{ borderBottom: `2px solid ${C.borderLight}`, textAlign: 'left' }}>
                  <th style={thStyle}>Nombre</th>
                  <th style={thStyle}>Tipo</th>
                  <th style={thStyle}>Teléfono</th>
                  <th style={thStyle}>Deuda (S/)</th>
                </tr>
              </thead>
              <tbody>
                {clientes.map((c) => (
                  <tr key={c.id} style={{ borderBottom: `1px solid ${C.borderLight}` }}>
                    <td style={tdStyle}>{c.nombre} {c.apellido}</td>
                    <td style={tdStyle}>{c.tipo === 'empresa' ? 'Empresa' : 'Individual'}</td>
                    <td style={tdStyle}>{c.telefono || '—'}</td>
                    <td style={{ ...tdStyle, color: C.alert, fontWeight: 700 }}>{c.deuda_total.toFixed(2)}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </>
      )}
    </div>
  );
}

// ─── Inventario ───────────────────────────────────────────────────────────────

function ReporteInventarioTab() {
  const [desde, setDesde] = useState('');
  const [hasta, setHasta] = useState('');
  const [reporte, setReporte] = useState<ReporteMovimientos | null>(null);
  const [loading, setLoading] = useState(true);

  const cargar = async () => {
    setLoading(true);
    try {
      const { data } = await reportesApi.inventario(desde || undefined, hasta || undefined);
      setReporte(data);
    } catch (error) {
      console.error('Error al cargar movimientos:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { cargar(); }, []);

  const descargar = async (formato: 'pdf' | 'excel') => {
    try {
      const { data } = await reportesApi.inventarioArchivo(formato, desde || undefined, hasta || undefined);
      descargarArchivo(data, `reporte-inventario.${formato === 'pdf' ? 'pdf' : 'xlsx'}`);
    } catch {
      alert('Error al descargar el reporte');
    }
  };

  return (
    <div style={{ backgroundColor: C.surface, borderRadius: '20px', border: `1px solid ${C.borderLight}`, padding: '24px' }}>
      <div style={{ display: 'flex', flexWrap: 'wrap', justifyContent: 'space-between', alignItems: 'center', gap: '16px', marginBottom: '20px' }}>
        <div style={{ display: 'flex', gap: '10px', flexWrap: 'wrap', alignItems: 'center' }}>
          <input type="date" value={desde} onChange={(e) => setDesde(e.target.value)} style={selectStyle} />
          <span style={{ color: C.textSecondary, fontSize: '0.85rem' }}>a</span>
          <input type="date" value={hasta} onChange={(e) => setHasta(e.target.value)} style={selectStyle} />
          <button onClick={cargar} style={botonSecundario}>Filtrar</button>
        </div>
        <BotonesDescarga onPDF={() => descargar('pdf')} onExcel={() => descargar('excel')} />
      </div>

      {loading ? (
        <p style={{ color: C.textSecondary }}>Cargando...</p>
      ) : !reporte ? (
        <p style={{ color: C.textSecondary }}>No se pudo cargar el reporte</p>
      ) : (
        <>
          <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(180px, 1fr))', gap: '16px', marginBottom: '24px' }}>
            <StatTile label="Total entradas" valor={reporte.total_entradas.toFixed(2)} color={C.secondary} />
            <StatTile label="Total salidas" valor={reporte.total_salidas.toFixed(2)} color={C.alert} />
            <StatTile label="Periodo" valor={`${reporte.desde} a ${reporte.hasta}`} color={C.primary} />
          </div>

          {reporte.movimientos.length === 0 ? (
            <p style={{ color: C.textSecondary, textAlign: 'center', padding: '24px' }}>Sin movimientos en este periodo</p>
          ) : (
            <div style={{ overflowX: 'auto' }}>
              <table style={{ width: '100%', borderCollapse: 'collapse', fontSize: '0.9rem' }}>
                <thead>
                  <tr style={{ borderBottom: `2px solid ${C.borderLight}`, textAlign: 'left' }}>
                    <th style={thStyle}>Fecha</th>
                    <th style={thStyle}>Producto</th>
                    <th style={thStyle}>Tipo</th>
                    <th style={thStyle}>Cantidad</th>
                  </tr>
                </thead>
                <tbody>
                  {reporte.movimientos.map((m) => (
                    <tr key={m.id} style={{ borderBottom: `1px solid ${C.borderLight}` }}>
                      <td style={tdStyle}>{m.fecha}</td>
                      <td style={tdStyle}>{m.producto?.nombre || '—'}</td>
                      <td style={{ ...tdStyle, color: m.tipo === 'entrada' ? C.secondary : C.alert, fontWeight: 600 }}>
                        {m.tipo === 'entrada' ? 'Entrada' : 'Salida'}
                      </td>
                      <td style={tdStyle}>{m.cantidad}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </>
      )}
    </div>
  );
}

function StatTile({ label, valor, color }: { label: string; valor: string; color: string }) {
  return (
    <div style={{ backgroundColor: `${color}10`, border: `1px solid ${color}30`, borderRadius: '16px', padding: '16px' }}>
      <p style={{ fontSize: '0.75rem', color: C.textSecondary, textTransform: 'uppercase', letterSpacing: '0.05em', margin: '0 0 6px' }}>{label}</p>
      <p style={{ fontSize: '1.3rem', fontWeight: 700, color, margin: 0 }}>{valor}</p>
    </div>
  );
}

const selectStyle: React.CSSProperties = {
  padding: '10px 14px', border: `1px solid ${C.borderLight}`, borderRadius: '16px',
  fontSize: '0.85rem', backgroundColor: C.surface, color: C.textPrimary, outline: 'none',
};

const botonSecundario: React.CSSProperties = {
  display: 'flex', alignItems: 'center', gap: '6px', padding: '10px 16px', borderRadius: '40px',
  border: `1px solid ${C.borderLight}`, backgroundColor: C.surface, color: C.textPrimary,
  fontSize: '0.85rem', fontWeight: 600, cursor: 'pointer',
};

const thStyle: React.CSSProperties = { padding: '10px 12px', color: C.textSecondary, fontWeight: 600, fontSize: '0.8rem', textTransform: 'uppercase', letterSpacing: '0.03em' };
const tdStyle: React.CSSProperties = { padding: '10px 12px', color: C.textPrimary };
