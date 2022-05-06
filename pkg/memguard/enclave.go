package memguard

import "github.com/awnumar/memguard"

func invert(key *memguard.Enclave) *memguard.Enclave {
	// Decrypt the key into a local copy
	b, err := key.Open()
	if err != nil {
		memguard.SafePanic(err)
	}
	defer b.Destroy() // Destroy the copy when we return

	// Open returns the data in an immutable buffer, so make it mutable
	b.Melt()

	// Set every element to its complement
	for i := range b.Bytes() {
		b.Bytes()[i] = ^b.Bytes()[i]
	}

	// Return the new data in encrypted form
	return b.Seal() // <- sealing also destroys b
}
