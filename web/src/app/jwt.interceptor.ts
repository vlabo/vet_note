import { HttpInterceptorFn } from '@angular/common/http';
import { inject } from '@angular/core';
import { AuthService } from './auth.service';

export const jwtInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const currentUser = authService.currentUserValue;

  if (currentUser && currentUser.authdata) {
    req = req.clone({
      setHeaders: {
        Authorization: `Basic ${currentUser.authdata}`
      }
    });
  }

  return next(req);
};