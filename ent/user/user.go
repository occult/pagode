// Code generated by ent, DO NOT EDIT.

package user

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldPassword holds the string denoting the password field in the database.
	FieldPassword = "password"
	// FieldVerified holds the string denoting the verified field in the database.
	FieldVerified = "verified"
	// FieldAdmin holds the string denoting the admin field in the database.
	FieldAdmin = "admin"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeOwner holds the string denoting the owner edge name in mutations.
	EdgeOwner = "owner"
	// EdgePaymentCustomer holds the string denoting the payment_customer edge name in mutations.
	EdgePaymentCustomer = "payment_customer"
	// Table holds the table name of the user in the database.
	Table = "users"
	// OwnerTable is the table that holds the owner relation/edge.
	OwnerTable = "password_tokens"
	// OwnerInverseTable is the table name for the PasswordToken entity.
	// It exists in this package in order to avoid circular dependency with the "passwordtoken" package.
	OwnerInverseTable = "password_tokens"
	// OwnerColumn is the table column denoting the owner relation/edge.
	OwnerColumn = "user_id"
	// PaymentCustomerTable is the table that holds the payment_customer relation/edge.
	PaymentCustomerTable = "users"
	// PaymentCustomerInverseTable is the table name for the PaymentCustomer entity.
	// It exists in this package in order to avoid circular dependency with the "paymentcustomer" package.
	PaymentCustomerInverseTable = "payment_customers"
	// PaymentCustomerColumn is the table column denoting the payment_customer relation/edge.
	PaymentCustomerColumn = "payment_customer_user"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldEmail,
	FieldPassword,
	FieldVerified,
	FieldAdmin,
	FieldCreatedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "users"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"payment_customer_user",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "github.com/occult/pagode/ent/runtime"
var (
	Hooks [1]ent.Hook
	// NameValidator is a validator for the "name" field. It is called by the builders before save.
	NameValidator func(string) error
	// EmailValidator is a validator for the "email" field. It is called by the builders before save.
	EmailValidator func(string) error
	// PasswordValidator is a validator for the "password" field. It is called by the builders before save.
	PasswordValidator func(string) error
	// DefaultVerified holds the default value on creation for the "verified" field.
	DefaultVerified bool
	// DefaultAdmin holds the default value on creation for the "admin" field.
	DefaultAdmin bool
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
)

// OrderOption defines the ordering options for the User queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByEmail orders the results by the email field.
func ByEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmail, opts...).ToFunc()
}

// ByPassword orders the results by the password field.
func ByPassword(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPassword, opts...).ToFunc()
}

// ByVerified orders the results by the verified field.
func ByVerified(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldVerified, opts...).ToFunc()
}

// ByAdmin orders the results by the admin field.
func ByAdmin(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldAdmin, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByOwnerCount orders the results by owner count.
func ByOwnerCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newOwnerStep(), opts...)
	}
}

// ByOwner orders the results by owner terms.
func ByOwner(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOwnerStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByPaymentCustomerField orders the results by payment_customer field.
func ByPaymentCustomerField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPaymentCustomerStep(), sql.OrderByField(field, opts...))
	}
}
func newOwnerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OwnerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, true, OwnerTable, OwnerColumn),
	)
}
func newPaymentCustomerStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PaymentCustomerInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, true, PaymentCustomerTable, PaymentCustomerColumn),
	)
}
