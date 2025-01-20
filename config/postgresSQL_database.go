// Add and upgrade to use multiple connections to optimize performance
// Use different connections to same DB for Crud, Analytics etc. (specific tasks)
// Be sure to manage concurrency while using multiple connections (use locks, transactions and isolation levels)
package config

import (
	"fmt"
	"log"
	"os"
	"time"

	common_tables "server/models/common"
	question_hierarchy "server/models/question_bank/question_hierarchy"
	question_type "server/models/question_bank/question_type"
	student_tables "server/models/student_psql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var postgresDBConnection *gorm.DB

// ensureSchemaExists ensures that a schema exists and creates it if it doesn't
func ensureSchemaExists(tx *gorm.DB, schemaName string) error {
	var count int64
	// Check if the schema exists
	if err := tx.Raw(`
		SELECT COUNT(*) 
		FROM information_schema.schemata 
		WHERE schema_name = ?`, schemaName).Scan(&count).Error; err != nil {
		return fmt.Errorf("failed to check schema existence: %w", err)
	}

	// If the schema does not exist, create it
	if count == 0 {
		if err := tx.Exec(fmt.Sprintf("CREATE SCHEMA %s", schemaName)).Error; err != nil {
			return fmt.Errorf("failed to create schema %s: %w", schemaName, err)
		}
		log.Printf("Schema %s created successfully", schemaName)
	}

	// Grant necessary access to the schema
	if err := grantSchemaAccess(tx, schemaName); err != nil {
		return fmt.Errorf("failed to grant access to schema %s: %w", schemaName, err)
	}

	// Set Default Grant privileges to the tables
	if err := grantDefaultTablePrivileges(tx, schemaName); err != nil {
		return fmt.Errorf("failed to grant access to table %s: %w", schemaName, err)
	}

	return nil
}

// grantSchemaAccess grants necessary access privileges to a schema
func grantSchemaAccess(tx *gorm.DB, schemaName string) error {
	// Define roles with their respective privileges
	adminAndCoordinators := []string{"admin", "tnp_coordinator"} // Usage and Create
	usageOnlyRoles := []string{"tnp_volunteer", "student"}       // Usage only

	// Grant USAGE and CREATE for admin and coordinators
	for _, role := range adminAndCoordinators {
		if err := tx.Exec(fmt.Sprintf(`
            GRANT USAGE, CREATE ON SCHEMA %s TO %s;
        `, schemaName, role)).Error; err != nil {
			return fmt.Errorf("failed to grant USAGE and CREATE access to schema %s for role %s: %w", schemaName, role, err)
		}
	}

	// Grant only USAGE for volunteers and students
	for _, role := range usageOnlyRoles {
		if err := tx.Exec(fmt.Sprintf(`
            GRANT USAGE ON SCHEMA %s TO %s;
        `, schemaName, role)).Error; err != nil {
			return fmt.Errorf("failed to grant USAGE access to schema %s for role %s: %w", schemaName, role, err)
		}
	}

	return nil
}

// grantDefaultTablePrivileges sets default privileges for new tables in the schema
func grantDefaultTablePrivileges(tx *gorm.DB, schemaName string) error {
	// Define the roles and privileges for new tables
	defaultAccessRoles := map[string]string{
		"admin":           "SELECT, INSERT, UPDATE, DELETE, REFERENCES, TRIGGER",
		"tnp_coordinator": "SELECT, INSERT, UPDATE, DELETE, REFERENCES, TRIGGER",
		"tnp_volunteer":   "SELECT, INSERT, UPDATE",
	}
	// Handle student role separately for "question_schema" schema
	if schemaName == "question_schema" {
		defaultAccessRoles["student"] = "SELECT" // Student role only gets SELECT privilege
	} else {
		defaultAccessRoles["student"] = "SELECT, UPDATE" // Student role gets both SELECT and UPDATE privileges
	}

	// Loop through the roles and grant default privileges for each role
	for role, privileges := range defaultAccessRoles {
		if err := tx.Exec(fmt.Sprintf(`
			ALTER DEFAULT PRIVILEGES IN SCHEMA %s
			GRANT %s ON TABLES TO %s;
		`, schemaName, privileges, role)).Error; err != nil {
			return fmt.Errorf("failed to set default privileges for role %s in schema %s: %w", role, schemaName, err)
		}
	}

	return nil
}

// // TODO: Update and include for partitioning hierarchy
// // setupPartitioning creates partitioned tables for the question hierarchy
// func setupPartitioning(db *gorm.DB) error {
// 	// Define partition configurations for each level
// 	partitions := []struct {
// 		TableName      string // Parent table name
// 		PartitionField string // Field used for partitioning
// 	}{
// 		{"question_schema.question_sub_domains_table", "question_domain_id"},
// 		{"question_schema.question_niches_table", "question_sub_domain_id"},
// 		{"question_schema.question_difficulty_level_table", "question_niche_id"},
// 		{"question_schema.question_formats_table", "question_difficulty_id"},
// 	}
// 	// Loop through and create partitions for each level
// 	for _, partition := range partitions {
// 		cmd := fmt.Sprintf(`
// 			CREATE TABLE IF NOT EXISTS %s_partition (
// 				LIKE %s INCLUDING ALL
// 			) PARTITION BY LIST (%s);`,
// 			partition.TableName,      // Partition table name
// 			partition.TableName,      // Parent table name
// 			partition.PartitionField, // Field used for partitioning
// 		)
// 		// Execute the SQL command
// 		if err := db.Exec(cmd).Error; err != nil {
// 			return fmt.Errorf("failed to create partition for table %s: %w", partition.TableName, err)
// 		}
// 	}
// 	log.Println("Partitioning setup completed successfully.")
// 	return nil
// }

// runAutoMigrations runs auto migrations to create or update the database schema based on the models
func runAutoMigrations() error {
	// Define the schemas you need to ensure exist
	schemas := []string{"public", "student_schema", "question_schema"}

	// Ensure schemas exist
	for _, schema := range schemas {
		err := postgresDBConnection.Transaction(func(tx *gorm.DB) error {
			if err := ensureSchemaExists(tx, schema); err != nil {
				return fmt.Errorf("failed to ensure schema %s exists: %w", schema, err)
			}
			// The transaction will be committed automatically if no error occurs
			return nil
		})
		if err != nil {
			return err
		}
	}

	// Migrate common schema tables
	err := postgresDBConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(&common_tables.DocumentPrivilegeTable{}); err != nil {
			return fmt.Errorf("failed to auto migrate common models: %w", err)
		}
		// The transaction will be committed automatically if no error occurs
		return nil
	})
	if err != nil {
		return err
	}

	// Migrate student schema tables
	err = postgresDBConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(
			&student_tables.StudentDocumentTable{},
			&student_tables.StudentFamilyDetailsTable{},
			&student_tables.StudentLogInDetailsTable{},
			&student_tables.StudentScholarshipDetailsTable{},
			&student_tables.StudentLeaderboardRecordTable{},
			&student_tables.StudentPracticeSessionRecordTable{},
		); err != nil {
			return fmt.Errorf("failed to auto migrate student models: %w", err)
		}
		// The transaction will be committed automatically if no error occurs
		return nil
	})
	if err != nil {
		return err
	}

	// Migrate dependent student schema tables
	err = postgresDBConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(
			&student_tables.StudentAcademicDetailsTable{},
			&student_tables.StudentCertificationDetailsTable{},
			&student_tables.StudentProfileDetailsTable{},
			&student_tables.StudentCertificationLookup{},
			&student_tables.EnrollmentMasterLookupTable{},
			&student_tables.StudentPracticeSessionLookupTable{},
			&student_tables.StudentLeaderboardLookupTable{},
		); err != nil {
			return fmt.Errorf("failed to auto migrate dependent student models: %w", err)
		}
		// The transaction will be committed automatically if no error occurs
		return nil
	})
	if err != nil {
		return err
	}

	// Migrate question hierarchy schema tables
	err = postgresDBConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(
			&question_hierarchy.QuestionDomainsTable{},
			&question_hierarchy.QuestionSubDomainsTable{},
			&question_hierarchy.QuestionNicheTable{},
			&question_hierarchy.QuestionDifficultyLevelTable{},
			&question_hierarchy.QuestionFormatTable{},
		); err != nil {
			return fmt.Errorf("failed to auto migrate question hierarchy models: %w", err)
		}
		// The transaction will be committed automatically if no error occurs
		return nil
	})
	if err != nil {
		return err
	}

	// Migrate question type schema tables
	err = postgresDBConnection.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(
			&question_type.TextBasedQuestion{},
			&question_type.TrueFalseQuestion{},
			&question_type.FillInTheBlankQuestion{},
			&question_type.MCQQuestion{},
		); err != nil {
			return fmt.Errorf("failed to auto migrate question type models: %w", err)
		}
		// The transaction will be committed automatically if no error occurs
		return nil
	})
	if err != nil {
		return err
	}

	// // TODO: Update and include for partitioning hierarchy
	// // Add partitioning hierarchy using raw SQL
	// err = postgresDBConnection.Transaction(func(tx *gorm.DB) error {
	// 	if err := setupPartitioning(tx); err != nil {
	// 		return fmt.Errorf("failed to setup partitioning: %w", err)
	// 	}
	// 	// The transaction will be committed automatically if no error occurs
	// 	return nil
	// })
	// if err != nil {
	// 	return err
	// }

	// Log migration success
	log.Println("Auto migration completed successfully for all models.")
	return nil
}

// ConnectPostgresDB initializes the connection to PostgreSQL database.
func ConnectPostgresDB() error {
	var dsn string

	// Run migrations and localhost in the development environment
	env := os.Getenv("ENV")
	if env == "development" {
		// Validate required environment variables for manual configuration
		requiredEnvVars := []string{"POSTGRES_DB_HOST", "POSTGRES_DB_USER_ADMIN",
			"POSTGRES_DB_PASSWORD_ADMIN", "POSTGRES_DB_NAME", "POSTGRES_DB_PORT"}
		for _, envVar := range requiredEnvVars {
			if value := os.Getenv(envVar); value == "" {
				return fmt.Errorf("environment variable %s is not set or is empty", envVar)
			}
		}

		// Build the DSN from environment variables
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata",
			os.Getenv("POSTGRES_DB_HOST"),           // e.g., "localhost"
			os.Getenv("POSTGRES_DB_USER_ADMIN"),     // e.g., "tnp_user"
			os.Getenv("POSTGRES_DB_PASSWORD_ADMIN"), // e.g., "securepassword"
			os.Getenv("POSTGRES_DB_NAME"),           // e.g., "tnp_database"
			os.Getenv("POSTGRES_DB_PORT"),           // e.g., "5432"
		)

		log.Println("Using manual database configuration for connection.")

	} else {

		// Check for an external database URL (e.g., `POSTGRES_DB_URL`)
		externalDBURL := os.Getenv("POSTGRES_DB_URL_EXTERNAL")
		internalDBURL := os.Getenv("POSTGRES_DB_URL_INTERNAL")

		if internalDBURL != "" {
			// Use the external URL if provided
			dsn = internalDBURL
			log.Println("Using internal database URL for connection.")
		} else if externalDBURL != "" {
			// Use the external URL if provided
			dsn = externalDBURL
			log.Println("Using external database URL for connection.")
		} else {
			return fmt.Errorf("no database URL found in environment variables")
		}
	}

	// Attempt to open a connection to the database
	var err error
	// Initialize GORM with logging
	postgresDBConnection, err = gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.New(
		log.New(os.Stdout, "\n", log.LstdFlags), // log.New creates a new logger
		logger.Config{
			LogLevel: logger.Info, // or logger.Debug for more verbose logging
			Colorful: true,        // Enable color output
		},
	)})
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	// Ran Automigration after initializing postgresdb
	if env == "development" {
		if err := runAutoMigrations(); err != nil {
			return fmt.Errorf("failed to run migrations: %w", err)
		}

		log.Println("Auto migration completed successfully.")

		// PostgreSQL version check for diagnostics (useful for development)
		var version string
		err = postgresDBConnection.Raw("SELECT version()").Scan(&version).Error
		if err == nil {
			log.Printf("Connected to PostgreSQL, version: %s", version)
		} else {
			log.Printf("Failed to retrieve PostgreSQL version: %v", err)
		}
	}

	sqlDB, err := postgresDBConnection.DB()
	if err != nil {
		panic("Failed to get native database handle")
	}

	// Configure connection pooling
	sqlDB.SetMaxOpenConns(100)          // Maximum open connections
	sqlDB.SetMaxIdleConns(20)           // Maximum idle connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum connection lifetime

	// Log the successful connection to PostgreSQL
	log.Println("Connected to PostgreSQL!")

	return nil
}

// DisconnectPostgresDB disconnects from the PostgreSQL database
func DisconnectPostgresDB() error {
	// If postgresdb is nil, return early since there is no active connection
	if postgresDBConnection == nil {
		log.Println("No active connection to PostgreSQL.")
		return nil
	}

	// Get the underlying *sql.DB instance
	sqlDB, err := postgresDBConnection.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB instance: %w", err)
	}

	// Close the connection
	if err := sqlDB.Close(); err != nil {
		return fmt.Errorf("failed to close PostgreSQL connection: %w", err)
	}

	// Log the successful disconnection from PostgreSQL
	log.Println("Disconnected from PostgreSQL")
	return nil
}

// GetPostgresDBConnection returns the global postgresdb connection (use with caution)
func GetPostgresDBConnection() *gorm.DB {
	return postgresDBConnection
}

// GetPostgresTable returns the GORM model instance for the specific table (similar to GetCollectionMongo)
func GetPostgresTable(table interface{}) *gorm.DB {
	return postgresDBConnection.Model(table)
}
