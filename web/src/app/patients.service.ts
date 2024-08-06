import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { catchError, Observable, Subject, throwError } from 'rxjs';
import { ViewListPatient, ViewPatient, ViewProcedure } from 'src/app/types';
import { ServerURL } from 'src/environments/environment';
import { AuthService } from './auth.service';


@Injectable({
  providedIn: 'root',
})
export class PatientsService {
  private patientListSubject: Subject<ViewListPatient[]> = new Subject();

  constructor(private http: HttpClient, private authService: AuthService) {
    this.handleError = this.handleError.bind(this);
  }

  // Main list
  public triggerPatientListReload() {
    this.http.get<ViewListPatient[]>(`${ServerURL}/patient-list`).pipe(catchError(this.handleError)).subscribe({
      next: patientsList => {
        this.patientListSubject.next(patientsList);
      }
    })
  }

  public getPatientListObservable(): Observable<ViewListPatient[]> {
    return this.patientListSubject.asObservable();
  }

  // Patient
  public getPatient(key: string): Observable<ViewPatient> {
    return this.http.get<ViewPatient>(`${ServerURL}/patient/` + key).pipe(catchError(this.handleError));
  }

  public updatePatient(view: ViewPatient): Observable<ViewPatient> {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.http.post<ViewPatient>(`${ServerURL}/patient`, view, { headers }).pipe(catchError(this.handleError));
  }

  public deletePatient(id: string) {
    return this.http.delete(`${ServerURL}/patient/${id}`).pipe(catchError(this.handleError));
  }

  // Procedure
  public getProcedure(key: string): Observable<ViewProcedure> {
    return this.http.get<ViewProcedure>(`${ServerURL}/procedure/` + key).pipe(catchError(this.handleError));
  }

  public updateProcedure(patientId: string, procedure: ViewProcedure) {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.http.post<ViewProcedure>(`${ServerURL}/procedure/` + patientId, procedure, { headers }).pipe(catchError(this.handleError));
  }

  public deleteProcedure(id: string) {
    return this.http.delete(`${ServerURL}/procedure/` + id).pipe(catchError(this.handleError));
  }

  // Settings
  public getPatientTypes(): Observable<string[]> {
    return this.http.get<string[]>(`${ServerURL}/patient-types`).pipe(catchError(this.handleError));
  }

  public updatePatientTypes(types: string[]): Observable<string[]> {
    return this.http.post<string[]>(`${ServerURL}/patient-types`, types).pipe(catchError(this.handleError));
  }

  public getProcedureTypes(): Observable<string[]> {
    return this.http.get<string[]>(`${ServerURL}/procedure-types`).pipe(catchError(this.handleError));
  }

  public updateProcedureTypes(types: string[]): Observable<string[]> {
    return this.http.post<string[]>(`${ServerURL}/procedure-types`, types).pipe(catchError(this.handleError));
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    // Customize the error handling logic as needed
    console.error('An error occurred:', error);
    if (error.status == 401) {
      this.authService.logout();
      window.location.reload();
    }
    // Return an observable with a user-facing error message
    return throwError(() => 'Something bad happened; please try again later.');
  }
}
