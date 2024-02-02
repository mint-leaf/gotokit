package main

func Filter[S ~[]T, T any](s S, fs ...func(ele T) bool) S {
	c := make([]T, 0)
	for index := range s {
		for i := range fs {
			if fs[i](s[index]) {
				c = append(c, s[index])
			}
		}
	}
	return c
}

func Mapper[S ~[]T, N []V, T any, V any](s S, f func(ele T) V) N {
	c := make([]V, len(s))
	for index := range s {
		c[index] = f(s[index])
	}
	return c
}
