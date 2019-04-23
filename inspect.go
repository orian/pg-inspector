package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/gocraft/dbr"
	_ "github.com/lib/pq"
)

type (
	CardinalNumber dbr.NullInt64
	CharacterData dbr.NullString
	SQLIdentifier dbr.NullString
	TimeStamp dbr.NullTime
	YesOrNo dbr.NullString
)

// Data types:
// CardinalNumber
// A nonnegative integer.
//
// CharacterData
// A character string (without specific maximum length).
//
// SQLIdentifier
// A character string. This type is used for SQL identifiers, the type
// CharacterData is used for any other kind of text data.
//
// time_stamp
// A domain over the type timestamp with time zone
//
// YesOrNo
// A character string domain that contains either YES or NO. This is used
// to represent Boolean (true/false) data in the information schema.
// (The information schema was invented before the type boolean
// was added to the SQL standard, so this convention is necessary
// to keep the information schema backward compatible.)

// https://www.postgresql.org/docs/9.6/infoschema-schemata.html
type TSchemata struct {
	CatalogName                SQLIdentifier `db:"catalog_name"` // Name of the database that the schema is contained in (always the current database)
	SchemaName                 SQLIdentifier                     // Name of the schema
	SchemaOwner                SQLIdentifier                     // Name of the owner of the schema
	DefaultCharacterSetCatalog SQLIdentifier                     // Applies to a feature not available in PostgreSQL
	DefaultCharacterSetSchema  SQLIdentifier                     // Applies to a feature not available in PostgreSQL
	DefaultCharacterSetName    SQLIdentifier                     // Applies to a feature not available in PostgreSQL
	SQLPath                    CharacterData                     // Applies to a feature not available in PostgreSQL
}

// https://www.postgresql.org/docs/9.6/infoschema-tables.html
type TTables struct {
	TableCatalog              SQLIdentifier `db:"table_catalog"`                // Name of the database that contains the table (always the current database)
	TableSchema               SQLIdentifier `db:"table_schema"`                 // Name of the schema that contains the table
	TableName                 SQLIdentifier `db:"table_name"`                   // Name of the table
	TableType                 CharacterData `db:"table_type"`                   // Type of the table: BASE TABLE for a persistent base table (the normal table type), VIEW for a view, FOREIGN TABLE for a foreign table, or LOCAL TEMPORARY for a temporary table
	SelfReferencingColumnName SQLIdentifier `db:"self_referencing_column_name"` // Applies to a feature not available in PostgreSQL
	ReferenceGeneration       CharacterData `db:"reference_generation"`         // Applies to a feature not available in PostgreSQL
	UserDefinedTypeCatalog    SQLIdentifier `db:"user_defined_type_catalog"`    // If the table is a typed table, the name of the database that contains the underlying data type (always the current database), else null.
	UserDefinedTypeSchema     SQLIdentifier `db:"user_defined_type_schema"`     // If the table is a typed table, the name of the schema that contains the underlying data type, else null.
	UserDefinedTypeName       SQLIdentifier `db:"user_defined_type_name"`       // If the table is a typed table, the name of the underlying data type, else null.
	IsInsertableInto          YesOrNo       `db:"is_insertable_into"`           // YES if the table is insertable into, NO if not (Base tables are always insertable into, views not necessarily.)
	IsTyped                   YesOrNo       `db:"is_typed"`                     // YES if the table is a typed table, NO if not
	CommitAction              CharacterData `db:"commit_action"`                // Not yet implemented
}

type TColumns struct {
	TableCatalog           SQLIdentifier  `db:"table_catalog"`            // Name of the database containing the table (always the current database)
	TableSchema            SQLIdentifier  `db:"table_schema"`             // Name of the schema containing the table
	TableName              SQLIdentifier  `db:"table_name"`               // Name of the table
	ColumnName             SQLIdentifier  `db:"column_name"`              // Name of the column
	OrdinalPosition        CardinalNumber `db:"ordinal_position"`         // Ordinal position of the column within the table (count starts at 1)
	ColumnDefault          CharacterData  `db:"column_default"`           // Default expression of the column
	IsNullable             YesOrNo        `db:"is_nullable"`              // YES if the column is possibly nullable, NO if it is known not nullable. A not-null constraint is one way a column can be known not nullable, but there can be others.
	DataType               CharacterData  `db:"data_type"`                // Data type of the column, if it is a built-in type, or ARRAY if it is some array (in that case, see the view element_types), else USER-DEFINED (in that case, the type is identified in udt_name and associated columns). If the column is based on a domain, this column refers to the type underlying the domain (and the domain is identified in domain_name and associated columns).
	CharacterMaximumLength CardinalNumber `db:"character_maximum_length"` // If data_type identifies a character or bit string type, the declared maximum length; null for all other data types or if no maximum length was declared.
	CharacterOctetLength   CardinalNumber `db:"character_octet_length"`   // If data_type identifies a character type, the maximum possible length in octets (bytes) of a datum; null for all other data types. The maximum octet length depends on the declared character maximum length (see above) and the server encoding.
	NumericPrecision       CardinalNumber `db:"numeric_precision"`        // If data_type identifies a numeric type, this column contains the (declared or implicit) precision of the type for this column. The precision indicates the number of significant digits. It can be expressed in decimal (base 10) or binary (base 2) terms, as specified in the column numeric_precision_radix. For all other data types, this column is null.
	NumericPrecisionRadix  CardinalNumber `db:"numeric_precision_radix"`  // If data_type identifies a numeric type, this column indicates in which base the values in the columns numeric_precision and numeric_scale are expressed. The value is either 2 or 10. For all other data types, this column is null.
	NumericScale           CardinalNumber `db:"numeric_scale"`            // If data_type identifies an exact numeric type, this column contains the (declared or implicit) scale of the type for this column. The scale indicates the number of significant digits to the right of the decimal point. It can be expressed in decimal (base 10) or binary (base 2) terms, as specified in the column numeric_precision_radix. For all other data types, this column is null.
	DatetimePrecision      CardinalNumber `db:"datetime_precision"`       // If data_type identifies a date, time, timestamp, or interval type, this column contains the (declared or implicit) fractional seconds precision of the type for this column, that is, the number of decimal digits maintained following the decimal point in the seconds value. For all other data types, this column is null.
	IntervalType           CharacterData  `db:"interval_type"`            // If data_type identifies an interval type, this column contains the specification which fields the intervals include for this column, e.g., YEAR TO MONTH, DAY TO SECOND, etc. If no field restrictions were specified (that is, the interval accepts all fields), and for all other data types, this field is null.
	IntervalPrecision      CardinalNumber `db:"interval_precision"`       // Applies to a feature not available in PostgreSQL (see datetime_precision for the fractional seconds precision of interval type columns)
	CharacterSetCatalog    SQLIdentifier  `db:"character_set_catalog"`    // Applies to a feature not available in PostgreSQL
	CharacterSetSchema     SQLIdentifier  `db:"character_set_schema"`     // Applies to a feature not available in PostgreSQL
	CharacterSetName       SQLIdentifier  `db:"character_set_name"`       // Applies to a feature not available in PostgreSQL
	CollationCatalog       SQLIdentifier  `db:"collation_catalog"`        // Name of the database containing the collation of the column (always the current database), null if default or the data type of the column is not collatable
	CollationSchema        SQLIdentifier  `db:"collation_schema"`         // Name of the schema containing the collation of the column, null if default or the data type of the column is not collatable
	CollationName          SQLIdentifier  `db:"collation_name"`           // Name of the collation of the column, null if default or the data type of the column is not collatable
	DomainCatalog          SQLIdentifier  `db:"domain_catalog"`           // If the column has a domain type, the name of the database that the domain is defined in (always the current database), else null.
	DomainSchema           SQLIdentifier  `db:"domain_schema"`            // If the column has a domain type, the name of the schema that the domain is defined in, else null.
	DomainName             SQLIdentifier  `db:"domain_name"`              // If the column has a domain type, the name of the domain, else null.
	UdtCatalog             SQLIdentifier  `db:"udt_catalog"`              // Name of the database that the column data type (the underlying type of the domain, if applicable) is defined in (always the current database)
	UdtSchema              SQLIdentifier  `db:"udt_schema"`               // Name of the schema that the column data type (the underlying type of the domain, if applicable) is defined in
	UdtName                SQLIdentifier  `db:"udt_name"`                 // Name of the column data type (the underlying type of the domain, if applicable)
	ScopeCatalog           SQLIdentifier  `db:"scope_catalog"`            // Applies to a feature not available in PostgreSQL
	ScopeSchema            SQLIdentifier  `db:"scope_schema"`             // Applies to a feature not available in PostgreSQL
	ScopeName              SQLIdentifier  `db:"scope_name"`               // Applies to a feature not available in PostgreSQL
	MaximumCardinality     CardinalNumber `db:"maximum_cardinality"`      // Always null, because arrays always have unlimited maximum cardinality in PostgreSQL
	DtdIdentifier          SQLIdentifier  `db:"dtd_identifier"`           // An identifier of the data type descriptor of the column, unique among the data type descriptors pertaining to the table. This is mainly useful for joining with other instances of such identifiers. (The specific format of the identifier is not defined and not guaranteed to remain the same in future versions.)
	IsSelfReferencing      YesOrNo        `db:"is_self_referencing"`      // Applies to a feature not available in PostgreSQL
	IsIdentity             YesOrNo        `db:"is_identity"`              // Applies to a feature not available in PostgreSQL
	IdentityGeneration     CharacterData  `db:"identity_generation"`      // Applies to a feature not available in PostgreSQL
	IdentityStart          CharacterData  `db:"identity_start"`           // Applies to a feature not available in PostgreSQL
	IdentityIncrement      CharacterData  `db:"identity_increment"`       // Applies to a feature not available in PostgreSQL
	IdentityMaximum        CharacterData  `db:"identity_maximum"`         // Applies to a feature not available in PostgreSQL
	IdentityMinimum        CharacterData  `db:"identity_minimum"`         // Applies to a feature not available in PostgreSQL
	IdentityCycle          YesOrNo        `db:"identity_cycle"`           // Applies to a feature not available in PostgreSQL
	IsGenerated            CharacterData  `db:"is_generated"`             // Applies to a feature not available in PostgreSQL
	GenerationExpression   CharacterData  `db:"generation_expression"`    // Applies to a feature not available in PostgreSQL
	IsUpdatable            YesOrNo        `db:"is_updatable"`             // YES if the column is updatable, NO if not (Columns in base tables are always updatable, columns in views not necessarily)
}

func main() {
	connStr := flag.String("db", "", "PostgreSQL connection string.")
	flag.Parse()

	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "Jan 02, 15:04:06",
	}

	dbConn, err := dbr.Open("postgres", *connStr, nil)
	if err != nil {
		log.WithError(err).Fatal("Cannot connect to db")
	}
	defer dbConn.Close()

	dbS := dbConn.NewSession(nil)
	var dbName string
	if err = dbS.Select("*").From("information_schema.information_schema_catalog_name").LoadOne(&dbName); err != nil {
		log.WithError(err).Fatal("load database name")
	}

	log.Infof("db name: %s", dbName)

	var schemaWhitelist []string = []string{
		"swipe",
	}

	var schemas []TSchemata
	n, err := dbS.SelectBySql("SELECT * FROM information_schema.schemata WHERE schema_name IN ?", schemaWhitelist).Load(&schemas)
	if err != nil {
		log.WithError(err).Fatal("select schemas")
	}
	if n == 0 {
		log.Warn("no schemas available")
		return
	}
	for _, v := range schemas {
		log.Debugf("schema %s owned by %s", v.SchemaName.String, v.SchemaOwner.String)
	}

	var tables []TTables
	n, err = dbS.SelectBySql("SELECT * FROM information_schema.tables WHERE table_schema IN ?", schemaWhitelist).Load(&tables)
	if err != nil {
		log.WithError(err).Fatal("select tables")
	}
	if n == 0 {
		log.Warn("no tables available")
		return
	}
	for _, v := range tables {
		log.Debugf("table %s.%s type %s", v.TableSchema.String, v.TableName.String, v.TableType.String)
	}

	var columns []TColumns
	n, err = dbS.SelectBySql("SELECT * FROM information_schema.columns WHERE table_schema IN ?", schemaWhitelist).Load(&columns)
	if err != nil {
		log.WithError(err).Fatal("select columns")
	}
	if n == 0 {
		log.Warn("no columns available")
		return
	}
	for _, v := range columns {
		log.Debugf("column %s.%s.%s", v.TableSchema.String, v.TableName.String, v.ColumnName.String)
	}
}

type Column struct {
	Name       string
	ParseValue interface{}
}

type Table struct {
	Schema string
	Name   string

	Columns []Column
	FKs     []ForeignKey
	PK      PrimaryKey
}

type ForeignKey struct{}
type PrimaryKey struct{}

// /:schema/:table
