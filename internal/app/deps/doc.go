// Package deps is the composition root: it wires the entire dependency graph in
// one place (samber/do/v2) and hands the runtime a Container it can resolve
// services from and tear down in reverse order.
//
// Providers are registered lazily; New then forces every long-lived service to
// be instantiated so wiring, configuration and connection errors surface at
// startup rather than on the first Telegram update. samber/do records the
// resulting invocation order and calls each service's Shutdown hook in reverse
// on Container.Shutdown — there is no hand-rolled shutdown registry.
package deps
