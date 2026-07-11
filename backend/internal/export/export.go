package export

import (
	"fmt"
	"io"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"

	"restaurantos/internal/models"
	"restaurantos/internal/services"
)

// ─── Ventas ───────────────────────────────────────────────────────────────────

func VentasExcel(w io.Writer, r services.ReporteVentas) error {
	f := excelize.NewFile()
	defer f.Close()
	sheet := "Ventas"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "Reporte de Ventas")
	f.SetCellValue(sheet, "A2", fmt.Sprintf("Periodo: %s (%s a %s)", r.Periodo, r.Desde, r.Hasta))
	f.SetCellValue(sheet, "A4", "Fecha")
	f.SetCellValue(sheet, "B4", "Total")
	f.SetCellValue(sheet, "C4", "Cantidad de pedidos")

	row := 5
	for _, d := range r.PorDia {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), d.Fecha)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), d.Total)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), d.CantidadPedidos)
		row++
	}
	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.TotalVentas)
	f.SetCellValue(sheet, fmt.Sprintf("C%d", row), r.TotalPedidos)

	return f.Write(w)
}

func VentasPDF(w io.Writer, r services.ReporteVentas) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Reporte de Ventas")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Periodo: %s (%s a %s)", r.Periodo, r.Desde, r.Hasta))
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(60, 8, "Fecha", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 8, "Total (S/)", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 8, "Pedidos", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 11)
	for _, d := range r.PorDia {
		pdf.CellFormat(60, 8, d.Fecha, "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 8, fmt.Sprintf("%.2f", d.Total), "1", 0, "C", false, 0, "")
		pdf.CellFormat(60, 8, fmt.Sprintf("%d", d.CantidadPedidos), "1", 1, "C", false, 0, "")
	}

	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Total ventas: S/ %.2f   -   Total pedidos: %d", r.TotalVentas, r.TotalPedidos))

	return pdf.Output(w)
}

// ─── Clientes con deuda ────────────────────────────────────────────────────────

func DeudoresExcel(w io.Writer, clientes []models.Cliente) error {
	f := excelize.NewFile()
	defer f.Close()
	sheet := "Deudores"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "Reporte de Clientes con Deuda")
	f.SetCellValue(sheet, "A3", "Nombre")
	f.SetCellValue(sheet, "B3", "Tipo")
	f.SetCellValue(sheet, "C3", "Teléfono")
	f.SetCellValue(sheet, "D3", "Deuda total (S/)")

	row := 4
	var totalDeuda float64
	for _, c := range clientes {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("%s %s", c.Nombre, c.Apellido))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), c.Tipo)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), c.Telefono)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), c.DeudaTotal)
		totalDeuda += c.DeudaTotal
		row++
	}
	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total deuda")
	f.SetCellValue(sheet, fmt.Sprintf("D%d", row), totalDeuda)

	return f.Write(w)
}

func DeudoresPDF(w io.Writer, clientes []models.Cliente) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Reporte de Clientes con Deuda")
	pdf.Ln(14)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(70, 8, "Nombre", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, "Tipo", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, "Teléfono", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, "Deuda (S/)", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 11)
	var totalDeuda float64
	for _, c := range clientes {
		pdf.CellFormat(70, 8, fmt.Sprintf("%s %s", c.Nombre, c.Apellido), "1", 0, "", false, 0, "")
		pdf.CellFormat(40, 8, c.Tipo, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, c.Telefono, "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 8, fmt.Sprintf("%.2f", c.DeudaTotal), "1", 1, "C", false, 0, "")
		totalDeuda += c.DeudaTotal
	}

	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Total deuda: S/ %.2f", totalDeuda))

	return pdf.Output(w)
}

// ─── Movimientos de inventario ────────────────────────────────────────────────

func MovimientosExcel(w io.Writer, r services.ReporteMovimientos) error {
	f := excelize.NewFile()
	defer f.Close()
	sheet := "Movimientos"
	f.SetSheetName("Sheet1", sheet)

	f.SetCellValue(sheet, "A1", "Reporte de Movimientos de Inventario")
	f.SetCellValue(sheet, "A2", fmt.Sprintf("Del %s al %s", r.Desde, r.Hasta))
	f.SetCellValue(sheet, "A4", "Fecha")
	f.SetCellValue(sheet, "B4", "Producto")
	f.SetCellValue(sheet, "C4", "Tipo")
	f.SetCellValue(sheet, "D4", "Cantidad")

	row := 5
	for _, m := range r.Movimientos {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), m.Fecha.Format("2006-01-02"))
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), m.Producto.Nombre)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), m.Tipo)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), m.Cantidad)
		row++
	}
	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total entradas")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.TotalEntradas)
	row++
	f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "Total salidas")
	f.SetCellValue(sheet, fmt.Sprintf("B%d", row), r.TotalSalidas)

	return f.Write(w)
}

func MovimientosPDF(w io.Writer, r services.ReporteMovimientos) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(0, 10, "Reporte de Movimientos de Inventario")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Del %s al %s", r.Desde, r.Hasta))
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(35, 8, "Fecha", "1", 0, "C", false, 0, "")
	pdf.CellFormat(75, 8, "Producto", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Tipo", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 8, "Cantidad", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 11)
	for _, m := range r.Movimientos {
		pdf.CellFormat(35, 8, m.Fecha.Format("2006-01-02"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(75, 8, m.Producto.Nombre, "1", 0, "", false, 0, "")
		pdf.CellFormat(30, 8, m.Tipo, "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("%.2f", m.Cantidad), "1", 1, "C", false, 0, "")
	}

	pdf.Ln(6)
	pdf.SetFont("Arial", "B", 11)
	pdf.Cell(0, 8, fmt.Sprintf("Total entradas: %.2f   -   Total salidas: %.2f", r.TotalEntradas, r.TotalSalidas))

	return pdf.Output(w)
}
