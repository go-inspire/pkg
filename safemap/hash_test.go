/*
 * Copyright 2024 Enoch <lanxenet@gmail.com>. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 */

package safemap

import (
	"strings"
	"testing"
)

func TestHashString_EmptyString(t *testing.T) {
	result := hashString("")
	if result == 0 {
		t.Errorf("expected non-zero hash for empty string, got %d", result)
	}
}

func TestHashString_SameStringSameHash(t *testing.T) {
	str := "test"
	hash1 := hashString(str)
	hash2 := hashString(str)
	if hash1 != hash2 {
		t.Errorf("expected same hash for same string, got %d and %d", hash1, hash2)
	}
}

func TestHashString_DifferentStringsDifferentHashes(t *testing.T) {
	hash1 := hashString("test1")
	hash2 := hashString("test2")
	if hash1 == hash2 {
		t.Errorf("expected different hashes for different strings, got %d and %d", hash1, hash2)
	}
}

func TestHashString_LongString(t *testing.T) {
	longStr := "a" + strings.Repeat("b", 1000) + "c"
	result := hashString(longStr)
	if result == 0 {
		t.Errorf("expected non-zero hash for long string, got %d", result)
	}
}

var instids = []string{
	"AP410",
	"AP411",
	"AP412",
	"AP501",
	"AP503",
	"AP504",
	"AP505",
	"CF409",
	"CF411",
	"CF501",
	"CF503",
	"CF505",
	"CF507",
	"CJ409",
	"CJ412",
	"CJ501",
	"CJ503",
	"CJ505",
	"CJ507",
	"CY408",
	"CY409",
	"CY410",
	"CY411",
	"CY412",
	"CY501",
	"CY502",
	"CY503",
	"CY504",
	"CY505",
	"CY506",
	"CY507",
	"FG408",
	"FG409",
	"FG410",
	"FG411",
	"FG412",
	"FG501",
	"FG502",
	"FG503",
	"FG504",
	"FG505",
	"FG506",
	"FG507",
	"IC2407",
	"IC2408",
	"IC2409",
	"IC2412",
	"IC2503",
	"IF2407",
	"IF2408",
	"IF2409",
	"IF2412",
	"IF2503",
	"IH2407",
	"IH2408",
	"IH2409",
	"IH2412",
	"IH2503",
	"IM2407",
	"IM2408",
	"IM2409",
	"IM2412",
	"IM2503",
	"JR409",
	"JR411",
	"JR501",
	"JR503",
	"JR505",
	"JR507",
	"LR409",
	"LR411",
	"LR501",
	"LR503",
	"LR505",
	"LR507",
	"MA408",
	"MA409",
	"MA410",
	"MA411",
	"MA412",
	"MA501",
	"MA502",
	"MA503",
	"MA504",
	"MA505",
	"MA506",
	"MA507",
	"OI409",
	"OI411",
	"OI501",
	"OI503",
	"OI505",
	"OI507",
	"PF408",
	"PF409",
	"PF410",
	"PF411",
	"PF412",
	"PF501",
	"PF502",
	"PF503",
	"PF504",
	"PF505",
	"PF506",
	"PF507",
	"PK410",
	"PK411",
	"PK412",
	"a2409",
	"a2411",
	"a2501",
	"a2503",
	"a2505",
	"a2507",
	"ag2407",
	"ag2408",
	"ag2409",
	"ag2410",
	"ag2411",
	"ag2412",
	"ag2501",
	"ag2502",
	"ag2503",
	"ag2504",
	"ag2505",
	"ag2506",
	"ag2507",
	"agefp",
	"al2407",
	"al2408",
	"al2409",
	"al2410",
	"al2411",
	"al2412",
	"al2501",
	"al2502",
	"al2503",
	"al2504",
	"al2505",
	"al2506",
	"al2507",
	"alefp",
	"ao2407",
	"ao2408",
	"ao2409",
	"ao2410",
	"ao2411",
	"ao2412",
	"ao2501",
	"ao2502",
	"ao2503",
	"ao2504",
	"ao2505",
	"ao2506",
	"ao2507",
	"aoefp",
	"au2407",
	"au2408",
	"au2409",
	"au2410",
}

func TestHashString(t *testing.T) {
	stats := make(map[uint64]int, 10)
	for _, instid := range instids {
		hash := hashString(instid)
		stats[hash%10] = stats[hash%10] + 1
		//t.Logf("hash of %s is %d, %d", instid, hash, hash%10)
	}
	t.Log(stats)
}
