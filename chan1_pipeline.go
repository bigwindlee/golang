package main

func main() {
	minSize, maxSize, suffixes, files := handleCommandLine()
	sink(filterSize(minSize, maxSize, filterSuffixes(suffixes, source(files))))
}

// the type <-chan Type is a receive-only channel
// this way we have better expressed our intentions
func source(files []string) <-chan string {
	out := make(chan string, 1000)
	go func() {
		for _, filename := range files {
			out <- filename
		}
		close(out)
	}()
	return out
}

// by including directions we have precisely expressed the semantics
// we want the function to have —— and ensured that the compiler enforces them
func filterSuffixes(suffixes []string, in <-chan string) <-chan string {
	out := make(chan string, cap(in))
	go func() {
		for filename := range in {
			if len(suffixes) == 0 {
				out <- filename
				continue
			}
			ext := strings.ToLower(filepath.Ext(filename))
			for _, suffix := range suffixes {
				if ext == suffix {
					out <- filename
					break
				}
			}
		}
		close(out)
	}()
	return out
}

func sink(in <-chan string) {
	for filename := range in {
		fmt.Println(filename)
	}
}
