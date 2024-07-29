import { Injectable } from '@angular/core';
import { ViewListPatient, ViewPatient, ViewProcedure } from './types';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';


@Injectable({
  providedIn: 'root',
})
export class PatientsService {
  private types: String[] = ["Куче", "Котка", "Птица", "Заек"];

  private Url: string = "http://localhost:8080/v1";

  constructor(private http: HttpClient) { }

  public getPatient(key: string): Observable<ViewPatient> {
      return this.http.get<ViewPatient>(`${this.Url}/patient/` + key);
  }

  public getProcedure(key: string): Observable<ViewProcedure> {
      return this.http.get<ViewProcedure>(`${this.Url}/procedure/` + key);
  }

  public getPatientList(): Observable<ViewListPatient[]> {
      return this.http.get<ViewListPatient[]>(`${this.Url}/patient-list`);
  }

  public updatePatient(view: ViewPatient): Observable<ViewPatient> {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.http.post<ViewPatient>(`${this.Url}/patient`, view, { headers });
  }

  public deletePatient(id: string) {
    return this.http.delete(`${this.Url}/patient/${id}`);
  }

  public updateProcedure(patientId: string, procedure: ViewProcedure) {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.http.post<ViewProcedure>(`${this.Url}/procedure/` + patientId, procedure, { headers });
  }
  public deleteProcedure(id: string) {
    return this.http.delete(`${this.Url}/procedure/` + id);
  }

  public getTypes(): String[] {
    return this.types;
  }
}
