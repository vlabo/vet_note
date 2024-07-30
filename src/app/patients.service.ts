import { Injectable } from '@angular/core';
import { ViewListPatient, ViewPatient, ViewProcedure } from './types';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';


@Injectable({
  providedIn: 'root',
})
export class PatientsService {
  private Url: string = "http://192.168.88.11:8080/v1";

  constructor(private http: HttpClient) { }

  // Main list
  public getPatientList(): Observable<ViewListPatient[]> {
      return this.http.get<ViewListPatient[]>(`${this.Url}/patient-list`);
  }

  // Patient
  public getPatient(key: string): Observable<ViewPatient> {
      return this.http.get<ViewPatient>(`${this.Url}/patient/` + key);
  }

  public updatePatient(view: ViewPatient): Observable<ViewPatient> {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.http.post<ViewPatient>(`${this.Url}/patient`, view, { headers });
  }

  public deletePatient(id: string) {
    return this.http.delete(`${this.Url}/patient/${id}`);
  }

  // Procedure
  public getProcedure(key: string): Observable<ViewProcedure> {
      return this.http.get<ViewProcedure>(`${this.Url}/procedure/` + key);
  }

  public updateProcedure(patientId: string, procedure: ViewProcedure) {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.http.post<ViewProcedure>(`${this.Url}/procedure/` + patientId, procedure, { headers });
  }

  public deleteProcedure(id: string) {
    return this.http.delete(`${this.Url}/procedure/` + id);
  }

  // Settings
  public getPatientTypes(): Observable<string[]> {
      return this.http.get<string[]>(`${this.Url}/patient-types`);
  }

  public updatePatientTypes(types: string[]): Observable<string[]> {
      return this.http.post<string[]>(`${this.Url}/patient-types`, types);
  }

  public getProcedureTypes(): Observable<string[]> {
      return this.http.get<string[]>(`${this.Url}/procedure-types`);
  }

  public updateProcedureTypes(types: string[]): Observable<string[]> {
      return this.http.post<string[]>(`${this.Url}/procedure-types`, types);
  }
}
