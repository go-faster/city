// Package cityhash provides hash functions for strings.  The functions mix
// the input bits thoroughly but are not suitable for cryptography.  See
// "Hash Quality," below, for details on how CityHash was tested and so on.
//
// All members of the CityHash family were designed with heavy reliance
// on previous work by Austin Appleby, Bob Jenkins, and others.
// For example, Hash32 has many similarities with Murmur3a.
package cityhash
