import { Component, OnDestroy, OnInit, TemplateRef, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Patient, PatientsService, Procedure } from '../patients.service';
import { CommonModule } from '@angular/common';
import { IonModal, IonicModule } from '@ionic/angular';
import { FormsModule } from '@angular/forms';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { faCakeCandles, faUser, faCalendarDays, faMarsAndVenus, faMicrochip, faMars, faVenus, faPhone } from '@fortawesome/free-solid-svg-icons';

import { add, create, checkmark, close, arrowBack, chevronForward, paw, person, calendar, heartHalf } from 'ionicons/icons';
import { addIcons } from "ionicons";
import { Location } from '@angular/common';
import { formatDate } from '@angular/common';

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


  patient: Patient | undefined = undefined;
  procedures: Procedure[] = [];
  originalPatient: any;
  isViewProcedure = false;

  selectedProcedure: Procedure = new Procedure();

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
    this.route.paramMap.subscribe(params => {
      let id = params.get('id');
      if (id === null) {
        return;
      }
      this.patient = this.patientService.getPatient(id);
      this.procedures = this.patientService.getProcedures(this.patient!.Procedures);
    });
  }

  ngOnDestroy() {
    this.isViewProcedure = false;
  }

  openEdit() {
    this.router.navigate(["/patient", this.patient?.Id, "edit"]);
  }

  goBack(): void {
    this.location.back();
  }

  setOpen(isOpen: boolean) {
    this.isViewProcedure = isOpen;
  }

  openViewProcedure(procedure: Procedure) {
    if (this.patient) {
      this.router.navigate(['/procedure', procedure.Id]);
    }
  }

  createProcedure(): void {
    if (this.patient) {
      this.router.navigate(['/procedure', this.patient.Id, 'new']);
      this.isViewProcedure = false;
    }
  }

  formatBirthDate(data: Date): string {
    return formatDate(data, "dd.MM.yyyy", "en-US");
  }


  getAge(): string {
    let now = new Date()
    let years = now.getFullYear() - this.patient!.BirthDate.getFullYear();
    let months = now.getMonth() - this.patient!.BirthDate.getMonth();

    // Adjust if the month difference is negative
    if (months < 0) {
      years--;
      months += 12;
    }
    let result = years + "г. ";
    if(months > 0) {
      result += months + "м."
    }
    return result; 
  };
}
