import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { catchError, Observable, Subject, throwError } from 'rxjs';
import { ViewListPatient, ViewPatient, ViewProcedure, ViewSetting } from 'src/app/types';
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

  public deletePatient(id: number) {
    return this.http.delete(`${ServerURL}/patient/${id}`).pipe(catchError(this.handleError));
  }

  // Procedure
  public getProcedure(key: string): Observable<ViewProcedure> {
    return this.http.get<ViewProcedure>(`${ServerURL}/procedure/` + key).pipe(catchError(this.handleError));
  }

  public updateProcedure(procedure: ViewProcedure) {
    const headers = new HttpHeaders({ 'Content-Type': 'application/json' });
    return this.http.post<ViewProcedure>(`${ServerURL}/procedure`, procedure, { headers }).pipe(catchError(this.handleError));
  }

  public deleteProcedure(id: number) {
    return this.http.delete(`${ServerURL}/procedure/` + id).pipe(catchError(this.handleError));
  }

  // Settings
  public getPatientTypes(): Observable<ViewSetting[]> {
    return this.http.get<ViewSetting[]>(`${ServerURL}/patient-types`).pipe(catchError(this.handleError));
  }

  public getProcedureTypes(): Observable<ViewSetting[]> {
    return this.http.get<ViewSetting[]>(`${ServerURL}/procedure-types`).pipe(catchError(this.handleError));
  }

  public updateSettings(settings: ViewSetting[]): Observable<void> {
    return this.http.post<void>(`${ServerURL}/settings`, settings).pipe(catchError(this.handleError));
  }

  public updateSetting(setting: ViewSetting): Observable<void> {
    return this.http.post<void>(`${ServerURL}/setting`, setting).pipe(catchError(this.handleError));
  }

  public deleteSetting(setting: ViewSetting): Observable<any> {
    return this.http.delete(`${ServerURL}/setting/` + setting.id).pipe(catchError(this.handleError));
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
