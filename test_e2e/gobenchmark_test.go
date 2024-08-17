package runner_test

import (
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/leonardinius/goloxvm/internal/tests"
)

func BenchmarkAll(b *testing.B) {
	workDir, err := tests.ProjectDir()
	if err != nil {
		b.Fatalf("Failed to get absolute path: %v", err)
	}
	gobin1 := buildGobin1(b, workDir)

	benchmarks := []string{
		"testdata/benchmark/binary_trees.lox",
		"testdata/benchmark/equality.lox",
		"testdata/benchmark/fib.lox",
		"testdata/benchmark/instantiation.lox",
		"testdata/benchmark/invocation.lox",
		"testdata/benchmark/method_call.lox",
		"testdata/benchmark/properties.lox",
		// "testdata/benchmark/string_equality.lox", // go constant limit
		"testdata/benchmark/trees.lox",
		// "testdata/benchmark/zoo_batch.lox", // always take 10 seconds and reports throughput
		"testdata/benchmark/zoo.lox",
	}

	b.ResetTimer()
	for _, bench := range benchmarks {
		b.Run("GOBIN1/"+bench, func(b *testing.B) {
			runBenchN(b, workDir, gobin1, bench)
		})
	}
}

func runBenchN(b *testing.B, workDir string, args ...string) {
	b.Helper()
	for n := 0; n < b.N; n++ {
		runBench(b, workDir, args...)
	}
}

func runBench(b *testing.B, workDir string, args ...string) {
	b.Helper()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = workDir
	stdout := new(strings.Builder)
	stderr := new(strings.Builder)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	func() {
		defer func() {
			if err := recover(); err != nil {
				b.Errorf("Execute error %v: %v", cmd, err)
				return
			}
		}()
		if err := cmd.Run(); err != nil {
			b.Errorf("Execute error %v: %v", cmd, err)
			return
		}
	}()

	// Run binary
	exitCode := cmd.ProcessState.ExitCode()
	outputLines := strings.Split(stdout.String(), "\n")
	errorLines := strings.Split(stderr.String(), "\n")
	for len(outputLines) > 0 && outputLines[len(outputLines)-1] == "" {
		outputLines = outputLines[:len(outputLines)-1]
	}
	for len(errorLines) > 0 && errorLines[len(errorLines)-1] == "" {
		errorLines = errorLines[:len(errorLines)-1]
	}

	if exitCode != 0 || len(errorLines) > 0 {
		b.Errorf("Command %v exited with code %v\nerror:\n%v", cmd, exitCode, stderr)
		return
	}

	elapsedTimeString := outputLines[len(outputLines)-1]
	elapsedTimeSeconds, err := strconv.ParseFloat(elapsedTimeString, 64)
	if err != nil {
		b.Errorf("Failed to parse elapsed time %v", elapsedTimeString)
		return
	}
	b.ReportMetric(elapsedTimeSeconds, "elapsed/op")
}

func buildGobin1(b *testing.B, workDir string) string {
	b.Helper()
	mainGo := workDir + "/main.go"
	gobin1 := workDir + "/bin/golox"
	cmd := exec.Command("go", "build", "-o", gobin1, mainGo)
	if outbytes, err := cmd.CombinedOutput(); err != nil {
		out := string(outbytes)
		b.Fatalf("go build failed with %v: %#v\n", err, out)
	}
	return gobin1
}
