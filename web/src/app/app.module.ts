import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { RouteReuseStrategy } from '@angular/router';

import { IonicModule, IonicRouteStrategy } from '@ionic/angular';

import { AppComponent } from './app.component';
import { AppRoutingModule } from './app-routing.module';
import { provideHttpClient, withInterceptorsFromDi } from '@angular/common/http';
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
  ],
  imports: [BrowserModule, IonicModule.forRoot(), FormsModule, AppRoutingModule, FontAwesomeModule],
  providers: [
    { provide: RouteReuseStrategy, useClass: IonicRouteStrategy },
    provideHttpClient(withInterceptorsFromDi()),
    { provide: LocationStrategy, useClass: HashLocationStrategy }
  ],
  bootstrap: [AppComponent],
})
export class AppModule { }
