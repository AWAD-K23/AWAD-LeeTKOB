package main

import "leetkob/backend/internal/localenv"

// Keep the server entry point while sharing automatic loading with the worker.
func loadLocalEnv() error { return localenv.Load() }
