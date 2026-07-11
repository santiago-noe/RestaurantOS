package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"restaurantos/internal/middleware"
	"restaurantos/internal/models"
	"restaurantos/internal/repository"
)

// ─── Mocks ────────────────────────────────────────────────────────────────────

type mockPagoRepo struct{ mock.Mock }

func (m *mockPagoRepo) Create(p *models.Pago) error { return m.Called(p).Error(0) }
func (m *mockPagoRepo) FindByCliente(id int) ([]models.Pago, error) {
	args := m.Called(id)
	return args.Get(0).([]models.Pago), args.Error(1)
}
func (m *mockPagoRepo) SumaPagadoPorCliente(id int) (float64, error) {
	args := m.Called(id)
	return args.Get(0).(float64), args.Error(1)
}

type mockClienteRepoC struct{ mock.Mock }

func (m *mockClienteRepoC) Create(c *models.Cliente) error { return m.Called(c).Error(0) }
func (m *mockClienteRepoC) FindByID(id int) (*models.Cliente, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cliente), args.Error(1)
}
func (m *mockClienteRepoC) FindAll(p, pp int, t string) ([]models.Cliente, int64, error) {
	args := m.Called(p, pp, t)
	return args.Get(0).([]models.Cliente), args.Get(1).(int64), args.Error(2)
}
func (m *mockClienteRepoC) Update(id int, f map[string]interface{}) (*models.Cliente, error) {
	args := m.Called(id, f)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Cliente), args.Error(1)
}
func (m *mockClienteRepoC) Deactivate(id int) error          { return m.Called(id).Error(0) }
func (m *mockClienteRepoC) EmailExists(e string, id int) bool { return m.Called(e, id).Bool(0) }

var _ repository.PagoRepo   = (*mockPagoRepo)(nil)
var _ repository.ClienteRepo = (*mockClienteRepoC)(nil)

// ─── Router ───────────────────────────────────────────────────────────────────

func routerCreditos(pr repository.PagoRepo, cr repository.ClienteRepo) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	h := NewCreditosHandler(pr, cr)
	const secret = "test-secret"

	adm := r.Group("/api/admin", middleware.JWTMiddleware(secret), middleware.RequireRole("admin"))
	adm.GET("/creditos", h.ListarDeudores)
	adm.GET("/creditos/:cliente_id", h.EstadoCuenta)
	adm.POST("/pagos", h.RegistrarPago)
	return r
}

// ─── Tests ListarDeudores ────────────────────────────────────────────────────

func TestListarDeudores_RetornaClientesConDeuda(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)
	cr.On("FindAll", 1, 100, "").Return([]models.Cliente{
		{ID: 1, Nombre: "Juan", DeudaTotal: 120.00},
		{ID: 2, Nombre: "Ana", DeudaTotal: 0.00},
	}, int64(2), nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/creditos", nil)
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	// Solo debe aparecer Juan (DeudaTotal > 0)
	deudores := resp["deudores"].([]interface{})
	assert.Len(t, deudores, 1)
}

func TestListarDeudores_ErrorBDRetorna500(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)
	cr.On("FindAll", 1, 100, "").Return([]models.Cliente{}, int64(0), errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/creditos", nil)
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

// ─── Tests EstadoCuenta ───────────────────────────────────────────────────────

func TestEstadoCuenta_ClienteExisteRetornaDatos(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)
	cr.On("FindByID", 1).Return(&models.Cliente{ID: 1, Nombre: "Juan", DeudaTotal: 70.00}, nil)
	pr.On("FindByCliente", 1).Return([]models.Pago{{ID: 10, Monto: 50.00}}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/creditos/1", nil)
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.NotNil(t, resp["cliente"])
	assert.NotNil(t, resp["pagos_realizados"])
}

func TestEstadoCuenta_ClienteInexistenteRetorna404(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)
	cr.On("FindByID", 99).Return(nil, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/creditos/99", nil)
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestEstadoCuenta_IDInvalidoRetorna400(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/creditos/abc", nil)
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ─── Tests RegistrarPago ─────────────────────────────────────────────────────

func TestRegistrarPago_DatosValidosRetorna201(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)
	cr.On("FindByID", 1).Return(&models.Cliente{ID: 1, DeudaTotal: 120.00}, nil)
	pr.On("Create", mock.AnythingOfType("*models.Pago")).Return(nil)

	body := `{"cliente_id":1,"monto":50.00,"metodo":"yape"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/pagos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestRegistrarPago_MontoMayorADeudaRetorna422(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)
	cr.On("FindByID", 1).Return(&models.Cliente{ID: 1, DeudaTotal: 50.00}, nil)

	body := `{"cliente_id":1,"monto":200.00,"metodo":"efectivo"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/pagos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRegistrarPago_ClienteSinDeudaRetorna422(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)
	cr.On("FindByID", 1).Return(&models.Cliente{ID: 1, DeudaTotal: 0.00}, nil)

	body := `{"cliente_id":1,"monto":10.00,"metodo":"efectivo"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/pagos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
}

func TestRegistrarPago_ClienteInexistenteRetorna404(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)
	cr.On("FindByID", 99).Return(nil, nil)

	body := `{"cliente_id":99,"monto":10.00,"metodo":"efectivo"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/pagos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestRegistrarPago_BodyInvalidoRetorna400(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/pagos", bytes.NewBufferString(`{bad}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegistrarPago_ErrorBDRetorna500(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)
	cr.On("FindByID", 1).Return(&models.Cliente{ID: 1, DeudaTotal: 100.00}, nil)
	pr.On("Create", mock.AnythingOfType("*models.Pago")).Return(errors.New("db error"))

	body := `{"cliente_id":1,"monto":50.00,"metodo":"efectivo"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/admin/pagos", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestEstadoCuenta_ErrorBDRetorna500(t *testing.T) {
	pr, cr := new(mockPagoRepo), new(mockClienteRepoC)
	cr.On("FindByID", 1).Return(nil, errors.New("db error"))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/admin/creditos/1", nil)
	req.Header.Set("Authorization", tokAdm())

	routerCreditos(pr, cr).ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
