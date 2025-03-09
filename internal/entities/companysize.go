package entities

import "math/rand/v2"

type CompanySize string

const (
	Micro CompanySize = "Micro"
	SME   CompanySize = "SME"
	Large CompanySize = "Large"
)

func (r CompanySize) GetBaseRevenue() float64 {
	if r == Micro {
		return 800_000 + rand.NormFloat64()*400_000 // 400K - 1.2M for micro businesses
	}
	if r == SME {
		return 3_000_000 + rand.NormFloat64()*1_000_000 // 2M - 4M for SME businesses
	}
	return 7_500_000 + rand.NormFloat64()*2_500_000 // 5M - 10M for large businesses
}

func (r CompanySize) GetBaseJobs() map[CareerLevel]int {
	if r == Micro { // Micro companies have < 15 jobs
		return map[CareerLevel]int{
			EntryLevel:     rand.IntN(3) + 3, // 2-5 jobs
			MidLevel:       rand.IntN(3) + 1, // 2-3 jobs
			SeniorLevel:    rand.IntN(2),     // 0-1 jobs
			ExecutiveLevel: 0,                // 0 jobs
		}
	}
	if r == SME { // SME companies have < 50 jobs
		return map[CareerLevel]int{
			EntryLevel:     rand.IntN(11) + 10, // 10-20 jobs
			MidLevel:       rand.IntN(11) + 5,  // 5-15 jobs
			SeniorLevel:    rand.IntN(5) + 4,   // 4-8 jobs
			ExecutiveLevel: rand.IntN(2) + 1,   // 1-2 jobs
		}
	} // Large companies have <100 jobs
	return map[CareerLevel]int{
		EntryLevel:     rand.IntN(31) + 15, // 15-45 jobs
		MidLevel:       rand.IntN(15) + 11, // 11-25 jobs
		SeniorLevel:    rand.IntN(10) + 3,  // 3-12 jobs
		ExecutiveLevel: rand.IntN(3) + 2,   // 2-4 jobs
	}
}

var companysizes = []CompanySize{
	Micro, SME, Large,
}

func GetRandomCompanySize() CompanySize {
	return companysizes[rand.IntN(len(companysizes))]
}
