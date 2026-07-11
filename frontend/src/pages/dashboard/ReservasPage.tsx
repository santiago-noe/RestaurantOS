import { useEffect, useState } from 'react';
import { reservasApi } from '../../services/api';
import type { Reserva } from '../../types';
import { CalendarClock, Users, Phone, Check, X as XIcon, CalendarX2 } from 'lucide-react';

const C = {
  bg: '#fefaf5',
  surface: '#ffffff',
  primary: '#b4532b',
  primaryDark: '#8b3a1a',
  secondary: '#3f6b5c',
  textPrimary: '#2c2c2a',
  textSecondary: '#6b5e55',
  borderLight: '#efe5dc',
  success: '#2e7d64',
  successBg: '#e8f5ef',
  alert: '#e53e3e',
  alertBg: '#fff5f5',
  warning: '#b38f2c',
  warningBg: '#fdf6e3',
};

const estadoStyle: Record<Reserva['estado'], { bg: string; color: string; label: string }> = {
  pendiente: { bg: C.warningBg, color: C.warning, label: 'Pendiente' },
  confirmada: { bg: C.successBg, color: C.success, label: 'Confirmada' },
  cancelada: { bg: C.alertBg, color: C.alert, label: 'Cancelada' },
};

export default function ReservasPage() {
  const [reservas, setReservas] = useState<Reserva[]>([]);
  const [total, setTotal] = useState(0);
  const [filtro, setFiltro] = useState<'' | Reserva['estado']>('');
  const [loading, setLoading] = useState(true);
  const [actualizando, setActualizando] = useState<number | null>(null);

  const cargar = async () => {
    setLoading(true);
    try {
      const { data } = await reservasApi.listar(1, filtro);
      setReservas(data.data);
      setTotal(data.total);
    } catch (error) {
      console.error('Error al cargar reservas:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    cargar();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [filtro]);

  const cambiarEstado = async (id: number, estado: Reserva['estado']) => {
    setActualizando(id);
    try {
      await reservasApi.actualizarEstado(id, estado);
      await cargar();
    } catch (err: any) {
      alert(err.response?.data?.error || 'Error al actualizar la reserva');
    } finally {
      setActualizando(null);
    }
  };

  return (
    <div
      style={{
        maxWidth: '1400px',
        margin: '0 auto',
        padding: '24px 20px',
        backgroundColor: C.bg,
        minHeight: '100vh',
        fontFamily: '"Inter", system-ui, sans-serif',
      }}
    >
      <div style={{ marginBottom: '32px' }}>
        <h1 style={{ fontSize: 'clamp(1.75rem, 5vw, 2rem)', fontWeight: 700, color: C.textPrimary, margin: '0 0 6px 0', letterSpacing: '-0.02em' }}>
          Reservas
        </h1>
        <p style={{ fontSize: '0.9rem', color: C.textSecondary, margin: 0 }}>
          {total} reservas recibidas desde la página pública
        </p>
      </div>

      {/* Filtro por estado */}
      <div style={{ display: 'flex', gap: '8px', marginBottom: '24px', flexWrap: 'wrap' }}>
        {(['', 'pendiente', 'confirmada', 'cancelada'] as const).map((valor) => (
          <button
            key={valor || 'todas'}
            onClick={() => setFiltro(valor)}
            style={{
              padding: '8px 16px',
              borderRadius: '40px',
              border: `1px solid ${filtro === valor ? C.primary : C.borderLight}`,
              backgroundColor: filtro === valor ? C.primary : C.surface,
              color: filtro === valor ? '#fff' : C.textSecondary,
              fontSize: '0.85rem',
              fontWeight: 500,
              cursor: 'pointer',
              transition: 'all 0.2s',
            }}
          >
            {valor === '' ? 'Todas' : estadoStyle[valor].label}
          </button>
        ))}
      </div>

      {loading ? (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {[...Array(4)].map((_, i) => (
            <div key={i} style={{ height: '90px', backgroundColor: C.surface, borderRadius: '20px', border: `1px solid ${C.borderLight}` }} />
          ))}
        </div>
      ) : reservas.length === 0 ? (
        <div style={{ textAlign: 'center', padding: '48px 24px', backgroundColor: C.surface, borderRadius: '28px', border: `1px solid ${C.borderLight}` }}>
          <CalendarX2 size={48} style={{ color: C.textSecondary, opacity: 0.5, marginBottom: '16px' }} />
          <p style={{ color: C.textSecondary, fontSize: '1rem', margin: 0 }}>
            No hay reservas {filtro ? `en estado "${estadoStyle[filtro].label.toLowerCase()}"` : 'registradas aún'}
          </p>
        </div>
      ) : (
        <div style={{ display: 'flex', flexDirection: 'column', gap: '12px' }}>
          {reservas.map((reserva) => {
            const estado = estadoStyle[reserva.estado];
            return (
              <div
                key={reserva.id}
                style={{
                  backgroundColor: C.surface,
                  borderRadius: '20px',
                  border: `1px solid ${C.borderLight}`,
                  padding: '16px 20px',
                  display: 'flex',
                  flexWrap: 'wrap',
                  justifyContent: 'space-between',
                  alignItems: 'center',
                  gap: '12px',
                }}
              >
                <div style={{ display: 'flex', alignItems: 'center', gap: '14px', flex: '2', minWidth: '220px' }}>
                  <div
                    style={{
                      width: '48px',
                      height: '48px',
                      borderRadius: '24px',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      backgroundColor: `${C.primary}15`,
                    }}
                  >
                    <CalendarClock size={22} style={{ color: C.primary }} />
                  </div>
                  <div>
                    <p style={{ fontWeight: 600, color: C.textPrimary, margin: 0, fontSize: '1rem' }}>
                      {reserva.nombre}
                    </p>
                    <div style={{ display: 'flex', gap: '12px', marginTop: '4px', flexWrap: 'wrap' }}>
                      <span style={{ fontSize: '0.7rem', color: C.textSecondary, display: 'flex', alignItems: 'center', gap: '4px' }}>
                        <Phone size={12} /> {reserva.whatsapp}
                      </span>
                      <span style={{ fontSize: '0.7rem', color: C.textSecondary, display: 'flex', alignItems: 'center', gap: '4px' }}>
                        <Users size={12} /> {reserva.personas} personas
                      </span>
                      <span style={{ fontSize: '0.7rem', color: C.textSecondary }}>
                        {reserva.fecha?.split('T')[0]}
                      </span>
                    </div>
                  </div>
                </div>

                <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
                  <span
                    style={{
                      backgroundColor: estado.bg,
                      color: estado.color,
                      fontWeight: 700,
                      fontSize: '0.75rem',
                      padding: '6px 12px',
                      borderRadius: '40px',
                      border: `1px solid ${estado.color}20`,
                    }}
                  >
                    {estado.label}
                  </span>

                  {reserva.estado === 'pendiente' && (
                    <>
                      <button
                        disabled={actualizando === reserva.id}
                        onClick={() => cambiarEstado(reserva.id, 'confirmada')}
                        title="Confirmar reserva"
                        style={{
                          display: 'flex', alignItems: 'center', gap: '6px',
                          backgroundColor: C.success, color: '#fff', border: 'none',
                          padding: '8px 14px', borderRadius: '40px', fontSize: '0.8rem',
                          fontWeight: 600, cursor: 'pointer', opacity: actualizando === reserva.id ? 0.6 : 1,
                        }}
                      >
                        <Check size={14} /> Confirmar
                      </button>
                      <button
                        disabled={actualizando === reserva.id}
                        onClick={() => cambiarEstado(reserva.id, 'cancelada')}
                        title="Cancelar reserva"
                        style={{
                          display: 'flex', alignItems: 'center', gap: '6px',
                          backgroundColor: 'transparent', color: C.alert, border: `1px solid ${C.alert}40`,
                          padding: '8px 14px', borderRadius: '40px', fontSize: '0.8rem',
                          fontWeight: 600, cursor: 'pointer', opacity: actualizando === reserva.id ? 0.6 : 1,
                        }}
                      >
                        <XIcon size={14} /> Cancelar
                      </button>
                    </>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
