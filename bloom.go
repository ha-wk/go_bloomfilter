package main

import (
	"fmt"
	"hash/crc32"
	"hash/fnv"
	"math"

	bitarray "github.com/Workiva/go-datastructures/bitarray" //TO IMPORT BIT ARRAY PACKAGE
)

type bf struct {
	set bitarray.BitArray
	k   uint64  //STORE NUMBER OF HASH FUNCTIONS
	n   uint64  //STORE NUMBER OF ITEMS
	m   uint64  //TO STORE SIZE OF BIT ARRAY
	p   float64 //TO STORE FALSE POSITIVE PROBABILITY
}

//CREATING NEW BLOOM FILTER BY CALCULATING K,M based on N,P

func BloomFilter(size uint64, falseprob float64) *bf {
	_m := -(float64(size) * math.Log(falseprob) / (math.Log(2) * math.Log(2)))
	m1 := uint64(_m)
	_k := (_m / float64(size)) * math.Log(2)
	k1 := uint64(_k)
	_set := bitarray.NewBitArray(m1)

	return &bf{
		set: _set,
		k:   k1,
		m:   m1,
		p:   falseprob,
		n:   size,
	}
}

type Interface interface { //BASIC OPERATIONS OF BLOOM FILTERS
	Add(item []byte)
	Check(item []byte) bool
	Details()
}

func (BF bf) Add(item string) { //TO ADD NEW DATA TO OUR EXISTING SET
	digests := []int{}
	m := BF.m
	k := BF.k
	set := BF.set

	for i := 0; i < int(k); i++ {
		digest := getIndex(item, i, m)
		digests = append(digests, int(digest))
		set.SetBit(digest)
	}
}

func (BF bf) Check(item string) bool { //TO CHECK WHETHER THE GIVEN ITEM'S HASHVALUE IS PRESENT IN BIT ARRAY OR NOT
	k := BF.k
	m := BF.m
	set := BF.set

	for i := 0; i < int(k); i++ {
		digest := getIndex(item, i, m)
		valAtPosition, error := set.GetBit(digest)
		if error != nil {
			fmt.Errorf("error occured at %d", digest)
		}
		if valAtPosition == false {
			return false //IF NOT,WE RETURN FALSE
		}
	}

	return true //ELSE,WE RETURN TRUE
}
func (BF bf) Details() { //THIS FUNC PRINT ALL DETAILS LIKE N,M,K,P
	fmt.Printf("n=%d\nm=%d\nk=%d\np=%f\n", BF.n, BF.m, BF.k, BF.p)
}

//IMPORTANT FUNCTION!
//IT GENERALLY TAKES TWO HASH FUNCTION(FNV AND CRC) TO GENERATE HASH VALUES

func getIndex(item string, i int, m uint64) uint64 { //IMPORTANT FUNCTION

	hash := fnv.New32()
	hash.Write([]byte(item)) //HASH FUNC 1
	val := hash.Sum32()

	hash2 := crc32.NewIEEE() //HASH FUNC 2
	hash2.Write([]byte(item))
	val2 := hash2.Sum32()

	return uint64((int(val) + int(val2)*i) % int(m))
}

func main() {
	bft := BloomFilter(200, 0.0001)
	bft.Details()

	bft.Add("http://fraud.com")
	bft.Add("http://insecureweb.com") //STORING DUMMY DATA(WEBSITES)
	bft.Add("http://webhack.com")

	fmt.Printf("Is Malicious ? %t\n", bft.Check("https://www.maniartech.com"))
	fmt.Printf("Is Malicious ? %t\n", bft.Check("http://insecureweb.com")) //CHECKING IT'S PRESENCE

}
