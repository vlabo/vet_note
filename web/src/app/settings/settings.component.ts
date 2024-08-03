import { Component, OnInit } from '@angular/core';
import { ItemReorderEventDetail } from '@ionic/angular';
import { Location } from '@angular/common';

import { addIcons } from "ionicons";
import { arrowBack, add, close, checkmark, trash } from 'ionicons/icons';
import { PatientsService } from '../patients.service';
import { AuthService } from '../auth.service';
@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss'],
})
export class SettingsComponent implements OnInit {

  patientTypes: string[] = [];
  procedureTypes: string[] = [];


  addingNewPatient: boolean = false;
  newPatient: string = "";
  addingNewProcedure: boolean = false;
  newProcedure: string = "";

  constructor(private location: Location, private patientsService: PatientsService, private authService: AuthService) { 
    addIcons({"add": add});
  }

  ngOnInit() {
    addIcons({
      "arrow-back": arrowBack, "close": close, "checkmark": checkmark, "trash": trash
    });

    this.patientsService.getPatientTypes().subscribe({
      next: types => {
        if(types) {
          this.patientTypes = types;
        }
      }
    });

    this.patientsService.getProcedureTypes().subscribe({
      next: types => {
        if(types) {
          this.procedureTypes = types;
        }
      }
    });
  }

  handleReorderPatients(ev: CustomEvent<ItemReorderEventDetail>) {
    var toValue = this.patientTypes[ev.detail.to]
    this.patientTypes[ev.detail.to] = this.patientTypes[ev.detail.from];
    this.patientTypes[ev.detail.from] = toValue;
    this.patientsService.updatePatientTypes(this.patientTypes).subscribe({});
    ev.detail.complete();
  }

  handleReorderProcedures(ev: CustomEvent<ItemReorderEventDetail>) {
    var toValue = this.procedureTypes[ev.detail.to]
    this.procedureTypes[ev.detail.to] = this.procedureTypes[ev.detail.from];
    this.procedureTypes[ev.detail.from] = toValue;
    this.patientsService.updateProcedureTypes(this.procedureTypes).subscribe({});
    ev.detail.complete();
  }

  goBack() {
    this.location.back();
  }

  cancelAddPatient() {
    this.addingNewPatient = false;
    this.newPatient = "";
  }

  cancelAddProcedure() {
    this.addingNewProcedure = false;
    this.newProcedure = "";
  }

  pushNewPatient() {
    this.patientTypes.push(this.newPatient);
    this.patientsService.updatePatientTypes(this.patientTypes).subscribe({});
    this.addingNewPatient = false;
    this.newPatient = "";
  }

  pushNewProcedure() {
    this.procedureTypes.push(this.newProcedure);
    this.patientsService.updateProcedureTypes(this.procedureTypes).subscribe({});
    this.addingNewProcedure = false;
    this.newProcedure = "";
  }

  deletePatientType(index: number) {
    this.patientTypes.splice(index, 1);
    this.patientsService.updatePatientTypes(this.patientTypes).subscribe({});
  }

  deleteProcedureType(index: number) {
    this.procedureTypes.splice(index, 1);
    this.patientsService.updateProcedureTypes(this.procedureTypes).subscribe({});
  }

  logout() {
    this.authService.logout();
    window.location.reload();
  }
}
