package economy

import (
	"fmt"
	"slices"

	"github.com/janithl/citylyf/internal/entities"
)

// Employment handles job assignments
type Employment struct {
	CompanyService *CompanyService
}

// AssignJobs assigns unemployed people to jobs
func (e *Employment) AssignJobs() {
	for _, person := range entities.Sim.People.People {
		retirement := person.ConsiderRetirement(e.CompanyService.RemoveEmployeeFromCompany)
		if retirement {
			fmt.Printf("[  Job ] %s %s (%d) has retired\n", person.FirstName, person.FamilyName, person.Age())
			continue
		}

		if person.IsEmployable() && !person.IsEmployed() {
			if companyID, remaining := e.findSuitableJob(*person); companyID != 0 {
				e.CompanyService.AddEmployeeToCompany(companyID, person.ID)
				person.EmployerID = companyID
				fmt.Printf("[  Job ] %s %s has accepted a job as %s, %d jobs remain\n",
					person.FirstName, person.FamilyName, person.Occupation, remaining)
			}
		}
	}
}

// findSuitableJob finds an appropriate job for a person based on their industry and career level
func (e *Employment) findSuitableJob(p entities.Person) (companyID int, remaining int) {
	for _, company := range entities.Sim.Companies {
		if company.Industry == p.Industry {
			for i, opening := range company.JobOpenings {
				if opening.Job == p.Occupation {
					company.JobOpenings = slices.Delete(company.JobOpenings, i, i+1)
					entities.Sim.Companies[company.ID] = company
					return company.ID, len(company.JobOpenings)
				}
			}
		}
	}

	return 0, 0 // No suitable job found
}
