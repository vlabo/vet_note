import { Injectable } from '@angular/core';
import * as uuid from 'uuid';


export class Procedure {
  Id: string = "";
  Type: 'Examination' | 'Surgery' | 'Vaccine' | 'Castration' | 'Blood Test' = "Examination";
  Date: Date = new Date();
  Details: String = "";
}

export class Patient {
  Id: string = ""
  Type: string = ""
  Name: string = ""
  Owner: string = ""
  Procedures: string[] = [];
}

@Injectable({
  providedIn: 'root'
})
export class PatientsService {
  private patients = new Map<string, Patient>()
  private procedures = new Map<string, Procedure>()
  private types: String[] = ["Dog", "Cat", "Bird", "Rabbit"];

  private generateMockData() {
    this.procedures.set("0510ee92-b30e-4bff-a6d9-6af70b0e6acc", {
      Id: "0510ee92-b30e-4bff-a6d9-6af70b0e6acc",
      Type: 'Examination',
      Date: new Date('2023-01-15'),
      Details: 'Routine check-up. All vitals are normal.'
    });
    this.procedures.set("b7306c26-8b48-4a8b-83c9-a2425c117364", {
      Id: "b7306c26-8b48-4a8b-83c9-a2425c117364",
      Type: 'Vaccine',
      Date: new Date('2023-02-20'),
      Details: 'Administered rabies vaccine.'
    });
    this.procedures.set("661a3ed6-c2a8-4b55-a3e2-d2f11afedbd1", {
      Id: "661a3ed6-c2a8-4b55-a3e2-d2f11afedbd1",
      Type: 'Blood Test',
      Date: new Date('2023-03-10'),
      Details: 'Blood test for heartworm. Results are negative.'
    });

    let procedures = Array.from(this.procedures.keys());

    this.patients.set('1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p', { Id: '1a2b3c4d-5e6f-7g8h-9i0j-1k2l3m4n5o6p', Type: "Cat", Name: "Luna", Owner: "Karen", Procedures: procedures });
    this.patients.set('2b3c4d5e-6f7g8h-9i0j1k-2l3m4n5o6p7q', { Id: '2b3c4d5e-6f7g8h-9i0j1k-2l3m4n5o6p7q', Type: "Bird", Name: "Sky", Owner: "Leo", Procedures: procedures });
    this.patients.set('3c4d5e6f-7g8h9i-0j1k2l-3m4n5o6p7q8r', { Id: '3c4d5e6f-7g8h9i-0j1k2l-3m4n5o6p7q8r', Type: "Dog", Name: "Bella", Owner: "Mia", Procedures: procedures });
    this.patients.set('4d5e6f7g-8h9i0j-1k2l3m-4n5o6p7q8r9s', { Id: '4d5e6f7g-8h9i0j-1k2l3m-4n5o6p7q8r9s', Type: "Cat", Name: "Oliver", Owner: "Nina", Procedures: procedures });
    this.patients.set('5e6f7g8h-9i0j1k-2l3m4n-5o6p7q8r9s0t', { Id: '5e6f7g8h-9i0j1k-2l3m4n-5o6p7q8r9s0t', Type: "Bird", Name: "Kiwi", Owner: "Oscar", Procedures: procedures });
    this.patients.set('6f7g8h9i-0j1k2l-3m4n5o-6p7q8r9s0t1u', { Id: '6f7g8h9i-0j1k2l-3m4n5o-6p7q8r9s0t1u', Type: "Dog", Name: "Daisy", Owner: "Paul", Procedures: procedures });
    this.patients.set('7g8h9i0j-1k2l3m-4n5o6p-7q8r9s0t1u2v', { Id: '7g8h9i0j-1k2l3m-4n5o6p-7q8r9s0t1u2v', Type: "Cat", Name: "Chloe", Owner: "Quinn", Procedures: procedures });
    this.patients.set('8h9i0j1k-2l3m4n-5o6p7q-8r9s0t1u2v3w', { Id: '8h9i0j1k-2l3m4n-5o6p7q-8r9s0t1u2v3w', Type: "Rabbit", Name: "Coco", Owner: "Rose", Procedures: procedures });
    this.patients.set('9i0j1k2l-3m4n5o-6p7q8r-9s0t1u2v3w4x', { Id: '9i0j1k2l-3m4n5o-6p7q8r-9s0t1u2v3w4x', Type: "Dog", Name: "Molly", Owner: "Steve", Procedures: procedures });
    this.patients.set('0j1k2l3m-4n5o6p-7q8r9s-0t1u2v3w4x5y', { Id: '0j1k2l3m-4n5o6p-7q8r9s-0t1u2v3w4x5y', Type: "Cat", Name: "Tiger", Owner: "Tina", Procedures: procedures });
    this.patients.set('1k2l3m4n-5o6p7q-8r9s0t-1u2v3w4x5y6z', { Id: '1k2l3m4n-5o6p7q-8r9s0t-1u2v3w4x5y6z', Type: "Bird", Name: "Sunny", Owner: "Uma", Procedures: procedures });
    this.patients.set('2l3m4n5o-6p7q8r-9s0t1u-2v3w4x5y6z7a', { Id: '2l3m4n5o-6p7q8r-9s0t1u-2v3w4x5y6z7a', Type: "Dog", Name: "Bailey", Owner: "Vince", Procedures: procedures });
    this.patients.set('3m4n5o6p-7q8r9s-0t1u2v-3w4x5y6z7a8b', { Id: '3m4n5o6p-7q8r9s-0t1u2v-3w4x5y6z7a8b', Type: "Cat", Name: "Loki", Owner: "Wendy", Procedures: procedures });
    this.patients.set('4n5o6p7q-8r9s0t-1u2v3w-4x5y6z7a8b9c', { Id: '4n5o6p7q-8r9s0t-1u2v3w-4x5y6z7a8b9c', Type: "Rabbit", Name: "Snowball", Owner: "Xander", Procedures: procedures });
    this.patients.set('5o6p7q8r-9s0t1u-2v3w4x-5y6z7a8b9c0d', { Id: '5o6p7q8r-9s0t1u-2v3w4x-5y6z7a8b9c0d', Type: "Dog", Name: "Lucy", Owner: "Yara", Procedures: procedures });
    this.patients.set('6p7q8r9s-0t1u2v-3w4x5y-6z7a8b9c0d1e', { Id: '6p7q8r9s-0t1u2v-3w4x5y-6z7a8b9c0d1e', Type: "Cat", Name: "Nala", Owner: "Zane", Procedures: procedures });
    this.patients.set('7q8r9s0t-1u2v3w-4x5y6z-7a8b9c0d1e2f', { Id: '7q8r9s0t-1u2v3w-4x5y6z-7a8b9c0d1e2f', Type: "Bird", Name: "Peach", Owner: "Amy", Procedures: procedures });
    this.patients.set('8r9s0t1u-2v3w4x-5y6z7a-8b9c0d1e2f3g', { Id: '8r9s0t1u-2v3w4x-5y6z7a-8b9c0d1e2f3g', Type: "Dog", Name: "Cooper", Owner: "Ben", Procedures: procedures });
    this.patients.set('9s0t1u2v-3w4x5y-6z7a8b-9c0d1e2f3g4h', { Id: '9s0t1u2v-3w4x5y-6z7a8b-9c0d1e2f3g4h', Type: "Cat", Name: "Leo", Owner: "Cathy", Procedures: procedures });
    this.patients.set('0t1u2v3w-4x5y6z-7a8b9c-0d1e2f3g4h5i', { Id: '0t1u2v3w-4x5y6z-7a8b9c-0d1e2f3g4h5i', Type: "Rabbit", Name: "Fluffy", Owner: "Dan", Procedures: procedures });
    this.patients.set('1u2v3w4x-5y6z7a-8b9c0d-1e2f3g4h5i6j', { Id: '1u2v3w4x-5y6z7a-8b9c0d-1e2f3g4h5i6j', Type: "Dog", Name: "Toby", Owner: "Ella", Procedures: procedures });
    this.patients.set('2v3w4x5y-6z7a8b-9c0d1e-2f3g4h5i6j7k', { Id: '2v3w4x5y-6z7a8b-9c0d1e-2f3g4h5i6j7k', Type: "Cat", Name: "Milo", Owner: "Fred", Procedures: procedures });
    this.patients.set('3w4x5y6z-7a8b9c-0d1e2f-3g4h5i6j7k8l', { Id: '3w4x5y6z-7a8b9c-0d1e2f-3g4h5i6j7k8l', Type: "Bird", Name: "Blue", Owner: "Gina", Procedures: procedures });
    this.patients.set('4x5y6z7a-8b9c0d-1e2f3g-4h5i6j7k8l9m', { Id: '4x5y6z7a-8b9c0d-1e2f3g-4h5i6j7k8l9m', Type: "Dog", Name: "Jack", Owner: "Holly", Procedures: procedures });
    this.patients.set('5y6z7a8b-9c0d1e-2f3g4h-5i6j7k8l9m0n', { Id: '5y6z7a8b-9c0d1e-2f3g4h-5i6j7k8l9m0n', Type: "Cat", Name: "Oscar", Owner: "Ian", Procedures: procedures });
    this.patients.set('6z7a8b9c-0d1e2f-3g4h5i-6j7k8l9m0n1o', { Id: '6z7a8b9c-0d1e2f-3g4h5i-6j7k8l9m0n1o', Type: "Rabbit", Name: "Patches", Owner: "Jill", Procedures: procedures });
  }
  constructor() {
    this.generateMockData();
  }


  public getPatients() {
    return Array.from(this.patients.values())
  }

  public getPatient(key: string): Patient | undefined {
    const patient = this.patients.get(key);
    return patient ? { ...patient } : undefined;
  }

  public getProcedure(key: string): Procedure | undefined {
    const procedure = this.procedures.get(key);
    return procedure ? { ...procedure } : undefined;
  }

  public getProcedures(keys: string[]): Procedure[] {
    var procedures: Procedure[] = [];
    var service = this;
    keys.forEach(function(key) {
      const procedure = service.procedures.get(key);
      if (procedure) {
        procedures.push(procedure);
      }
    })

    return procedures;
  }

  public addPatient(patient: Patient) {
    patient.Id = uuid.v4();
    this.patients.set(patient.Id, patient);
  }

  public addProcedure(patientId: string, procedure: Procedure) {
    procedure.Id = uuid.v4();
    this.procedures.set(procedure.Id, procedure);
    let patient = this.patients.get(patientId);
    patient?.Procedures.unshift(procedure.Id);
  }

  public updatePatient(patient: Patient) {
    this.patients.set(patient.Id, patient);
  }


  public updateProcedure(procedure: Procedure) {
    this.procedures.set(procedure.Id, procedure);
  }

  public getTypes(): String[] {
    return this.types;
  }
}
