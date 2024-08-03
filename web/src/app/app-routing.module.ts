import { NgModule } from '@angular/core';
import { PreloadAllModules, RouterModule, Routes } from '@angular/router';
import { MainComponent } from './main/main.component';
import { SettingsComponent } from './settings/settings.component';
import { EditPatientComponent } from './edit-patient/edit-patient.component';
import { PatientComponent } from './patient/patient.component';
import { ProcedureComponent } from './procedure/procedure.component';
import { LoginComponent } from './login/login.component';
import { authGuard } from './auth.guard';

const routes: Routes = [
  { path: '', component: MainComponent, canActivate: [authGuard] },
  { path: 'login', component: LoginComponent },
  { path: 'settings',component: SettingsComponent, canActivate: [authGuard] },
  { path: 'patient/new', component: EditPatientComponent, canActivate: [authGuard] },
  { path: 'patient/:id', component: PatientComponent, canActivate: [authGuard] },
  { path: 'patient/:id/edit', component: EditPatientComponent, canActivate: [authGuard] },
  { path: 'procedure/:patientId/new', component: ProcedureComponent, data: { newMode: true }, canActivate: [authGuard] }, 
  { path: 'procedure/:procedureId', component: ProcedureComponent , data: { newMode: false }, canActivate: [authGuard] },
];

@NgModule({
  imports: [
    RouterModule.forRoot(routes, { preloadingStrategy: PreloadAllModules })
  ],
  exports: [RouterModule]
})
export class AppRoutingModule { }
