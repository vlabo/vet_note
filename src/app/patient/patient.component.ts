import { Component, OnDestroy, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { PatientsService } from '../patients.service';
import { CommonModule } from '@angular/common';
import { IonModal, IonicModule } from '@ionic/angular';
import { FormsModule } from '@angular/forms';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { faCakeCandles, faUser, faCalendarDays, faMarsAndVenus, faMicrochip, faMars, faVenus, faPhone, faWeightHanging, faCheck, faX } from '@fortawesome/free-solid-svg-icons';

import { add, create, checkmark, close, arrowBack, chevronForward, paw, person, calendar, heartHalf } from 'ionicons/icons';
import { addIcons } from "ionicons";
import { Location } from '@angular/common';
import { formatDate } from '@angular/common';
import { ViewPatient, ViewProcedure } from '../types';

@Component({
  selector: 'app-patient',
  templateUrl: './patient.component.html',
  styleUrls: ['./patient.component.scss'],
  imports: [IonicModule, CommonModule, FormsModule, FontAwesomeModule],
  standalone: true,
})
export class PatientComponent implements OnInit, OnDestroy {

  // Icons
  cake = faCakeCandles;
  user = faUser;
  calendar = faCalendarDays;
  gender = faMarsAndVenus;
  microchip = faMicrochip;
  mars = faMars;
  venus = faVenus;
  phone = faPhone;
  weight = faWeightHanging;
  check = faCheck;
  xIcon = faX;

  patient: ViewPatient | undefined = undefined;
  procedures: ViewProcedure[] = [];
  originalPatient: any;
  isViewProcedure = false;

  selectedProcedure: ViewProcedure | null = null;

  @ViewChild('procedureDetailsModal', { static: true }) modal!: IonModal;

  constructor(private route: ActivatedRoute,
    private router: Router,
    private location: Location,
    private patientService: PatientsService,
  ) {
    addIcons({
      "add": add,
      "create": create,
      "checkmark": checkmark,
      "close": close,
      "arrow-back": arrowBack,
      "chevron-forward": chevronForward,
      "paw": paw,
      "person": person,
      "calendar": calendar,
      "heart-half": heartHalf,
    });
  }

  ngOnInit() {
    this.route.paramMap.subscribe(async params => {
      let id = params.get('id');
      if (id === null) {
        return;
      }
      this.patient = await this.patientService.getPatient(id);
      console.log(this.patient);
      // this.procedures = this.patientService.getProcedures(this.patient!.Procedures);
    });
  }

  ngOnDestroy() {
    this.isViewProcedure = false;
    this.updatePatient();
  }

  openEdit() {
    this.router.navigate(["/patient", this.patient?.id, "edit"]);
  }

  goBack(): void {
    this.location.back();
  }

  setOpen(isOpen: boolean) {
    this.isViewProcedure = isOpen;
  }

  openViewProcedure(procedure: ViewProcedure) {
    if (this.patient) {
      this.router.navigate(['/procedure', procedure.id]);
    }
  }

  createProcedure(): void {
    if (this.patient) {
      this.router.navigate(['/procedure', this.patient.id, 'new']);
      this.isViewProcedure = false;
    }
  }

  formatBirthDate(data: Date): string {
    return formatDate(data, "dd.MM.yyyy", "en-US");
  }

  updatePatient() {
    if (this.patient) {
      this.patientService.updatePatient(this.patient).subscribe({});
    }
  }

  getAge(): string {
    let now = new Date()
    let date = new Date(this.patient!.birthDate)
    let years = now.getFullYear() - date.getFullYear();
    let months = now.getMonth() - date.getMonth();

    // Adjust if the month difference is negative
    if (months < 0) {
      years--;
      months += 12;
    }
    let result = years + "г";
    if (months > 0) {
      result += " " + months + "м"
    }
    return result;
  };
}
