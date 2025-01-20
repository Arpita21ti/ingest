package triggers

import (
	"gorm.io/gorm"
)

func PreventPrivilegesDocumentIdUpdate(db *gorm.DB) error {
	// Define the trigger and the associated function
	return db.Exec(`
        -- Function to prevent updates to DocumentID in privileges
        CREATE OR REPLACE FUNCTION prevent_privileges_documentid_update()
        RETURNS TRIGGER AS $$
        BEGIN
            IF TG_OP = 'UPDATE' AND NEW.document_id != OLD.document_id THEN
                RAISE EXCEPTION 'Updating document_id in privileges is not allowed';
            END IF;
            RETURN NEW;
        END;
        $$ LANGUAGE plpgsql;

        -- Trigger to enforce the function
        CREATE TRIGGER prevent_documentid_update
        BEFORE UPDATE ON privileges
        FOR EACH ROW
        WHEN (OLD.document_id IS DISTINCT FROM NEW.document_id)
        EXECUTE FUNCTION prevent_privileges_documentid_update();
    `).Error
}

func PreventPrivilegesDocumentIdUpdateRollBack(db *gorm.DB) error {
	// Drop the trigger and function if rolling back
	return db.Exec(`
        DROP TRIGGER IF EXISTS prevent_documentid_update ON privileges;
        DROP FUNCTION IF EXISTS prevent_privileges_documentid_update;
    `).Error
}
