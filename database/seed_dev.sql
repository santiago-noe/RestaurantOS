-- Seed para desarrollo y pruebas
-- Ejecutar DESPUÉS de que el servidor haya corrido al menos una vez (AutoMigrate)

-- Limpiar datos previos (en orden por FK)
TRUNCATE pagos, pedido_items, pedidos, movimiento_stocks, menu_publicos, productos, clientes, users RESTART IDENTITY CASCADE;

-- Usuarios de prueba (contraseña: "password123")
INSERT INTO users (nombre, apellido, email, password, rol) VALUES
('Admin',  'Sistema', 'admin@restaurante.com',  '$2a$10$/WSTt/70w4Cv1B/7WRw9ZuOiEImMhCqxPOGoJsXPN4ljru5KHSc6a', 'admin'),
('Carlos', 'Gomez',   'carlos@restaurante.com', '$2a$10$/WSTt/70w4Cv1B/7WRw9ZuOiEImMhCqxPOGoJsXPN4ljru5KHSc6a', 'empleado');

-- Clientes
INSERT INTO clientes (nombre, apellido, tipo, telefono, direccion) VALUES
('Juan',         'Pérez',    'individual', '987001001', 'Jr. Las Flores 123'),
('María',        'López',    'individual', '987001002', 'Av. Principal 456'),
('Pedro',        'Sánchez',  'individual', '987001003', NULL),
('Constructora', 'Norte SAC','empresa',    '01-234567',  'Av. Industrial 789'),
('Obra',         'Central',  'empresa',    '987001005', 'Sector 3, Lote 12');

-- Productos (inventario)
INSERT INTO productos (nombre, unidad, stock_actual, stock_minimo, precio_venta) VALUES
('Arroz',        'kg',     15.0,  5.0,  3.50),
('Pollo',        'kg',      8.5,  3.0, 12.00),
('Papa',         'kg',     20.0,  5.0,  2.00),
('Tomate',       'kg',      6.0,  2.0,  3.00),
('Aceite',       'litro',   0.5,  2.0,  8.00),
('Pan',          'unidad', 40.0, 10.0,  0.50),
('Huevo',        'unidad', 60.0, 20.0,  0.80),
('Jugo naranja', 'litro',   5.0,  2.0,  5.00);

-- Menú público
INSERT INTO menu_publicos (categoria, nombre, descripcion, precio, disponible, orden) VALUES
('Desayunos',  'Desayuno Completo',  'Pan, huevo frito, jugo natural y café', 8.00,  true, 1),
('Desayunos',  'Sánguche de Pollo',  'Pan artesanal con pollo a la plancha',  5.00,  true, 2),
('Almuerzos',  'Menú del Día',       'Sopa + segundo + refresco',            12.00, true, 1),
('Almuerzos',  'Almuerzo Especial',  'Segundo doble + ensalada + refresco',  16.00, true, 2),
('Bebidas',    'Refresco Natural',   'Chicha morada o limonada',              3.00,  true, 1),
('Bebidas',    'Jugo de Naranja',    'Jugo natural de naranja',               4.00,  true, 2),
('Promociones','Combo Familiar x4',  '4 menús del día con descuento',        40.00, true, 1);

-- Pedidos
INSERT INTO pedidos (cliente_id, user_id, fecha, tipo_comida, estado, forma_pago, total) VALUES
(1, 2, CURRENT_DATE - 5, 'almuerzo', 'entregado', 'credito', 12.00),
(1, 2, CURRENT_DATE - 3, 'almuerzo', 'entregado', 'credito', 24.00),
(2, 2, CURRENT_DATE - 4, 'desayuno', 'entregado', 'contado',  8.00),
(3, 2, CURRENT_DATE - 2, 'almuerzo', 'entregado', 'credito', 16.00),
(4, 1, CURRENT_DATE - 1, 'almuerzo', 'entregado', 'credito', 48.00),
(5, 1, CURRENT_DATE,     'almuerzo', 'pendiente', 'credito', 60.00);

-- Items con subtotal calculado explícitamente
INSERT INTO pedido_items (pedido_id, producto_id, cantidad, precio_unitario, subtotal) VALUES
(1, 3, 1, 12.00, 12.00),
(2, 3, 2, 12.00, 24.00),
(3, 1, 1,  8.00,  8.00),
(4, 4, 1, 16.00, 16.00),
(5, 3, 4, 12.00, 48.00),
(6, 3, 5, 12.00, 60.00);

-- Pago parcial de Juan Pérez
INSERT INTO pagos (cliente_id, monto, metodo, fecha) VALUES
(1, 20.00, 'efectivo', CURRENT_DATE - 1);

-- Actualizar deuda_total de clientes con crédito (pedidos crédito - pagos)
UPDATE clientes SET deuda_total = 16.00 WHERE id = 1;  -- 36 - 20 pagado
UPDATE clientes SET deuda_total = 16.00 WHERE id = 3;
UPDATE clientes SET deuda_total = 48.00 WHERE id = 4;
UPDATE clientes SET deuda_total = 60.00 WHERE id = 5;