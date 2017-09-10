package main

type Job struct {
	filename string
	results  chan<- Result
}

func (job Job) Do(lineRx *regexp.Regexp) {
	file, err := os.Open(job.filename)
	if err != nil {
		log.Printf("error: %s\n", err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for lino := 1; ; lino++ {
		line, err := reader.ReadBytes('\n')
		line = bytes.TrimRight(line, "\n\r")
		if lineRx.Match(line) {
			job.results <- Result{job.filename, lino, string(line)}
		}
		if err != nil {
			if err != io.EOF {
				log.Printf("error:%d: %s\n", lino, err)
			}
			break
		}
	}
}

type Result struct {
	filename string
	lino     int
	line     string
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Use all the machine's cores
	if len(os.Args) < 3 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s <regexp> <files>\n",
			filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	if lineRx, err := regexp.Compile(os.Args[1]); err != nil {
		log.Fatalf("invalid regexp: %s\n", err)
	} else {
		grep(lineRx, commandLineFiles(os.Args[2:]))
	}
}

var workers = runtime.NumCPU()

func grep(lineRx *regexp.Regexp, filenames []string) {
	jobs := make(chan Job, workers)
	results := make(chan Result, minimum(1000, len(filenames)))
	done := make(chan struct{}, workers)
	go addJobs(jobs, filenames, results) // Executes in its own goroutine
	for i := 0; i < workers; i++ {
		go doJobs(done, lineRx, jobs) // Each executes in its own goroutine
	}
	go awaitCompletion(done, results) // Executes in its own goroutine
	processResults(results)           // Blocks until the work is done
}

/*
Send all the jobs into the "jobs" channel and close the channel, which will
notify the "doJobs" worker goroutine to break the for...range loop
*/
func addJobs(jobs chan<- Job, filenames []string, results chan<- Result) {
	for _, filename := range filenames {
		jobs <- Job{filename, results}
	}
	close(jobs)
}

/*
it is a Go convention that functions that have channel parameters
have the destination channels first, followed by the source channels.
*/
func doJobs(done chan<- struct{}, lineRx *regexp.Regexp, jobs <-chan Job) {
	for job := range jobs {
		job.Do(lineRx)
	}
	// since all that matters is whether something is sent, we have made the
	// channel’s type chan struct{} (i.e., an empty struct) to more clearly
	// express our semantics
	done <- struct{}{}
}

/*
Waiting for all the workers completing their jobs and close the results channel.
*/
func awaitCompletion(done <-chan struct{}, results chan Result) {
	for i := 0; i < workers; i++ {
		<-done
	}
	close(results)
}

func processResults(results <-chan Result) {
	for result := range results {
		fmt.Printf("%s:%d:%s\n", result.filename, result.lino, result.line)
	}
}

func waitAndProcessResults_1(done <-chan struct{}, results <-chan Result) {
	for working := workers; working > 0; {
		select { // Blocking
		case result := <-results:
			fmt.Printf("%s:%d:%s\n", result.filename, result.lino,
				result.line)
		case <-done:
			working--
		}
	}
DONE:
	for {
		select { // Nonblocking
		case result := <-results:
			fmt.Printf("%s:%d:%s\n", result.filename, result.lino, result.line)
		default:
			break DONE
		}
	}
}

func waitAndProcessResults_2(timeout int64, done <-chan struct{},
	results <-chan Result) {
	finish := time.After(time.Duration(timeout))
	for working := workers; working > 0; {
		select { // Blocking
		case result := <-results:
			fmt.Printf("%s:%d:%s\n", result.filename, result.lino, result.line)
		case <-finish:
			fmt.Println("timed out")
			return // Time's up so finish with what results there were
		case <-done:
			working--
		}
	}
	for {
		select { // Nonblocking
		case result := <-results:
			fmt.Printf("%s:%d:%s\n", result.filename, result.lino,
				result.line)
		case <-finish:
			fmt.Println("timed out")
			return // Time's up so finish with what results there were
		default:
			return
		}
	}
}
