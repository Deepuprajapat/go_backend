package migration_jobs

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/VI-IM/im_backend_go/shared/logger"
)

var (
	cityIDToCityNameMap                                       = make(map[int64]LCity)
	localityIDToLocalityMap                                   = make(map[int64]LLocality)
	developerIDToDeveloperMap                                 = make(map[int64]LDeveloper)
	projectIDToProjectMap                                     = make(map[int64]LProject)
	propertyIDToPropertyMap                                   = make(map[int64]LProperty)
	propertyConfigurationIDToPropertyConfigurationMap         = make(map[int64]LPropertyConfiguration)
	propertyConfigurationTypeIDToPropertyConfigurationTypeMap = make(map[int64]LPropertyConfigurationType)
	projectAmenityIDToProjectAmenityMap                       = make(map[int64]LProjectAmenity)
	PropertyConfigurationIDToPropertyConfigurationMap         = make(map[int64]LPropertyConfiguration)
	ProjectConfigurationIDToProjectConfigurationMap           = make(map[int64]LPropertyConfiguration)
	DeveloperIDToDeveloperMap                                 = make(map[int64]LDeveloper)
	ProjectImageIDToProjectImageMap                           = make(map[int64]LProjectImage)
	FloorPlanIDToFloorPlanMap                                 = make(map[int64]LFloorPlan)
	PaymentPlanIDToPaymentPlanMap                             = make(map[int64]LPaymentPlan)
	FAQIDToFAQMap                                             = make(map[int64]LFAQ)
	ProjectIDToReraMap                                        = make(map[int64][]LRera)
	ProjectIDToFloorPlanMap                                   = make(map[int64][]LFloorPlan)
	ProjectIDToPaymentPlanMap                                 = make(map[int64][]LPaymentPlan)
	ProjectIDToFAQMap                                         = make(map[int64][]LFAQ)
	ProjectIDToProjectConfigurationMap                        = make(map[int64]LPropertyConfiguration)
	ProjectIDToPropertyConfigurationMap                       = make(map[int64]LPropertyConfiguration)
	ProjectIDToPropertyConfigurationTypeMap                   = make(map[int64]LPropertyConfigurationType)
	ProjectIDToProjectImageMap                                = make(map[int64]LProjectImage)
	ProjectIDToProjectAmenityMap                              = make(map[int64][]LProjectAmenity)
)

// TableInfo represents information about a database table
type TableInfo struct {
	Name         string `json:"table_name"`
	RowCount     int64  `json:"row_count"`
	ExportedAt   string `json:"exported_at"`
	ExportedFile string `json:"exported_file"`
}

// DatabaseExportSummary contains summary information about the export
type DatabaseExportSummary struct {
	DatabaseName   string      `json:"database_name"`
	ExportedAt     string      `json:"exported_at"`
	TotalTables    int         `json:"total_tables"`
	TotalRows      int64       `json:"total_rows"`
	ExportDir      string      `json:"export_directory"`
	Tables         []TableInfo `json:"tables"`
	ExportDuration string      `json:"export_duration"`
}

// ExportAllTablesToJSON connects to the database and exports all tables to JSON files
func ExportAllTablesToJSON(ctx context.Context, exportDir string) error {
	startTime := time.Now()

	// Connect to legacy database
	db, err := LegacyDBConnection()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to connect to legacy database")
		return fmt.Errorf("failed to connect to legacy database: %w", err)
	}
	defer db.Close()

	// Create export directory if it doesn't exist
	if exportDir == "" {
		exportDir = "./migration_jobs/database_export"
	}

	if err := os.MkdirAll(exportDir, 0755); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create export directory")
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	logger.Get().Info().Msgf("Starting database export to directory: %s", exportDir)

	// Get all table names
	tables, err := getAllTableNames(ctx, db)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get table names")
		return fmt.Errorf("failed to get table names: %w", err)
	}

	logger.Get().Info().Msgf("Found %d tables to export", len(tables))

	// Export each table
	var exportSummary DatabaseExportSummary
	exportSummary.DatabaseName = "investmango"
	exportSummary.ExportedAt = time.Now().Format(time.RFC3339)
	exportSummary.TotalTables = len(tables)
	exportSummary.ExportDir = exportDir
	exportSummary.Tables = make([]TableInfo, 0, len(tables))

	var totalRows int64
	for _, tableName := range tables {
		logger.Get().Info().Msgf("Exporting table: %s", tableName)

		tableInfo, err := exportTableToJSON(ctx, db, tableName, exportDir)
		if err != nil {
			logger.Get().Error().Err(err).Msgf("Failed to export table %s", tableName)
			continue
		}

		exportSummary.Tables = append(exportSummary.Tables, *tableInfo)
		totalRows += tableInfo.RowCount
	}

	exportSummary.TotalRows = totalRows
	exportSummary.ExportDuration = time.Since(startTime).String()

	// Save export summary
	summaryFile := filepath.Join(exportDir, "export_summary.json")
	if err := saveJSON(summaryFile, exportSummary); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to save export summary")
		return fmt.Errorf("failed to save export summary: %w", err)
	}

	logger.Get().Info().Msgf("Database export completed successfully in %s", exportSummary.ExportDuration)
	logger.Get().Info().Msgf("Exported %d tables with %d total rows", exportSummary.TotalTables, exportSummary.TotalRows)
	logger.Get().Info().Msgf("Export summary saved to: %s", summaryFile)

	return nil
}

// getAllTableNames retrieves all table names from the database
func getAllTableNames(ctx context.Context, db *sql.DB) ([]string, error) {
	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema = 'investmango'
		AND table_type = 'BASE TABLE'
		ORDER BY table_name
	`

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query table names: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, tableName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating table names: %w", err)
	}

	return tables, nil
}

// exportTableToJSON exports a single table to a JSON file
func exportTableToJSON(ctx context.Context, db *sql.DB, tableName string, exportDir string) (*TableInfo, error) {
	// Get table structure
	columns, err := getTableColumns(ctx, db, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to get table columns: %w", err)
	}

	// Query all data from the table
	query := fmt.Sprintf("SELECT * FROM `%s`", tableName)
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query table %s: %w", tableName, err)
	}
	defer rows.Close()

	// Prepare result slice
	var results []map[string]interface{}
	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to get column types: %w", err)
	}

	// Process each row
	for rows.Next() {
		// Create a slice of interface{} to hold the row data
		columnValues := make([]interface{}, len(columns))
		columnPointers := make([]interface{}, len(columns))

		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		// Scan the row into the column pointers
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Create a map for this row
		rowData := make(map[string]interface{})
		for i, column := range columns {
			// Handle different data types
			value := columnValues[i]
			if value != nil {
				switch v := value.(type) {
				case []byte:
					// Convert byte arrays to strings for JSON compatibility
					rowData[column] = string(v)
				case time.Time:
					// Convert time to RFC3339 string
					rowData[column] = v.Format(time.RFC3339)
				default:
					rowData[column] = v
				}
			} else {
				rowData[column] = nil
			}
		}

		results = append(results, rowData)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Create table export data
	tableExport := map[string]interface{}{
		"table_name":   tableName,
		"exported_at":  time.Now().Format(time.RFC3339),
		"row_count":    len(results),
		"columns":      columns,
		"column_types": getColumnTypeNames(columnTypes),
		"data":         results,
	}

	// Save to JSON file
	fileName := fmt.Sprintf("%s.json", tableName)
	filePath := filepath.Join(exportDir, fileName)

	if err := saveJSON(filePath, tableExport); err != nil {
		return nil, fmt.Errorf("failed to save table data: %w", err)
	}

	tableInfo := &TableInfo{
		Name:         tableName,
		RowCount:     int64(len(results)),
		ExportedAt:   time.Now().Format(time.RFC3339),
		ExportedFile: fileName,
	}

	logger.Get().Info().Msgf("Exported table %s: %d rows -> %s", tableName, len(results), fileName)
	return tableInfo, nil
}

// getTableColumns retrieves column names for a table
func getTableColumns(ctx context.Context, db *sql.DB, tableName string) ([]string, error) {
	query := `
		SELECT column_name 
		FROM information_schema.columns 
		WHERE table_schema = DATABASE() 
		AND table_name = ?
		ORDER BY ordinal_position
	`

	rows, err := db.QueryContext(ctx, query, tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to query columns for table %s: %w", tableName, err)
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var columnName string
		if err := rows.Scan(&columnName); err != nil {
			return nil, fmt.Errorf("failed to scan column name: %w", err)
		}
		columns = append(columns, columnName)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating columns: %w", err)
	}

	return columns, nil
}

// getColumnTypeNames extracts column type names from ColumnType slice
func getColumnTypeNames(columnTypes []*sql.ColumnType) []string {
	var typeNames []string
	for _, ct := range columnTypes {
		typeNames = append(typeNames, ct.DatabaseTypeName())
	}
	return typeNames
}

// saveJSON saves data to a JSON file with proper formatting
func saveJSON(filePath string, data interface{}) error {
	// Marshal with indentation for readability
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Write to file
	if err := ioutil.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	return nil
}

// ExportSpecificTablesToJSON exports only specific tables to JSON
func ExportSpecificTablesToJSON(ctx context.Context, tableNames []string, exportDir string) error {
	if len(tableNames) == 0 {
		return fmt.Errorf("no table names provided")
	}

	startTime := time.Now()

	// Connect to legacy database
	db, err := LegacyDBConnection()
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to connect to legacy database")
		return fmt.Errorf("failed to connect to legacy database: %w", err)
	}
	defer db.Close()

	// Create export directory if it doesn't exist
	if exportDir == "" {
		exportDir = fmt.Sprintf("database_export_specific_%s", time.Now().Format("2006-01-02_15-04-05"))
	}

	if err := os.MkdirAll(exportDir, 0755); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create export directory")
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	logger.Get().Info().Msgf("Starting export of %d specific tables to directory: %s", len(tableNames), exportDir)

	// Validate table names exist
	allTables, err := getAllTableNames(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get table names: %w", err)
	}

	tableSet := make(map[string]bool)
	for _, table := range allTables {
		tableSet[table] = true
	}

	// Filter valid table names
	var validTables []string
	var invalidTables []string

	for _, tableName := range tableNames {
		if tableSet[tableName] {
			validTables = append(validTables, tableName)
		} else {
			invalidTables = append(invalidTables, tableName)
		}
	}

	if len(invalidTables) > 0 {
		logger.Get().Warn().Msgf("Invalid table names (will be skipped): %s", strings.Join(invalidTables, ", "))
	}

	if len(validTables) == 0 {
		return fmt.Errorf("no valid table names found")
	}

	// Export each valid table
	var exportSummary DatabaseExportSummary
	exportSummary.DatabaseName = "investmango"
	exportSummary.ExportedAt = time.Now().Format(time.RFC3339)
	exportSummary.TotalTables = len(validTables)
	exportSummary.ExportDir = exportDir
	exportSummary.Tables = make([]TableInfo, 0, len(validTables))

	var totalRows int64
	for _, tableName := range validTables {
		logger.Get().Info().Msgf("Exporting table: %s", tableName)

		tableInfo, err := exportTableToJSON(ctx, db, tableName, exportDir)
		if err != nil {
			logger.Get().Error().Err(err).Msgf("Failed to export table %s", tableName)
			continue
		}

		exportSummary.Tables = append(exportSummary.Tables, *tableInfo)
		totalRows += tableInfo.RowCount
	}

	exportSummary.TotalRows = totalRows
	exportSummary.ExportDuration = time.Since(startTime).String()

	// Save export summary
	summaryFile := filepath.Join(exportDir, "export_summary.json")
	if err := saveJSON(summaryFile, exportSummary); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to save export summary")
		return fmt.Errorf("failed to save export summary: %w", err)
	}

	logger.Get().Info().Msgf("Specific tables export completed successfully in %s", exportSummary.ExportDuration)
	logger.Get().Info().Msgf("Exported %d tables with %d total rows", exportSummary.TotalTables, exportSummary.TotalRows)
	logger.Get().Info().Msgf("Export summary saved to: %s", summaryFile)

	return nil
}
