import { Component, OnInit } from '@angular/core';
import { ItemReorderEventDetail } from '@ionic/angular';
import { Location } from '@angular/common';

import { addIcons } from "ionicons";
import { arrowBack, add, close, checkmark, trash } from 'ionicons/icons';
import { PatientsService } from '../patients.service';
import { AuthService } from '../auth.service';
import { ViewSetting } from '../types';
@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss'],
})
export class SettingsComponent implements OnInit {

  patientTypes: ViewSetting[] = [];
  procedureTypes: ViewSetting[] = [];

  newPatient?: ViewSetting;
  newProcedure?: ViewSetting;

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
    const reorderedItems = ev.detail.complete(this.patientTypes);

    // Update the index field of each item
    reorderedItems.forEach((item: ViewSetting, index: number) => {
      item.index = index;
    });

    // Update the items array with the reordered items
    this.patientTypes = reorderedItems;
    this.patientsService.updateSettings(this.patientTypes).subscribe({});
  }

  handleReorderProcedures(ev: CustomEvent<ItemReorderEventDetail>) {
    const reorderedItems = ev.detail.complete(this.procedureTypes);

    // Update the index field of each item
    reorderedItems.forEach((item: ViewSetting, index: number) => {
      item.index = index;
    });

    // Update the items array with the reordered items
    this.procedureTypes = reorderedItems;
    this.patientsService.updateSettings(this.procedureTypes).subscribe({});
  }

  goBack() {
    this.location.back();
  }

  cancelAddPatientType() {
    this.newPatient = undefined;
  }

  cancelAddProcedureType() {
    this.newProcedure = undefined;
  }

  openNewPatient() {
    let index = 0;
    if(this.patientTypes.length > 0) {
      index = this.patientTypes[this.patientTypes.length - 1].index + 1;
    }
    this.newPatient = {
      type: 'PatientType',
      value: "",
      index: index,
    };
  }

  pushNewPatientType() {
    if(!this.newPatient) {
      return;
    }
    this.patientTypes.push(this.newPatient);
    this.patientsService.updateSetting(this.newPatient).subscribe({});
    this.newPatient = undefined;
  }

  openNewProcedure() {
    let index = 0;
    if(this.procedureTypes.length > 0) {
      index = this.procedureTypes[this.procedureTypes.length - 1].index + 1;
    }
    this.newProcedure = {
      type: 'ProcedureType',
      value: "",
      index: index,
    };
  }

  pushNewProcedureType() {
    if(!this.newProcedure) {
      return;
    }
    this.procedureTypes.push(this.newProcedure);
    this.patientsService.updateSetting(this.newProcedure).subscribe({});
    this.newProcedure = undefined;
  }

  deletePatientType(index: number) {
    let deleted = this.patientTypes.splice(index, 1)[0];
    this.patientsService.deleteSetting(deleted).subscribe({});
  }

  deleteProcedureType(index: number) {
    let deleted = this.procedureTypes.splice(index, 1)[0];
    this.patientsService.deleteSetting(deleted).subscribe({});
  }

  logout() {
    this.authService.logout();
    window.location.reload();
  }
}
