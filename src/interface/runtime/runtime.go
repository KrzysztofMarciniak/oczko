// Package runtime defines the contract for an executable parsed program.
package runtime

// Program is an executable, parsed oczko program.
type Program interface{ Execute() error }
