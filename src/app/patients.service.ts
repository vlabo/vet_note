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

  private Url: string = "http://localhost:8080/v1";

  constructor(private http: HttpClient) { }

  public getPatient(key: string): Observable<ViewPatient> {
      return this.http.get<ViewPatient>(this.Url + '/patient/' + key);
  }

  public getProcedure(key: string): Observable<ViewProcedure> {
      return this.http.get<ViewProcedure>(this.Url + '/procedure/' + key);
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

  public getPatientList(): Observable<ViewListPatient[]> {
      return this.http.get<ViewListPatient[]>(this.Url + '/patient-list');
  }

  public addPatient(patient: ViewPatient) {
    this.patients.set(patient.id, patient);
  }

  public updatePatient(view: ViewPatient): Observable<ViewPatient> {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.http.post<ViewPatient>('http://localhost:8080/v1/patient', view, { headers });
  }

  public deletePatient(id: string) {
    return this.http.delete('http://localhost:8080/v1/patient/' + id);
  }

  public updateProcedure(patientId: string, procedure: ViewProcedure) {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.http.post<ViewProcedure>('http://localhost:8080/v1/procedure/' + patientId, procedure, { headers });
  }
  public deleteProcedure(id: string) {
    return this.http.delete('http://localhost:8080/v1/procedure/' + id);
  }

  public getTypes(): String[] {
    return this.types;
  }
}
