package economy

import (
	"fmt"

	"github.com/janithl/citylyf/internal/entities"
)

// Employment handles job assignments
type Employment struct{}

// AssignJobs assigns unemployed people to jobs
func (e *Employment) AssignJobs() {
	for i := range entities.Sim.People.Households {
		for j := range entities.Sim.People.Households[i].Members {
			person := &entities.Sim.People.Households[i].Members[j]
			retirement := person.ConsiderRetirement(RemoveEmployeeFromCompany)
			if retirement {
				fmt.Printf("[  Job ] %s %s (%d) has retired\n", person.FirstName, person.FamilyName, person.Age())
				continue
			}

			if person.IsEmployable() && !person.IsEmployed() {
				if companyID, remaining := e.findSuitableJob(*person); companyID != 0 {
					AddEmployeeToCompany(companyID, person.ID)
					person.EmployerID = companyID
					fmt.Printf("[  Job ] %s %s has accepted a job as %s, %d jobs remain\n",
						person.FirstName, person.FamilyName, person.Occupation, remaining)
				}
			}
		}
	}
}

// findSuitableJob finds an appropriate job for a person based on their industry and career level
func (e *Employment) findSuitableJob(p entities.Person) (companyID int, remaining int) {
	for i := range entities.Sim.Companies {
		company := &entities.Sim.Companies[i]
		if company.Industry == p.Industry {
			if openings, exists := company.JobOpenings[p.CareerLevel]; exists && openings > 0 {
				company.JobOpenings[p.CareerLevel]--
				return company.ID, company.GetNumberOfJobOpenings()
			}
		}
	}
	return 0, 0 // No suitable job found
}
