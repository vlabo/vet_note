import { Injectable } from '@angular/core';
import { ViewListPatient, ViewPatient, ViewProcedure } from './types';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';


@Injectable({
  providedIn: 'root',
})
export class PatientsService {
  private patients = new Map<string, ViewPatient>()
  private procedures = new Map<string, ViewProcedure>()
  private types: String[] = ["Куче", "Котка", "Птица", "Заек"];

  constructor(private http: HttpClient) { }

  public async getPatient(key: string): Promise<ViewPatient | undefined> {
    console.log("key", key);
    try {
      const response = await fetch('http://localhost:8080/v1/patient/' + key);
      if (!response.ok) {
        throw new Error(`Error: ${response.status}`);
      }
      const patient: ViewPatient = await response.json();
      return patient;
    } catch (error) {
      console.error("Failed to fetch patient list:", error);
      return undefined;
    }
  }

  public getProcedure(key: string): ViewProcedure | undefined {
    const procedure = this.procedures.get(key);
    return procedure ? { ...procedure } : undefined;
  }

  public getProcedures(keys: string[]): ViewProcedure[] {
    var procedures: ViewProcedure[] = [];
    var service = this;
    keys.forEach(function(key) {
      const procedure = service.procedures.get(key);
      if (procedure) {
        procedures.push(procedure);
      }
    })

    return procedures;
  }

  public async getPatientList(): Promise<ViewListPatient[]> {
    try {
      const response = await fetch('http://localhost:8080/v1/patient-list');
      if (!response.ok) {
        throw new Error(`Error: ${response.status}`);
      }
      const list: ViewListPatient[] = await response.json();
      return list;
    } catch (error) {
      console.error("Failed to fetch patient list:", error);
      return []; // Return an empty array in case of error
    }
  }

  public addPatient(patient: ViewPatient) {
    this.patients.set(patient.id, patient);
  }

  public addProcedure(patientId: string, procedure: ViewProcedure) {
    this.procedures.set(procedure.id, procedure);
  }

  public updatePatient(view: ViewPatient): Observable<ViewPatient> {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    console.log("patient", view);
    return this.http.post<ViewPatient>('http://localhost:8080/v1/patient', view, { headers });
  }

  public updateProcedure(procedure: ViewProcedure) {
    this.procedures.set(procedure.id, procedure);
  }

  public getTypes(): String[] {
    return this.types;
  }
}
