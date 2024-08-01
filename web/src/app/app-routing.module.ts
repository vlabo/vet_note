import { NgModule } from '@angular/core';
import { PreloadAllModules, RouterModule, Routes } from '@angular/router';
import { MainComponent } from './main/main.component';
import { SettingsComponent } from './settings/settings.component';
import { EditPatientComponent } from './edit-patient/edit-patient.component';
import { PatientComponent } from './patient/patient.component';
import { ProcedureComponent } from './procedure/procedure.component';

const routes: Routes = [
  { path: '', component: MainComponent },
  { path: 'settings',component: SettingsComponent },
  { path: 'patient/new', component: EditPatientComponent },
  { path: 'patient/:id', component: PatientComponent },
  { path: 'patient/:id/edit', component: EditPatientComponent },
  { path: 'procedure/:patientId/new', component: ProcedureComponent, data: { newMode: true } }, 
  { path: 'procedure/:procedureId', component: ProcedureComponent , data: { newMode: false } },
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, { preloadingStrategy: PreloadAllModules })
  ],
  exports: [RouterModule]
})
export class AppRoutingModule { }
