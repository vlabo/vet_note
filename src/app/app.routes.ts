import { Routes } from '@angular/router';
import { PatientComponent } from './patient/patient.component';
import { MainComponent } from './main/main.component';
import { EditPatientComponent } from './edit-patient/edit-patient.component';
import { ProcedureComponent } from './procedure/procedure.component';
import { SettingsComponent } from './settings/settings.component';

export const routes: Routes = [
  { path: '', component: MainComponent },
  { path: 'patient/new', component: EditPatientComponent },
  { path: 'patient/:id', component: PatientComponent },
  { path: 'patient/:id/edit', component: EditPatientComponent },
  { path: 'procedure/:patientId/new', component: ProcedureComponent, data: { newMode: true } }, 
  { path: 'procedure/:procedureId', component: ProcedureComponent , data: { newMode: false } },
  { path: 'settings', component: SettingsComponent},
];
