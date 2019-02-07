// Package powerwalk concurrently walks file trees.
// Aside from SkipDir functionality not working and the fact that the
// WalkFunc gets run concurrently, this is a drop-in replacement
// for filepath.Walk.
package powerwalk
