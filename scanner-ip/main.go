package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type ScanResult struct {
	IP       string
	Port     int
	Status   string
	Duration time.Duration
	Error    error
}

func tcping(ip string, port int, timeout time.Duration) ScanResult {
	address := net.JoinHostPort(ip, strconv.Itoa(port))
	start := time.Now()

	conn, err := net.DialTimeout("tcp", address, timeout)
	duration := time.Since(start)

	result := ScanResult{
		IP:       ip,
		Port:     port,
		Duration: duration,
	}

	if err != nil {
		result.Status = "FAILED"
		result.Error = err
		return result
	}

	conn.Close()
	result.Status = "SUCCESS"
	return result
}

func scanPorts(ip string, ports []int, timeout time.Duration, verbose bool) {
	fmt.Printf("Scanning %s...\n\n", ip)

	var successCount, failedCount int
	var totalDuration time.Duration

	for _, port := range ports {
		result := tcping(ip, port, timeout)

		if result.Status == "SUCCESS" {
			successCount++
			totalDuration += result.Duration
			fmt.Printf("Probing %s:%d/tcp - Port is open - time=%s\n",
				result.IP, result.Port, formatDuration(result.Duration))
		} else {
			failedCount++
			if verbose {
				fmt.Printf("Probing %s:%d/tcp - Port is closed - %s\n",
					result.IP, result.Port, result.Error.Error())
			} else {
				fmt.Printf("Probing %s:%d/tcp - Port is closed\n",
					result.IP, result.Port)
			}
		}
	}

	// Summary report
	fmt.Printf("\n--- %s tcping statistics ---\n", ip)
	fmt.Printf("%d ports scanned, %d open, %d closed\n",
		len(ports), successCount, failedCount)

	if successCount > 0 {
		avgDuration := totalDuration / time.Duration(successCount)
		fmt.Printf("Average response time: %s\n", formatDuration(avgDuration))
	}
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.3fms", float64(d.Nanoseconds())/1000000.0)
	} else if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000.0)
	} else {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
}

func parsePorts(portStr string) ([]int, error) {
	var ports []int

	// Handle comma-separated or space-separated ports
	portStr = strings.ReplaceAll(portStr, ",", " ")
	parts := strings.Fields(portStr)

	for _, part := range parts {
		// Handle port ranges (e.g., 8080-8090)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start port: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end port: %s", rangeParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("start port must be <= end port")
			}

			for p := start; p <= end; p++ {
				if p < 1 || p > 65535 {
					return nil, fmt.Errorf("port %d out of range (1-65535)", p)
				}
				ports = append(ports, p)
			}
		} else {
			port, err := strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}

			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("port %d out of range (1-65535)", port)
			}

			ports = append(ports, port)
		}
	}

	return ports, nil
}

func readPortsFromFile(filename string) ([]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	var portStrs []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		portStrs = append(portStrs, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	if len(portStrs) == 0 {
		return nil, fmt.Errorf("no ports found in file")
	}

	portStr := strings.Join(portStrs, " ")
	return parsePorts(portStr)
}

func main() {
	var (
		timeout  = flag.Duration("timeout", 3*time.Second, "Connection timeout")
		verbose  = flag.Bool("v", false, "Verbose output (show error messages)")
		help     = flag.Bool("h", false, "Show help message")
		portFile = flag.String("f", "", "Read ports from file (one port per line or comma/space separated)")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <IP> <port1> [port2] [port3] ...\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Examples:\n")
		fmt.Fprintf(os.Stderr, "  %s 1.1.1.1 8081\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s 1.1.1.1 8081 8082 8083\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s 1.1.1.1 8080-8090\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s 1.1.1.1 8081,8082,8083\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f port.txt 1.1.1.1\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f port.txt 1.1.1.1 8081 8082\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "Error: IP address is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	ip := flag.Arg(0)

	// Validate IP address
	if net.ParseIP(ip) == nil {
		fmt.Fprintf(os.Stderr, "Error: Invalid IP address: %s\n", ip)
		os.Exit(1)
	}

	// Collect ports from file and/or command-line arguments
	var allPorts []int

	// Read ports from file if specified
	if *portFile != "" {
		filePorts, err := readPortsFromFile(*portFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading ports from file: %v\n", err)
			os.Exit(1)
		}
		allPorts = append(allPorts, filePorts...)
	}

	// Parse ports from command-line arguments
	if flag.NArg() > 1 {
		var portStrs []string
		for i := 1; i < flag.NArg(); i++ {
			portStrs = append(portStrs, flag.Arg(i))
		}

		portStr := strings.Join(portStrs, " ")
		cmdPorts, err := parsePorts(portStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing ports: %v\n", err)
			os.Exit(1)
		}
		allPorts = append(allPorts, cmdPorts...)
	}

	// Remove duplicates
	portMap := make(map[int]bool)
	var ports []int
	for _, port := range allPorts {
		if !portMap[port] {
			portMap[port] = true
			ports = append(ports, port)
		}
	}

	if len(ports) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No valid ports specified\n")
		os.Exit(1)
	}

	scanPorts(ip, ports, *timeout, *verbose)
}
