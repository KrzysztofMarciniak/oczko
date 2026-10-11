// Package runner defines the application-facing interpreter contract.
package runner

// Runner parses and executes one oczko program.
type Runner interface{ Run(source string) error }
