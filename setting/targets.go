package setting

import (
	"bufio"
	"github.com/vpxuser/proxy"
	"os"
)

func loadTargets(path string, targets *map[string]struct{}) {
	file, err := os.Open(path)
	if err != nil {
		proxy.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		(*targets)[scanner.Text()] = struct{}{}
	}

	if err = scanner.Err(); err != nil {
		proxy.Fatal(err)
	}
}
