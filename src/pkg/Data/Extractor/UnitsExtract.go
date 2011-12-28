package Extractor


func ExtractUnits(path string, outpath string, NttExtractDone chan bool) {
	defer Panic()
	defer func() {
		NttExtractDone <- true
	}()
	
}