import { Injectable } from '@angular/core';
import * as uuid from 'uuid';
import { Patient } from '../../server/bindings/Patient';
import { Procedure } from '../../server/bindings/Procedure';
import { ListPatient } from '../../server/bindings/ListPatient';


@Injectable({
  providedIn: 'root'
})
export class PatientsService {
  private patients = new Map<string, Patient>()
  private procedures = new Map<string, Procedure>()
  private types: String[] = ["Куче", "Котка", "Птица", "Заек"];

  constructor() { }

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

  public async getPatientList(): Promise<ListPatient[]> {
    try {
      const response = await fetch('http://localhost:8080/v1/patient-list');
      if (!response.ok) {
        throw new Error(`Error: ${response.status}`);
      }
      const list: ListPatient[] = await response.json();
      return list;
    } catch (error) {
      console.error("Failed to fetch patient list:", error);
      return []; // Return an empty array in case of error
    }
  }

  public addPatient(patient: Patient) {
    patient.id = uuid.v4();
    this.patients.set(patient.id, patient);
  }

  public addProcedure(patientId: string, procedure: Procedure) {
    procedure.id = uuid.v4();
    this.procedures.set(procedure.id, procedure);
  }

  public updatePatient(view: Patient) {
    // var patient = this.patients.get(view.Id);
    // if (!patient) {
    //   patient = {};
    //   patient.id = uuid.v4();
    //   this.patients.set(patient.Id, patient);
    // }

    //    // Iterate over the specified fields
    // ViewPatient.fieldsToCheck.forEach(key => {
    //   // @ts-ignore
    //   if (patient[key] !== view[key]) {
    //     // @ts-ignore
    //     patient[key] = view[key];
    //   }
    // });

    // // Update the LastModified date
    // patient.LastModified = new Date();
  }

  public updateProcedure(procedure: Procedure) {
    this.procedures.set(procedure.id, procedure);
  }

  public getTypes(): String[] {
    return this.types;
  }
}
