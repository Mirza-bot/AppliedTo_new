package utils

func Patch[T any](target *T, patch *T) {
    if patch != nil {
        *target = *patch
    }
}

func PatchRef[T any](target **T, patch *T) {
	if patch != nil {
		*target = patch
	}
}
