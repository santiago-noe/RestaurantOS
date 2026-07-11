package models

import "time"

type User struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Nombre    string    `gorm:"size:100;not null" json:"nombre"`
	Apellido  string    `gorm:"size:100;not null" json:"apellido"`
	Email     string    `gorm:"size:150;uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"size:255;not null" json:"-"`
	Rol       string    `gorm:"size:20;not null;check:rol IN ('admin','empleado')" json:"rol"`
	Activo    bool      `gorm:"default:true" json:"activo"`
	CreatedAt time.Time `json:"created_at"`
}

type Cliente struct {
	ID         int       `gorm:"primaryKey" json:"id"`
	Nombre     string    `gorm:"size:100;not null" json:"nombre"`
	Apellido   string    `gorm:"size:100" json:"apellido"`
	Tipo       string    `gorm:"size:20;not null;check:tipo IN ('individual','empresa')" json:"tipo"`
	Telefono   string    `gorm:"size:20" json:"telefono"`
	Direccion  string    `json:"direccion"`
	Email      string    `gorm:"size:150" json:"email"`
	DeudaTotal float64   `gorm:"type:decimal(10,2);default:0" json:"deuda_total"`
	Activo     bool      `gorm:"default:true" json:"activo"`
	CreatedAt  time.Time `json:"created_at"`
}

type Producto struct {
	ID           int       `gorm:"primaryKey" json:"id"`
	Nombre       string    `gorm:"size:150;not null" json:"nombre"`
	Unidad       string    `gorm:"size:20;not null" json:"unidad"`
	StockActual  float64   `gorm:"type:decimal(10,3);default:0" json:"stock_actual"`
	StockMinimo  float64   `gorm:"type:decimal(10,3);default:0" json:"stock_minimo"`
	PrecioVenta  float64   `gorm:"type:decimal(10,2);default:0" json:"precio_venta"`
	Activo       bool      `gorm:"default:true" json:"activo"`
	CreatedAt    time.Time `json:"created_at"`
}

type Pedido struct {
	ID         int          `gorm:"primaryKey" json:"id"`
	ClienteID  int          `gorm:"not null" json:"cliente_id"`
	UserID     int          `gorm:"not null" json:"user_id"`
	Fecha      time.Time    `gorm:"type:date;not null" json:"fecha"`
	TipoComida string       `gorm:"size:30;not null" json:"tipo_comida"`
	Estado     string       `gorm:"size:20;not null;default:'pendiente'" json:"estado"`
	FormaPago  string       `gorm:"size:20;not null" json:"forma_pago"`
	Total      float64      `gorm:"type:decimal(10,2);not null;default:0" json:"total"`
	Notas      string       `json:"notas"`
	CreatedAt  time.Time    `json:"created_at"`
	Cliente    Cliente      `gorm:"foreignKey:ClienteID" json:"cliente,omitempty"`
	User       User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Items      []PedidoItem `gorm:"foreignKey:PedidoID" json:"items,omitempty"`
}

type PedidoItem struct {
	ID             int      `gorm:"primaryKey" json:"id"`
	PedidoID       int      `gorm:"not null" json:"pedido_id"`
	ProductoID     int      `gorm:"not null" json:"producto_id"`
	Cantidad       float64  `gorm:"type:decimal(10,3);not null" json:"cantidad"`
	PrecioUnitario float64  `gorm:"type:decimal(10,2);not null" json:"precio_unitario"`
	Subtotal       float64  `gorm:"type:decimal(10,2);not null" json:"subtotal"`
	Producto       Producto `gorm:"foreignKey:ProductoID" json:"producto,omitempty"`
}

type Pago struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	ClienteID int       `gorm:"not null" json:"cliente_id"`
	PedidoID  *int      `json:"pedido_id"`
	Monto     float64   `gorm:"type:decimal(10,2);not null" json:"monto"`
	Metodo    string    `gorm:"size:20;not null" json:"metodo"`
	Fecha     time.Time `gorm:"type:date;not null" json:"fecha"`
	CreatedAt time.Time `json:"created_at"`
}

type MovimientoStock struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	ProductoID  int       `gorm:"not null" json:"producto_id"`
	Tipo        string    `gorm:"size:10;not null" json:"tipo"`
	Cantidad    float64   `gorm:"type:decimal(10,3);not null" json:"cantidad"`
	ReferenciaID *int     `json:"referencia_id"`
	Notas       string    `json:"notas"`
	Fecha       time.Time `gorm:"type:date;not null" json:"fecha"`
	CreatedAt   time.Time `json:"created_at"`
	Producto    Producto  `gorm:"foreignKey:ProductoID" json:"producto,omitempty"`
}

type Reserva struct {
	ID        int       `gorm:"primaryKey" json:"id"`
	Nombre    string    `gorm:"size:100;not null" json:"nombre"`
	Whatsapp  string    `gorm:"size:20;not null" json:"whatsapp"`
	Fecha     time.Time `gorm:"type:date;not null" json:"fecha"`
	Personas  string    `gorm:"size:10;not null" json:"personas"`
	Ocasion   string    `gorm:"size:50" json:"ocasion"`
	Estado    string    `gorm:"size:20;not null;default:'pendiente';check:estado IN ('pendiente','confirmada','cancelada')" json:"estado"`
	PedidoID  *int      `json:"pedido_id"`
	Pedido    *Pedido   `gorm:"foreignKey:PedidoID" json:"pedido,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type MenuPublico struct {
	ID          int    `gorm:"primaryKey" json:"id"`
	Categoria   string `gorm:"size:50;not null" json:"categoria"`
	Nombre      string `gorm:"size:150;not null" json:"nombre"`
	Descripcion string `json:"descripcion"`
	Precio      float64 `gorm:"type:decimal(10,2)" json:"precio"`
	ImagenURL   string `gorm:"size:500" json:"imagen_url"`
	Disponible  bool   `json:"disponible"`
	Orden       int    `gorm:"default:0" json:"orden"`
	ProductoID  *int      `json:"producto_id"`
	Producto    *Producto `gorm:"foreignKey:ProductoID" json:"producto,omitempty"`
}
