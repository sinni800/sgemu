package Extractor


func ExtractUnits(path string, outpath string, UnitExtractDone chan bool) {
	defer Panic()
	defer func() {
		UnitExtractDone <- true
	}()
	
}