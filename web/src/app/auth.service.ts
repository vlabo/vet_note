import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, catchError, map, Observable, throwError } from 'rxjs';
import { ServerURL } from 'src/environments/environment';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private currentUserSubject: BehaviorSubject<any> = new BehaviorSubject(null);
  public currentUser: Observable<any> = new Observable();

  constructor(private http: HttpClient) {
    let user = localStorage.getItem('currentUser')
    if (user) {
      this.currentUserSubject = new BehaviorSubject<any>(JSON.parse(user));
      this.currentUser = this.currentUserSubject.asObservable();
    }
  }

  public get currentUserValue(): any {
    return this.currentUserSubject.value;
  }

  public login(username: string, password: string) {
    const headers = new HttpHeaders({
      Authorization: 'Basic ' + btoa(username + ':' + password)
    });

    return this.http.get<any>(`${ServerURL}/authenticate`, { headers: headers })
      .pipe(catchError(this.handleError),
        map(_ => {
          // Store user details and basic auth credentials in local storage
          let user: any = {};
          user.authdata = btoa(username + ':' + password);
          localStorage.setItem('currentUser', JSON.stringify(user));
          this.currentUserSubject.next(user);
          return user;
        }));
  }

  public logout() {
    // Remove user from local storage and set current user to null
    localStorage.removeItem('currentUser');
    this.currentUserSubject.next(null);
  }

  private handleError(error: HttpErrorResponse): Observable<never> {
    let errorMessage = 'Възникна неизвестна грешка!';
    if (error.error instanceof ErrorEvent) {
      // Грешка от страна на клиента или мрежата
      errorMessage = `Грешка: ${error.error.message}`;
    } else {
      // Грешка от страна на сървъра
      switch (error.status) {
        case 401:
          errorMessage = 'Невалидно потребителско име или парола.';
          break;
        case 403:
          errorMessage = 'Достъпът е отказан.';
          break;
        case 404:
          errorMessage = 'Услугата не е намерена.';
          break;
        case 500:
          errorMessage = 'Вътрешна грешка на сървъра.';
          break;
        default:
          errorMessage = `Код на грешката: ${error.status}\nСъобщение: ${error.message}`;
      }
    }
    // Връщане на observable с потребителско съобщение за грешка
    return throwError(() => new Error(errorMessage));
  }
}