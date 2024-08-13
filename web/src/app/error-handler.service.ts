import { HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, throwError } from 'rxjs';
import { AuthService } from './auth.service';
import { ToastController } from '@ionic/angular';

@Injectable({
  providedIn: 'root'
})
export class ErrorHandlerService {

  constructor(
    private toastController: ToastController,
    private authService: AuthService
  ) {}

  private async presentToast(message: string) {
    const toast = await this.toastController.create({
      message: message,
      duration: 3000,
      position: 'bottom'
    });
    toast.present();
  }

  public handleError = (error: HttpErrorResponse): Observable<never> => {
    if (error.status == 401) {
      this.authService.logout();
      window.location.reload();
    } else {
      // Show a toast notification to the user
      this.presentToast('Възникна грешка: ' + error.status + " " + error.statusText);
    }
    // Return an observable with a user-facing error message
    return throwError(() => new Error('Възникна грешка: ' + error.status + " " + error.statusText));
  }
}
