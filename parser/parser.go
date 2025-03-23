// implementation is based on https://r2p.dev/b/2024-03-18-1brc-go/ all credits to the author

package main

import (
	"fmt"
	"io"
	"maps"
	"os"
	"runtime/pprof"
	"slices"
	"sync"
	"time"
)

const READ_BUFFER_SIZE = 2048 * 2048
const N_WORKERS = 75

type TrashItem struct {
	Idx     int
	Value   []byte
	Initial bool
}

var lock = &sync.Mutex{}
var lockIdx = 0

func hash(name []byte) uint64 {
	var h uint64 = 5381
	for _, b := range name {
		h = (h << 5) + h + uint64(b)
	}
	return h
}

func consumer(file *os.File, trash chan *TrashItem, output chan map[uint64]any, wg *sync.WaitGroup) {
	defer wg.Done()
	data := make(map[uint64]any, 1024)

	readBuffer := make([]byte, READ_BUFFER_SIZE)
	for {
		lock.Lock()
		lockIdx++
		idx := lockIdx
		n, err := file.Read(readBuffer)
		lock.Unlock()

		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		// ignoring first line
		start := 0
		for i := 0; i < n; i++ {
			if readBuffer[i] == 10 {
				start = i + 1
				break
			}
		}
		trash <- &TrashItem{idx - 1, readBuffer[:start], false}

		// ignoring last line
		final := 0
		for i := n - 1; i >= 0; i-- {
			if readBuffer[i] == 10 {
				final = i
				break
			}
		}
		trash <- &TrashItem{idx, readBuffer[final+1 : n], true}

		readingIndex := start
		for readingIndex < final {
			next := nextEnterpriseLine(readingIndex, readBuffer, data)
			readingIndex = next
		}
	}

	output <- data
}

func saveCan(can []*TrashItem, data map[uint64]any, buffer []byte) []*TrashItem {
	for i, ref := range can {
		if ref.Idx == 0 {
			_ = nextEnterpriseLine(0, ref.Value, data)
			return slices.Delete(can, i, i+1)
		}

		for j, oth := range can {
			if ref.Idx == oth.Idx && i != j {
				if ref.Initial {
					copy(buffer[:len(ref.Value)], ref.Value)
					copy(buffer[len(ref.Value):], oth.Value)
				} else {
					copy(buffer[:len(oth.Value)], oth.Value)
					copy(buffer[len(oth.Value):], ref.Value)
				}
				total := len(ref.Value) + len(oth.Value)

				end := nextEnterpriseLine(0, buffer, data)
				if end < total {
					_ = nextEnterpriseLine(end, buffer, data)
				}

				if i > j {
					can = slices.Delete(can, i, i+1)
					can = slices.Delete(can, j, j+1)
				} else {
					can = slices.Delete(can, j, j+1)
					can = slices.Delete(can, i, i+1)
				}

				return can
			}
		}
	}

	return can
}

func run() {
	// Read file
	file, err := os.Open("./testdata/TEST.EMPRECSV")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	outputChannels := make([]chan map[uint64]any, N_WORKERS+1)

	var wg sync.WaitGroup
	var wgTrash sync.WaitGroup

	wg.Add(N_WORKERS)
	wgTrash.Add(1)
	trash := make(chan *TrashItem, N_WORKERS*2)
	output := make(chan map[uint64]any, 1)
	go trashBin(trash, output, &wgTrash)
	outputChannels[0] = output

	for i := range N_WORKERS {
		output := make(chan map[uint64]any, 1)
		go consumer(file, trash, output, &wg)
		outputChannels[i+1] = output
	}

	wg.Wait()
	close(trash)
	wgTrash.Wait()

	for i := range N_WORKERS + 1 {
		close(outputChannels[i])
	}

	data := make(map[uint64]any, 1000)
	for i := range N_WORKERS + 1 {
		maps.Copy(data, <-outputChannels[i])
	}

	printResult(data)
}

func trashBin(input chan *TrashItem, output chan map[uint64]any, wg *sync.WaitGroup) {
	defer wg.Done()
	data := make(map[uint64]any, 1024)

	can := []*TrashItem{}
	buffer := make([]byte, 1024)

	for item := range input {
		can = append(can, item)
		can = saveCan(can, data, buffer)
	}

	output <- data
}

func printResult(data map[uint64]any) {
	print("[\n")
	for _, k := range data {
		switch v := (k).(type) {
		case *EnterpriseData:
			fmt.Printf(`  {"basic_cnpj":"%s","corporate_name":"%s","legal_nature":"%s","responsible_qualification":"%s","social_capital":"%s","company_size":"%s","federative_entity":"%s"}`+"\n", v.BasicCNPJ, v.CorporateName, v.LegalNature, v.ResponsibleQualification, v.SocialCapital, v.CompanySize, v.FederativeEntity)
		}
	}
	print("]\n")
}

func process() {
	f, err := os.Create("cpu_profile.prof")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := pprof.StartCPUProfile(f); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()

	_ = time.Now()
	run()
	// fmt.Printf("%0.6f\n", time.Since(started).Seconds())
}

func nextEnterpriseLine(readingIndex int, reading []byte, data map[uint64]any) int {
	temp := &EnterpriseData{} // FIXME use dynamic type
	size := temp.Size()
	fields := make([]int, size)

	i := readingIndex + 1 // skip first "

	for j := 0; j < size-2; j += 2 {
		fields[j] = i
		for reading[i] != 59 { // ;
			i++
		}
		fields[j+1] = i - 1 // skip last field "
		i += 2              // skip ;"
	}

	fields[size-2] = i
	for i < len(reading) && reading[i] != 10 { // \n
		i++
	}
	fields[size-1] = i - 1 // skip last field "

	id := temp.ID(reading, fields)
	data[id] = temp.Data(reading, fields)

	readingIndex = i + 1
	return readingIndex
}
