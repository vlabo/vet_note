import { NgModule, isDevMode } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { RouteReuseStrategy } from '@angular/router';

import { IonicModule, IonicRouteStrategy } from '@ionic/angular';

import { AppComponent } from './app.component';
import { AppRoutingModule } from './app-routing.module';
import { HTTP_INTERCEPTORS, provideHttpClient, withInterceptors, withInterceptorsFromDi } from '@angular/common/http';
import { MainComponent } from './main/main.component';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { SettingsComponent } from './settings/settings.component';
import { FormsModule } from '@angular/forms';
import { PatientComponent } from './patient/patient.component';
import { EditPatientComponent } from './edit-patient/edit-patient.component';
import { ProcedureComponent } from './procedure/procedure.component';
import { DatePickerModalComponent } from './date-picker-modal/date-picker-modal.component';
import { ViewListPatientComponentComponent } from './view-list-patient-component/view-list-patient-component.component';
import { HashLocationStrategy, LocationStrategy } from '@angular/common';
import { ServiceWorkerModule } from '@angular/service-worker';
import { LoginComponent } from './login/login.component';
import { jwtInterceptor } from './jwt.interceptor';

@NgModule({
  declarations: [
    AppComponent,
    MainComponent,
    SettingsComponent,
    PatientComponent,
    EditPatientComponent,
    ProcedureComponent,
    DatePickerModalComponent,
    ViewListPatientComponentComponent,
    LoginComponent,
  ],
  imports: [BrowserModule, IonicModule.forRoot(), FormsModule, AppRoutingModule, FontAwesomeModule, ServiceWorkerModule.register('ngsw-worker.js', {
  enabled: !isDevMode(),
  // Register the ServiceWorker as soon as the application is stable
  // or after 30 seconds (whichever comes first).
  registrationStrategy: 'registerWhenStable:30000'
})],
  providers: [
    { provide: RouteReuseStrategy, useClass: IonicRouteStrategy },
    provideHttpClient(withInterceptors([jwtInterceptor])),
    { provide: LocationStrategy, useClass: HashLocationStrategy },
  ],
  bootstrap: [AppComponent],
})
export class AppModule { }
