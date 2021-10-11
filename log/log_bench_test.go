package log

import (
	"fmt"
	"io"
	"log"
	"testing"
)

func Benchmark_Logger_Print(b *testing.B) {
	logger := std.(*zapLogger)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print("No context.")
		}
	})
}

func Benchmark_Logger_Print_With3StringgArgs(b *testing.B) {
	logger := std.(*zapLogger)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print("No context.", "1", "2")
		}
	})
}

func Benchmark_Logger_Print_With5StringArgs(b *testing.B) {
	logger := std.(*zapLogger)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print("No context.", "1", "2", "1", "2")
		}
	})
}

func Benchmark_Logger_Print_With3StringArgsAnd2IntArgs(b *testing.B) {
	logger := std.(*zapLogger)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Print("No context.", "1", "2", 1, 2)
		}
	})
}

func Benchmark_Logger_Printf(b *testing.B) {
	logger := std.(*zapLogger)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("No context.")
		}
	})
}

func Benchmark_Logger_Printf_With2StringArgs(b *testing.B) {
	logger := std.(*zapLogger)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("No context.%s%s", "1", "2")
		}
	})
}

func Benchmark_Logger_Printf_With4StringArgs(b *testing.B) {
	logger := std.(*zapLogger)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("No context.%s%s%s%s", "1", "2", "1", "2")
		}
	})
}

func Benchmark_Logger_Printf_With2StringArgsAnd2IntArgs(b *testing.B) {
	logger := std.(*zapLogger)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Printf("No context.%s%s%d%d", "1", "2", 1, 2)
		}
	})
}

func Benchmark_SysLog_Print(b *testing.B) {
	log.SetOutput(io.Discard)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("No context.")
		}
	})
}

func Benchmark_SysLog_Print_With3StringArgs(b *testing.B) {
	log.SetOutput(io.Discard)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("No context.", "1", "2")
		}
	})
}

func Benchmark_SysLog_Print_With3StringArgsAnd2IntArgs(b *testing.B) {
	log.SetOutput(io.Discard)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Print("No context.", "1", "2", 1, 2)
		}
	})
}

func Benchmark_SysLog_Printf(b *testing.B) {
	log.SetOutput(io.Discard)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Printf("No context.")
		}
	})
}

func Benchmark_SysLog_Printf_With2StringArgs(b *testing.B) {
	log.SetOutput(io.Discard)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Printf("No context.%s%s", "1", "2")
		}
	})
}

func Benchmark_SysLog_Printf_With2StringArgsAnd2IntArgs(b *testing.B) {
	log.SetOutput(io.Discard)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			log.Printf("No context.%s%s%d%d", "1", "2", 1, 2)
		}
	})
}

func Benchmark_fmt_Print(b *testing.B) {
	log.SetOutput(io.Discard)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fmt.Fprint(io.Discard, "No context.")
		}
	})

}

func Benchmark_fmt_Printf(b *testing.B) {
	log.SetOutput(io.Discard)
	b.ResetTimer()
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			fmt.Fprintf(io.Discard, "No context.")
		}
	})

}
