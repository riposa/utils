package utils

import (
	"testing"
	"time"
)

func TestCheckFieldLength(t *testing.T) {
	var test = "test content"
	CheckFieldLength(test, 5, "test")
}

func TestNewTimeIt(t *testing.T) {
	timeIt := NewTimeIt("test")
	defer timeIt.End()
}

func TestCombinationResult(t *testing.T) {
	CombinationResult(10, 12)
}

func TestFindNumsByIndexes(t *testing.T) {
	FindNumsByIndexes([]int{1, 2, 3}, [][]int{{1, 2, 3}})
}

func TestSubstr(t *testing.T) {
	Substr("test-content", 2, 3)
}

func TestGetWxaCode(t *testing.T) {
	t1 := time.Now().UnixNano()
	token := "13_kHs8yUEtvzWaBwDMT7cg_INelQD1_aJIG4gxngOcz5TwLNRSHYTi-lDaAR1oE8FgTTxqS71uRK3f-0jN834jvvXfKEncUaWek5izGilDV0niKuq7ANS1nRNz9RPs1Be7w9qkdrMTuLoYZR9dNOHhAFAJKD"
	scene := "S111"
	page := ""
	utilsLogger.Info(GetWxaCode(token, scene, page))
	utilsLogger.Infof("cost: %d", (time.Now().UnixNano()-t1)/1e6)
}
